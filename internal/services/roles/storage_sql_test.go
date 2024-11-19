package roles

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

	t.Run("GetByID not found", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		// Run
		res, err := store.GetByName(ctx, "some-invalid-id")

		// Asserts
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("Save success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role := NewFakeRole(t).
			WithName("admin").
			WithPermissions(UploadPost, Moderation).
			Build()

		// Run
		err := store.Save(ctx, role)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByID success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role := NewFakeRole(t).
			WithName("admin").
			WithPermissions(UploadPost, Moderation).
			BuildAndStore(ctx, db)

		// Run
		res, err := store.GetByName(ctx, role.Name())

		// Asserts
		require.NoError(t, err)
		require.Equal(t, role, res)
	})

	t.Run("GetAll success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role := NewFakeRole(t).
			WithName("admin").
			WithPermissions(UploadPost, Moderation).
			BuildAndStore(ctx, db)

		// Run
		res, err := store.GetAll(ctx)

		// Asserts
		require.NoError(t, err)
		require.Equal(t, []Role{*role}, res)
	})
}
