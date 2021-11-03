package result

import (
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/drive"
	"git.xenonstack.com/util/drive-portal/src/types"
)

var poolWise bool

func UsersResult(driveId, page int) (map[string]interface{}, error) {
	mapd := make(map[string]interface{})
	poolWise = false
	db := config.DB

	//check whether to send result or only user list
	// if drive is ongoing then only user list otherwise result
	var count int64
	db.Model(&types.UserSession{}).Where("drive_id=?", driveId).Count(&count)
	users, total, err := drive.ListUsers(strconv.Itoa(driveId), page)
	if err != nil {
		zap.S().Error(err)
		return mapd, err
	}
	mapd["users"] = users
	mapd["total"] = total
	mapd["poolWise"] = poolWise
	if count == 0 {
		// when no user start the test then also users list will be shown
		return mapd, nil
	}
	db.Model(&types.UserSession{}).Where("drive_id=? AND expire>?", driveId, time.Now().Unix()).Count(&count)
	if count != 0 {
		// users list
		return mapd, nil
	}
	// result
	result := make([]bodyTypes.Result, 0)

	// create view if not exists
	db.Exec("CREATE OR REPLACE VIEW result_drive_view AS SELECT email,drive_id,count(ques_id) as total,count(marked_id) filter (where marked_id>0) as attempted,count(marked_id) filter (where marked_id<=0) as not_attempted,count(marked_id) filter (where marked_id=answer_id) as correct,count(marked_id) filter (where marked_id>0 AND marked_id<>answer_id) as wrong from answers group by email,drive_id;")

	// fetch each user result
	for i := 0; i < len(users); i++ {
		one, err := FetchUserResult(users[i].Email, driveId, "false")
		if err != nil {
			zap.S().Error(err)
			return mapd, err
		}
		result = append(result, one)
	}
	mapd["users"] = result
	mapd["total"] = total
	mapd["poolWise"] = poolWise
	return mapd, nil
}

func PoolResult(driveId int, email, results string) (map[string]interface{}, error) {
	mapd := make(map[string]interface{})
	db := config.DB
	if results != "true" {
		//check whether to send result or only user list
		// if drive is ongoing then only user list otherwise result
		var count int64
		db.Model(&types.UserSession{}).Where("drive_id=? AND expire>?", driveId, time.Now().Unix()).Count(&count)
		if count != 0 {
			// users list
			mapd["error"] = "true"
			mapd["message"] = "Drive is ongoing! Please wait for the drive to finish."
			return mapd, errors.New("Drive is ongoing! Please wait for the drive to finish")
		}

		// create view if not exists
		db.Exec("CREATE OR REPLACE VIEW result_pool_view AS SELECT email,drive_id,pool_id,count(ques_id) as total,count(marked_id) filter (where marked_id>0) as attempted,count(marked_id) filter (where marked_id<=0) as not_attempted,count(marked_id) filter (where marked_id=answer_id) as correct,count(marked_id) filter (where marked_id>0 AND marked_id<>answer_id) as wrong from answers group by email,drive_id,pool_id;")
	}

	result := bodyTypes.UserResult{}

	//fetch user result
	user, err := FetchUserResult(email, driveId, results)
	if err != nil {
		zap.S().Error("error in fetching users...", err)
		// users list
		mapd["error"] = true
		mapd["message"] = err.Error()
		return mapd, err
	}
	if user.Total == 0 {
		zap.S().Error("error in fetching users...", err)
		// users list
		mapd["error"] = false
		mapd["result"] = result
		return mapd, err
	}

	//fetch user assigned pools
	pools := []types.TestPool{}
	db.Where("test_id IN (?)", db.Model(types.Drive{}).Select("test_id").Where("id=?", driveId).QueryExpr()).Find(&pools)

	// fetch each pool result
	poolResult := make([]bodyTypes.PoolResult, 0)
	for i := 0; i < len(pools); i++ {
		one, err := fetchPoolResult(email, pools[i].PoolId, driveId, results)
		if err != nil {
			zap.S().Error(err)
			continue
		}
		one.Total = pools[i].NoOfQuestions
		poolResult = append(poolResult, one)
	}

	//send final result to websocket connection
	result.Pool = poolResult
	result.Result = user
	mapd["result"] = result
	return mapd, nil
}

type Result struct {
	Total        int `gorm:"column:total"`
	Attempted    int `gorm:"column:attempted"`
	NotAttempted int `gorm:"column:not_attempted"`
	Correct      int `gorm:"column:correct"`
	Wrong        int `gorm:"column:wrong"`
}

func FetchUserResult(email string, driveId int, results string) (bodyTypes.Result, error) {
	db := config.DB
	var result Result

	if results != "true" {
		db.Raw("select * from result_drive_view where  drive_id=? AND email=?", driveId, email).Find(&result)
		if result.Total != 0 {
			poolWise = true
		} else {
			db.Raw("select count(ques_id) as total,count(marked_id) filter (where marked_id>0) as attempted,count(marked_id) filter (where marked_id=answer_id) as correct,count(marked_id) filter (where marked_id>0 AND marked_id<>answer_id) as wrong from results where  drive_id=? AND email=? group by email,drive_id;", driveId, email).Find(&result)
		}
	} else {
		db.Raw("select count(ques_id) as total,count(marked_id) filter (where marked_id>0) as attempted,count(marked_id) filter (where marked_id=answer_id) as correct,count(marked_id) filter (where marked_id>0 AND marked_id<>answer_id) as wrong from results where  drive_id=? AND email=? group by email,drive_id;", driveId, email).Find(&result)
	}
	//calculate test status
	userSession := []types.UserSession{}
	db.Where("drive_id=? AND email=?", driveId, email).Find(&userSession)
	status := ""
	var restart, browser, timeTaken int64
	if len(userSession) == 0 {
		status = "Not Started"
		restart = 0
		timeTaken = 0
		browser = 0
	} else {
		if userSession[0].Expire < time.Now().Unix() {
			status = "Completed"
		} else {
			status = "Ongoing"
		}
		restart = int64(userSession[0].Restart)
		browser = int64(userSession[0].Browser)
		timeTaken = int64(userSession[0].TimeTaken)
	}

	return bodyTypes.Result{
		Email:      email,
		Attempted:  result.Attempted,
		Correct:    result.Correct,
		Wrong:      result.Attempted - result.Correct,
		Total:      result.Total,
		TestStatus: status,
		Restart:    restart,
		Browser:    browser,
		TimeTaken:  timeTaken,
	}, nil
}

func fetchPoolResult(email, pool string, driveId int, results string) (bodyTypes.PoolResult, error) {
	db := config.DB
	db = db.Debug()
	poolDetails := types.Pool{}
	db.Where("id=?", pool).Find(&poolDetails)

	var result Result
	if results != "true" {
		db.Raw("select * from result_pool_view where  drive_id=? AND email=? AND pool_id=?", driveId, email, pool).Find(&result)
		if result.Total == 0 {
			db.Raw("select count(ques_id) as total,count(marked_id) filter (where marked_id>0) as attempted,count(marked_id) filter (where marked_id=answer_id) as correct from results where  drive_id=? AND email=? AND pool_id=? group by email,drive_id,pool_id;", driveId, email, pool).Find(&result)
		}
	} else {
		db.Raw("select count(ques_id) as total,count(marked_id) filter (where marked_id>0) as attempted,count(marked_id) filter (where marked_id=answer_id) as correct from results where  drive_id=? AND email=? AND pool_id=? group by email,drive_id,pool_id;", driveId, email, pool).Find(&result)
	}
	return bodyTypes.PoolResult{
		PoolId:    pool,
		PoolName:  poolDetails.Name,
		Attempted: result.Attempted,
		Correct:   result.Correct,
		Wrong:     result.Attempted - result.Correct,
		Total:     result.Total,
	}, nil
}
