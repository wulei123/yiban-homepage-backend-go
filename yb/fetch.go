package yb

import (
	"time"
	"../configuration"
	"../ybtempl"
	"io/ioutil"
	"log"
	"encoding/json"
	"os"
)
func updateData(notice_1 []ybtempl.NoticeContent,notice_2 []ybtempl.NoticeContent,schoolIntro ybtempl.SchoolIntroTempl){
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
}
func UpdateYBData(){
	for true{
		c := Login(configuration.Account,configuration.Password)
		if CheckLogin(c) {
			notice_1 := getNotices(c,notice1_url)
			notice_2 := getNotices(c,notice2_url)
			schoolIntro := getSchoolIntro(c)
			updateData(notice_1,notice_2,schoolIntro)
			time.Sleep(configuration.IntervalHours*time.Hour)
		}
	}
}