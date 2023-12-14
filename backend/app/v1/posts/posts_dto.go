package posts

// Post representa la estructura de un post
type Post struct {
	UserID int    `json:"userId" validate:"required"`
	ID     int    `json:"id" `
	Title  string `json:"title" validate:"required,min=4,max=15"`
	Body   string `json:"body" validate:"required,min=4,max=200"`
}
