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
			permissions: []Permission{},
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

func (f *FakeRoleBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *Role {
	f.t.Helper()

	storage := newSqlStorage(db)

	role := f.Build()

	err := storage.Save(ctx, role)
	require.NoError(f.t, err)

	return role
}
