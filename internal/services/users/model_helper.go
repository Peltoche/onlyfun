package users

import (
	"context"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/tools/secret"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type FakeUserBuilder struct {
	t             testing.TB
	user          *User
	roleBuilder   *perms.FakePermissionsBuilder
	avatarBuilder *medias.FakeFileMetaBuilder
}

func NewFakeUser(t testing.TB) *FakeUserBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()
	createdAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	return &FakeUserBuilder{
		t:             t,
		roleBuilder:   perms.NewFakePermissions(t),
		avatarBuilder: medias.NewFakeFileMeta(t),

		user: &User{
			id:                uuidProvider.New(),
			createdAt:         createdAt,
			passwordChangedAt: createdAt,
			role:              nil,
			username:          gofakeit.Username(),
			avatar:            uuidProvider.New(),
			password:          secret.NewText(gofakeit.Password(true, true, true, false, false, 8)),
			status:            Active,
			createdBy:         uuidProvider.New(),
		},
	}
}

func (f *FakeUserBuilder) WithUsername(username string) *FakeUserBuilder {
	f.user.username = username

	return f
}

func (f *FakeUserBuilder) WithPassword(password string) *FakeUserBuilder {
	f.user.password = secret.NewText(password)

	return f
}

func (f *FakeUserBuilder) WithAvatar(media *medias.FileMeta) *FakeUserBuilder {
	f.user.avatar = media.ID()
	f.avatarBuilder = nil

	return f
}

func (f *FakeUserBuilder) CreatedBy(user *User) *FakeUserBuilder {
	f.user.createdBy = user.ID()

	return f
}

func (f *FakeUserBuilder) WithRole(role *perms.Role) *FakeUserBuilder {
	f.user.role = role
	f.roleBuilder = nil

	return f
}

func (f *FakeUserBuilder) WithStatus(status Status) *FakeUserBuilder {
	f.user.status = status

	return f
}

func (f *FakeUserBuilder) Build() *User {
	if f.roleBuilder != nil {
		role, _ := f.roleBuilder.Build()
		f.user.role = role
	}

	if f.avatarBuilder != nil {
		avatar := f.avatarBuilder.Build()
		f.user.avatar = avatar.ID()
	}

	return f.user
}

func (f *FakeUserBuilder) Store(ctx context.Context, db sqlstorage.Querier) {
	f.t.Helper()

	storage := newSqlStorage(db)

	if f.avatarBuilder != nil {
		f.avatarBuilder.Store(ctx, db)
	}

	if f.roleBuilder != nil {
		f.roleBuilder.Store(ctx, db)
	}

	err := storage.Save(ctx, f.user)
	require.NoError(f.t, err)
}

func (f *FakeUserBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *User {
	f.t.Helper()

	user := f.Build()

	f.Store(ctx, db)

	return user
}
