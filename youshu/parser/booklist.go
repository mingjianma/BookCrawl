package parser

import (
    "BookCrawl/model"
    "regexp"
 //    "fmt"
    "strconv"
)

const bookListRe = `<a[^>]*href="(/book/\d+)"[^>]*class="bookname" [^>]*>[^<]*</a>`
const urlRe = `page=(\d+)`
const totalNumRe = `"total":([1-9]\d+),`


// cityList 的 ParserFunc func([]byte) ParseResult
// 解析种子页面 - 获取小说列表
func ParseBookList(contents []byte, url string) model.ParseResult {
    result := model.ParseResult{}

    // 计算小说库页数
    totalNumg := regexp.MustCompile(totalNumRe)
    totalNum, _ := strconv.Atoi(extractString(contents, totalNumg))
    pages := totalNum / PREPAGE

    // 增加下一页到队列
    urlg := regexp.MustCompile(urlRe)
    urlSubmatch := urlg.FindStringSubmatch(url)
    page, _ := strconv.Atoi(urlSubmatch[1])
    //判断是否已经处理完所有小说列表
    if page > pages {
        result.EndFlag = true
        return result
    }
    page += 1
    next_url := PAGELISTURL + strconv.Itoa(page)
    result.Items = append(result.Items, next_url)
    result.Requests = append(result.Requests, model.Request{
       Url:        next_url,
       ParserFunc: ParseBookList,
    })

    // 通过正则表达式生成正则对象：()用于提取
    rg := regexp.MustCompile(bookListRe)
    allSubmatch := rg.FindAllSubmatch(contents, -1)

    // FindAllSubmatch获取的是一个byte数组，m[0]为匹配到的内容, m[1]为第一个()匹配到的内容
    // 遍历列表的每一个小说详情地址，并且将 Url 和小说解析器封装为一个 Request
    // 最后将该 Request 添加到 ParseResult 中
    for _, m := range allSubmatch {
        url := DOMAIN + string(m[1])
        result.Items = append(result.Items, url)
        result.Requests = append(result.Requests, model.Request{
           Url:        url,
           ParserFunc: ParseBook,
        })
    }

    // 返回 ParseResult
    return result
}