package yb

import (
	"fmt"
	"net/http/cookiejar"
	"net/url"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"github.com/PuerkitoBio/goquery"
	"errors"
)

const (
	login_page_url 	string = "https://www.yiban.cn/login"
	login_post_url 	string = "https://www.yiban.cn/login/doLoginAjax"
)
type LoginBackMsg struct {
	Code int32 `json:"code"`
	Message string `json:"message"`
	Data DataTemplate `json:"data"`
}
type DataTemplate struct {
	Url string `json:"url"`
}
type CheckLoginMsg struct {
	Code int32 `json:"code"`
	Message string `json:"message"`
	Data DataCheckTemplate `json:"data"`
}
type DataCheckTemplate struct{
	IsLogin bool `json:"isLogin"`
}
func getLoginPage(login_page_url string) *http.Response{
c := &http.Client{}
req, _ := http.NewRequest("GET",login_page_url,nil)
res, _ := c.Do(req)
fmt.Printf("%v\n",res.Cookies())
return res
}
func getDataKeysAndKeysTime(res *http.Response) (data_keys string, data_keys_time string, err error){
	doc, _ := goquery.NewDocumentFromResponse(res)
	data_keys, exists := doc.Find("#login-pr").Attr("data-keys")
	if !exists {
		err = errors.New("can not find the attribute named data-keys")
		return data_keys,data_keys_time,err
	}
	data_keys_time, exists = doc.Find("#login-pr").Attr("data-keys-time")
	if !exists {
		err = errors.New("can not find the attribute named data_keys_time")
		return data_keys,data_keys_time,err
	}
	return data_keys,data_keys_time,err
}
func CheckLogin(c *http.Client) bool{
	var checkLoginMsg CheckLoginMsg
	res, _ := c.Get("http://www.yiban.cn/ajax/my/getLogin")
	data , _ := ioutil.ReadAll(res.Body)
	fmt.Printf("%v\n", string(data))
	fmt.Printf("%v\n", c.Jar)
	json.Unmarshal([]byte(data),&checkLoginMsg)
	fmt.Printf("%v\n",checkLoginMsg)
	return checkLoginMsg.Data.IsLogin
}
func Login (account string,password string) *http.Client{
	c := &http.Client{}
	res := getLoginPage(login_page_url)
	data_keys , data_keys_time , err:= getDataKeysAndKeysTime(res)
	if err != nil{
		fmt.Printf("%v\n",err)
		return c
	}
	pwdRsa := RsaEncrypt([]byte(password),[]byte(data_keys))

	Jar, _ := cookiejar.New(nil)
	loginPostUrl,_ := url.Parse(login_post_url)
	temp_cookies := res.Cookies()
	Jar.SetCookies(loginPostUrl,temp_cookies)
	c.Jar = Jar
	postValues := url.Values{}
	postValues.Add("account",account)
	postValues.Add("password",pwdRsa)
	postValues.Add("captcha","")
	postValues.Add("keysTime",data_keys_time)
	response, _ := c.PostForm(login_post_url,postValues)
	decoder := json.NewDecoder(response.Body)
	var msg LoginBackMsg
	err = decoder.Decode(&msg)
	if err != nil{
		fmt.Printf("error :%v",err)
	}
	fmt.Printf("%v\n",msg)
	data, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Println(string(data))
	fmt.Printf("%v\n",response.Header)
	return c
}