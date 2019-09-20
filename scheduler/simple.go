package scheduler

import (
    "BookCrawl/model"
)

type SimpleScheduler struct {
    workerChan chan model.Request
}

// 使用指针接收者，改变 SimpleScheduler 内部的 workerChan
func (s *SimpleScheduler) ConfigureMasterWorkerChan(in chan model.Request) {
    s.workerChan = in
}

func (s *SimpleScheduler) Submit(request model.Request) {
    // 每个 Request 一个 Goroutine
    go func() { s.workerChan <- request }()
}