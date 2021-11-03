package health

import (
	"git.xenonstack.com/util/drive-portal/config"
	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Healthz() error {

	//checking health of cockroachdb
	// connecting to db
	db, err := gorm.Open("postgres", config.DBConfig())
	if err != nil {
		zap.S().Error(err)
		return err
	}
	// close db instance whenever whole work completed
	defer db.Close()
	return nil
}
