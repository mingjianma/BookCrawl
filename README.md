# BookCrawl

本项目通过使用Go，完成并发爬虫，爬取优书网的书库，获取书库列表的每本小说的详细信息

目录结构
```js
├─engine                       //引擎目录
│      ConcurrentEngine.go	   //单进程并发引擎
│      simpleEngine.go         //单进程顺序引擎
│  
│
├─fetcher	                   //分析器目录
│	   fetcher.go	           //分析器
│
├─models	                   //基础数据结构定义目录
│      book.go                 //小说详情结构
│	   types.go.go             //请求和请求结果结构
│
├─scheduler                    //协程调用目录
│      scheduler.go            //协程调用接口
│      simple.go               //协程处理
│
└─youshu                       //处理优书网数据目录
	│
	│	 main.go               //优书网处理程序入口
	│
    └─parser                   //优书网数据处理目录
    	book.go                //处理小说详情页规则
    	booklist.go            //处理小说列表规则
    	booklist_test.go       //处理小说列表规则测试
    	config.go              //具体配置
```

执行main.go即可获取到小说数据
数据处理待续