package questions

import (
	"fmt"
	"strconv"
	"strings"

	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
)

//read data from csv and make a structure
func CSVQuestions(data [][]string, pool string) (int, string, error) {
	// check pool exists or not
	code, err := CheckPool(pool)
	if err != nil {
		return code, "", err
	}

	// dry run to check all questions are correct
	not := make([]int, 0)
	for i := 1; i < len(data); i++ {
		if data[i][0] == "" {
			not = append(not, i+1)
		} else {
			question := bodyTypes.Questions{}
			question.Title = data[i][0]
			question.Type = "string"
			options := make([]bodyTypes.Options, 0)
			for j := 1; j < len(data[i]); j = j + 2 {
				if data[i][j] != "" && data[i][j+1] != "" {
					val, err := strconv.ParseBool(data[i][j+1])
					if err == nil {
						options = append(options, bodyTypes.Options{
							Value:     data[i][j],
							Type:      "string",
							IsCorrect: val,
						})
					}
				}
			}
			if len(options) == 0 {
				not = append(not, i+1)
			} else {
				question.Options = options
				_, err := validateQuestion(question)
				if err != nil {
					not = append(not, i+1)
				}
			}
		}
	}

	strNot := make([]string, 0)
	for i := 0; i < len(not); i++ {
		strNot = append(strNot, strconv.Itoa(not[i]))
	}

	if len(not) == 0 {
		// if there is no incorrect question add question in database
		go func() {
			not = make([]int, 0)
			for i := 1; i < len(data); i++ {
				if data[i][0] == "" {
					not = append(not, i+1)
				} else {
					question := bodyTypes.Questions{}
					question.Title = data[i][0]
					question.Type = "string"
					options := make([]bodyTypes.Options, 0)
					for j := 1; j < len(data[i]); j = j + 2 {
						if data[i][j] != "" && data[i][j+1] != "" {
							val, err := strconv.ParseBool(data[i][j+1])
							if err == nil {
								options = append(options, bodyTypes.Options{
									Value:     data[i][j],
									Type:      "string",
									IsCorrect: val,
								})
							}
						}
					}
					if len(options) == 0 {
						not = append(not, i+1)
					} else {
						question.Options = options
						_, err := AddQuestion(pool, question, false)
						if err != nil {
							not = append(not, i+1)
						}
					}
				}
			}
		}()
		return 200, "All questions inserted successfully", nil
	}
	str := fmt.Sprint(len(data)-len(not)-1, " questions inserted and these questions (", strings.Join(strNot, ","), ") are not inserted.")
	return 200, str, nil
}
