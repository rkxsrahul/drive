package college

import (
	"fmt"
	"strconv"
	"strings"

	"git.xenonstack.com/util/drive-portal/src/bodyTypes"
	"go.uber.org/zap"
)

//read data from csv and make a structure
func CSVColleges(data [][]string) string {
	not := make([]string, 0)
	for i := 1; i < len(data); i++ {
		if strings.TrimSpace(data[i][0]) == "" {
			not = append(not, strconv.Itoa(i+1))
			continue
		}
		if strings.TrimSpace(data[i][1]) == "" {
			not = append(not, strconv.Itoa(i+1))
			continue
		}
		college := bodyTypes.College{
			Name:     data[i][0],
			Location: data[i][1],
		}
		_, err := Add(college)
		if err != nil {
			zap.S().Error("error in adding college ", err)
			not = append(not, strconv.Itoa(i+1))
		}
	}
	if len(not) == 0 {
		return "All colleges inserted successfully."
	}
	return fmt.Sprint(len(data)-len(not)-1, " college inserted and these colleges (", strings.Join(not, ","), ") are not inserted.")
}
