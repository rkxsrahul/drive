package api

import (
	"log"
	"strings"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/questions"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	// fetch file from form-data
	file, fh, err := c.Request.FormFile("image")
	if err != nil {
		zap.S().Error(err)
		c.JSON(400, gin.H{"error": true, "message": err.Error()})
		return
	}
	if file != nil {
		// fetch filename and check file is correct
		fileNameParts := strings.Split(fh.Filename, ".")
		if len(fileNameParts) < 2 {
			c.JSON(400, gin.H{"message": "invalid filename"})
			return
		}

		//creating aws session
		var sess *session.Session
		if config.Region == "" {
			sess, err = session.NewSession(&aws.Config{
				Region: aws.String("ap-south-1")},
			)
			if err != nil {
				zap.S().Error("session error .... ", err)
				c.JSON(500, gin.H{"error": true, "message": err.Error()})
				return
			}
		} else {
			sess, err = session.NewSession(&aws.Config{
				Region: aws.String(config.Region)},
			)
			if err != nil {
				zap.S().Error("session error .... ", err)
				c.JSON(500, gin.H{"error": true, "message": err.Error()})
				return
			}
		}
		// s3 uploader
		uploader := s3manager.NewUploader(sess)

		// uploading file in aws s3 bucekt
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(config.Bucket),
			Key:    aws.String(fh.Filename),
			Body:   file,
		})
		if err != nil {
			zap.S().Error("uploading error .... ", err)
			c.JSON(500, gin.H{"error": true, "message": err.Error()})
			return
		}
		log.Printf("Successfully uploaded %q to %q\n", fh.Filename, result.Location)

		// fetch location of file
		var link string
		if config.AssetLink == "" {
			link = result.Location
		} else {
			link = config.AssetLink + "/" + fh.Filename
		}

		c.JSON(200, gin.H{
			"error": false,
			"link":  link,
		})
	}

}

type Image struct {
	Url string `json:"url" binding:"required"`
}

func DeleteImage(c *gin.Context) {
	var url Image
	err := c.BindJSON(&url)
	if err != nil {
		// if there is some error passing bad status code
		c.JSON(400, gin.H{"error": true, "message": "url is required field."})
		return
	}

	err = questions.DeleteImage(url.Url)
	if err != nil {
		zap.S().Error(err)
		c.JSON(500, gin.H{"error": true, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{
		"error": false,
		"link":  "Image succesfully deleted",
	})
}
