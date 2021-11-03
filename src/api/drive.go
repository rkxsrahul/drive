package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/drive"
)

func CreateDrive(c *gin.Context) {
	var data bodyTypes.Drive
	if err := c.BindJSON(&data); err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass all required details",
		})
		return
	}
	code, err := drive.Add(data)
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": data.Name + " drive created succesfully",
	})
}

func ListDrives(c *gin.Context) {
	source := c.Query("value")
	list, err := drive.List(source)
	if err != nil {
		zap.S().Error(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Not able to fetch drives",
		})
		return
	}
	c.JSON(200, gin.H{
		"error":  false,
		"drives": list,
	})
}

func ViewDrive(c *gin.Context) {
	drive, err := drive.View(c.Param("id"), "")
	if err != nil {
		zap.S().Error(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"drive": drive,
	})
}

func DriveSummary(c *gin.Context) {
	summary, err := drive.Summary(c.Param("id"))
	if err != nil {
		zap.S().Error(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"summary": summary,
	})
}

func UpdateDrive(c *gin.Context) {
	// bind data from request body
	var data bodyTypes.Drive
	if err := c.BindJSON(&data); err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{
			"error":   true,
			"message": "Please pass all required details",
		})
		return
	}

	//update drive data
	code, err := drive.Edit(c.Param("id"), data)
	if err != nil {
		zap.S().Error(err)
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": c.Param("id") + " drive updated successfully.",
	})
}

func DeleteDrive(c *gin.Context) {
	code, err := drive.Delete(c.Param("id"))
	if err != nil {
		zap.S().Error(err)
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": c.Param("id") + " drive deleted successfully.",
	})
}
