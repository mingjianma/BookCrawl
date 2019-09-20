package parser

import (
    "mike/crawler/model"
    "regexp"
    "strconv"
    "fmt"
)

var nameRe = regexp.MustCompile(`<h1 class="book-name"[^>]*>(.*)</h1>`)
var authorRe = regexp.MustCompile(`<p class="book-author hidden-md-and-up" [^>]*><a [^>]*>(.*)</a>`)
var statusRe = regexp.MustCompile(`<span class="status hidden-sm-and-down" data-v-38b3f138>·(.*)</span>`)
var wordageRe = regexp.MustCompile(`"countWord":([1-9]\d+),`)
var scoreRe = regexp.MustCompile(`"score":(.*?),`)
var scoreCountRe = regexp.MustCompile(`"scorerCount":(.*?),`)
var addListCountRe = regexp.MustCompile(`<span class="addListCount" [^>]*>(.*?)</span>次`)
var tagRe = regexp.MustCompile(`<a href="/bookstore/\?classId=\d+" target="_blank">(.*?)</a>`)
var lastUpdateRe = regexp.MustCompile(`"updateAt":"(.*?)",`)
//var lastUpdateRe = regexp.MustCompile(`更新时间：(.*)\s*<span`)
var scoreDetailRe = regexp.MustCompile(`"scoreDetail":(\[[^>]*?\]),`)

var lastUpdateTemRe = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}).*`)

// 解析小说详情
func ParseBook(contents []byte, _ string) model.ParseResult {
    book := model.Book{}

    // 书名
    book.Book_name = extractString(contents, nameRe)

    // 作者
    book.Author = extractString(contents, authorRe)

    // 字数
    book.Wordage , _ = strconv.Atoi(extractString(contents, wordageRe))

    // 状态
    book.Status = extractString(contents, statusRe)

    // 评分
    book.Score, _ = strconv.ParseFloat(extractString(contents, scoreRe), 64)
    book.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", book.Score), 64)

    // 评分人数
    book.Score_count, _ = strconv.Atoi(extractString(contents, scoreCountRe))

    // 评分人数
    book.Score_detail = extractString(contents, scoreDetailRe)

    // 加入书单次数
    book.AddListCount, _ = strconv.Atoi(extractString(contents, addListCountRe))

    // 小说标签
    book.Type = extractString(contents, tagRe)

    // 上次更新时间
    book.LastUpdate = extractString(contents, lastUpdateRe)
    match := lastUpdateTemRe.FindStringSubmatch(book.LastUpdate)
    if len(match) >= 3 {
        book.LastUpdate = match[1] + " " + match[2]
    }
    
    fmt.Println(book)
    result := model.ParseResult{
        Items: []interface{}{book},
    }
    
    return result
}


func extractString(body []byte, re *regexp.Regexp) string {
    match := re.FindSubmatch(body) // 只找到第一个match的
    if len(match) >= 2 {
        return string(match[1])
    }
    return ""
}