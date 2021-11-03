package pool

import (
	"time"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/methods"
	"git.xenonstack.com/util/drive-portal/src/questions"
	"git.xenonstack.com/util/drive-portal/src/types"

	"errors"

	"go.uber.org/zap"
)

func CreatePool(name string) error {
	// convert pool name in slug
	slug := methods.SlugOfString(name)
	if slug == "" {
		return errors.New("Please enter valid pool name contains atleast one alphabet or number")
	}

	// saving data according to pool structure
	poolData := types.Pool{
		Id:        slug,
		Name:      name,
		Date:      time.Now().Unix(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := config.DB

	dberr := db.Create(&poolData)

	if dberr.Error != nil {
		println("Error Created ")
		zap.S().Error(dberr.Error)
		return dberr.Error
	}
	return nil
}

//===========================================================================//

func ListPool() ([]bodyTypes.PoolList, error) {
	// intialise return result variable
	pools := []bodyTypes.PoolList{}

	db := config.DB
	db.Raw("select id,name,date from pools").Scan(&pools)
	// check user had not submitted more question then no of question in a pool
	for i := 0; i < len(pools); i++ {
		list, err := questions.ListQuestions(pools[i].Id)
		if err != nil {
			zap.S().Error(err)
			return pools, err
		}
		pools[i].TotalQuestion = len(list)
	}
	return pools, nil
}

//==============================================================================//

func PoolDetails(id string) (types.Pool, error) {
	// intialise return result variable
	pool := []types.Pool{}

	db := config.DB

	db.Raw("select id, name, created_at from pools where id=?", id).Scan(&pool)
	if len(pool) == 0 {
		return types.Pool{}, errors.New("No pool exists")
	}
	return pool[0], nil
}

//===================================================================================//

func EditPool(name, id string) error {
	// convert pool name in slug
	slug := methods.SlugOfString(name)
	if slug == "" {
		return errors.New("Please enter valid pool name contains atleast one alphabet or number")
	}

	db := config.DB
	pool := types.Pool{
		Id:        slug,
		Name:      name,
		UpdatedAt: time.Now(),
	}
	row := db.Model(types.Pool{}).Where("id=?", id).Update(&pool).RowsAffected
	if row == 0 {
		println("Row : ", row)
		return errors.New("Pool doesn't exists ")

	}
	// update in questions table also
	db.Exec("update questions set pool_id=? where pool_id=?", slug, id)

	// update in test table also
	db.Exec("update test_pools set pool_id=? where pool_id=?", slug, id)
	return nil
}

// //=============================================================================================//

func DeletePool(id string) error {
	db := config.DB

	var count int64
	db.Model(&types.TestPool{}).Where("pool_id=?", id).Count(&count)
	if count != 0 {
		return errors.New("You cannot delete this pool because it is assigned to a test. So first delete that test")
	}

	go func() {
		db.Exec("delete from results.pool_analytical where pool_id='" + id + "';")
	}()

	//remove pool
	row := db.Exec("delete from pools where id=?", id).RowsAffected
	if row == 0 {
		return errors.New("Pool doesn't exists")
	}

	return nil
}
