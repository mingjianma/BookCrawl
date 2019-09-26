package parser

import (
    "BookCrawl/model"
    "regexp"
    "strconv"
    "time"
    "fmt"
    "log"
)

//定义小说详情页的正则匹配器
var nameRe = regexp.MustCompile(`<h1 class="book-name"[^>]*>(.*)</h1>`)
var authorRe = regexp.MustCompile(`<p class="book-author hidden-md-and-up" [^>]*><a [^>]*>(.*)</a>`)
var statusRe = regexp.MustCompile(`<span class="status hidden-sm-and-down" [^>]*>·(.*)</span>`)
var wordageRe = regexp.MustCompile(`"countWord":([1-9]\d+),`)
var scoreRe = regexp.MustCompile(`"score":(.*?),`)
var scoreCountRe = regexp.MustCompile(`"scorerCount":(.*?),`)
var addListCountRe = regexp.MustCompile(`<span class="addListCount" [^>]*>(.*?)</span>次`)
var tagRe = regexp.MustCompile(`<a href="/bookstore/\?classId=\d+" target="_blank">(.*?)</a>`)
var lastUpdateRe = regexp.MustCompile(`"updateAt":"(.*?)",`)
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
    status := extractString(contents, statusRe)
    if status == "连载" {
        book.Status = "1"
    } else {
        book.Status = "0"
    }
    

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
    book.Tag = extractString(contents, tagRe)

    // 上次更新时间
    book.LastUpdate = extractString(contents, lastUpdateRe)
    match := lastUpdateTemRe.FindStringSubmatch(book.LastUpdate)
    if len(match) >= 3 {
        book.LastUpdate = match[1] + " " + match[2]
    }

    if book.Book_name != "" && book.Author != "" {
        save(book) 
    }
    
    result := model.ParseResult{
        Items: []interface{}{book},
    }
    
    return result
}

func save(book model.Book){
    t := time.Now()
    date := t.Format("200601")
    date_key := t.Format("2006-01-02")

    bookTable := model.MysqlDb{Dbname:MQ_DBNAME, Tblname:"book_info"}
    err := bookTable.NewMysqlDb()
    if err != nil {
        return 
    }
    defer bookTable.Close()
    //log.Println(book)
    bookWhere := map[string]string{"book_name":book.Book_name, "author":book.Author}
    book_info, _ := bookTable.SetWhere(bookWhere).Find()

    book_data :=  map[string]string{
        "book_name":book.Book_name,
        "author":book.Author,
        "tag":book.Tag,
        "status":book.Status,
        "score":strconv.FormatFloat(book.Score, 'e', -1, 64),
        "score_count":strconv.Itoa(book.Score_count),
        "score_detail":book.Score_detail,
        "add_list_count":strconv.Itoa(book.AddListCount),
        "last_update_time":book.LastUpdate,
    }
    var book_id string
    if book_info == nil {
        insert_id, err := bookTable.Insert(book_data)
        book_id = strconv.Itoa(insert_id)
        if err != nil {
            log.Printf("【%s】 book_info insert err：%s", book.Book_name, err)
        }
    } else {
        _, err := bookTable.Update(book_data)
        book_id = book_info["book_id"]
        if err != nil {
            log.Printf("【%s】 book_info update err：%s", book.Book_name, err)
        }
    }

    scoreTable := model.MysqlDb{Dbname:MQ_DBNAME, Tblname:("book_score_daily_" + date)}
    err = scoreTable.NewMysqlDb()
    if err != nil {
        return 
    }
    defer scoreTable.Close()
    scoreWhere := map[string]string{"book_id":book_id, "date_key":date_key}
    scoreInfo, _ := scoreTable.SetWhere(scoreWhere).Find()
    scorce_data := map[string]string{
        "book_id":book_id,
        "score":strconv.FormatFloat(book.Score, 'e', -1, 64),
        "score_count":strconv.Itoa(book.Score_count),
        "score_detail":book.Score_detail,
        "date_key":date_key,
    }
    if scoreInfo == nil {
        _, err := scoreTable.Insert(scorce_data)
        if err != nil {
            log.Printf("【%s】 score_info insert err：%s", book.Book_name, err)
        }
    } else {
        _, err := scoreTable.Update(scorce_data)
        if err != nil {
            log.Printf("【%s】 score_info update err：%s", book.Book_name, err)
        }
    }
}

//进行正则匹配，并获取第一个匹配到的match
func extractString(body []byte, re *regexp.Regexp) string {
    match := re.FindSubmatch(body)
    if len(match) >= 2 {
        return string(match[1])
    }

    return ""
}
