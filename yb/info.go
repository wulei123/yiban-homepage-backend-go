package yb

import (
	"net/http"
	"encoding/json"
	"../ybtempl"
	"github.com/PuerkitoBio/goquery"
	"fmt"
)

const (
	square_index_url string = "https://www.yiban.cn/square/index"
	my_school_url           = "http://www.yiban.cn/school/getMyGroupAjax?id=5370538"
)

func getMembers(memberTotal *goquery.Selection) (members string) {
	memberTotal.Children().Each(func(i int, span *goquery.Selection) {
		if i == 1 {
			members = span.Children().First().Next().Text()
		}
	})
	return members
}
func getSchoolIntro(c *http.Client) (schoolIntro ybtempl.SchoolIntroTempl){
	res, _ := c.Get(square_index_url)
	doc, _ := goquery.NewDocumentFromResponse(res)
	temp := doc.Find(".yiban-my-school")
	temp = temp.ChildrenFiltered(".school-intro").Children()
	temp = temp.ChildrenFiltered(".member-total")
	schoolIntro.Members= getMembers(temp)
	schoolIntro.Group = FetchMyGroup(c)
	fmt.Printf("%v", schoolIntro)
	return schoolIntro
}
func FetchMyGroup(c *http.Client) (group ybtempl.GroupTempl) {
	res, _ := c.Get(my_school_url)
	decoder := json.NewDecoder(res.Body)
	decoder.Decode(&group)
	return group
}
