package parser

import (
    "io/ioutil"
    "testing"
    "net/http"
//    "fmt"
)

func TestParseBookList(t *testing.T) {

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