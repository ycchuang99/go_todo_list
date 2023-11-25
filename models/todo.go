package models

type Todo struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	DoneAt      *string `json:"done_at"`
}
