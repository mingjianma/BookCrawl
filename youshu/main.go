package main

import (
    "BookCrawl/engine"
    "BookCrawl/model"
    "BookCrawl/scheduler"
    "BookCrawl/youshu/parser"
    "time"
    "log"
)

func main() {
    st := time.Now()
    defer func() {
        elapsed := time.Since(st)
        log.Println("App runtime: ", elapsed)
    }()
    engine.ConcurrentEngine{
        Scheduler:   &scheduler.SimpleScheduler{},
        WorkerCount: 100,
    }.Run(model.Request{
        // 种子 Url
        Url:        parser.SEED_URL,
        ParserFunc: parser.ParseBookList,
    })
}
