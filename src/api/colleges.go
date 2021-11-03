package api

import (
	"encoding/csv"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/college"
	"github.com/gin-gonic/gin"
)

func AddCollege(c *gin.Context) {
	var data bodyTypes.College
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass name and location",
		})
		return
	}

	if strings.TrimSpace(data.Name) == "" || strings.TrimSpace(data.Location) == "" {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass valid name and location",
		})
		return
	}

	code, err := college.Add(data)
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": "College added successfully",
	})
}

func ListColleges(c *gin.Context) {
	pageNo := 1
	pageStr := c.Query("page")
	if pageStr != "" {
		i, err := strconv.Atoi(pageStr)
		if err != nil {
			zap.S().Error(err)
			c.JSON(400, gin.H{
				"error":   true,
				"message": "Please pass valid current number",
			})
		}
		pageNo = i
	}

	total, list, err := college.List(pageNo, c.Query("search"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":    false,
		"colleges": list,
		"total":    total,
	})
}

func ViewCollege(c *gin.Context) {
	college, code, err := college.View(c.Param("id"))
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"college": college,
	})
}

func UpdateCollege(c *gin.Context) {
	var data bodyTypes.College
	if err := c.BindJSON(&data); err != nil {
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass name and location",
		})
	}
	code, err := college.Update(c.Param("id"), data)
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": "college updated successfully",
	})
}

func DeleteCollege(c *gin.Context) {
	code, err := college.Delete(c.Param("id"))
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": "college successfully deleted",
	})
}

func CSVCollege(c *gin.Context) {
	file, _, err := c.Request.FormFile("colleges")
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
		msg := college.CSVColleges(lines)
		c.JSON(200, gin.H{
			"error":   false,
			"message": msg,
		})
	}
}
