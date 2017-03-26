package feedback

import (
	"time"

	"github.com/WeCanHearYou/wechy/app"
)

//Idea represents an idea on a tenant board
type Idea struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdOn"`
	User        app.User  `json:"user"`
}

//Comment represents an user comment on an idea
type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedOn time.Time `json:"createdOn"`
	User      app.User  `json:"user"`
}
