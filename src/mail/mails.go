package mail

import (
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
)

func SendDriveMail(drive bodyTypes.DriveDetails, email string) {
	// parse start time to int
	i, err := strconv.ParseInt(drive.Start, 10, 64)
	if err != nil {
		zap.S().Error(err)
		return
	}
	// convert unix time to timestamp
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		loc, _ = time.LoadLocation("UTC")
	}
	tm := time.Unix(i, 0).In(loc)

	// verbose date
	suffix := "th"
	switch tm.Day() % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	}

	if tm.Day() > 10 && tm.Day() < 20 {
		suffix = "th"
	}

	//formating timestamp to specific format
	startTime := tm.Format("2" + suffix + " Jan 2006, 03:04 PM (MST)")

	// parse end time to int
	i, err = strconv.ParseInt(drive.End, 10, 64)
	if err != nil {
		zap.S().Error(err)
		return
	}
	tm = time.Unix(i, 0).In(loc)
	// verbose date
	suffix = "th"
	switch tm.Day() % 10 {
	case 1:
		suffix = "st"
	case 2:
		suffix = "nd"
	case 3:
		suffix = "rd"
	}
	//formating timestamp to specific format
	endTime := tm.Format("2" + suffix + " Jan 2006, 03:04 PM (MST)")

	// map saving name of user and verification code for email verification
	mapd := map[string]interface{}{
		"DriveName": drive.Name,
		"College":   drive.College.Name,
		"StartTime": startTime,
		"EndTime":   endTime,
		"Url":       os.Getenv("HP_ACS_FRONT_ADDR") + "login",
		"Website":   strings.TrimSuffix(strings.TrimPrefix(os.Getenv("HP_ACS_FRONT_ADDR"), "https://"), "/"),
	}
	// saving subject as string
	sub := "XenonStack Career Portal Drive Notification"
	// saving template as string by parsing above map
	tmpl := EmailTemplate(config.MailPath+"/mail/templates/drive.tmpl", mapd)
	//now sending mail
	SendMailV2(email, sub, tmpl)
}
