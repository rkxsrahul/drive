package result

import (
	"log"
	"strconv"
	"sync"
	"time"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/drive"
)

// DriveResult is a structure that contains information about drive wise result
type DriveResult struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Rank       string  `json:"rank"`
	Total      int     `json:"total"`
	Attempted  float32 `json:"attempted"`
	Correct    float32 `json:"correct"`
	Accuracy   float32 `json:"accuracy"`
	Percentage float32 `json:"percentage"`
	TimeTaken  float32 `json:"time_taken"`
}

type DriveRank struct {
	Ranking int `json:"ranking"`
}

func AllDrives() map[string]interface{} {
	mapd := make(map[string]interface{})
	var wg sync.WaitGroup

	db := config.DB

	drives := []DriveResult{}
	db.Raw("select drive_id as id,count(distinct email) as total,	(sum(pool_total_attempted)::decimal/greatest(count(distinct email),1)::decimal) as attempted, (sum(pool_total_correct)::decimal/greatest(count(distinct email),1)::decimal) as correct, (sum(pool_total_correct)::decimal/greatest(sum(pool_total_attempted),1)::decimal)*100 as accuracy, (sum(pool_total_correct)::decimal/greatest(sum(pool_total_question),1)::decimal)*100 as percentage from results.pool_analytical group by drive_id order by accuracy DESC;").Find(&drives)

	wg.Add(len(drives))
	for i := 0; i < len(drives); i++ {
		go driveDetails(&drives[i], &wg)
	}
	wg.Wait()

	mapd["error"] = "false"
	mapd["drives"] = drives
	return mapd
}

func driveDetails(drive *DriveResult, wg *sync.WaitGroup) {
	name := make(chan string)
	rank := make(chan string)
	time := make(chan float32)

	go driveName(drive.ID, name)
	go driveRank(drive.ID, rank)
	go driveTime(drive.ID, time)

	drive.Rank = <-rank
	drive.Name = <-name
	drive.TimeTaken = <-time

	wg.Done()
}

func driveName(driveID string, ch chan string) {
	// fetch drive details
	driveDetails, err := drive.View(driveID, "")
	if err != nil {
		ch <- ""
	} else {
		ch <- driveDetails.Name
	}
	return
}

func driveRank(driveID string, ch chan string) {
	db := config.DB

	rank := []DriveRank{}

	db.Raw("select ranking from results.drive_rank where drive_id=" + driveID + " order by calculated_at desc limit 2;").Find(&rank)

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

func driveTime(driveID string, ch chan float32) {
	db := config.DB

	type Time struct {
		TimeTaken float32 `json:"time_taken"`
	}
	var data Time
	db.Raw("select sum(time_taken)::decimal/count(email) as time_taken from user_sessions where drive_id=" + driveID + " group by drive_id;").Scan(&data)

	ch <- data.TimeTaken
	return
}

//============================================================================================================================================================================

type DriveCandidate struct {
	Email            string  `json:"email"`
	Questions        int     `json:"questions"`
	Attempted        int     `json:"attempted"`
	Correct          int     `json:"correct"`
	Accuracy         float32 `json:"accuracy"`
	Percentage       float32 `json:"percentage"`
	TimeTaken        int     `json:"time_taken"`
	AverageTimeTaken float32 `json:"average_time_taken"`
	Restart          int     `json:"restart"`
	Browser          int     `json:"browser"`
}

type DriveComp struct {
	ID      string    `json:"id" gorm:"column:drive_id"`
	Total   int       `json:"total"`
	Twenty  int       `json:"twenty"`
	Forty   int       `json:"forty"`
	Sixty   int       `json:"sixty"`
	Eighty  int       `json:"eighty"`
	Hundred int       `json:"hundred"`
	Date    time.Time `json:"date" gorm:"column:calculated_at"`
}

type Pie struct {
	Total        float32 `json:"total" gorm:"column:total"`
	Attempted    float32 `json:"attempted" gorm:"column:attempted"`
	NotAttempted float32 `json:"not_attempted"`
	Correct      float32 `json:"correct" gorm:"column:correct"`
	Wrong        float32 `json:"wrong"`
}

func DriveCandidates(driveID string, pageStr string) map[string]interface{} {
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	mapd := make(map[string]interface{})

	var wg sync.WaitGroup
	wg.Add(3)

	go fetchCandidates(driveID, page, mapd, &wg)
	go driveComparison(driveID, mapd, &wg)
	go drivePie(driveID, mapd, &wg)

	wg.Wait()

	log.Println(mapd)
	return mapd
}

func fetchCandidates(driveID string, page int, mapd map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	db := config.DB

	//total candidates
	var total int64
	db.Table("results.pool_analytical").Select("count(distinct email)").Where("drive_id=?", driveID).Count(&total)

	cand := []DriveCandidate{}

	db.Raw("select email,sum(pool_total_question) as questions, sum(pool_total_attempted) as attempted, sum(pool_total_correct) as correct, (sum(pool_total_correct)::decimal/greatest(sum(pool_total_question),1)*100) as percentage, (sum(pool_total_correct)::decimal/greatest(sum(pool_total_attempted),1)*100) as accuracy,time_taken, (time_taken::decimal/greatest(sum(pool_total_attempted),1)) as average_time_taken from results.pool_analytical where drive_id=" + driveID + " group by email,time_taken order by accuracy DESC LIMIT 10 offset " + strconv.Itoa((page-1)*10) + ";").Find(&cand)

	var cwg sync.WaitGroup
	cwg.Add(len(cand))

	for i := 0; i < len(cand); i++ {
		go candidateDetails(&cand[i], driveID, &cwg)
	}
	cwg.Wait()
	mapd["total"] = total
	mapd["candidates"] = cand
}

func candidateDetails(cand *DriveCandidate, driveID string, wg *sync.WaitGroup) {
	defer wg.Done()
	db := config.DB

	type CandMeta struct {
		TimeTaken int
		Restart   int
		Browser   int
	}

	var data CandMeta
	db.Raw("select time_taken, restart, browser from user_sessions where drive_id=" + driveID + " AND email='" + cand.Email + "';").Scan(&data)

	cand.TimeTaken = data.TimeTaken
	cand.Restart = data.Restart
	cand.Browser = data.Browser
	if cand.Attempted != 0 {
		cand.AverageTimeTaken = float32(data.TimeTaken) / float32(cand.Attempted)
	}
}

func driveComparison(driveID string, mapd map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	db := config.DB

	compare := []DriveComp{}

	db.Raw("select * from results.drive_week where drive_id=" + driveID + " order by calculated_at desc LIMIT 2;").Find(&compare)

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
func drivePie(driveID string, mapd map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	db := config.DB

	pie := Pie{}

	db.Raw("select sum(pool_total_question)::decimal/count(distinct email) as total, sum(pool_total_attempted)::decimal/count(distinct email) as attempted, sum(pool_total_correct)::decimal/count(distinct email) as correct from results.pool_analytical where drive_id=" + driveID + " group by drive_id;").Find(&pie)
	pie.NotAttempted = float32(pie.Total) - pie.Attempted
	pie.Wrong = pie.Attempted - pie.Correct

	mapd["pie"] = pie
}
