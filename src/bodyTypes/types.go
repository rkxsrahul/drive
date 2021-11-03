// file contains structures that to be used in fetching information from
// request and send body during api call in json format
package bodyTypes

// structure for binding job data
type Jobs struct {
	Id       string   `json:"id"`
	Name     string   `json:"name" binding:"required"`
	Summary  string   `json:"summary" binding:"required"`
	Location string   `json:"location" binding:"required"`
	Body     string   `json:"body" binding:"required"`
	TeamId   string   `json:"teamId" binding:"required"`
	TeamName string   `json:"teamName" binding:"required"`
	Skills   []string `json:"skills" binding:"required"`
}

// structure for sending list of jobs
type JobList struct {
	TeamName string        `json:"teamName"`
	TeamId   string        `json:"teamId"`
	Jobs     []TeamJobList `json:"jobs"`
}

type TeamJobList struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Summary  string   `json:"summary"`
	Location string   `json:"location"`
	Skills   []string `json:"skills"`
}

// structure for binding team job data
type JobTeam struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

//structure for binding college data
type College struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}

//structure for sending list of colleges
type CollegeList struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// structure for binding pool data
type Pool struct {
	Name string `json:"name" binding:"required"`
}

// structure for sending list of pools
type PoolList struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Date          int64  `json:"date"`
	TotalQuestion int    `json:"totalQuestion"`
}

// structures for binding questions data
type Questions struct {
	Title     string    `json:"title" binding:"required"`
	Type      string    `json:"type" binding:"required"`
	Options   []Options `json:"options" binding:"required"`
	ImagesUrl string    `json:"images_url"`
}

//structre for binding options data
type Options struct {
	Value     string `json:"value" binding:"required"`
	Type      string `json:"type" binding:"required"`
	ImagesUrl string `json:"images_url"`
	IsCorrect bool   `json:"is_correct"`
}

// structure for sending question list in output
type QuestionList struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

// structure for binding test data
type Test struct {
	Name           string     `json:"name" binding:"required"`
	Duration       string     `json:"duration" binding:"required"`
	TotalQuestions int        `json:"totalQuestions"`
	Pools          []TestPool `json:"pools" binding:"required"`
}

type TestPool struct {
	PoolName      string `json:"poolName"`
	PoolId        string `json:"poolId" binding:"required"`
	TotalQuestion int    `json:"totalQuestions"`
	NoOfQuestions int    `json:"noOfQuestions" binding:"required"`
}

// structure for sending test list in output
type TestList struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Duration int64    `json:"duration"`
	Total    int      `json:"total"`
	Date     int64    `json:"date"`
	Pools    []string `json:"pools"`
}

// input structure for binding drive data
type Drive struct {
	Type      string `json:"type" binding:"required"`
	Name      string `json:"name" binding:"required"`
	TestId    string `json:"test_id" binding:"required"`
	Start     string `json:"start" binding:"required"`
	End       string `json:"end" binding:"required"`
	StartStr  string `json:"startStr"`
	EndStr    string `json:"endStr"`
	CollegeId string `json:"college_id"`
}

// output structure for drive details
type DriveDetails struct {
	Id        string  `json:"id"`
	Type      string  `json:"type"`
	Name      string  `json:"name"`
	TestId    string  `json:"test_id"`
	Start     string  `json:"start"`
	End       string  `json:"end"`
	StartStr  string  `json:"startStr"`
	EndStr    string  `json:"endStr"`
	CollegeId string  `json:"college_id"`
	College   College `json:"college"`
}

// output structure for drive list
type DriveList struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Start string `json:"start"`
	End   string `json:"end"`
	Users int64  `json:"users"`
}

// output structure for drive summary
type SummaryDrive struct {
	DriveDetails
	Ongoing    int64 `json:"ongoing"`
	NotStarted int64 `json:"not_started"`
	Completed  int64 `json:"completed"`
}

// input structure for drive users
type DriveUser struct {
	Email string `json:"email" binding:"required"`
}

// output structure for result
type Result struct {
	Email      string `json:"email"`
	Correct    int    `json:"correct"`
	Wrong      int    `json:"wrong"`
	Attempted  int    `json:"attempted"`
	Total      int    `json:"total"`
	TimeTaken  int64  `json:"time_taken"`
	Restart    int64  `json:"restart"`
	Browser    int64  `json:"browser"`
	TestStatus string `json:"test_status"`
}

// output structure for pool wise result
type PoolResult struct {
	PoolName  string `json:"pool_name"`
	PoolId    string `json:"pool_id"`
	Correct   int    `json:"correct"`
	Wrong     int    `json:"wrong"`
	Attempted int    `json:"attempted"`
	Total     int    `json:"total"`
	TimeTaken int64  `json:"time_taken"`
}

// output structure for user detail Result
type UserResult struct {
	Result Result       `json:"result"`
	Pool   []PoolResult `json:"pool_result"`
}

// output structre for completed test
type Completetest struct {
	Drive     string `json:"drive"`
	Test      string `json:"test"`
	Completed int64  `json:"completed"`
}
