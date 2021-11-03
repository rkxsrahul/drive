package config

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

var (
	// logout other sessions of user
	IsLogoutOthers string = os.Getenv("IS_LOGOUT_OTHER")
	PrivateKey     string = os.Getenv("HIRING_PRIVATE_KEY")
	//HTTP Server Configuration ==========
	AuthServerPort string = os.Getenv("DRIVE_PORTAL_HTTP_PORT")

	//bucket to store images
	Bucket    = os.Getenv("AWS_BUCKET")
	Region    = os.Getenv("AWS_REGION")
	AssetLink = os.Getenv("ASSET_SERVER_LINK")

	//redis configuration
	//Database of redis
	RedisDB string = os.Getenv("HIRING_REDIS_DB")
	//host address of redis
	RedisHost string = os.Getenv("HIRING_REDIS_HOST")
	//port number of redis
	RedisPort string = os.Getenv("HIRING_REDIS_PORT")
	// password of redis
	RedisPass string = os.Getenv("HIRING_REDIS_PASS")

	//configuration of cockroach db
	DBName string = os.Getenv("DRIVE_PORTAL_DB_NAME")
	DbUser string = os.Getenv("DRIVE_PORTAL_DB_USER")
	DbPass string = os.Getenv("DRIVE_PORTAL_DB_PASS")
	DbHost string = os.Getenv("DRIVE_PORTAL_DB_HOST")
	DbPort string = os.Getenv("DRIVE_PORTAL_DB_PORT")
	DBType string = os.Getenv("DRIVE_PORTAL_DB_TYPE")

	//private key to commuicate with test portal
	TestPortalKey string = os.Getenv("TEST_PORTAL_PRIVATE_KEY")

	DisableMail string = os.Getenv("MAIL_DISABLE")

	// mail templates and images path
	MailPath string = os.Getenv("MAIL_TEMPLATES_PATH")

	DB *gorm.DB
)

const (
	DBConnSSL     string        = "disable"
	JWTExpireTime time.Duration = time.Hour * 24
)

func init() {
	if MailPath == "" {
		MailPath = "./src"
	}
}

func DBConfig() string {
	if DBType == "postgres" {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", DbHost, DbPort, DbUser, DbPass, DBName, DBConnSSL)
	}

	// creating db connection string
	str := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		DbUser,
		DbPass,
		DbHost,
		DbPort,
		DBName, DBConnSSL)

	return str
}
