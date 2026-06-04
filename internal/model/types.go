package model

type Todo struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Response struct {
	Message string `json:"message"`
}
