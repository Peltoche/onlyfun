package moderations

import (
	"time"

	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	v "github.com/go-ozzo/ozzo-validation"
)

type Moderation struct {
	id        uint
	postID    uint
	reason    string
	createdAt time.Time
	createdBy uuid.UUID
}

func (m *Moderation) ID() uint             { return m.id }
func (m *Moderation) PostID() uint         { return m.postID }
func (m *Moderation) Reason() string       { return m.reason }
func (m *Moderation) CreatedAt() time.Time { return m.createdAt }
func (m *Moderation) CreatedBy() uuid.UUID { return m.createdBy }

type PostModerationCmd struct {
	User   *users.User
	Post   *posts.Post
	Reason string
}

func (t PostModerationCmd) Validate() error {
	return v.ValidateStruct(&t,
		v.Field(&t.User, v.Required),
		v.Field(&t.Post, v.Required),
		v.Field(&t.Reason, v.Required, v.Length(5, 300)),
	)
}
