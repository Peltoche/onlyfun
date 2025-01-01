package taskrunner

import (
	"context"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/require"
)

func TestTasksSqlStorage(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("Save success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		task := newFakeTask(t).Build()

		err := store.Save(ctx, task)
		require.NoError(t, err)
	})

	t.Run("GetNext success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		now := time.Now().UTC()

		task1 := newFakeTask(t).
			WithStatus(queuing).
			WithPriority(3).
			RegisteredAt(now).
			BuildAndStore(ctx, db)
		_ = newFakeTask(t).
			WithStatus(queuing).
			WithPriority(3).
			RegisteredAt(now.Add(time.Second)).
			BuildAndStore(ctx, db)

		res, err := store.GetNext(ctx)
		require.NoError(t, err)
		require.Equal(t, task1, res)
	})

	t.Run("GetNext with no queuing", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		_ = newFakeTask(t).
			WithStatus(failed).
			WithPriority(3).
			BuildAndStore(ctx, db)
		_ = newFakeTask(t).
			WithStatus(failed).
			WithPriority(3).
			BuildAndStore(ctx, db)

		res, err := store.GetNext(ctx)
		require.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})

	t.Run("Update success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		task := newFakeTask(t).
			WithStatus(queuing).
			WithPriority(3).
			BuildAndStore(ctx, db)

		updatedTask := *task
		updatedTask.Status = failed

		err := store.Update(ctx, &updatedTask)
		require.NoError(t, err)

		res, err := store.GetByID(ctx, task.ID)
		require.NoError(t, err)
		require.Equal(t, &updatedTask, res)
	})

	t.Run("Delete success", func(t *testing.T) {
		t.Parallel()

		db := sqlstorage.NewTestStorage(t)
		store := newSqlStorage(db)

		task := newFakeTask(t).
			WithStatus(queuing).
			WithPriority(3).
			BuildAndStore(ctx, db)

		err := store.Delete(ctx, task.ID)
		require.NoError(t, err)

		res, err := store.GetByID(ctx, task.ID)
		require.Nil(t, res)
		require.ErrorIs(t, err, errNotFound)
	})
}
