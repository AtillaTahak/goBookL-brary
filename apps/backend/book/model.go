package book

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title" validate:"required"`
	Author string `json:"author" validate:"required"`
	Year   int    `json:"year" validate:"required"`
}
