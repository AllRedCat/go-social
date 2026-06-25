package posts

import "time"

type Post struct {
	id          uint
	Title       string
	Description string
	ImageURL    string
	CreatedAt   time.Time
}
