package api

import (
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/job"
	"github.com/gin-gonic/gin"
)

func AddJob(c *gin.Context) {
	var data bodyTypes.Jobs

	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Please pass the required fields."})
		return
	}

	err := job.AddJob(data)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "job succesfully created",
	})
}

func ListJobs(c *gin.Context) {
	// fetch pool list from db
	list, err := job.ListJobs()
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"list":  list,
	})
}

func ListJobsByTeam(c *gin.Context) {
	// fetch pool list from db
	list, err := job.ListJobByTeam(c.Param("team"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"jobs":  list,
	})
}

func JobDetails(c *gin.Context) {
	// fetch pool list from db
	job, err := job.JobDetails(c.Param("team"), c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"job":   job,
	})
}

func DeleteJob(c *gin.Context) {
	// fetch pool list from db
	err := job.DeleteJob(c.Param("team"), c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "job succesfully deleted",
	})
}

func UpdateJob(c *gin.Context) {
	var data bodyTypes.Jobs
	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Please pass the required fields."})
		return
	}
	// fetch pool list from db
	err := job.UpdateJob(c.Param("team"), c.Param("id"), data)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "job succesfully edited",
	})
}
