package drive

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/mail"
	"git.xenonstack.com/util/drive-portal/src/methods"
	"git.xenonstack.com/util/drive-portal/src/types"
)

func AddUser(drive, email string, check bool) (int, error) {
	//removing spaces from email id
	email = strings.TrimSpace(email)
	//parse drive string to int
	driveInt, err := strconv.Atoi(drive)
	if err != nil {
		zap.S().Error(err)
		return 400, errors.New("please pass valid drive id only int")
	}
	var driveDetails bodyTypes.DriveDetails
	// check drive is there or not
	driveDetails, err = View(drive, "")
	if err != nil {
		return 400, err
	}

	db := config.DB

	// insert in db table and along with check already there is user or not
	row := db.Create(&types.DriveUser{
		DriveId:   driveInt,
		UserEmail: email,
	}).RowsAffected
	if row == 0 {
		return 400, errors.New("Cannot add duplicate users in one drive")
	}

	// send email to users for notification
	go mail.SendDriveMail(driveDetails, email)
	return 200, nil
}

func DeleteUser(drive, email string) (int, error) {
	//removing spaces from email id
	email = strings.TrimSpace(email)
	//parse drive string to int
	driveInt, err := strconv.Atoi(drive)
	if err != nil {
		zap.S().Error(err)
		return 400, errors.New("please pass valid drive id only int")
	}
	// check drive is there or not
	_, err = View(drive, "")
	if err != nil {
		return 400, err
	}

	db := config.DB
	//delete from drive user table
	row := db.Exec("delete from drive_users where drive_id=? and user_email= ?", driveInt, email).RowsAffected
	if row == 0 {
		return 400, errors.New("There is no user with this email id in drive")
	}

	//delete user session
	db.Exec("delete from user_sessions where drive_id=? and email= ?", driveInt, email)
	// delete user answers
	db.Exec("delete from answers where drive_id=? and email=?", driveInt, email)

	// delete user result
	db.Exec("delete from results where drive_id=? and email=?", driveInt, email)
	return 200, nil
}

type Users struct {
	Email  string `json:"email"`
	Status string `json:"test_status"`
}

func ListUsers(drive string, current int) ([]Users, int64, error) {

	emails := make([]Users, 0)
	//parse drive string to int
	driveInt, err := strconv.Atoi(drive)
	if err != nil {
		zap.S().Error(err)
		return emails, 0, errors.New("please pass valid drive id only int")
	}
	// check drive is there or not
	_, err = View(drive, "")
	if err != nil {
		return emails, 0, errors.New("this drive doesn't exists")
	}

	db := config.DB

	var count int64
	db.Raw("select count(drive_id) from drive_users where drive_id=?", driveInt).Count(&count)
	//users email from db
	users := []types.DriveUser{}
	if current <= 0 {
		db.Raw("select user_email from drive_users where drive_id=? order by user_email asc;", driveInt).Scan(&users)
	} else {
		db.Raw("select user_email from drive_users where drive_id=? order by user_email asc limit 10 offset (?-1)*10", driveInt, current).Scan(&users)
	}
	for i := 0; i < len(users); i++ {

		//calculate test status
		userSession := []types.UserSession{}
		db.Raw("select email,drive_id,token,expire from user_sessions where email=? and drive_id = ?", users[i].UserEmail, driveInt).Scan(&userSession)
		status := ""
		if len(userSession) == 0 {
			status = "Not Started"
		} else {
			if userSession[0].Expire < time.Now().Unix() {
				status = "Completed"
			} else {
				status = "Ongoing"
			}
		}

		emails = append(emails, Users{
			Email:  users[i].UserEmail,
			Status: status,
		})
	}
	return emails, count, nil
}

func CSVUser(drive string, data [][]string) (int, string) {
	// check drive is there or not
	_, err := View(drive, "")
	if err != nil {
		return 400, err.Error()
	}

	not := make([]string, 0)
	for i := 1; i < len(data); i++ {
		// if email is blank
		if data[i][0] == "" {
			not = append(not, strconv.Itoa(i+1))
			continue
		}

		if !methods.ValidateEmail(strings.ToLower(data[i][0])) {
			not = append(not, strconv.Itoa(i+1))
			continue
		}

		// add user in db
		_, err := AddUser(drive, strings.ToLower(data[i][0]), false)
		if err != nil {
			zap.S().Error(err)
			not = append(not, strconv.Itoa(i+1))
		}
	}
	if len(not) == 0 {
		return 200, "All users inserted successfully."
	}
	return 200, fmt.Sprint(len(data)-len(not)-1, " users inserted and these users (", strings.Join(not, ","), ") are not inserted.")
}

func UserDriveDetails(email string) ([]bodyTypes.DriveDetails, error) {
	//removing spaces from email id
	email = strings.TrimSpace(email)

	db := config.DB

	// fetch drive id list specific to that email id
	drives := []types.DriveUser{}
	db.Where("user_email= ? AND drive_id NOT IN (?)", email, db.Model(types.UserSession{}).Select("drive_id").Where("email= ? AND expire <= ?", email, time.Now().Unix()).QueryExpr()).Find(&drives)

	// creating final structure of list of drive details
	result := make([]bodyTypes.DriveDetails, 0)
	for i := 0; i < len(drives); i++ {
		oneDrive, err := View(strconv.Itoa(drives[i].DriveId), "ongoing")
		if err != nil {
			zap.S().Error(err)
			continue
		}
		result = append(result, oneDrive)
	}
	return result, nil
}
