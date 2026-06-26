package posts

import "time"

type Post struct {
	Id        uint
	UserId    uint
	Title     string
	Content   string
	ImageURL  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostRequest struct {
	Title   string
	Content string
}

type PostResponse struct {
	Id        uint
	UserId    uint
	Title     string
	Content   string
	ImageURL  string
	CreatedAt time.Time
}
