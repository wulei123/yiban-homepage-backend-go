package configuration

import (
	"log"
	"flag"
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

type Conf struct {
	ServerPort int `json:"server_port"`
	Admins []Admin `json:"admins"`
	ImageUrl string `json:"image_url"`
	SessionSecret string `json:"session_secret"`
	DataUrl string `json:"data_url"`
}
type Admin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	ServerPort string
	Admins map[string]string
	ImageUrl string
	SessionSecret string
	DataUrl string
)

func init(){
	config_file_name := flag.String("config","config.json","配置文件路径")
	flag.Parse()
	log.Printf("正在加载配置文件 : %s", *config_file_name)
	file , err := os.Open(*config_file_name)
	if err != nil{
		log.Printf("打开配置文件错误 : %v",err)
		return
	}
	defer file.Close()
	data , err := ioutil.ReadAll(file)
	if err != nil{
		log.Printf("读取配置文件错误 : %v",err)
		return
	}
	conf := Conf{}
	err = json.Unmarshal(data, &conf)
	if err != nil{
		log.Printf("解析配置文件错误 : %v", err)
		return
	}
	ServerPort = fmt.Sprintf(":%d",conf.ServerPort)
	Admins = make(map[string]string)
	for _,admin := range conf.Admins{
		Admins[admin.Username] = admin.Password
	}
	ImageUrl = conf.ImageUrl
	SessionSecret = conf.SessionSecret
	DataUrl = conf.DataUrl
	log.Println("配置文件加载成功")
}