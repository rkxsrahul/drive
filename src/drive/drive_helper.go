package drive

import (
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/college"
	"git.xenonstack.com/util/drive-portal/src/test"
	"git.xenonstack.com/util/drive-portal/src/types"
)

//function to check data is exist or not into drive
func checkData(data *bodyTypes.Drive) error {
	if data.Type == "open" {
		data.CollegeId = ""
	} else if data.Type == "college" {
		if data.CollegeId == "" {
			return errors.New("Please pass college id")
		}
		// check college exists
		_, _, err := college.View(data.CollegeId)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Please pass valid value of type i.e open or college")
	}
	//check test exists
	_, err := test.TestDetails(data.TestId)
	if err != nil {
		return err
	}

	//compare start and end time
	start, err := strconv.ParseInt(data.Start, 10, 64)
	if err != nil {
		return errors.New("Please pass valid start time in int only")
	}
	end, err := strconv.ParseInt(data.End, 10, 64)
	if err != nil {
		return errors.New("Please pass valid end time in int only")
	}
	if start > end {
		return errors.New("Please pass valid time start will be less then end")
	}
	return nil
}

//function to add data into drive
func Add(data bodyTypes.Drive) (int, error) {
	//validate data
	err := checkData(&data)
	if err != nil {
		return 400, err
	}
	db := config.DB

	//create drive structure with data
	drive := types.Drive{}
	drive.Name = data.Name
	drive.Type = data.Type
	drive.TestId = data.TestId
	drive.CollegeId = data.CollegeId
	drive.StartStr = data.StartStr
	drive.EndStr = data.EndStr
	//parse start and end time in int
	start, err := strconv.ParseInt(data.Start, 10, 64)
	if err != nil {
		return 400, errors.New("Please pass valid start time in int only")
	}
	end, err := strconv.ParseInt(data.End, 10, 64)
	if err != nil {
		return 400, errors.New("Please pass valid end time in int only")
	}
	drive.StartTime = start
	drive.EndTime = end
	//save data in db table
	err = db.Create(&drive).Error
	if err != nil {
		zap.S().Error(err)
		return 400, errors.New("Drive name (" + drive.Name + ") is already in use.")
	}
	return 200, nil
}

//function to find the list of drive
func List(source string) ([]bodyTypes.DriveList, error) {
	db := config.DB

	//fetch drives from db
	drive := []types.Drive{}
	if source == "ongoing" {
		db.Raw("select id, start_time,end_time,name,type from drives where end_time >= ? order by start_time desc", time.Now().Unix()).Scan(&drive)
	} else {
		db.Raw("select id,start_time,end_time,name,type from drives order by start_time DESC").Scan(&drive)
	}
	//make slice of result type
	drives := make([]bodyTypes.DriveList, 0)
	for i := 0; i < len(drive); i++ {
		one := bodyTypes.DriveList{}
		one.Id = strconv.Itoa(drive[i].Id)
		one.Start = strconv.FormatInt(drive[i].StartTime, 10)
		one.End = strconv.FormatInt(drive[i].EndTime, 10)
		one.Name = drive[i].Name
		one.Type = drive[i].Type
		// fetch no of users in a drive
		var count int64
		db.Raw("select count(drive_id) from drive_users where drive_id=?", drive[i].Id).Count(&count)
		one.Users = count
		//append drives in a slice
		drives = append(drives, one)
	}
	return drives, nil
}

//function to view the drive
func View(id, source string) (bodyTypes.DriveDetails, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error(err)
		return bodyTypes.DriveDetails{}, errors.New("Please pass valid drive id")
	}

	db := config.DB
	//fetch drive from db
	drive := []types.Drive{}

	if source == "ongoing" {
		db.Raw("select id,name,type,test_id,college_id,start_time,end_time,start_str,end_str from drives where id=? and end_time >= ?", idInt, time.Now().Unix()).Scan(&drive)
	} else {
		db.Raw("select id,name,type,test_id,college_id,start_time,end_time,start_str,end_str from drives where id=?", idInt).Scan(&drive)
	}
	if len(drive) == 0 {
		return bodyTypes.DriveDetails{}, errors.New("there is no such drive correspondance to this id")
	}

	if drive[0].Type == "college" {
		//fetch college details with college id
		college, _, _ := college.View(drive[0].CollegeId)
		//return final result in a structure if type is college
		return bodyTypes.DriveDetails{
			Id:        strconv.Itoa(drive[0].Id),
			Type:      drive[0].Type,
			Name:      drive[0].Name,
			TestId:    drive[0].TestId,
			CollegeId: drive[0].CollegeId,
			Start:     strconv.FormatInt(drive[0].StartTime, 10),
			End:       strconv.FormatInt(drive[0].EndTime, 10),
			StartStr:  drive[0].StartStr,
			EndStr:    drive[0].EndStr,
			College:   college,
		}, nil
	}

	//return final result in a structure if type is open
	return bodyTypes.DriveDetails{
		Id:       strconv.Itoa(drive[0].Id),
		Type:     drive[0].Type,
		Name:     drive[0].Name,
		TestId:   drive[0].TestId,
		Start:    strconv.FormatInt(drive[0].StartTime, 10),
		End:      strconv.FormatInt(drive[0].EndTime, 10),
		StartStr: drive[0].StartStr,
		EndStr:   drive[0].EndStr,
	}, nil
}

//function to find the summary
func Summary(driveID string) (bodyTypes.SummaryDrive, error) {
	drive, err := View(driveID, "")
	if err != nil {
		zap.S().Error(err)
		return bodyTypes.SummaryDrive{}, err
	}

	// db client
	db := config.DB

	// now time
	now := time.Now().Unix()

	// calculate summary
	var ongoing int64
	db.Raw("select count(drive_id) from user_sessions where drive_id=? and expire >=?", drive.Id, now).Count(&ongoing)
	var completed int64
	db.Raw("select count(drive_id) from user_sessions where drive_id=? and expire < ?", drive.Id, now).Count(&completed)
	var total int64
	db.Raw("select count(drive_id) from drive_users where drive_id=?", drive.Id).Count(&total)
	return bodyTypes.SummaryDrive{
		DriveDetails: drive,
		NotStarted:   total - ongoing - completed,
		Completed:    completed,
		Ongoing:      ongoing,
	}, nil
}

//function to edit drive
func Edit(id string, data bodyTypes.Drive) (int, error) {
	//parse id as int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error(err)
		return 400, errors.New("Please pass valid drive id")
	}

	db := config.DB

	//check drive is there or not
	var count int64
	db.Raw("select count(id) from drives where id=?", idInt).Count(&count)

	if count == 0 {
		return 400, errors.New("there is no such drive correspondance to this id")
	}

	//validate data
	err = checkData(&data)
	if err != nil {
		return 400, err
	}

	//create drive structure with data
	drive := types.Drive{}
	drive.Name = data.Name
	drive.Type = data.Type
	drive.TestId = data.TestId
	drive.CollegeId = data.CollegeId
	drive.StartStr = data.StartStr
	drive.EndStr = data.EndStr
	//parse start and end time in int
	start, err := strconv.ParseInt(data.Start, 10, 64)
	if err != nil {
		return 400, errors.New("Please pass valid start time in int only")
	}
	end, err := strconv.ParseInt(data.End, 10, 64)
	if err != nil {
		return 400, errors.New("Please pass valid end time in int only")
	}
	drive.StartTime = start
	drive.EndTime = end

	//update data in db table
	out := db.Model(&types.Drive{}).Where("id=?", idInt).Updates(&drive)

	if out.Error != nil {
		zap.S().Error(out.Error)
		return 400, out.Error
	}
	return 200, nil
}

//function to delete the drive
func Delete(id string) (int, error) {
	//parse id as int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error(err)
		return 400, errors.New("Please pass valid drive id")
	}
	db := config.DB

	//check any user is there
	var count int64
	db.Raw("select count(drive_id) from drive_users where drive_id=?", idInt).Count(&count)

	if count > 0 {
		return 400, errors.New("Delete all assigned users then only you can delete this drive")
	}

	go func() {
		db.Exec("delete from results.pool_analytical where drive_id=" + id + ";")
		db.Exec("delete from results.drive_analytical where drive_id=" + id + ";")
	}()

	//delete drive and check also drive exists
	row := db.Exec("delete from drives where id=?", idInt).RowsAffected
	if row == 0 {
		return 400, errors.New("there is no such drive correspondance to this id")
	}
	return 200, nil
}
