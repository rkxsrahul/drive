package api

import (
	"strings"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/redisdb"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// @Title CheckTokenValidity
// @Description middleware for checking token validity
// @Produce json
// @Param   Authorization	header	string	true	"Bearer <token>"
// @Failure 401 {object} string "Expired auth token"
// middleware for checking token validity
func CheckTokenValidity(c *gin.Context) {

	// fetch token from header
	token := c.Request.Header.Get("Authorization")
	// trim bearer from token
	token = strings.TrimPrefix(token, "Bearer ")

	// check token exist or not
	err := redisdb.CheckToken(token)
	if err != nil {
		// when token not exist
		c.Abort()
		c.JSON(401, gin.H{"error": true, "message": "Expired auth token"})
		return
	}
	c.Next()
}

// @Title CheckAdmin
// @Description middleware for checking person is admin or not
// @Produce json
// @Param   Authorization	header	string	true	"Bearer <token>"
// @Failure 401 {object} string "You are not allowed to perform these actions"
// middleware for checking token validity
func CheckAdmin(c *gin.Context) {
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	if claims["sys_role"].(string) != "admin" {
		// when person is not admin
		c.Abort()
		c.JSON(401, gin.H{"error": true, "message": "You are not allowed to perform these actions"})
		return
	}

	c.Next()
}

// @Title CheckUser
// @Description middleware for checking person is user or not
// @Produce json
// @Param   Authorization	header	string	true	"Bearer <token>"
// @Failure 401 {object} string "You are not allowed to perform these actions"
// middleware for checking token validity
func CheckUser(c *gin.Context) {
	// extracting jwt claims
	claims := jwt.ExtractClaims(c)

	if claims["sys_role"].(string) != "user" {
		// when person is not admin
		c.Abort()
		c.JSON(401, gin.H{"error": true, "message": "You are not allowed to perform these actions"})
		return
	}

	c.Next()
}

func ChangeMail(c *gin.Context) {
	config.DisableMail = c.Param("value")
}
