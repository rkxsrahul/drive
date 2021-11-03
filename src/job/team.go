package job

import (
	"errors"
	"time"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/methods"
	"git.xenonstack.com/util/drive-portal/src/types"
	"go.uber.org/zap"
)

func AddTeam(data bodyTypes.JobTeam) error {
	// convert pool name in slug
	slug := methods.SlugOfString(data.Name)
	if slug == "" {
		return errors.New("Please enter valid job team title contains atleast one alphabet or number")
	}

	db := config.DB

	// saving data according to job team structure
	teamData := types.JobTeam{
		Id:          slug,
		Name:        data.Name,
		Description: data.Description,
	}
	//save data in db
	dberr := db.Create(&teamData)
	if dberr.Error != nil {
		zap.S().Error(dberr.Error)
		return errors.New("This team is already exists")
	}
	return nil
}

func ListTeams() ([]types.JobTeam, error) {
	db := config.DB

	list := []types.JobTeam{}
	db.Raw("select id,name,description from job_teams").Scan(&list)
	return list, nil
}

func TeamDetails(id string) (types.JobTeam, error) {
	db := config.DB

	team := []types.JobTeam{}
	db.Raw("select id,name,description from job_teams where id=?", id).Scan(&team)
	if len(team) == 0 {
		return types.JobTeam{}, errors.New("This team is no longer exists")
	}
	return team[0], nil
}

func DeleteTeam(id string) error {
	db := config.DB
	row := db.Exec("delete from job_teams where id=?", id).RowsAffected
	if row == 0 {
		return errors.New("This team is no longer exists")
	}
	// delete jobs and data

	// row = db.Where("job_id IN (?)", db.Table("jobs").Select("id").Where("team_id=?", id).QueryExpr()).Delete(&types.JobSkill{}).RowsAffected
	// log.Println(row)
	// row = db.Where("team_id=?", id).Delete(&types.Jobs{}).RowsAffected
	// log.Println(row)

	return nil
}

func UpdateTeam(id string, data bodyTypes.JobTeam) error {
	db := config.DB
	// convert pool name in slug
	slug := methods.SlugOfString(data.Name)
	//update in jobs table
	db.Exec("update jobs set team_id='" + slug + "', team_name='" + data.Name + "',updated_at='" + time.Now().String() + "' where team_id='" + id + "';")
	//delete old team
	err := DeleteTeam(id)
	if err != nil {
		return err
	}
	// saving data according to job team structure
	teamData := types.JobTeam{
		Id:          slug,
		Name:        data.Name,
		Description: data.Description,
	}
	//save new team data in db
	dberr := db.Create(&teamData)
	if dberr.Error != nil {
		zap.S().Error(dberr.Error)
		return errors.New("This team is already exists")
	}
	return nil
}
