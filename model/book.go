package model

type Book struct {
    Book_id int `orm:"PK"`
    //书名
    Book_name string

    //作者
    Author string

    //类型
    Tag string

    //字数（万字）
    Wordage int

    //状态 
    Status string

    //评分
    Score float64

    //评分人数
    Score_count int

    //详细评分
    Score_detail string

    //收录书单次数
    Add_list_count int

    //上次更新时间（粗略）
    Last_update string
}


func (b *Book) TableName() string {
    return "book_info"
}