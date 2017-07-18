package main

import (
	"./configuration"
	"log"
	"net/http"
	"encoding/json"
	"github.com/gorilla/sessions"
	"fmt"
)

type NoticesStruct struct{
	Carousels [6]Carousel `json:"carousel"`
	Notice1 Notice `json:"notice_1"`
	Notice2 Notice `json:"notice_2"`
}

type Notice struct {
	Title string `json:"title"`
	Content [6]NoticeContent `json:"content"`
}
type NoticeContent struct {
	Content string `json:"content"`
	Url 	string `json:"url"`
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
	}else{
		fmt.Fprint(w,`{"code":2,"message":"你还没有登陆"}`)
	}
}

func FormSubmit(w http.ResponseWriter, r *http.Request){
	allowCORS(w)
	CheckLogin(w,r)
	decoder := json.NewDecoder(r.Body)
	var data Data
	err := decoder.Decode(&data)
	if err != nil{
		fmt.Fprintf(w,`{"code":4,"message":"解析失败","error":%v}`,err)
	}
	fmt.Println(data)
}

func main(){

	http.HandleFunc("/homepage/login",Login)
	http.HandleFunc("/homepage/login/check",CheckLogin)
	http.HandleFunc("/homepage/form",FormSubmit)
	//http.HandleFunc("/test/cookie",testCookie)
	//http.HandleFunc("/test/session",testSession)
	log.Printf("启动服务 : localhost%s\n",configuration.ServerPort)
	log.Printf("图片存储目录 : %s\n",configuration.ImageUrl)
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