package routes

import (
	"os"

	"git.xenonstack.com/util/drive-portal/src/api"
	"git.xenonstack.com/util/drive-portal/src/jwt"
	"git.xenonstack.com/util/drive-portal/src/userapis"
	"github.com/gin-gonic/gin"
)

func V1Routes(router *gin.Engine) {
	v1 := router.Group("/v1")

	v1.GET("/job", api.ListJobs)
	v1.GET("/job/:team", api.ListJobsByTeam)
	v1.GET("/job/:team/:id", api.JobDetails)
	v1.GET("/college", api.ListColleges)
	v1.GET("/pool_csv/:file", api.ServePoolCSV)
	v1.GET("/result_csv/:file", api.ServeResultCSV)
	//setting up middleware for protected apis
	authMiddleware := jwt.MwInitializer()
	//Protected resources ======================================
	v1.Use(authMiddleware.MiddlewareFunc())
	{
		// adding custom middleware for checking token validity
		v1.Use(api.CheckTokenValidity)
		{
			//user apis
			user := v1.Group("/user")

			user.GET("/drives", userapis.UserDriveList)
			user.Use(api.CheckUser)
			{
				user.GET("/completed", userapis.CompletedTest)
				user.GET("/drives/:id/:testid", userapis.CheckTest)
				user.GET("/test/:id", userapis.TestDetails)
			}
			//admin apis
			v1.Use(api.CheckAdmin)
			{
				// api endpoint related to job team
				v1.POST("/teamjob", api.AddTeamJob)
				v1.GET("/teamjob", api.ListTeamJob)
				v1.GET("/teamjob/:id", api.TeamDetail)
				v1.PUT("/teamjob/:id", api.UpdateTeam)
				v1.DELETE("/teamjob/:id", api.DeleteTeam)

				// api endpoints related to jobs
				v1.POST("/job", api.AddJob)
				v1.DELETE("/job/:team/:id", api.DeleteJob)
				v1.PUT("/job/:team/:id", api.UpdateJob)

				//	api endpoints related to colleges
				v1.POST("/college", api.AddCollege)
				v1.GET("/college/:id", api.ViewCollege)
				v1.PUT("/college/:id", api.UpdateCollege)
				v1.DELETE("/college/:id", api.DeleteCollege)
				v1.POST("/csv_college", api.CSVCollege)

				// api endpointss related to pool
				v1.POST("/pool", api.CreatePool)
				v1.GET("/pool", api.ListPools)
				v1.GET("/pool/:id", api.PoolDetail)
				v1.PUT("/pool/:id", api.EditPool)
				v1.DELETE("/pool/:id", api.DeletePool)

				// api endpoints related to questions
				v1.POST("/question/:pool", api.AddQuestion)
				v1.GET("/question/:pool", api.ListQuestions)
				v1.GET("/question/:pool/:id", api.QuestionDetails)
				v1.DELETE("/question/:pool/:id", api.DeleteQuestion)
				v1.PUT("/question/:pool/:id", api.UpdateQuestion)
				v1.POST("/csv_question/:pool", api.CSVQuestion)
				v1.GET("/csv_question/:pool", api.DownloadCSVQuestion)

				// api endpoints related to image on aws
				v1.POST("/upload_image", api.UploadImage)
				v1.DELETE("/delete_image", api.DeleteImage)

				// api endpoints related to test
				v1.POST("/test", api.CreateTest)
				v1.GET("/test", api.ListTests)
				v1.GET("/test/:id", api.TestDetail)
				v1.PUT("/test/:id", api.EditTest)
				v1.DELETE("/test/:id", api.DeleteTest)

				//api endpoints related to drive
				v1.POST("/drive", api.CreateDrive)
				v1.GET("/drive", api.ListDrives)
				v1.GET("/drive/:id", api.ViewDrive)
				v1.DELETE("/drive/:id", api.DeleteDrive)
				v1.GET("/drive/:id/summary", api.DriveSummary)

				//api endpoints related to drive Users
				v1.POST("/drive/:id/user", api.AddUser)
				v1.GET("/drive/:id/user", api.ListUsers)
				v1.DELETE("/drive/:id/user/:email", api.DeleteUser)
				v1.DELETE("/drive/:id/user", api.DeleteAllUsers)
				v1.POST("/drive/:id/csv_user", api.CSVAddUsers)

				//api endpoints related to result
				v1.GET("/result/drive/:id", api.Result)
				v1.GET("/result/drive/:id/csv", api.DownloadResult)
				v1.GET("/result/analytical/:drive", api.AnalyticalResult)
				v1.GET("/result/drivewise", api.DriveList)
				v1.GET("/result/drivewise/:drive", api.DriveCandidate)
				v1.GET("/result/poolwise", api.PoolList)
				v1.GET("/result/poolwise/:pool", api.PoolCandidate)

				//delete user details from all table api endpoint
				v1.DELETE("/user", api.DeleteUserDetails)

				// tester help endpoint
				if os.Getenv("ENVIRONMENT") != "production" {
					// toggle mail service
					v1.PUT("admin/mail/:value", api.ChangeMail)
				}
			}
		}
	}
}
