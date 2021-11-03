package userapis

import (
	"strings"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/src/drive"
	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

func UserDriveList(c *gin.Context) {
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	email, ok := claims["email"]
	if !ok {
		c.JSON(500, gin.H{
			"error":   true,
			"message": "Please login again",
		})
		return
	}
	// fetch drives detail list
	list, err := drive.UserDriveDetails(strings.ToLower(email.(string)))
	if err != nil {
		zap.S().Error(err)
		c.JSON(500, gin.H{
			"error":   true,
			"list":    list,
			"message": "Unable to fetch drives",
		})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"list":  list,
	})
	return
}
