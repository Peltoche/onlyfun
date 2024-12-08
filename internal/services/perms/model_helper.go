package perms

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools/ptr"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type FakePermissionsBuilder struct {
	t           testing.TB
	role        *Role
	permissions []Permission
}

func NewFakePermissions(t testing.TB) *FakePermissionsBuilder {
	t.Helper()

	randomDefaultRole := gofakeit.RandomMapKey(DefaultRoles).(Role)

	return &FakePermissionsBuilder{
		t:           t,
		role:        ptr.To(randomDefaultRole),
		permissions: DefaultRoles[randomDefaultRole],
	}
}

func (f *FakePermissionsBuilder) WithName(name string) *FakePermissionsBuilder {
	f.role = ptr.To(Role(name))

	return f
}

func (f *FakePermissionsBuilder) WithPermissions(perms ...Permission) *FakePermissionsBuilder {
	f.permissions = perms

	return f
}

func (f *FakePermissionsBuilder) Build() (*Role, []Permission) {
	return f.role, f.permissions
}

func (f *FakePermissionsBuilder) Store(ctx context.Context, db sqlstorage.Querier) {
	f.t.Helper()

	role, permissions := f.Build()

	storage := newSqlStorage(db)

	err := storage.Save(ctx, role, permissions)
	require.NoError(f.t, err)
}

func (f *FakePermissionsBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) (*Role, []Permission) {
	f.t.Helper()

	role, permissions := f.Build()

	f.Store(ctx, db)

	return role, permissions
}
