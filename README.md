# BookCrawl

# 本项目通过使用Go，完成并发爬虫，爬取优书网的书库，获取书库列表的每本小说的详细信息

# 目录结构
```js
├─config                       //全局配置目录
│      config.go               //全局配置   
│ 
├─data
│	   full.sql                //数据库建表文件
│
├─engine                       //引擎目录
│      ConcurrentEngine.go	   //单进程并发引擎
│      simpleEngine.go         //单进程顺序引擎
│  
│
├─fetcher	                   //分析器目录
│      fetcher.go	           //分析器
│
├─models	                   //基础数据结构定义目录
│      book.go                 //小说详情结构
│      types.go                //请求和请求结果结构
│      mysql.go                //mysql封装类
│
├─scheduler                    //协程调用目录
│      scheduler.go            //协程调用接口
│      simple.go               //协程处理
│
└─youshu                       //处理优书网数据目录
    │
    │  main.go                 //优书网处理程序入口
    │
    └─parser                   //优书网数据处理目录
    	book.go                //处理小说详情页规则
    	booklist.go            //处理小说列表规则
    	config.go              //爬虫网站具体配置
```
# 基础配置：
- 通过`data/full.sql`创建存储数据的相关数据表（book_score_daily_为每日积分统计分表，请按月份建表）
- 修改config/config.go的数据库配置
- mysql.go封装类是基于go-sql-driver/mysql的mysql的驱动，需要先下载相应的驱动
>   * 驱动下载命令：go get -u github.com/go-sql-driver/mysql

# 数据处理：
- 执行main.go即可获取到当天的数据