package engine

import (
    "BookCrawl/fetcher"
    "BookCrawl/model"
    "BookCrawl/scheduler"
    "log"
    "time"
    "sync"
)

// 并发引擎
type ConcurrentEngine struct {
    // 调度器
    Scheduler scheduler.Scheduler
    // 开启的 worker 数量
    WorkerCount int
}

var wg sync.WaitGroup

func (e ConcurrentEngine) Run(seeds ...model.Request) {
    wg.Add(e.WorkerCount)
    defer log.Printf("task end")

    in := make(chan model.Request)
    out := make(chan model.ParseResult)

    // 初始化调度器的 chann
    e.Scheduler.ConfigureMasterWorkerChan(in)

    // 创建 WorkerCount 个 worker
    for i := 0; i < e.WorkerCount; i++ {
        createWorker(in, out);
    }
    // 将 seeds 中的 Request 添加到调度器 chann
    for _, r := range seeds {
        e.Scheduler.Submit(r)
    }

    for {
        result, isClose := <-out // 阻塞获取
        //等待关闭的渠道数据处理完再退出
        if !isClose {
            //log.Printf("out channel closed!")
            break
        }
        
        for _, item := range result.Items {
            log.Printf("getItems, items: %v", item)
        }

        for _, r := range result.Requests {
            // 如果 submit 内部直接是 s.workerChan <- request，则阻塞等待发送，该方法阻塞在这里
            // 如果 submit 内部直接是 go func() { s.workerChan <- request }()，则为每个Request分配了一个Goroutine，这里不会阻塞在这里
            e.Scheduler.Submit(r)
        }

    }    
    wg.Wait()
}

func createWorker(in chan model.Request, out chan model.ParseResult) {
    go func() {
        defer func (){
            //恢复在关闭in渠道后返回的worker像out <- result
            if info := recover(); info != nil {
                //log.Printf("panic recover")
            }
            wg.Done()
        }()

        for {
            r, isClose := <-in // 阻塞等待获取
            if !isClose {
                //判断out渠道的数据是否都获取完，获取完关闭out渠道，并退出循环
                if len(out) == 0 {
                    _, isOutClose := <-out
                    if !isOutClose {
                        //log.Printf("out channel closed!")
                    } else {
                        close(out)
                        log.Printf("close out channel!")
                    }
                    break
                } else {
                    //log.Printf("in channel closed! wait out channel")
                    time.Sleep(time.Duration(2)*time.Second)
                    continue
                }
                
            }

            result, err := worker(r)
            if err != nil {
                continue
            }

            //判断任务是否爬完，爬完关闭in渠道
            if result.EndFlag == true {
                close(in)
                log.Printf("close in channel!")
            } else {
                out <- result // 阻塞发送
            }            
        }
    }()
}

func worker(r model.Request) (model.ParseResult, error) {
    body, err := fetcher.Fetch(r.Url)
    if err != nil {
        log.Printf("fetch error, url: %s, err: %v", r.Url, err)
        return model.ParseResult{}, nil
    }
    return r.ParserFunc(body, r.Url), nil
}