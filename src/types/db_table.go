package types

import (
	"fmt"
	"time"

	"git.xenonstack.com/util/drive-portal/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
)

//model define structure

// pool defined structure
type Pool struct {
	Id        string
	Name      string
	Date      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

//options structure
type Options struct {
	Id        int
	Type      string
	QuesId    int
	Value     string
	IsCorrect bool
	ImageUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// question structure
type Questions struct {
	Id        int
	Title     string
	Type      string
	PoolId    string
	AnswerId  int
	ImageUrl  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//to store image urls
type Images struct {
	Id        int
	SourceId  int
	Url       string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//to store tests
type Tests struct {
	Id        string
	Name      string
	Date      int64
	Duration  int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// to store test pools
type TestPool struct {
	Id            int
	TestId        string
	PoolId        string
	NoOfQuestions int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

//to store drive details
type Drive struct {
	Id        int
	Name      string
	Type      string
	TestId    string
	CollegeId string
	StartTime int64
	EndTime   int64
	StartStr  string
	EndStr    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//to store user related to that drive
type DriveUser struct {
	DriveId   int
	UserEmail string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// to store colleges
type College struct {
	Id        string
	Name      string
	Location  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// to store jobs
type Jobs struct {
	Id        string
	Name      string
	Summary   string
	Location  string
	Body      string
	TeamId    string
	TeamName  string
	Skills    []string `gorm:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// to store team of jobs
type JobTeam struct {
	Id          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// to store user session token
type UserSession struct {
	Email     string
	DriveId   int
	Token     string
	Expire    int64
	Start     int64 `gorm:"default:0"`
	TimeTaken int   `gorm:"default:0"`
	Browser   int   `gorm:"default:0"`
	Restart   int   `gorm:"default:0"`
}

// to store submitted answers and assigned questions
type Answers struct {
	Email     string
	DriveId   int
	PoolId    string
	QuesId    int
	MarkedId  int
	AnswerId  int
	QuesIndex int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// to store data after completing drive for result
type Result struct {
	Email     string
	DriveId   int
	PoolId    string
	QuesId    int
	MarkedId  int
	AnswerId  int
	QuesIndex int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// function create an table that doesn't exists
func CreateDBTablesIfNotExists() {

	// create database if database is not initialized
	createDatabaseIfNotExists()

	db := config.DB

	// creating all tables one by one but firstly checking whether table exists or not
	if !(db.HasTable(Jobs{})) {
		db.CreateTable(Jobs{})
	}
	if !(db.HasTable(JobTeam{})) {
		db.CreateTable(JobTeam{})
	}
	if !(db.HasTable(Pool{})) {
		db.CreateTable(Pool{})
	}
	if !(db.HasTable(Questions{})) {
		db.CreateTable(Questions{})
	}
	if !(db.HasTable(Options{})) {
		db.CreateTable(Options{})
	}
	if !(db.HasTable(Images{})) {
		db.CreateTable(Images{})
	}
	if !(db.HasTable(Tests{})) {
		db.CreateTable(Tests{})
	}
	if !(db.HasTable(TestPool{})) {
		db.CreateTable(TestPool{})
	}
	if !(db.HasTable(College{})) {
		db.CreateTable(College{})
	}
	if !(db.HasTable(Drive{})) {
		db.CreateTable(Drive{})
	}
	if !(db.HasTable(DriveUser{})) {
		db.CreateTable(DriveUser{})
	}
	if !(db.HasTable(UserSession{})) {
		db.CreateTable(UserSession{})
	}
	if !(db.HasTable(Answers{})) {
		db.CreateTable(Answers{})
	}
	if !(db.HasTable(Result{})) {
		db.CreateTable(Result{})
	}

	// Database migration
	dberr := db.AutoMigrate(&Jobs{},
		&JobTeam{},
		&Pool{},
		&Questions{},
		&Options{},
		&Images{},
		&TestPool{},
		&Tests{},
		&College{},
		&Drive{},
		&DriveUser{},
		&UserSession{},
		&Answers{},
		&Result{})

	zap.S().Error(dberr.Error)

	fmt.Println("Database initialized successfully.")
}

// create database is not exists
func createDatabaseIfNotExists() {
	// connecting with postgres db root db
	db, err := gorm.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbPass,
		"postgres", config.DBConnSSL))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	// executing create database query.
	db.Exec("create database " + config.DBName + ";")
}
