package moderations

import (
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/posts"
	"github.com/Peltoche/onlyfun/internal/services/users"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
)

type FakeModerationBuilder struct {
	t          testing.TB
	moderation *Moderation
}

func NewFakeModeration(t testing.TB) *FakeModerationBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()
	createdAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	return &FakeModerationBuilder{
		t: t,
		moderation: &Moderation{
			id:        gofakeit.Uint(),
			postID:    gofakeit.Uint(),
			reason:    gofakeit.LoremIpsumSentence(gofakeit.Number(1, 20)),
			createdAt: createdAt,
			createdBy: uuidProvider.New(),
		},
	}
}

func (f *FakeModerationBuilder) CreatedBy(user *users.User) *FakeModerationBuilder {
	f.moderation.createdBy = user.ID()

	return f
}

func (f *FakeModerationBuilder) WithPost(post *posts.Post) *FakeModerationBuilder {
	f.moderation.postID = post.ID()

	return f
}

func (f *FakeModerationBuilder) Build() *Moderation {
	return f.moderation
}

// func (f *FakeModerationBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *Moderation {
// 	f.t.Helper()
//
// 	storage := newSqlStorage(db)
//
// 	user := f.Build()
//
// 	err := storage.Save(ctx, user)
// 	require.NoError(f.t, err)
//
// 	return user
// }
