package posts

import (
	"context"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type FakePostBuilder struct {
	t    testing.TB
	post *Post
}

func NewFakePost(t testing.TB) *FakePostBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()
	createdAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	return &FakePostBuilder{
		t: t,
		post: &Post{
			id:        gofakeit.Uint64(),
			status:    Uploaded,
			title:     gofakeit.LoremIpsumSentence(4),
			fileID:    uuidProvider.New(),
			createdAt: createdAt,
			createdBy: uuidProvider.New(),
		},
	}
}

func (f *FakePostBuilder) WithMedia(media *medias.FileMeta) *FakePostBuilder {
	f.post.fileID = media.ID()

	return f
}

func (f *FakePostBuilder) WithStatus(status Status) *FakePostBuilder {
	f.post.status = status

	return f
}

func (f *FakePostBuilder) CreatedBy(user *users.User) *FakePostBuilder {
	f.post.createdBy = user.ID()

	return f
}

func (f *FakePostBuilder) Build() *Post {
	return f.post
}

func (f *FakePostBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *Post {
	f.t.Helper()

	storage := newSqlStorage(db)

	post := f.Build()

	err := storage.Save(ctx, post)
	require.NoError(f.t, err)

	return post
}
