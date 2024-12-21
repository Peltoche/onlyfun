package tasks

import (
	"context"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Tasks_SQLStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("Save success", func(t *testing.T) {
		db := sqlstorage.NewTestStorage(t)
		storage := newSqlStorage(db)

		task := NewFakeTask(t).Build()

		err := storage.Save(ctx, task)
		require.NoError(t, err)
	})

	t.Run("GetByID success", func(t *testing.T) {
		db := sqlstorage.NewTestStorage(t)
		storage := newSqlStorage(db)

		task := NewFakeTask(t).BuildAndStore(ctx, db)

		res, err := storage.GetByID(ctx, task.id)
		assert.Equal(t, task, res)
		require.NoError(t, err)
	})

	t.Run("GetNext success", func(t *testing.T) {
		db := sqlstorage.NewTestStorage(t)
		storage := newSqlStorage(db)

		task := NewFakeTask(t).BuildAndStore(ctx, db)

		res, err := storage.GetNext(ctx)
		require.NoError(t, err)
		assert.Equal(t, task, res)
	})

	t.Run("GetLastRegisteredTask success", func(t *testing.T) {
		db := sqlstorage.NewTestStorage(t)
		storage := newSqlStorage(db)

		task := NewFakeTask(t).BuildAndStore(ctx, db)

		res, err := storage.GetLastRegisteredTask(ctx, task.name)
		require.NoError(t, err)
		assert.Equal(t, task, res)
	})

	t.Run("GetLastRegisteredTask success", func(t *testing.T) {
		db := sqlstorage.NewTestStorage(t)
		storage := newSqlStorage(db)

		task := NewFakeTask(t).BuildAndStore(ctx, db)

		res, err := storage.GetLastRegisteredTask(ctx, task.name)
		require.NoError(t, err)
		assert.Equal(t, task, res)
	})

	t.Run("GetLastRegisteredTask with not tasks", func(t *testing.T) {
		db := sqlstorage.NewTestStorage(t)
		storage := newSqlStorage(db)

		res, err := storage.GetLastRegisteredTask(ctx, "unknown-task")
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("Delete success", func(t *testing.T) {
		db := sqlstorage.NewTestStorage(t)
		storage := newSqlStorage(db)

		task := NewFakeTask(t).BuildAndStore(ctx, db)

		err := storage.Delete(ctx, task.id)
		require.NoError(t, err)

		res, err := storage.GetByID(ctx, task.id)
		assert.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("Delete an already deleted task success", func(t *testing.T) {
		db := sqlstorage.NewTestStorage(t)
		storage := newSqlStorage(db)

		task := NewFakeTask(t).BuildAndStore(ctx, db)

		// Deleted by the previous test
		err := storage.Delete(ctx, task.id)
		require.NoError(t, err)
	})
}
