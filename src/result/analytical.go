package result

import (
	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/drive"
	"git.xenonstack.com/util/drive-portal/src/test"
	"go.uber.org/zap"
)

type Top10Drive struct {
	Email      string  `json:"email"`
	Questions  int     `json:"questions"`
	Attempted  int     `json:"attempted"`
	Correct    int     `json:"correct"`
	Percentage float64 `json:"percentage"`
	Accuracy   float64 `json:"accuracy"`
}

type Top10Pool struct {
	Email      string  `json:"email"`
	Percentage float64 `json:"percentage"`
	Accuracy   float64 `json:"accuracy"`
}

type PoolsWithTop struct {
	PoolName   string      `json:"pool_name"`
	Candidates []Top10Pool `json:"candidates"`
}

func DrivePoolResult(driveID string) (map[string]interface{}, error) {
	mapd := make(map[string]interface{})

	// fetch drive details
	driveDetails, err := drive.View(driveID, "")
	if err != nil {
		zap.S().Error(err)
		mapd["error"] = true
		mapd["message"] = err.Error()
		return mapd, err
	}

	// fetch test details for pool details
	testDetails, err := test.TestDetails(driveDetails.TestId)
	if err != nil {
		zap.S().Error(err)
		mapd["error"] = true
		mapd["message"] = err.Error()
		return mapd, err
	}

	db := config.DB

	// fetch overall top 10 students in a drive
	overall := []Top10Drive{}
	db.Raw("select email,drive_id,sum(pool_total_question) as questions,sum(pool_total_attempted) as attempted,sum(pool_total_correct) as correct,((sum(pool_total_correct))::decimal/((greatest(sum(pool_total_question),1)))::decimal*100) as percentage,((sum(pool_total_correct))::decimal/((greatest(sum(pool_total_attempted),1)))::decimal * 100) as accuracy from results.pool_analytical where drive_id=" + driveID + " group by email,drive_id order by accuracy desc limit 10;").Find(&overall)

	// fetch per pool top 10 students
	topPools := make([]PoolsWithTop, 0)
	for i := 0; i < len(testDetails.Pools); i++ {
		topPool := []Top10Pool{}
		db.Raw("select email,pool_percentage as percentage,pool_accuracy as accuracy from results.pool_analytical where drive_id=" + driveID + " AND pool_id='" + testDetails.Pools[i].PoolId + "' order by accuracy desc limit 10;").Find(&topPool)
		if testDetails.Pools[i].PoolName == "" {
			testDetails.Pools[i].PoolName = testDetails.Pools[i].PoolId
		}
		topPools = append(topPools, PoolsWithTop{
			PoolName:   testDetails.Pools[i].PoolName,
			Candidates: topPool,
		})
	}

	mapd["overall"] = overall
	mapd["pools"] = topPools
	mapd["error"] = false

	return mapd, nil
}
