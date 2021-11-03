package questions

import (
	"time"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/types"

	"errors"
	"strconv"

	"go.uber.org/zap"
)

func CheckPool(pool string) (int, error) {
	db := config.DB
	pools := []types.Pool{}
	db.Exec("select id from pools where id='" + pool + "';").Find(&pools)
	if len(pools) == 0 {
		return 400, errors.New("Pool not exists")
	}
	return 200, nil
}

func AddQuestion(pool string, data bodyTypes.Questions, check bool) (int, error) {
	if check {
		// check pool exists or not
		code, err := CheckPool(pool)
		if err != nil {
			return code, err
		}
	}

	// validate question data input by admin
	code, err := validateQuestion(data)
	if err != nil {
		return code, err
	}

	db := config.DB
	// creating question structure to save in db
	question := types.Questions{
		Title:    data.Title,
		Type:     data.Type,
		PoolId:   pool,
		ImageUrl: data.ImagesUrl,
	}

	db.Create(&question)

	// save options in db
	addOptions(question.Id, data.Options)

	// update correct option id in question table
	var opt []types.Options
	db.Raw("select id from options where ques_id=? and is_correct=? ", question.Id, true).Scan(&opt)
	if len(opt) == 1 {
		db.Exec("update questions set answer_id=" + strconv.Itoa(opt[0].Id) + " where id=" + strconv.Itoa(question.Id) + ";")
	}
	return 200, nil
}

func ListQuestions(pool string) ([]bodyTypes.QuestionList, error) {
	db := config.DB

	list := []types.Questions{}

	db.Raw("select id, title from questions where pool_id = ?", pool).Scan(&list)

	result := make([]bodyTypes.QuestionList, 0)

	for i := 0; i < len(list); i++ {
		result = append(result, bodyTypes.QuestionList{
			Id:    strconv.Itoa(list[i].Id),
			Title: list[i].Title,
		})
	}

	return result, nil
}

func QuestionDetail(pool, id string) (bodyTypes.Questions, error) {
	//convert id to int
	quesID, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error(err)
		return bodyTypes.Questions{}, err
	}
	db := config.DB

	// result question to be returned
	finalQuestion := bodyTypes.Questions{}

	//fetching question from db
	question := []types.Questions{}
	db.Raw("select title,type,image_url from questions where pool_id=? and id=?", pool, quesID).Scan(&question)
	if len(question) == 0 {
		// when no question exist
		return finalQuestion, errors.New("No question exist")
	}
	//setting type and title in final result question
	finalQuestion.Title = question[0].Title
	finalQuestion.Type = question[0].Type
	finalQuestion.ImagesUrl = question[0].ImageUrl

	// fetch options from db
	options := []types.Options{}
	db.Raw("select type,value,is_correct,image_url from options where ques_id=?", quesID).Scan(&options)
	finalOptions := make([]bodyTypes.Options, 0)
	for i := 0; i < len(options); i++ {
		// creating list of options
		finalOptions = append(finalOptions, bodyTypes.Options{
			Value:     options[i].Value,
			Type:      options[i].Type,
			IsCorrect: options[i].IsCorrect,
			ImagesUrl: options[i].ImageUrl,
		})
	}
	finalQuestion.Options = finalOptions
	//return final result
	return finalQuestion, nil
}

func EditQuestion(pool, id string, data bodyTypes.Questions) (int, error) {
	//convert id to int
	quesID, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error(err)
		return 400, err
	}

	// validate question data input by admin
	code, err := validateQuestion(data)
	if err != nil {
		return code, err
	}

	db := config.DB

	question := []types.Questions{}
	db.Raw("select type,image_url from questions where pool_id=? and id=?", pool, quesID).Scan(&question)
	if len(question) == 0 {
		return 400, errors.New("No question exist")
	}
	//update question details in db
	ques := types.Questions{
		Title:     data.Title,
		Type:      data.Type,
		ImageUrl:  data.ImagesUrl,
		UpdatedAt: time.Now(),
	}
	db.Model(types.Questions{}).Where("pool_id=? and id =?", pool, id).Update(&ques)

	//delete options
	err = deleteOptions(quesID)
	if err != nil {
		zap.S().Error(err)
		return 500, err
	}
	//add new options in db
	addOptions(quesID, data.Options)

	// update correct option id in question table
	var opt []types.Options
	db.Raw("select id from options where ques_id = ? and is_correct =?", id, true).Scan(&opt)
	if len(opt) == 1 {
		db.Exec("update questions set answer_id=" + strconv.Itoa(opt[0].Id) + " where id=" + id + ";")
	}
	return 200, nil
}

//Function to delete Question
func DeleteQuestion(pool, id string) error {
	quesID, err := strconv.Atoi(id)
	if err != nil {
		zap.S().Error(err)
		return err
	}

	db := config.DB

	// delete question
	info := db.Exec("delete from questions where pool_id=? and id=?", pool, quesID)

	return info.Error
}

//Function to Add options
func addOptions(quesID int, options []bodyTypes.Options) {
	db := config.DB

	for i := 0; i < len(options); i++ {
		// creating option structure to save in db
		option := types.Options{
			QuesId:    quesID,
			Type:      options[i].Type,
			Value:     options[i].Value,
			IsCorrect: options[i].IsCorrect,
			ImageUrl:  options[i].ImagesUrl,
		}
		db.Create(&option)
	}
}

// Function to Delete Options
func deleteOptions(id int) error {
	db := config.DB
	//delete options
	db.Exec("delete from options where ques_id = ?", id)
	return nil
}

func validateQuestion(data bodyTypes.Questions) (int, error) {
	// checking send data is correct or not
	// first check if type is images and length of image url array is zero
	if data.Type == "images" {
		if data.ImagesUrl == "" {
			return 400, errors.New("Your question type is image please send image url")
		}
	}

	// checking options
	count := 0
	mapd := make(map[string]int, 0)
	for i := 0; i < len(data.Options); i++ {
		// checking how many options are correct in one question
		if data.Options[i].IsCorrect {
			count++
		}
		// checking if options type is images and length of image url array is zero
		if data.Options[i].Type == "images" {
			if data.Options[i].ImagesUrl == "" {
				return 400, errors.New("Your option " + strconv.Itoa(i+1) + " type is image please send image url")
			}
		} else {
			// checking if two or more options have same value
			if mapd[data.Options[i].Value] == 0 {
				mapd[data.Options[i].Value]++
			} else {
				return 400, errors.New("Options should have different values")
			}
		}
	}
	if count != 1 {
		return 400, errors.New("Question can have only one correct option")
	}
	return 200, nil
}
