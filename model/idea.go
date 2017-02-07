package model

//Idea represents an idea on a tenant board
type Idea struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
