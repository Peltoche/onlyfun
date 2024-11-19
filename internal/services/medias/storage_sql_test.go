package medias

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserSqlStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := sqlstorage.NewTestStorage(t)
	store := newSqlStorage(db)

	// Data
	meta := NewFakeFileMeta(t).Build()

	t.Run("Save success", func(t *testing.T) {
		// Run
		err := store.Save(ctx, meta)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByID success", func(t *testing.T) {
		// Run
		res, err := store.GetByID(ctx, meta.ID())

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, meta, res)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		// Run
		res, err := store.GetByID(ctx, "some-invalid-id")

		// Asserts
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("Delete success", func(t *testing.T) {
		// Run
		err := store.Delete(ctx, meta.ID())

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByID a deleted file", func(t *testing.T) {
		// Run
		res, err := store.GetByID(ctx, meta.ID())

		// Asserts
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})
}
