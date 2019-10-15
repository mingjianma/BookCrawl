package model
import "time"

type BookScore struct {
    Log_id int `orm:"PK"`

    Book_id int

    //评分
    Score float64

    //评分人数
    Score_count int

    //详细评分
    Score_detail string

    //上次更新时间（粗略）
    Date_key string
}

func (bs *BookScore) TableName() string {
    t := time.Now()
    date := t.Format("200601")
    return "book_score_daily_"+date
}