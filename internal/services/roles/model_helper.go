package roles

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type FakeRoleBuilder struct {
	t    testing.TB
	role *Role
}

func NewFakeRole(t testing.TB) *FakeRoleBuilder {
	t.Helper()

	return &FakeRoleBuilder{
		t: t,
		role: &Role{
			name:        gofakeit.JobTitle(),
			permissions: DefaultUserRole.permissions,
		},
	}
}

func (f *FakeRoleBuilder) WithName(name string) *FakeRoleBuilder {
	f.role.name = name

	return f
}

func (f *FakeRoleBuilder) WithPermissions(perms ...Permission) *FakeRoleBuilder {
	f.role.permissions = append(f.role.permissions, perms...)

	return f
}

func (f *FakeRoleBuilder) Build() *Role {
	return f.role
}

func (f *FakeRoleBuilder) Store(ctx context.Context, db sqlstorage.Querier) {
	f.t.Helper()

	storage := newSqlStorage(db)

	err := storage.Save(ctx, f.role)
	require.NoError(f.t, err)
}

func (f *FakeRoleBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *Role {
	f.t.Helper()

	role := f.Build()

	f.Store(ctx, db)

	return role
}
