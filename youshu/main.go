package main

import (
    "BookCrawl/engine"
    "BookCrawl/model"
    "BookCrawl/scheduler"
    "BookCrawl/youshu/parser"
)

func main() {
    engine.ConcurrentEngine{
        Scheduler:   &scheduler.SimpleScheduler{},
        WorkerCount: 1000,
    }.Run(model.Request{
        // 种子 Url
        Url:        parser.SEED_URL,
        ParserFunc: parser.ParseBookList,
    })
}
