package api

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/questions"
	"git.xenonstack.com/util/drive-portal/src/types"
	"github.com/gin-gonic/gin"
)

func AddQuestion(c *gin.Context) {
	var data bodyTypes.Questions
	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Please pass the required fields."})
		return
	}

	code, err := questions.AddQuestion(c.Param("pool"), data, true)
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": "question succesfully added",
	})
}

func ListQuestions(c *gin.Context) {
	list, err := questions.ListQuestions(c.Param("pool"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":     false,
		"questions": list,
	})
}

func QuestionDetails(c *gin.Context) {
	question, err := questions.QuestionDetail(c.Param("pool"), c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":    false,
		"question": question,
	})
}

func DeleteQuestion(c *gin.Context) {
	// delete pool on basis of id pass in url
	err := questions.DeleteQuestion(c.Param("pool"), c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "question succesfully deleted",
	})
}

func UpdateQuestion(c *gin.Context) {
	// fetching pool name from request body
	var data bodyTypes.Questions
	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Please pass the required fields."})
		return
	}

	//edit pool on basis of id and send new name of pool
	code, err := questions.EditQuestion(c.Param("pool"), c.Param("id"), data)
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "question succesfully updated",
	})
}

func CSVQuestion(c *gin.Context) {
	file, _, err := c.Request.FormFile("questions")
	if err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{"error": true, "message": err.Error()})
		return
	}
	if file != nil {
		lines, err := csv.NewReader(file).ReadAll()
		if err != nil {
			c.JSON(500, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		code, msg, err := questions.CSVQuestions(lines, c.Param("pool"))
		if err != nil {
			c.JSON(code, gin.H{
				"error":   true,
				"message": err.Error(),
			})
			return
		}
		c.JSON(code, gin.H{
			"error":   false,
			"message": msg,
		})
	}
}

func DownloadCSVQuestion(c *gin.Context) {
	poolID := c.Param("pool")

	// directory
	dir := "./pool_csv/"
	// path of file
	file2 := dir + poolID + ".csv"

	//check file already exists
	if _, err := os.Stat(file2); !os.IsNotExist(err) {
		log.Println("file already exist")
		c.JSON(200, gin.H{
			"error": false,
			"file":  poolID + ".csv",
		})
		return
	}

	db := config.DB
	data := make([][]string, 0)
	one := make([]string, 0)
	head := []string{"title", "op1 value", "op1 correct", "op2 value", "op2 correct", "op3 value", "op3 correct", "op4 value", "op4 correct", "op5 value", "op5 correct", "op6 value", "op6 correct"}
	one = append(one, head...)
	data = append(data, one)
	// fetch questions
	questions := []types.Questions{}
	db.Raw("select id, title from questions where pool_id=? and type=?", poolID, "string").Scan(&questions)

	for i := 0; i < len(questions); i++ {
		one := make([]string, 0)
		// fetch options
		options := []types.Options{}
		db.Raw("select value,is_correct from options where ques_id=?", questions[i].Id).Scan(&options)
		one = append(one, questions[i].Title)
		for j := 0; j < len(options); j++ {
			one = append(one, options[j].Value)
			one = append(one, strconv.FormatBool(options[j].IsCorrect))
		}
		data = append(data, one)
	}

	//check directory and create directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Println("directory not exist")
		//create directory
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			zap.S().Error("error in creating directory...", err)
			c.JSON(500, gin.H{
				"error":   true,
				"message": "Please try again later...",
			})
			return
		}
	}

	// create file and write data to file
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

	c.JSON(200, gin.H{
		"error": false,
		"file":  poolID + ".csv",
	})
}

func ServePoolCSV(c *gin.Context) {
	go func() {
		time.Sleep(10 * time.Minute)
		os.Remove("./pool_csv/" + c.Param("file"))
	}()

	c.File("./pool_csv/" + c.Param("file"))
}
