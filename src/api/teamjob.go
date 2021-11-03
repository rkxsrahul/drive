package api

import (
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/job"
	"github.com/gin-gonic/gin"
)

func AddTeamJob(c *gin.Context) {
	var data bodyTypes.JobTeam

	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Please pass the required fields."})
		return
	}

	err := job.AddTeam(data)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "job team succesfully created",
	})
}

func ListTeamJob(c *gin.Context) {
	// fetch pool list from db
	list, err := job.ListTeams()
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"teams": list,
	})
}

func TeamDetail(c *gin.Context) {
	team, err := job.TeamDetails(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"team":  team,
	})
}

func DeleteTeam(c *gin.Context) {
	// fetch pool list from db
	err := job.DeleteTeam(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "team succesfully deleted",
	})
}

func UpdateTeam(c *gin.Context) {
	var data bodyTypes.JobTeam

	if err := c.BindJSON(&data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Please pass the required fields."})
		return
	}
	// fetch pool list from db
	err := job.UpdateTeam(c.Param("id"), data)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "team succesfully edited",
	})
}
