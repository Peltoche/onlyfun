package posts

import (
	"io"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	v "github.com/go-ozzo/ozzo-validation"
)

const (
	Uploaded  Status = "uploaded"
	Listed    Status = "listed"
	Moderated Status = "moderated"
)

type Status string

type Post struct {
	id        uint64
	status    Status
	title     string
	fileID    uuid.UUID
	createdAt time.Time
	createdBy uuid.UUID
}

func (p Post) ID() uint64           { return p.id }
func (p Post) Status() Status       { return p.status }
func (p Post) Title() string        { return p.title }
func (p Post) FileID() uuid.UUID    { return p.fileID }
func (p Post) CreatedAt() time.Time { return p.createdAt }
func (p Post) CreatedBy() uuid.UUID { return p.createdBy }

type CreateCmd struct {
	Title     string
	Media     io.Reader
	CreatedBy *users.User
}

func (t CreateCmd) Validate() error {
	return v.ValidateStruct(&t,
		v.Field(&t.Title, v.Required, v.Length(3, 280)),
		v.Field(&t.CreatedBy, v.Required),
		v.Field(&t.Media, v.Required),
	)
}
