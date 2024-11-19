package medias

import (
	"time"

	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

type MediaType string

const (
	Post   MediaType = "post"
	Avatar MediaType = "avatar"
)

type FileMeta struct {
	uploadedAt time.Time
	id         uuid.UUID
	mediaType  MediaType
	mimetype   string
	checksum   string
	size       uint64
}

func (f FileMeta) ID() uuid.UUID         { return f.id }
func (f FileMeta) Mimetype() string      { return f.mimetype }
func (f FileMeta) Checksum() string      { return f.checksum }
func (f FileMeta) Size() uint64          { return f.size }
func (f FileMeta) Type() MediaType       { return f.mediaType }
func (f FileMeta) UploadedAt() time.Time { return f.uploadedAt }
