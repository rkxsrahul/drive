package api

import (
	"log"

	"git.xenonstack.com/util/drive-portal/config"

	"github.com/gin-gonic/gin"
)

func DeleteUserDetails(c *gin.Context) {
	email := c.Query("email")

	row := config.DB.Exec("delete from drive_users where user_email=?", email).RowsAffected
	log.Println(row)
	row = config.DB.Exec("delete from answers where email=?", email).RowsAffected
	log.Println(row)
	row = config.DB.Exec("delete from user_sessions where email=?", email).RowsAffected
	log.Println(row)
	row = config.DB.Exec("delete from results.pool_analytical where email=?", email).RowsAffected
	log.Println(row)

	c.JSON(200, gin.H{
		"error":   false,
		"message": "user deleted",
	})
}
