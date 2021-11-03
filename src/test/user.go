package test

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/jwt"
	"git.xenonstack.com/util/drive-portal/src/types"
)

func CheckAndIntializeTest(email, name, driveid, testid string) (map[string]interface{}, error) {
	// Initialise result map
	mapd := make(map[string]interface{})

	//convert driveid string to int
	id, err := strconv.Atoi(driveid)
	if err != nil {
		zap.S().Error(err)
		mapd["error"] = true
		mapd["message"] = "Please pass valid driveid only integer"
		return mapd, err
	}

	db := config.DB

	var count int64
	//fetch test details
	test := []types.Tests{}
	db.Raw("select id,name,duration from tests where id=?", testid).Scan(&test)
	//check test exists
	if len(test) == 0 {
		mapd["error"] = true
		mapd["message"] = "No test exists with this id"
		return mapd, errors.New("No test exists with this id")
	}

	//fetch drive details
	drive := []types.Drive{}
	db.Raw("select id,name,end_time,test_id from drives where id=?", id).Scan(&drive)

	//check drive exists
	if len(drive) == 0 {
		mapd["error"] = true
		mapd["message"] = "No drive exists with this " + driveid
		return mapd, errors.New("No drive exists with this " + driveid)
	}

	//check drive testid is equal to passed testid
	if drive[0].TestId != testid {
		mapd["error"] = true
		mapd["message"] = "This test " + testid + " is not assign to this drive " + driveid
		return mapd, errors.New("This test " + testid + " is not assign to this drive " + driveid)
	}

	//check drive end time
	if drive[0].EndTime < time.Now().Unix() {
		mapd["error"] = true
		mapd["message"] = "This drive " + drive[0].Name + " is not active"
		return mapd, errors.New("This drive is not active")
	}

	//check user is assigned to that drive
	count = 0
	db.Model(&types.DriveUser{}).Where("drive_id=? AND user_email= ?", id, email).Count(&count)
	if count == 0 {
		mapd["error"] = true
		mapd["message"] = "You are not assigned to this drive"
		return mapd, errors.New("You are not assigned to this drive")
	}

	// fetch pool details
	poolDetails, err := PoolDetails(test[0].Id)
	if err != nil || len(poolDetails) == 0 {
		mapd["error"] = true
		mapd["message"] = "No pool is assigned to this test"
		return mapd, errors.New("No pool is assigned to this test")
	}

	//calculate total no. of questions in a test
	total := 0
	for i := 0; i < len(poolDetails); i++ {
		total += poolDetails[i].NoOfQuestions
	}

	//if total is zero
	if total == 0 {
		mapd["error"] = true
		mapd["message"] = "No pool is assigned to this test"
		return mapd, errors.New("No pool is assigned to this test")
	}

	// TODO check user college with drive college

	//check user already not given test
	sess := []types.UserSession{}
	db.Raw("select expire,token from user_sessions where email=? and drive_id=?", email, id).Scan(&sess)
	if len(sess) == 0 {

		// assign questions to user
		var wg sync.WaitGroup
		wg.Add(1)
		go assignQuestion1(poolDetails, drive[0].Id, email, &wg)

		//making map for jwt claims
		claims := make(map[string]interface{})
		claims["email"] = email
		claims["drive"] = driveid
		claims["test"] = testid
		claims["questions"] = total
		claims["name"] = name
		claims["drive_name"] = drive[0].Name
		claims["test_name"] = test[0].Name
		//generate new token
		mapd = jwt.NewToken(claims, test[0].Duration+120)
		mapd["start"] = time.Now().Format(time.RFC3339)
		mapd["email"] = email
		mapd["drive_name"] = drive[0].Name
		mapd["test_name"] = test[0].Name
		mapd["duration"] = test[0].Duration
		mapd["questions"] = total
		mapd["next_pool"] = poolDetails[0].PoolId
		mapd["next_index"] = 1
		mapd["progress"] = 0
		mapd["pools"] = poolDetails
		wg.Wait()
		time.Sleep(5 * time.Second)
		return mapd, nil
	}

	//check old TOKEN
	if sess[0].Expire < time.Now().Unix() {
		log.Println("expire time...", sess[0].Expire)
		mapd["error"] = true
		mapd["message"] = "Your test session has been expired"
		return mapd, errors.New("Your test session has been expired")
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go details(email, driveid, mapd, &wg)
	wg.Wait()
	mapd["token"] = sess[0].Token
	mapd["start"] = time.Now().Format(time.RFC3339)
	mapd["expire"] = time.Unix(sess[0].Expire, 0).Format(time.RFC3339)
	mapd["duration"] = (sess[0].Expire - time.Now().Unix())
	mapd["email"] = email
	mapd["drive_name"] = drive[0].Name
	mapd["test_name"] = test[0].Name
	mapd["questions"] = total
	mapd["pools"] = poolDetails
	return mapd, nil
}

func details(email, driveID string, mapd map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	db := config.DB
	db.Exec("update user_sessions set restart=restart+1 where drive_id=? AND email = ?", driveID, email)

	// find index where user left
	var answer types.Answers
	db.Raw("select pool_id,ques_index from answers where email='" + email + "' and drive_id = " + driveID + " order by updated_at desc limit 1;").Scan(&answer)

	var count int64
	db.Raw("select count(email) from answers where email='" + email + "' and drive_id = " + driveID + " AND marked_id<>0;").Count(&count)

	mapd["next_pool"] = answer.PoolId
	mapd["next_index"] = answer.QuesIndex
	if count > 0 {
		mapd["progress"] = count - 1
	} else {
		mapd["progress"] = 0
	}
	return
}

func PoolDetails(testid string) ([]bodyTypes.TestPool, error) {

	// save in a structure
	finalPools := make([]bodyTypes.TestPool, 0)
	db := config.DB
	//fetch test pool details from db
	testPools := []types.TestPool{}
	db.Raw("select pool_id,no_of_questions from test_pools where test_id='" + testid + "' order by pool_id").Scan(&testPools)

	for i := 0; i < len(testPools); i++ {
		one := bodyTypes.TestPool{}
		one.PoolId = testPools[i].PoolId
		one.NoOfQuestions = testPools[i].NoOfQuestions
		// fetch pool name from db
		pool := []types.Pool{}
		db.Select("name").Where("id=?", testPools[i].PoolId).Find(&pool)
		if len(pool) != 0 {
			one.PoolName = pool[0].Name
		}
		finalPools = append(finalPools, one)
	}
	return finalPools, nil
}

func CompletedUserTests(email string) ([]bodyTypes.Completetest, error) {
	result := make([]bodyTypes.Completetest, 0)
	// connecting to db
	db := config.DB

	drives := []types.UserSession{}
	db.Raw("select drive_id,expire from user_sessions where email='" + email + "' and expire <= " + strconv.FormatInt(time.Now().Unix(), 10) + " order by expire desc").Scan(&drives)
	for i := 0; i < len(drives); i++ {
		// fetch drive detail
		drive := []types.Drive{}
		db.Select("name,test_id").Where("id=?", drives[i].DriveId).Find(&drive)
		if len(drive) == 0 {
			continue
		}

		//fetch test name
		test := []types.Tests{}
		db.Select("name").Where("id=?", drive[0].TestId).Find(&test)
		if len(test) == 0 {
			continue
		}

		// creating final list
		result = append(result, bodyTypes.Completetest{
			Drive:     drive[0].Name,
			Test:      test[0].Name,
			Completed: drives[i].Expire,
		})
	}
	return result, nil
}

func assignQuestion1(pools []bodyTypes.TestPool, driveID int, email string, wg *sync.WaitGroup) {
	defer wg.Done()
	db := config.DB
	var count int64
	db.Model(&types.Answers{}).Where("email=? AND drive_id=?", email, driveID).Count(&count)
	if count != 0 {
		return
	}

	// prepare sql string for batch insertion
	sqlStr := "INSERT INTO answers(email,drive_id,pool_id,ques_id,answer_id,ques_index,marked_id) VALUES"
	for i := 0; i < len(pools); i++ {
		//fetch random question
		ques := []types.Questions{}
		db.Where("pool_id=?", pools[i].PoolId).Order("random()").Limit(pools[i].NoOfQuestions).Find(&ques)
		if len(ques) == 0 {
			continue
		}
		for j := 0; j < len(ques); j++ {
			sqlStr += "('" + email + "'," + strconv.Itoa(driveID) + ",'" + pools[i].PoolId + "'," + strconv.Itoa(ques[j].Id) + "," + strconv.Itoa(ques[j].AnswerId) + "," + strconv.Itoa(j+1) + "," + "0),"
		}
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	for i := 0; i < 5; i++ {
		time.Sleep(time.Duration(i) * time.Second)
		// begin transaction
		tx, err := db.DB().Begin()
		if err != nil {
			zap.S().Error(err)
			continue
		}
		// prepare sql query to execute
		sql, err := tx.Prepare(sqlStr)
		if err != nil {
			_ = tx.Rollback()
			continue
		}
		// execute query
		_, err = sql.Exec()
		if err != nil {
			_ = tx.Rollback()
			continue
		}

		// commit the changes
		err = tx.Commit()
		if err != nil {
			_ = tx.Rollback()
			continue
		} else {
			return
		}
	}
}

func AssignStartTime() {
	db := config.DB

	//delete drives
	db.Exec("delete from drives where test_id not in (select id from tests);")
	db.Exec("delete from drive_users where drive_id not in (select id from drives);")
	db.Exec("delete from user_sessions where drive_id not in (select id from drives);")
	db.Exec("delete from answers where drive_id not in(select id from drives);")
	db.Exec("delete from results where drive_id not in(select id from drives);")

	users := []types.UserSession{}
	db.Raw("select * from user_sessions where start=0 or time_taken=0").Scan(&users)
	for i := 0; i < len(users); i++ {
		driveDetails := types.Drive{}
		db.Raw("select test_id from drives where id=?", users[i].DriveId).Scan(&driveDetails)
		testDetails, err := TestDetails(driveDetails.TestId)
		if err != nil {
			continue
		}
		row := db.Exec("update user_sessions set start=expire-" + testDetails.Duration + " where drive_id=" + strconv.Itoa(users[i].DriveId) + " AND start=0 ").RowsAffected
		log.Println(row, "=-=-=-=-=")
		row = db.Exec("update user_sessions set time_taken=" + testDetails.Duration + " where drive_id=" + strconv.Itoa(users[i].DriveId) + " AND time_taken=0 ").RowsAffected
		log.Println(row, "=-=-=-=-=")
	}
}
