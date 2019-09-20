package parser

import (
    "io/ioutil"
    "testing"
    "net/http"
//    "fmt"
)

func TestParseBookList(t *testing.T) {
    //expectRequestsLen := 470
    //expectCitiesLen := 470
    // 表格驱动测试
    //expectRequestUrls := []string{
    //    "http://www.yousuu.com/book/141152",
    //    "http://www.yousuu.com/book/165",
    //    "http://www.yousuu.com/book/13899",
    //}
    /*expectRequestCities := []string{
        "city 阿坝",
        "city 阿克苏",
        "city 阿拉善盟",
    }*/

    //body, err := ioutil.ReadFile("citylist_test_data.html")
    resp, err := http.Get("http://www.yousuu.com/bookstore/?channel&classId&tag&countWord&status&update&sort&page=2")

    if err != nil {
        panic(err)
    }

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }
    result := ParseBookList(body, "http://www.yousuu.com/bookstore/?channel&classId&tag&countWord&status&update&sort&page=2")
    for _, data := range result.Requests {
        resp, _ := http.Get(data.Url)
        body, _ := ioutil.ReadAll(resp.Body)
        data.ParserFunc(body, data.Url)
    }
}