package main

import (
    "mike/crawler/engine"
    "mike/crawler/model"
    "mike/crawler/scheduler"
    "mike/crawler/youshu/parser"
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
