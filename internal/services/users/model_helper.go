package users

import (
	"context"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/roles"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type FakeUserBuilder struct {
	t    testing.TB
	user *User
}

func NewFakeUser(t testing.TB) *FakeUserBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()
	createdAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	return &FakeUserBuilder{
		t: t,
		user: &User{
			id:                uuidProvider.New(),
			createdAt:         createdAt,
			passwordChangedAt: createdAt,
			role:              "user",
			username:          gofakeit.Username(),
			avatar:            uuidProvider.New(),
			password:          secret.NewText(gofakeit.Password(true, true, true, false, false, 8)),
			status:            Active,
			createdBy:         uuidProvider.New(),
		},
	}
}

func (f *FakeUserBuilder) WithPassword(password string) *FakeUserBuilder {
	f.user.password = secret.NewText(password)

	return f
}

func (f *FakeUserBuilder) WithAvatar(media *medias.FileMeta) *FakeUserBuilder {
	f.user.avatar = media.ID()

	return f
}

func (f *FakeUserBuilder) CreatedBy(user *User) *FakeUserBuilder {
	f.user.createdBy = user.ID()

	return f
}

func (f *FakeUserBuilder) WithRole(role *roles.Role) *FakeUserBuilder {
	f.user.role = role.Name()

	return f
}

func (f *FakeUserBuilder) WithStatus(status Status) *FakeUserBuilder {
	f.user.status = status

	return f
}

func (f *FakeUserBuilder) Build() *User {
	return f.user
}

func (f *FakeUserBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *User {
	f.t.Helper()

	storage := newSqlStorage(db)

	user := f.Build()

	err := storage.Save(ctx, user)
	require.NoError(f.t, err)

	return user
}
