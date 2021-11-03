package api

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/drive"
	"git.xenonstack.com/util/drive-portal/src/result"
	"git.xenonstack.com/util/drive-portal/src/types"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
	var email bodyTypes.DriveUser
	if err := c.BindJSON(&email); err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass email address",
		})
		return
	}

	code, err := drive.AddUser(c.Param("id"), strings.ToLower(email.Email), true)
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": email.Email + " user added successfully to " + c.Param("id") + " drive",
	})
}

func ListUsers(c *gin.Context) {
	pageInt := 0
	pageInt, _ = strconv.Atoi(c.Param("page"))
	list, total, err := drive.ListUsers(c.Param("id"), pageInt)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"users": list,
		"total": total,
	})
}

func Result(c *gin.Context) {
	drive, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass valid drive id only int",
		})
	}

	results := c.Query("results")
	if results != "true" {
		results = "false"
	}
	email := strings.ToLower(c.Query("email"))
	pageStr := c.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	page, _ := strconv.Atoi(pageStr)

	if email == "" {
		// user brief result
		mapd, err := result.UsersResult(drive, page)
		if err != nil {

			c.JSON(500, mapd)
			return
		}
		c.JSON(200, mapd)
	} else {
		// user detail result -> pool wise result
		data, err := result.PoolResult(drive, email, results)
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"error":  false,
			"result": data,
		})
	}
}

func DownloadResult(c *gin.Context) {
	//check whether to send result or only user list
	// if drive is ongoing then only user list otherwise result
	db := config.DB
	var count int64
	db.Model(&types.UserSession{}).Where("drive_id=? AND expire>?", c.Param("id"), time.Now().Unix()).Count(&count)
	if count != 0 {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Test is ongoing",
		})
		return
	}

	dir := "./result/"
	//check directory and create directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println("directory not exist")
		//create directory
		err = os.Mkdir(dir, os.ModePerm)
		if err != nil {
			zap.S().Error("error in creating directory...", err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "Please try again later...",
			})
			return
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go poolResultCSV(c.Param("id"), &wg)

	file2 := dir + c.Param("id") + ".csv"

	//check file already exists
	if _, err := os.Stat(file2); !os.IsNotExist(err) {
		err = os.Remove(file2)
		zap.S().Error(err)
	}

	users, _, err := drive.ListUsers(c.Param("id"), -1)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}

	data := make([][]string, 0)
	// set header of csv
	head := []string{"CANDIDATE EMAIL", "TEST STATUS", "TOTAL QUESTIONS", "ATTEMPTED", "CORRECT", "WRONG", "PERCENTAGE", "ACCURACY", "TIME TAKEN", "RESTART", "BROWSER SWITCH"}
	oneData := make([]string, 0)
	oneData = append(oneData, head...)

	data = append(data, oneData)
	// fetch each user result
	driveID, _ := strconv.Atoi(c.Param("id"))
	for i := 0; i < len(users); i++ {
		one, err := result.FetchUserResult(users[i].Email, driveID, "false")
		if err != nil {
			continue
		}
		oneData := make([]string, 0)
		oneData = append(oneData, one.Email)
		oneData = append(oneData, one.TestStatus)
		oneData = append(oneData, strconv.Itoa(one.Total))
		oneData = append(oneData, strconv.Itoa(one.Attempted))
		oneData = append(oneData, strconv.Itoa(one.Correct))
		oneData = append(oneData, strconv.Itoa(one.Wrong))
		//percentage
		if one.Total == 0 {
			oneData = append(oneData, fmt.Sprintf("%0.2f", 0.00))
		} else {
			oneData = append(oneData, fmt.Sprintf("%0.2f", (float32)(one.Correct)/(float32)(one.Total)*100))
		}
		//accuracy
		if one.Attempted == 0 {
			oneData = append(oneData, fmt.Sprintf("%0.2f", 0.00))
		} else {
			oneData = append(oneData, fmt.Sprintf("%0.2f", (float32)(one.Correct)/(float32)(one.Attempted)*100))
		}
		oneData = append(oneData, fmt.Sprintf("%d", one.TimeTaken/60))
		oneData = append(oneData, fmt.Sprintf("%d", one.Restart))
		oneData = append(oneData, fmt.Sprintf("%d", one.Browser))

		data = append(data, oneData)
	}

	// write data in csv file
	file, err := os.Create(file2)
	if err != nil {
		zap.S().Error(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Please try again later...",
		})
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			zap.S().Error(err)
		}
	}
	wg.Wait()
	c.JSON(200, gin.H{
		"error": false,
		"file":  c.Param("id") + ".csv",
		"pool":  c.Param("id") + "_pool.csv",
	})
}

func poolResultCSV(drive string, wg *sync.WaitGroup) {
	defer wg.Done()

	dir := "./result/"
	file2 := dir + drive + "_pool.csv"
	//check file already exists
	if _, err := os.Stat(file2); !os.IsNotExist(err) {
		err = os.Remove(file2)
		zap.S().Error(err)
		return
	}

	// type to fetch data
	type PoolData struct {
		PoolID    string
		Email     string
		Total     int
		Attempted int
		Correct   int
	}

	data := make([][]string, 0)

	// set header of csv
	head := []string{"POOL ID", "CANDIDATE EMAIL", "TEST STATUS", "TOTAL QUESTIONS", "ATTEMPTED", "CORRECT", "RESTART", "BROWSER SWITCH"}
	oneData := make([]string, 0)
	oneData = append(oneData, head...)
	data = append(data, oneData)

	// fetch data pool wise
	db := config.DB

	pdata := []PoolData{}
	db.Raw("SELECT rank_filter.* FROM (select pool_id,email,total,attempted,correct,rank() over (partition by pool_id order by correct::decimal/(greatest(sum(attempted),1)) desc, correct::decimal/(greatest(sum(total),1)) desc) from result_pool_view where drive_id=" + drive + " group by pool_id,email,total,attempted,correct order by pool_id) rank_filter where rank<=5;").Scan(&pdata)
	if len(pdata) == 0 {
		db.Raw("SELECT rank_filter.* FROM (select pool_id,email,count(ques_id) as total, count(marked_id) filter (where marked_id>0) as attempted, count(marked_id) filter (where marked_id=answer_id) as correct,rank() over (partition by pool_id order by count(marked_id) filter (where marked_id=answer_id)::decimal/(greatest(count(ques_id),1))*100 desc, count(marked_id) filter (where marked_id=answer_id)::decimal/(greatest(count(marked_id) filter (where marked_id>0),1))*100 desc) from results where drive_id=" + drive + " group by pool_id,email order by pool_id) rank_filter where rank<=5;").Scan(&pdata)
	}

	// fetch each user test status
	for i := 0; i < len(pdata); i++ {
		userSession := []types.UserSession{}
		db.Select("browser,restart,expire").Where("drive_id=? AND email=?", drive, pdata[i].Email).Find(&userSession)
		status := ""
		var restart, browser int64
		if len(userSession) == 0 {
			status = "Not Started"
			restart = 0
			browser = 0
		} else {
			if userSession[0].Expire < time.Now().Unix() {
				status = "Completed"
			} else {
				status = "Ongoing"
			}
			restart = int64(userSession[0].Restart)
			browser = int64(userSession[0].Browser)

			oneData := make([]string, 0)
			oneData = append(oneData, pdata[i].PoolID)
			oneData = append(oneData, pdata[i].Email)
			oneData = append(oneData, status)
			oneData = append(oneData, strconv.Itoa(pdata[i].Total))
			oneData = append(oneData, strconv.Itoa(pdata[i].Attempted))
			oneData = append(oneData, strconv.Itoa(pdata[i].Correct))
			oneData = append(oneData, fmt.Sprintf("%d", restart))
			oneData = append(oneData, fmt.Sprintf("%d", browser))

			data = append(data, oneData)
		}

	}

	// write data in csv file
	file, err := os.Create(file2)
	if err != nil {
		zap.S().Error(err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range data {
		err := writer.Write(value)
		if err != nil {
			zap.S().Error(err)
		}
	}

}

func DeleteUser(c *gin.Context) {
	code, err := drive.DeleteUser(c.Param("id"), strings.ToLower(c.Param("email")))
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": c.Param("email") + " user deleted successfully from " + c.Param("id") + " drive",
	})
}

func DeleteAllUsers(c *gin.Context) {
	drive := c.Param("id")
	db := config.DB
	//delete from drive user table
	db.Exec("delete from drive_users where drive_id=" + drive + ";")
	//delete user session
	db.Exec("delete from user_sessions where drive_id=" + drive + ";")
	// delete user answers
	db.Exec("delete from answers where drive_id=" + drive + ";")

	// delete user results
	db.Exec("delete from results where drive_id=" + drive + ";")
	// delete from pool analytical
	db.Exec("delete from results.pool_analytical where drive_id=" + drive + ";")

	c.JSON(200, gin.H{
		"error":   false,
		"message": "all users deleted successfully from " + drive + " drive",
	})
}

func CSVAddUsers(c *gin.Context) {
	file, _, err := c.Request.FormFile("users")
	if err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{"error": true, "message": err.Error()})
		return
	}
	if file != nil {
		lines, err := csv.NewReader(file).ReadAll()
		if err != nil {
			return
		}
		code, msg := drive.CSVUser(c.Param("id"), lines)
		if code != 200 {
			c.JSON(code, gin.H{
				"error":   true,
				"message": msg,
			})
			return
		}
		c.JSON(code, gin.H{
			"error":   false,
			"message": msg,
		})
	}
}

func ServeResultCSV(c *gin.Context) {
	go func() {
		time.Sleep(10 * time.Minute)
		os.Remove("./result/" + c.Param("file"))
	}()
	c.File("./result/" + c.Param("file"))
}
