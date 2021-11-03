package college

import (
	"errors"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/methods"
	"git.xenonstack.com/util/drive-portal/src/types"
)

func Add(data bodyTypes.College) (int, error) {
	//create slug of college name + location
	id := methods.SlugOfString(data.Name + "-" + data.Location)

	db := config.DB

	if id == "" {
		return 400, errors.New("Name and location cannot be blank.")
	}
	//check there is already same college or not
	colleges := []types.College{}
	db.Raw("select id from colleges where id=?", id).Scan(&colleges)
	if len(colleges) != 0 {
		return 400, errors.New("Cannot add duplicate college")
	}

	//saving data according to college structure
	college := types.College{
		Id:        id,
		Name:      data.Name,
		Location:  data.Location,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// save college details in table
	err := db.Create(&college).Error
	if err != nil {
		zap.S().Error(err)
	}
	return 200, err
}

func List(current int, search string) (int64, []bodyTypes.CollegeList, error) {
	search = strings.Replace(search, " ", "%' AND name ilike '%", -1)
	search = "name ilike '%" + search
	search += "%'"

	db := config.DB
	//find total number of colleges
	var total int64
	db.Model(types.College{}).Where(search).Count(&total)

	if total == 0 {
		log.Println("there is no college")
		return total, []bodyTypes.CollegeList{}, nil
	}

	//check current page is valid
	totalPages := int(total / 10)
	if total%10 != 0 {
		totalPages = totalPages + 1
	}
	if current <= 0 || current > totalPages {
		log.Println("wrong input")
		return total, []bodyTypes.CollegeList{}, errors.New("please pass valid current page number")
	}

	// fetch colleges detail from db
	colleges := []types.College{}
	db.Where(search).Order("name").Limit(10).Offset((current - 1) * 10).Find(&colleges)

	// save data according to college list body type
	list := make([]bodyTypes.CollegeList, 0)
	for i := 0; i < len(colleges); i++ {
		list = append(list, bodyTypes.CollegeList{
			Id:       colleges[i].Id,
			Name:     colleges[i].Name,
			Location: colleges[i].Location,
		})
	}
	return total, list, nil
}

func View(id string) (bodyTypes.College, int, error) {

	db := config.DB

	// fetch college detail from db
	colleges := []types.College{}
	db.Raw("select id, name, location from colleges where id=?", id).Scan(&colleges)

	if len(colleges) == 0 {
		return bodyTypes.College{}, 400, errors.New("No college exist with this id")
	}
	return bodyTypes.College{
		Name:     colleges[0].Name,
		Location: colleges[0].Location,
	}, 200, nil
}

func Update(id string, data bodyTypes.College) (int, error) {
	//create slug of college name + location
	newId := methods.SlugOfString(data.Name + "-" + data.Location)

	if newId == "" {
		return 400, errors.New("Please enter valid college details contains atleast one alphabet or number")
	}

	db := config.DB

	//check there is already same college or not
	colleges := []types.College{}
	db.Raw("select id from colleges  where id=?", newId).Scan(&colleges)
	if len(colleges) != 0 {
		return 400, errors.New("Cannot add duplicate college")
	} else {
		//saving data according to college structure
		college := types.College{
			Id:        newId,
			Name:      data.Name,
			Location:  data.Location,
			UpdatedAt: time.Now(),
		}
		//update data in db
		row := db.Model(types.College{}).Where("id=?", id).Update(&college).RowsAffected
		if row == 0 {
			return 400, errors.New("No college exist with this id")
		}
		return 200, nil
	}
}

func Delete(id string) (int, error) {
	db := config.DB

	// delete college detail from db
	row := db.Exec("delete from colleges where id=?", id).RowsAffected
	if row == 0 {
		return 400, errors.New("No college exist with this id")
	}
	return 200, nil
}
