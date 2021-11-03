package api

import (
	"git.xenonstack.com/util/drive-portal/src/result"

	"github.com/gin-gonic/gin"
)

// AnalyticalResult is an API handler
func AnalyticalResult(c *gin.Context) {
	mapd, err := result.DrivePoolResult(c.Param("drive"))
	if err != nil {
		c.JSON(400, mapd)
	} else {
		c.JSON(200, mapd)
	}
}

// DriveList is an API handler
func DriveList(c *gin.Context) {
	mapd := result.AllDrives()
	c.JSON(200, mapd)
}

// DriveCandidate is an API handler
func DriveCandidate(c *gin.Context) {
	mapd := result.DriveCandidates(c.Param("drive"), c.Query("page"))
	c.JSON(200, mapd)
}

// PoolList is an API handler
func PoolList(c *gin.Context) {
	mapd := result.AllPools()
	c.JSON(200, mapd)
}

// PoolCandidate is an API handler
func PoolCandidate(c *gin.Context) {
	mapd := result.PoolCandidates(c.Param("pool"), c.Query("page"))
	c.JSON(200, mapd)
}
