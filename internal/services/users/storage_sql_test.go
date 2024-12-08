package users

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/services/medias"
	"github.com/Peltoche/onlyfun/internal/services/perms"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Users_SqlStorage(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("GetAll with nothing", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		// Run
		res, err := store.GetAll(ctx, &sqlstorage.PaginateCmd{Limit: 10})

		// Asserts
		require.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("Save success", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		role, _ := perms.NewFakePermissions(t).BuildAndStore(ctx, db)
		avatar := medias.NewFakeFileMeta(t).BuildAndStore(ctx, db)
		user := NewFakeUser(t).WithRole(role).WithAvatar(avatar).Build()

		// Run
		err := store.Save(ctx, user)

		// Asserts
		require.NoError(t, err)
	})

	t.Run("GetByID success", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		user := NewFakeUser(t).BuildAndStore(ctx, db)

		// Run
		res, err := store.GetByID(ctx, user.ID())

		// Asserts
		assert.NotNil(t, res)
		require.NoError(t, err)
		assert.Equal(t, user, res)
	})

	t.Run("GetByID not found", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		// Run
		res, err := store.GetByID(ctx, "some-invalid-id")

		// Asserts
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("Patch success", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		user := NewFakeUser(t).WithUsername("old-username").BuildAndStore(ctx, db)

		// Run
		err := store.Patch(ctx, user.ID(), map[string]any{"username": "new-username"})
		require.NoError(t, err)

		// Asserts
		res, err := store.GetByID(ctx, user.ID())
		require.NoError(t, err)
		assert.Equal(t, "new-username", res.username)
	})

	t.Run("GetByUsername success", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		user := NewFakeUser(t).BuildAndStore(ctx, db)

		// Run
		res, err := store.GetByUsername(ctx, user.Username())

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, user, res)
	})

	t.Run("GetByUsername not found", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		// Run
		res, err := store.GetByUsername(ctx, "some-invalid-username")

		// Asserts
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("GetAll success", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		user := NewFakeUser(t).BuildAndStore(ctx, db)

		// Run
		res, err := store.GetAll(ctx, &sqlstorage.PaginateCmd{Limit: 10})

		// Asserts
		require.NoError(t, err)
		assert.Equal(t, []User{*user}, res)
	})

	t.Run("HardDelete success", func(t *testing.T) {
		t.Parallel()
		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		user := NewFakeUser(t).BuildAndStore(ctx, db)

		// Run
		err := store.HardDelete(ctx, user.ID())
		require.NoError(t, err)

		// Asserts
		res, err := store.GetAll(ctx, nil)
		require.NoError(t, err)
		assert.Empty(t, res)
	})
}
