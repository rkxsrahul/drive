package userapis

import (
	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/src/test"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func CheckTest(c *gin.Context) {
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	// extract email from claims
	email, ok := claims["email"]
	if !ok {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}

	// extract name from claims
	name, ok := claims["name"]
	if !ok {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}

	//check test exists, user assigned to it and other things then intialize the test
	mapd, err := test.CheckAndIntializeTest(email.(string), name.(string), c.Param("id"), c.Param("testid"))
	if err != nil {
		zap.S().Error(err)
		c.JSON(500, mapd)
		return
	}
	c.JSON(200, mapd)
}

func TestDetails(c *gin.Context) {
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

func CompletedTest(c *gin.Context) {
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	// extract email from claims
	email, ok := claims["email"]
	if !ok {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}

	result, err := test.CompletedUserTests(email.(string))
	if err != nil {
		zap.S().Error(err)
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Right now not able to fetch completed test list",
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"list":  result,
	})
}
