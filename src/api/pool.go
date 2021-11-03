package api

import (
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/pool"
	"github.com/gin-gonic/gin"
)

// api handler for creating pool
func CreatePool(c *gin.Context) {
	// fetching pool name from request body
	var request_data bodyTypes.Pool
	if err := c.BindJSON(&request_data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Name is required field."})
		return
	}

	//creating pool and checking error
	err := pool.CreatePool(request_data.Name)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "pool succesfully created",
	})
}

//=======================================================================================================//

// api handler for listing pool
func ListPools(c *gin.Context) {

	// fetch pool list from db
	list, err := pool.ListPool()
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"pools": list,
	})
}

//========================================================================================================//

// api handler for pool details
func PoolDetail(c *gin.Context) {

	// fetch pool details on basis of id pass in url
	pool, err := pool.PoolDetails(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"pool":  pool,
	})
}

//========================================================================================================//

// api handler for editing pool
func EditPool(c *gin.Context) {
	// fetching pool name from request body
	var request_data bodyTypes.Pool
	if err := c.BindJSON(&request_data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Name is required field."})
		return
	}

	//edit pool on basis of id and send new name of pool
	err := pool.EditPool(request_data.Name, c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "pool succesfully edited",
	})
}

//=======================================================================================================//

// api handler for deleting pool
func DeletePool(c *gin.Context) {

	// delete pool on basis of id pass in url
	err := pool.DeletePool(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "pool succesfully deleted",
	})
}
