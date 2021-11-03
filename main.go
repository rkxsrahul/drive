package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/health"
	"git.xenonstack.com/util/drive-portal/src/logger"
	"git.xenonstack.com/util/drive-portal/src/redisdb"
	"git.xenonstack.com/util/drive-portal/src/routes"
	"git.xenonstack.com/util/drive-portal/src/scheduler"
	"git.xenonstack.com/util/drive-portal/src/test"
	"git.xenonstack.com/util/drive-portal/src/types"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

func main() {
	// setup zap logger
	level, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = -1
	}
	err = logger.Init(level, os.Getenv("LOG_TYPE"), os.Getenv("LOG_ENVIRONMENT"))
	log.Println("zap logger", err)
	if err != nil {
		zap.S().Error(err)
		return
	}

	// initialize redis
	err = redisdb.Initialise()
	if err != nil {
		zap.S().Error(err)
		zap.S().Error("error in connecting to redis")
		os.Exit(1)
	}

	//number of ideal connections
	var ideal int
	idealStr := os.Getenv("IDEAL_CONNECTIONS")
	if idealStr == "" {
		ideal = 100
	} else {
		ideal, _ = strconv.Atoi(idealStr)
	}
	// connecting db using connection string
	db, err := gorm.Open("postgres", config.DBConfig())
	if err != nil {
		zap.S().Error(err)
		return
	}
	// close db instance whenever whole work completed
	defer db.Close()
	db.DB().SetMaxIdleConns(ideal)
	db.DB().SetConnMaxLifetime(1 * time.Hour)
	config.DB = db

	//setting up Database
	types.CreateDBTablesIfNotExists()
	test.AssignStartTime()
	scheduler.Start()
	defer scheduler.Cron.Stop()

	// intialize gin router
	router := gin.Default()
	//set zap logger as std logger
	router.Use(ginzap.Ginzap(logger.Log, time.RFC3339, true))
	//allowing CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization")
	corsConfig.AddAllowMethods("DELETE")
	router.Use(cors.New(corsConfig))
	staticRoutes := router.Group("/")
	staticRoutes.Static("/images", "./images")
	// creating healthz end point, logs end point and end end point
	router.GET("/healthz", apiHealthz)
	if os.Getenv("ENVIRONMENT") != "production" {
		router.GET("/logs", checkToken, readLogs)
		router.GET("/end", checkToken, readEnv)
	}
	//v1 Protected and non protected routes
	routes.V1Routes(router)
	// run gin router on specific port
	zap.S().Info("server is running on ... ", config.AuthServerPort)
	router.Run(":" + config.AuthServerPort)
}

// @Title apiHealthz
// @Description retrieves health of drive portal service
// @Produce json
// @Success 200 {object} string "ok"
// @Failure 503 {object} string "Service Unavailable"
// @Router /healthz [get]
func apiHealthz(c *gin.Context) {
	err := health.Healthz()
	if err != nil {
		c.JSON(500, gin.H{
			"error": err,
		})
		return
	}

	c.Writer.Write([]byte("ok"))
}

// serving info file at browser
func readLogs(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "info.txt")
}

// read all environment variables set and pass to browser
func readEnv(c *gin.Context) {
	env := make([]string, 0)
	for _, pair := range os.Environ() {
		env = append(env, pair)
	}
	c.JSON(200, gin.H{
		"environments": env,
	})
}

// check header is set or not for secured api
func checkToken(c *gin.Context) {
	xt := c.Request.Header.Get("X-TOKEN")
	if xt != "xyz" {
		c.Abort()
		c.JSON(401, gin.H{"message": "You are not authorised."})
		return
	}
	c.Next()
}
