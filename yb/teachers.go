package yb

import (
	"net/http"
	"fmt"
	"../ybtempl"
	"encoding/json"
)

const teachers_url = "http://www.yiban.cn/school/getTeacherAjax?id=5370538&page="
func getTeachers(c *http.Client) ( teachers []ybtempl.TeacherTempl){
	i := 0
	for {
		var teacher ybtempl.TeacherTempl
		res, _ := c.Get(fmt.Sprintf("%s%d",teachers_url,i))
		decoder := json.NewDecoder(res.Body)
		decoder.Decode(&teacher)
		if len(teacher.Data) == 0{
			break
		}
		for i := range teacher.Data{
			teacher.Data[i].Url = "http://www.yiban.cn"+teacher.Data[i].Url
		}
		teachers = append(teachers,teacher)
		i += 1
	}
	return teachers
}