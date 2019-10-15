package parser

import (
    "BookCrawl/model"
    "regexp"
    "strconv"
    "time"
    "fmt"
    "log"
    "github.com/astaxie/beego/orm"
)

//定义小说详情页的正则匹配器
var nameRe = regexp.MustCompile(`<h1 class="book-name"[^>]*>(.*)</h1>`)
var authorRe = regexp.MustCompile(`<p class="book-author hidden-md-and-up" [^>]*><a [^>]*>(.*)</a>`)
var statusRe = regexp.MustCompile(`<span class="status hidden-sm-and-down" [^>]*>·(.*?)</span>`)
var wordageRe = regexp.MustCompile(`"countWord":(.*?),`)
var scoreRe = regexp.MustCompile(`"score":(.*?),`)
var scoreCountRe = regexp.MustCompile(`"scorerCount":(.*?),`)
var addListCountRe = regexp.MustCompile(`<span class="addListCount" [^>]*>(.*?)</span>次`)
var tagRe = regexp.MustCompile(`<a href="/bookstore/\?classId=\d+" target="_blank">(.*?)</a>`)
var lastUpdateRe = regexp.MustCompile(`"updateAt":"(.*?)",`)
var scoreDetailRe = regexp.MustCompile(`"scoreDetail":(\[[^>]*?\]),`)

var lastUpdateTemRe = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})T(\d{2}:\d{2}:\d{2}).*`)

// 解析小说详情
func ParseBook(contents []byte, _ string) model.ParseResult {
    defer func (){
        //恢复在关闭in渠道后返回的worker像out <- result
        if info := recover(); info != nil {
            //log.Printf("panic recover")
        }
    }()
    o := orm.NewOrm()
    book := model.Book{}

    // 书名
    book.Book_name = extractString(contents, nameRe)

    // 作者
    book.Author = extractString(contents, authorRe)

    insert_flag := false
    if book.Book_name != "" && book.Author != "" {
        err := o.Read(&book, "book_name", "author")
        if err == orm.ErrNoRows {
            insert_flag = true 
        }

        // 状态
        status := extractString(contents, statusRe)
        if status == "连载" {
            book.Status = "1"
        } else {
            book.Status = "0"
        }

        // 小说标签
        book.Tag = extractString(contents, tagRe)

        // 字数
        book.Wordage , _ = strconv.Atoi(extractString(contents, wordageRe))

        // 评分
        book.Score, _ = strconv.ParseFloat(extractString(contents, scoreRe), 64)
        book.Score, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", book.Score), 64)

        // 评分人数
        book.Score_count, _ = strconv.Atoi(extractString(contents, scoreCountRe))

        // 评分人数
        book.Score_detail = extractString(contents, scoreDetailRe)

        // 加入书单次数
        book.Add_list_count, _ = strconv.Atoi(extractString(contents, addListCountRe))

        // 上次更新时间
        book.Last_update = extractString(contents, lastUpdateRe)
        match := lastUpdateTemRe.FindStringSubmatch(book.Last_update)
        if len(match) >= 3 {
            book.Last_update = match[1] + " " + match[2]
        }
        log.Println(book)
        //更新小说详情
        if insert_flag == true {
            _, err = o.Insert(&book)
        } else {
            _, err = o.Update(&book)
        }

        if err == nil {
            //更新小说每日评分
            log.Printf("==【%s】 book_info update Success", book.Book_name)
            book_score := model.BookScore{}
            t := time.Now()
            date_key := t.Format("2006-01-02")

            book_score.Book_id = book.Book_id
            book_score.Date_key = date_key
            err = o.Read(&book_score, "book_id", "date_key")
            insert_flag = false
            if err == orm.ErrNoRows {
                insert_flag = true
            }
            book_score.Score = book.Score
            book_score.Score_count = book.Score_count
            book_score.Score_detail = book.Score_detail
            if insert_flag == true {
                _, err = o.Insert(&book_score)
            } else {
                _, err = o.Update(&book_score)
            }

            if err != nil {
                log.Printf("==【%s】 book_score update err：%s", book.Book_name, err)
            } else {
                log.Printf("==【%s】 book_score update Success", book.Book_name)
            }

        } else {
           log.Printf("==【%s】 book_info update err：%s", book.Book_name, err)
        }
    } 

    result := model.ParseResult{
        Items: []interface{}{book},
    }
    
    return result
}

//进行正则匹配，并获取第一个匹配到的match
func extractString(body []byte, re *regexp.Regexp) string {
    match := re.FindSubmatch(body)
    if len(match) >= 2 {
        return string(match[1])
    }

    return ""
}
