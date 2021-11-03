package api

import (
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/test"
	"github.com/gin-gonic/gin"
)

// api handler for creating test
func CreateTest(c *gin.Context) {
	// fetching test data from request body
	var request_data bodyTypes.Test
	if err := c.BindJSON(&request_data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Please pass the required fields."})
		return
	}

	//creating test and checking error
	code, err := test.CreateTest(request_data)
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": "Test succesfully created",
	})
}

//=======================================================================================================//

// api handler for listing test
func ListTests(c *gin.Context) {

	// fetch test list from db
	list, err := test.ListTest()
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"tests": list,
	})
}

//========================================================================================================//

// api handler for test details
func TestDetail(c *gin.Context) {

	// fetch test details on basis of id pass in url
	test, err := test.TestDetails(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"test":  test,
	})
}

//========================================================================================================//

// api handler for editing test
func EditTest(c *gin.Context) {
	// fetching test name from request body
	var request_data bodyTypes.Test
	if err := c.BindJSON(&request_data); err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "Please pass the required fields."})
		return
	}

	//edit test on basis of id and send new name of test
	code, err := test.EditTest(request_data, c.Param("id"))
	if err != nil {
		c.JSON(code, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(code, gin.H{
		"error":   false,
		"message": "Test succesfully edited",
	})
}

//=======================================================================================================//

// api handler for deleting test
func DeleteTest(c *gin.Context) {

	// delete test on basis of id pass in url
	err := test.DeleteTest(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{
			"error":   true,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"error":   false,
		"message": "Test succesfully deleted",
	})
}
