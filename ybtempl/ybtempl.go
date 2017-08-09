package ybtempl

type NoticesStruct struct{
	Carousels []Carousel `json:"carousels"`
	Notice1 Notice `json:"notice_1"`
	Notice2 Notice `json:"notice_2"`
}

type Notice struct {
	Title string `json:"title"`
	Content [6]NoticeContent `json:"content"`
}
type NoticeContent struct {
	Text string `json:"text"`
	Href string `json:"href"`
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
type YBData struct {
	Notices NoticesStruct `json:"notices"`
	Apps [5][]App `json:"apps"`
	SchoolIntro SchoolIntroTempl `json:"school_intro"`
}
type SchoolIntroTempl struct {
	Group GroupTempl `json:"group"`
	Members string `json:"members"`
}
type Organization struct {
	Href string `json:"href"`
	Src string `json:"src"`
}
type GroupTempl struct{
	Code int `json:"code"`
	Message string `json:"message"`
	Data []GroupData `json:"data"`
}
type GroupData struct {
	Id string `json:"id"`
	UserId string `json:"user_id"`
	Name string `json:"name"`
	Brief string `json:"brief"`
	Kind string `json:"kind"`
	Img string `json:"img"`
	QrCode string `json:"qrCode"`
	Auth string `json:"auth"`
	Type string `json:"type"`
	OriginId string `json:"originId"`
	OldClassId string `json:"oldClassId"`
	UpdateTime string `json:"updateTime"`
	CreateTime string `json:"createTime"`
	Sort string `json:"sort"`
	Label string `json:"label"`
	Top string `json:"top"`
	IsMember int `json:"isMember"`
	Avatar string `json:"avatar"`
	Url string `json:"url"`
}
//assign new value's new value to old value
func AssignNoticesAndApps(old *YBData,new *YBData){
	old.Notices = new.Notices
	old.Apps = new.Apps
}