package job

import (
	"errors"
	"time"

	"github.com/lib/pq"
	"go.uber.org/zap"

	"git.xenonstack.com/util/drive-portal/config"
	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"git.xenonstack.com/util/drive-portal/src/methods"
	"git.xenonstack.com/util/drive-portal/src/types"
)

func AddJob(data bodyTypes.Jobs) error {
	//check team exists
	_, err := TeamDetails(data.TeamId)
	if err != nil {
		zap.S().Error(err)
		return err
	}
	// convert pool name in slug
	slug := methods.SlugOfString(data.Name)
	if slug == "" {
		return errors.New("Please enter valid job title contains atleast one alphabet or number")
	}

	db := config.DB

	// saving data according to job team structure
	jobData := types.Jobs{
		Id:        slug,
		Name:      data.Name,
		Summary:   data.Summary,
		Location:  data.Location,
		Body:      data.Body,
		TeamId:    data.TeamId,
		TeamName:  data.TeamName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	//save data in db
	dberr := db.Create(&jobData)
	if dberr.Error != nil {
		zap.S().Error(dberr.Error)
		return errors.New("This job is already exists")
	}

	pdb := db.DB()
	_, err = pdb.Exec("UPDATE jobs set skills=$2 where id=$1", slug, pq.Array(data.Skills))
	if err != nil {
		zap.S().Error(err)
		return errors.New("This job is already exists")
	}

	return nil
}

func ListJobs() ([]bodyTypes.JobList, error) {
	db := config.DB

	//fetch all teams
	teams := []types.JobTeam{}
	db.Raw("select id, name, description from job_teams").Scan(&teams)

	finalList := make([]bodyTypes.JobList, 0)
	// fetch jobs according to team
	for i := 0; i < len(teams); i++ {
		//fetch job list by team
		jobsList, err := ListJobByTeam(teams[i].Id)
		if err != nil {
			zap.S().Error(err)
			return []bodyTypes.JobList{}, err
		}
		finalList = append(finalList, bodyTypes.JobList{
			TeamName: teams[i].Name,
			TeamId:   teams[i].Id,
			Jobs:     jobsList,
		})
	}
	return finalList, nil
}

func ListJobByTeam(teamId string) ([]bodyTypes.TeamJobList, error) {
	db := config.DB
	pdb := db.DB()
	jobs := []types.Jobs{}
	db.Raw("select id,name,summary,location from jobs where team_id=?", teamId).Scan(&jobs)
	jobsList := make([]bodyTypes.TeamJobList, 0)
	for j := 0; j < len(jobs); j++ {
		var skills []string
		err := pdb.QueryRow(`SELECT skills from jobs WHERE id=$1`, jobs[j].Id).Scan(pq.Array(&skills))
		if err != nil {
			continue
		}
		//append jobs in a slice
		jobsList = append(jobsList, bodyTypes.TeamJobList{
			Id:       jobs[j].Id,
			Name:     jobs[j].Name,
			Summary:  jobs[j].Summary,
			Skills:   skills,
			Location: jobs[j].Location,
		})
	}
	return jobsList, nil
}

func JobDetails(teamId string, jobId string) (bodyTypes.Jobs, error) {
	db := config.DB
	jobs := []types.Jobs{}
	db.Raw("select name,summary,body,location from jobs where team_id=? AND id=?", teamId, jobId).Scan(&jobs)

	if len(jobs) == 0 {
		return bodyTypes.Jobs{}, errors.New("This job is no longer exists")
	}

	team := []types.JobTeam{}
	db.Select("name").Where("id=?", teamId).Find(&team)

	pdb := db.DB()
	var skills []string
	_ = pdb.QueryRow(`select skills from jobs where team_id= $1 and id= $2`, teamId, jobId).Scan(pq.Array(&skills))

	return bodyTypes.Jobs{
		Id:       jobId,
		Name:     jobs[0].Name,
		Summary:  jobs[0].Summary,
		Location: jobs[0].Location,
		Body:     jobs[0].Body,
		TeamName: team[0].Name,
		TeamId:   teamId,
		Skills:   skills,
	}, nil
}

func DeleteJob(teamId string, jobId string) error {
	db := config.DB
	row := db.Exec("delete from jobs where team_id=? and id=?", teamId, jobId).RowsAffected

	if row == 0 {
		return errors.New("This job is no longer exists")
	}

	return nil
}

func UpdateJob(teamId string, jobId string, data bodyTypes.Jobs) error {
	db := config.DB

	err := DeleteJob(teamId, jobId)
	if err != nil {
		return err
	}
	// convert pool name in slug
	slug := methods.SlugOfString(data.Name)
	// saving data according to job team structure
	jobData := types.Jobs{
		Id:        slug,
		Name:      data.Name,
		Summary:   data.Summary,
		Location:  data.Location,
		Body:      data.Body,
		TeamId:    data.TeamId,
		TeamName:  data.TeamName,
		Skills:    data.Skills,
		UpdatedAt: time.Now(),
	}
	//save data in db
	dberr := db.Create(&jobData)
	if dberr.Error != nil {
		zap.S().Error(dberr.Error)
		return errors.New("This job is already exists")
	}
	return nil
}
