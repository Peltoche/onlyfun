package perms

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleSqlStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("GetPermissions not found", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role := Role("some-invalid-id")

		// Run
		res, err := store.GetPermissions(ctx, &role)

		// Asserts
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("Save success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role := Role("admin")
		permissions := []Permission{UploadPost, Moderation}

		// Run
		err := store.Save(ctx, &role, permissions)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByID success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role, permissions := NewFakePermissions(t).
			WithName("admin").
			WithPermissions(UploadPost, Moderation).
			BuildAndStore(ctx, db)

		// Run
		res, err := store.GetPermissions(ctx, role)

		// Asserts
		require.NoError(t, err)
		require.Equal(t, permissions, res)
	})

	t.Run("GetAll success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role, permissions := NewFakePermissions(t).
			WithName("admin").
			WithPermissions(UploadPost, Moderation).
			BuildAndStore(ctx, db)

		// Run
		res, err := store.GetAll(ctx)

		// Asserts
		require.NoError(t, err)
		require.Equal(t, map[Role][]Permission{
			*role: permissions,
		}, res)
	})
}
