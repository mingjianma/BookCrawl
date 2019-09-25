package model

type Book struct {
	//书名
	Book_name string

	//作者
	Author string

	//类型
	Tag string

	//字数（万字）
	Wordage	int

	//状态 
	Status string

	//评分
	Score float64

	//评分人数
	Score_count int

	//详细评分
	Score_detail string

	//收录书单次数
	AddListCount int

	//上次更新时间（粗略）
	LastUpdate string
}