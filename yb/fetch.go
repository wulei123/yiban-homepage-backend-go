package yb

import (
	"time"
	"../configuration"
	"../ybtempl"
	"io/ioutil"
	"log"
	"encoding/json"
	"os"
	"net/http"
)

func updateData(notice_1 []ybtempl.NoticeContent,notice_2 []ybtempl.NoticeContent,schoolIntro ybtempl.SchoolIntroTempl,teachers []ybtempl.TeacherTempl){
	var ybData ybtempl.YBData
	fileData, err := ioutil.ReadFile(configuration.DataUrl)
	if err != nil{
		log.Fatalf("open data file error , notice not changed: %v\n",err)
		return
	}
	err = json.Unmarshal(fileData,&ybData)
	if err != nil{
		log.Fatalf("decode data.json failed : %v\n",err)
		return
	}
	for i := 0; i < len(ybData.Notices.Notice1.Content);i++{
		ybData.Notices.Notice1.Content[i] = notice_1[i]
	}
	for i := 0; i < len(ybData.Notices.Notice2.Content);i++{
		ybData.Notices.Notice2.Content[i] = notice_2[i]
	}
	ybData.SchoolIntro = schoolIntro
	ybData.Teachers = teachers
	file, err := os.OpenFile(configuration.DataUrl, os.O_WRONLY|os.O_TRUNC, 0666)
	defer file.Close()

	if err != nil{
		log.Fatalf("write in data.json failed : %v\n",err)
		return
	}
	dataJson, err := json.Marshal(ybData)
	if err != nil{
		log.Fatalf("write in data.json failed : %v\n",err)
		return
	}
	file.Write(dataJson)
	log.Println(string(dataJson))
}
func UpdateData(c *http.Client){
	log.Println("获取最新公告")
	notice_1 := getNotices(c,notice1_url)
	log.Println("获取校园活动")
	notice_2 := getNotices(c,notice2_url)
	log.Println("获取机构群和本校成员人数")
	schoolIntro := getSchoolIntro(c)
	log.Println("获取名师推荐")
	teachers := getTeachers(c)
	log.Printf("数据获取完毕，更新%s\n",configuration.DataUrl)
	updateData(notice_1,notice_2,schoolIntro,teachers)
	log.Println("数据更新完成")
}
func UpdateYBData(){
	for true{
		log.Println("登陆公共账号")
		c := Login(configuration.Account,configuration.Password)
		if CheckLogin(c) {
			log.Println("登陆成功")
			UpdateData(c)
			time.Sleep(configuration.IntervalHours*time.Hour)
		}
	}
}