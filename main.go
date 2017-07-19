package main

import (
	"./configuration"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/sessions"
	"fmt"
	"os"
	"time"
	"io"
)
var Data_Origin_Path = configuration.DataUrl + "data_copy.json"
type NoticesStruct struct{
	Carousels [6]Carousel `json:"carousels"`
	Notice1 Notice `json:"notice_1"`
	Notice2 Notice `json:"notice_2"`
}

type Notice struct {
	Title string `json:"title"`
	Content [6]NoticeContent `json:"content"`
}
type NoticeContent struct {
	Text string `json:"text"`
	Href 	string `json:"href"`
}
type Carousel struct {
	Name string `json:"name"`
	Href string `json:"href"`
	Src string `json:"src"`
}
type App struct {
	Name string `json:"name"`
	Href string `json:"href"`
	Icon string `json:"icon"`
}
type Data struct {
	Notices NoticesStruct `json:"notices"`
	Apps [5][]App `json:"apps"`
}
var store = sessions.NewCookieStore([]byte("my-secret"))
func allowCORS(w http.ResponseWriter){
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8086")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, If-Modified-Since") //header的类型
	//w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials","true")
}

func Login(w http.ResponseWriter, r *http.Request){
	allowCORS(w)
	if sess,_ := store.Get(r,"JsSession");sess.Values["login"] == true{
		fmt.Fprint(w,`{"code:1,"message":"你已经登陆了"}`)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var admin configuration.Admin
	err := decoder.Decode(&admin)
	if err != nil{
		http.Error(w,fmt.Sprintf("解析json错误:%v",err),4)
		return
	}
	defer r.Body.Close()
	sess, err := store.New(r,"JsSession")
	if err != nil{
		log.Printf("session error : %v \n",err)
		return
	}
	sess.Options.MaxAge = 12 * 3600
	if configuration.Admins[admin.Username] == admin.Password{
		sess.Values["login"] = true
		sess.Save(r,w)
		fmt.Fprint(w,`{"code":1,"message":"登陆成功"}`)
	}else{
		sess.Values["login"] = false
		sess.Save(r,w)
		fmt.Fprint(w,`{"code":0,"message":"用户名或密码错"}`)
	}
}

func CheckLogin(w http.ResponseWriter, r *http.Request){
	allowCORS(w)
	//fmt.Println(r.Cookie("JsSession"))
	sess, err := store.Get(r,"JsSession")
	if err != nil{
		log.Printf("session error : %v \n",err)
		return
	}
	//println(session.Values["login"].(bool))
	if sess.Values["login"] == true{
		fmt.Fprint(w,`{"code":1,"message":"你已经登陆"}`)
		return
	}else{
		fmt.Fprint(w,`{"code":2,"message":"你还没有登陆"}`)
		return
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
func checkLogin(w http.ResponseWriter, r *http.Request)(login bool, err error){
	sess, err := store.Get(r,"JsSession")
	if err != nil{
		log.Printf("session error : %v \n",err)
		return false, err
	}
	if sess.Values["login"] == true{
		return true,err
	}else{
		fmt.Fprint(w,`{"code":2,"message":"你还没有登陆"}`)
		return false,err
	}
	return false,err
}


func FormSubmit(w http.ResponseWriter, r *http.Request){
	allowCORS(w)
	if login,_ := checkLogin(w,r);login == true{
		decoder := json.NewDecoder(r.Body)
		var data Data
		err := decoder.Decode(&data)
		if err != nil{
			fmt.Fprintf(w,`{"code":4,"message":"解析失败","error":%v}`,err)
		}
		fmt.Println(data)
		bd,err := json.Marshal(data)
		if err != nil{
			log.Printf("save failed : %v\n",err)
			return
		}
		exist, err := pathExists(Data_Origin_Path)
		if exist{
			newPath := fmt.Sprintf("%s%v_data.json",configuration.DataUrl,time.Now().Format("2006-01-02_15_04_05"))
			err = os.Rename(Data_Origin_Path,newPath)
			if err != nil{
				log.Printf("save failed : %v\n",err)
				return
			}
		}

		df, err := os.OpenFile(Data_Origin_Path, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil{
			log.Printf("save failed : %v\n",err)
			return
		}
		defer df.Close()
		_ , err = df.Write(bd)
		if err != nil{
			log.Printf("save failed : %v\n",err)
			return
		}
		fmt.Fprint(w,`{"code":5,"message":"上传成功,数据已存储"}`)
	}
}
func HandleFileUpload(w http.ResponseWriter, r *http.Request){
	allowCORS(w)
	if login,_ := checkLogin(w,r); login == true{
		r.ParseMultipartForm(32 << 20)
		file, handler , err := r.FormFile("icon")
		if err != nil{
			log.Println(err)
			return
		}
		defer file.Close()
		//fmt.Fprintf(w, "%v", handler.Header)
		newFileName := time.Now().Format("2006-01-02_15_04_05")+handler.Filename
		f, err := os.OpenFile(configuration.ImageUrlServer+newFileName, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil{
			fmt.Printf("打开文件失败 : %v\n",err)
			return
		}
		defer f.Close()
		io.Copy(f,file)
		fmt.Fprintf(w,`{"code":6,"message":"图片上传成功","image_url":"%s"}`,configuration.ImageUrlFront+newFileName)
	}
}
func main(){

	http.HandleFunc("/homepage/login",Login)
	http.HandleFunc("/homepage/login/check",CheckLogin)
	http.HandleFunc("/homepage/form",FormSubmit)
	http.HandleFunc("/homepage/image",HandleFileUpload)
	//http.HandleFunc("/test/cookie",testCookie)
	//http.HandleFunc("/test/session",testSession)
	log.Printf("启动服务 : localhost%s\n",configuration.ServerPort)
	log.Printf("图片存储目录 : %s\n",configuration.ImageUrlServer)
	err := http.ListenAndServe(configuration.ServerPort, nil)
	if err != nil{
		log.Printf("服务启动失败 : %v \n", err)
		return
	}
}

//func testSession(w http.ResponseWriter, r *http.Request){
//	allowCORS(w)
//	sess, err := store.New(r,"JsSession")
//	if err != nil{
//		log.Println(err)
//	}
//	sess.Values["login"] = true
//	sess.Save(r,w)
//	//sess := session.GobalSessions.SessionStart(w,r)
//	//sess.Set("username","sdffsfds")
//}
//func testCookie (w http.ResponseWriter, r *http.Request){
//	allowCORS(w)
//	//fmt.Println(r.Cookie("username"))
//	expiration := time.Now()
//	expiration = expiration.AddDate(1,0,0)
//	cookie := http.Cookie{Name: "username", Value: "astaxie", Expires: expiration}
//	http.SetCookie(w,&cookie)
//}