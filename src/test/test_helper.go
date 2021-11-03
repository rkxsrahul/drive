package test

import (
	"log"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/methods"
	"git.xenonstack.com/util/drive-portal/src/questions"
	"git.xenonstack.com/util/drive-portal/src/types"

	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"
)

func validateData(data bodyTypes.Test) (int, error) {
	//check on pools
	mapd := make(map[string]int, 0)
	for i := 0; i < len(data.Pools); i++ {
		// check pool exists
		code, err := questions.CheckPool(data.Pools[i].PoolId)
		if err != nil {
			return code, errors.New("Please enter existing pools only")
		}
		// check user had not submitted more question then no of question in a pool
		list, err := questions.ListQuestions(data.Pools[i].PoolId)
		if err != nil {
			return 500, err
		}
		if len(list) < data.Pools[i].NoOfQuestions {
			return 400, errors.New("Please enter valid number of questions")
		}
		// check test does not contains duplicate pools
		if mapd[data.Pools[i].PoolId] == 0 {
			mapd[data.Pools[i].PoolId]++
		} else {
			return 400, errors.New("Please don't enter duplicate pools in a single test")
		}
	}
	return 200, nil
}

func CreateTest(data bodyTypes.Test) (int, error) {
	//check data is valid or not
	code, err := validateData(data)
	if err != nil {
		zap.S().Error(err)
		return code, err
	}

	db := config.DB
	// convert test name in slug
	slug := methods.SlugOfString(data.Name)
	if slug == "" {
		return 400, errors.New("Please enter valid test name contains atleast one alphabet or number")
	}
	//convert duration to int
	duration, err := strconv.ParseInt(data.Duration, 10, 0)
	if err != nil {
		zap.S().Error(err)
		return 400, err
	}

	//saving data according to test structure
	testData := types.Tests{
		Id:       slug,
		Name:     data.Name,
		Date:     time.Now().Unix(),
		Duration: duration * 60,
	}
	//save data in db
	err = db.Create(&testData).Error
	if err != nil {
		zap.S().Error(err)
		return 500, err
	}

	// save test pools in db
	for i := 0; i < len(data.Pools); i++ {
		testPoolData := types.TestPool{
			PoolId:        data.Pools[i].PoolId,
			NoOfQuestions: data.Pools[i].NoOfQuestions,
			TestId:        slug,
		}
		db.Create(&testPoolData)
	}
	return 200, nil
}

func ListTest() ([]bodyTypes.TestList, error) {
	db := config.DB
	//fetch all tests from db
	list := []types.Tests{}
	db.Find(&list)
	//creating result array of list in a defined structure
	result := make([]bodyTypes.TestList, 0)
	for i := 0; i < len(list); i++ {
		one := bodyTypes.TestList{
			Id:       list[i].Id,
			Name:     list[i].Name,
			Duration: list[i].Duration / 60,
			Date:     list[i].Date,
		}
		// fech test pools name from dbs
		testPools := []types.TestPool{}
		// db.Where("test_id=?", list[i].Id).Order("updated_at DESC").Find(&testPools)
		db.Raw("select pool_id, no_of_questions from test_pools where test_id= '" + list[i].Id + "' order by updated_at desc").Scan(&testPools)
		count := 0
		pools := make([]string, 0)
		for j := 0; j < len(testPools); j++ {
			// count total no. of questions in a test.
			count += testPools[j].NoOfQuestions
			// fetch pool name from db
			pool := []types.Pool{}
			db.Select("name").Where("id=?", testPools[j].PoolId).Find(&pool)
			if len(pool) != 0 {
				pools = append(pools, pool[0].Name)
			}
		}
		one.Total = count
		one.Pools = pools
		// append result list
		result = append(result, one)
	}
	return result, nil
}

//functin to find the detail of test
func TestDetails(id string) (bodyTypes.Test, error) {
	db := config.DB
	// final result structure variable
	finalTest := bodyTypes.Test{}
	//fetching test from db
	tests := []types.Tests{}
	db.Raw("select name,duration from tests where id=?", id).Scan(&tests)
	if len(tests) == 0 {
		// when no test exist
		return finalTest, errors.New("No test exist")
	}
	finalTest.Name = tests[0].Name
	finalTest.Duration = strconv.FormatInt(tests[0].Duration/60, 10)

	//fetch test pool details from db
	testPools := []types.TestPool{}
	db.Raw("select pool_id, no_of_questions from test_pools where test_id=?", id).Scan(&testPools)
	count := 0
	// save in a structure
	finalPools := make([]bodyTypes.TestPool, 0)
	for i := 0; i < len(testPools); i++ {
		one := bodyTypes.TestPool{}
		one.PoolId = testPools[i].PoolId
		one.NoOfQuestions = testPools[i].NoOfQuestions
		count += one.NoOfQuestions
		// fetch pool name from db
		pool := []types.Pool{}
		db.Select("name").Where("id=?", testPools[i].PoolId).Find(&pool)
		if len(pool) != 0 {
			one.PoolName = pool[0].Name
			// total no. of questions in a pool
			list, err := questions.ListQuestions(testPools[i].PoolId)
			if err == nil {
				one.TotalQuestion = len(list)
			}
		}
		finalPools = append(finalPools, one)
	}
	finalTest.Pools = finalPools
	finalTest.TotalQuestions = count
	return finalTest, nil
}

//function to edit test
func EditTest(data bodyTypes.Test, id string) (int, error) {
	//check data is valid or not
	code, err := validateData(data)
	if err != nil {
		zap.S().Error(err)
		return code, err
	}

	db := config.DB
	// convert test name in slug
	slug := methods.SlugOfString(data.Name)

	//convert duration to int
	duration, err := strconv.ParseInt(data.Duration, 10, 0)
	if err != nil {
		zap.S().Error(err)
		return 400, err
	}

	//update slug or test id in drives before deleting
	db.Model(&types.Drive{}).Where("test_id=?", id).Update("test_id", slug)

	// remove test
	row := db.Exec("delete from tests where id=?", id).RowsAffected
	if row == 0 {
		return 400, errors.New("Test doesn't exists")
	}

	//saving data according to test structure
	testData := types.Tests{
		Id:       slug,
		Name:     data.Name,
		Date:     time.Now().Unix(),
		Duration: duration * 60,
	}
	//save data in db
	err = db.Create(&testData).Error
	if err != nil {
		zap.S().Error(err)
		return 500, err
	}
	err = db.Exec("DELETE FROM test_pools WHERE test_id = '" + id + "';").Error
	if err != nil {
		log.Println(err)
	}
	// save test pools in db
	for i := 0; i < len(data.Pools); i++ {
		testPoolData := types.TestPool{
			PoolId:        data.Pools[i].PoolId,
			NoOfQuestions: data.Pools[i].NoOfQuestions,
			TestId:        slug,
		}
		db.Create(&testPoolData)

	}
	return 200, nil
}

func DeleteTest(id string) error {
	db := config.DB
	// check in drive
	var count int64
	db.Model(&types.Drive{}).Where("test_id=?", id).Count(&count)

	zap.S().Error(count)
	if count != 0 {
		return errors.New("You cannot delete this test because this test is assigned to some drives")
	}
	// remove test
	row := db.Exec("delete from tests where id=?", id).RowsAffected

	zap.S().Error(row)
	if row == 0 {
		return errors.New("Test doesn't exists")
	}
	//delete test pools
	db.Exec("delete from test_pools where test_id=?", id)
	return nil
}
