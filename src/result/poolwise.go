package result

import (
	"log"
	"strconv"
	"sync"
	"time"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/pool"
)

// PoolsResult is a structure that contains information about pool wise result
type PoolsResult struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Rank       string  `json:"rank"`
	Total      int     `json:"total"`
	Attempted  float32 `json:"attempted"`
	Correct    float32 `json:"correct"`
	Accuracy   float32 `json:"accuracy"`
	Percentage float32 `json:"percentage"`
}

type PoolRank struct {
	Ranking int `json:"ranking"`
}

func AllPools() map[string]interface{} {

	var wg sync.WaitGroup

	mapd := make(map[string]interface{})

	db := config.DB

	pools := []PoolsResult{}
	db.Raw("select pool_id as id,count(distinct email) as total,(sum(pool_total_attempted)::decimal/greatest(count(distinct email),1)::decimal) as attempted, (sum(pool_total_correct)::decimal/greatest(count(distinct email),1)::decimal) as correct, (sum(pool_total_correct)::decimal/greatest(sum(pool_total_attempted),1)::decimal)*100 as accuracy, (sum(pool_total_correct)::decimal/greatest(sum(pool_total_question),1)::decimal)*100 as percentage from results.pool_analytical group by pool_id  order by accuracy DESC;").Find(&pools)

	wg.Add(len(pools))
	for i := 0; i < len(pools); i++ {
		go poolDetails(&pools[i], &wg)
	}
	wg.Wait()

	mapd["error"] = "false"
	mapd["pools"] = pools
	log.Println(mapd)
	return mapd
}

func poolDetails(pool *PoolsResult, wg *sync.WaitGroup) {
	name := make(chan string)
	rank := make(chan string)

	go poolName(pool.ID, name)
	go poolRank(pool.ID, rank)

	pool.Rank = <-rank
	pool.Name = <-name

	wg.Done()
}

func poolName(poolID string, ch chan string) {
	// fetch drive details
	details, err := pool.PoolDetails(poolID)
	if err != nil {
		ch <- ""
	} else {
		if details.Name == "" {
			ch <- poolID
		} else {
			ch <- details.Name
		}
	}
	return
}

func poolRank(poolID string, ch chan string) {
	db := config.DB

	rank := []PoolRank{}

	db.Raw("select ranking from results.pool_rank where pool_id='" + poolID + "' order by calculated_at desc limit 2;").Find(&rank)

	if len(rank) != 2 {
		ch <- "same"
	} else {
		if rank[0].Ranking == rank[1].Ranking {
			ch <- "same"
		} else if rank[0].Ranking < rank[1].Ranking {
			ch <- "down"
		} else {
			ch <- "up"
		}
	}
	return
}

//============================================================================================================================================================================

type PoolCandidate struct {
	Email            string  `json:"email"`
	Questions        int     `json:"questions"`
	Attempted        int     `json:"attempted"`
	Correct          int     `json:"correct"`
	Accuracy         float32 `json:"accuracy"`
	Percentage       float32 `json:"percentage"`
	TimeTaken        int     `json:"time_taken"`
	AverageTimeTaken float32 `json:"average_time_taken"`
}

type PoolComp struct {
	ID      string    `json:"id" gorm:"column:pool_id"`
	Total   int       `json:"total"`
	Twenty  int       `json:"twenty"`
	Forty   int       `json:"forty"`
	Sixty   int       `json:"sixty"`
	Eighty  int       `json:"eighty"`
	Hundred int       `json:"hundred"`
	Date    time.Time `json:"date" gorm:"column:calculated_at"`
}

func PoolCandidates(poolID string, pageStr string) map[string]interface{} {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	mapd := make(map[string]interface{})

	var wg sync.WaitGroup
	wg.Add(3)

	go poolCandidates(poolID, page, mapd, &wg)
	go poolComparison(poolID, mapd, &wg)
	go poolPie(poolID, mapd, &wg)

	wg.Wait()

	return mapd
}

func poolCandidates(poolID string, page int, mapd map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	db := config.DB
	//total candidates
	var total int64
	db.Table("results.pool_analytical").Select("count(distinct email)").Where("pool_id=?", poolID).Count(&total)

	cand := []PoolCandidate{}

	db.Raw("select email,sum(pool_total_question) as questions, sum(pool_total_attempted) as attempted, sum(pool_total_correct) as correct, (sum(pool_total_correct)::decimal/greatest(sum(pool_total_question),1)*100) as percentage, (sum(pool_total_correct)::decimal/greatest(sum(pool_total_attempted),1)*100) as accuracy from results.pool_analytical where pool_id='" + poolID + "' group by email order by accuracy DESC LIMIT 10 offset " + strconv.Itoa((page-1)*10) + ";").Find(&cand)

	mapd["candidates"] = cand
	mapd["total"] = total
}
func poolComparison(poolID string, mapd map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	db := config.DB

	compare := []PoolComp{}

	db.Raw("select * from results.pool_week where pool_id='" + poolID + "' order by calculated_at desc LIMIT 2;").Find(&compare)

	result := make([][]interface{}, 0)
	one := make([]interface{}, 0)
	one = append(one, "Date")
	one = append(one, "0-20%")
	one = append(one, "20-40%")
	one = append(one, "40-60%")
	one = append(one, "60-80%")
	one = append(one, "80-100%")
	result = append(result, one)
	for i := 0; i < len(compare); i++ {
		one = make([]interface{}, 0)
		one = append(one, compare[i].Date.Format("2 Jan"))
		one = append(one, compare[i].Twenty)
		one = append(one, compare[i].Forty)
		one = append(one, compare[i].Sixty)
		one = append(one, compare[i].Eighty)
		one = append(one, compare[i].Hundred)
		result = append(result, one)
	}

	mapd["compare"] = result
}
func poolPie(poolID string, mapd map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	db := config.DB

	pie := Pie{}

	db.Raw("select sum(pool_total_question)::decimal/count(distinct email) as total, sum(pool_total_attempted)::decimal/count(distinct email) as attempted, sum(pool_total_correct)::decimal/count(distinct email) as correct from results.pool_analytical where pool_id='" + poolID + "' group by pool_id;").Find(&pie)
	pie.NotAttempted = float32(pie.Total) - pie.Attempted
	pie.Wrong = pie.Attempted - pie.Correct

	mapd["pie"] = pie
}
