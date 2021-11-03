package questions

import (
	"log"
	"strings"

	"git.xenonstack.com/util/drive-portal/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
)

func DeleteImage(url string) error {
	// fetching file name from url
	key := strings.TrimPrefix(url, "https://")
	key = strings.TrimPrefix(key, "http://")
	keys := strings.Split(key, "/")

	//creating aws session
	var sess *session.Session
	var err error
	if config.Region == "" {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String("ap-south-1")},
		)
		if err != nil {
			zap.S().Error("session error .... ", err)
			return err
		}
	} else {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String(config.Region)},
		)
		if err != nil {
			zap.S().Error("session error .... ", err)
			return err
		}
	}
	// Create S3 service client
	svc := s3.New(sess)
	// Delete the item
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(config.Bucket), Key: aws.String(keys[len(keys)-1])})
	if err != nil {
		zap.S().Error(err)
		return err
	}

	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(config.Bucket),
		Key:    aws.String(keys[len(keys)-1]),
	})
	if err != nil {
		zap.S().Error(err)
		return err
	}
	log.Println("Object ", keys[len(keys)-1], " succesfully deleted")
	return nil
}
