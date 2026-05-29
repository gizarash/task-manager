package model

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Store struct {
	CurrentId int    `json:"current_id"`
	Todos     []Todo `json:"todos"`
}
