package tasks

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/ptr"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTasksStorage(t *testing.T) {
	t.Run("Run success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		taskRunner := newMockTaskRunner(t)

		task := NewFakeTask(t).WithTaksName("some-task").Build()

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		// First loop
		storage.On("GetNext", mock.Anything).Return(task, nil).Once()
		taskRunner.On("Run", mock.Anything, task.args).Return(nil).Once()
		storage.On("Delete", mock.Anything, task.id).Return(nil).Once()

		// Second loop
		storage.On("GetNext", mock.Anything).Return(nil, errNotFound).Once()

		err := svc.Run(context.Background())
		require.NoError(t, err)
	})

	t.Run("Run with a GetNext error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		taskRunner := newMockTaskRunner(t)

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		storage.On("GetNext", mock.Anything).Return(nil, fmt.Errorf("some-error")).Once()

		err := svc.Run(context.Background())
		require.ErrorContains(t, err, "some-error")
	})

	t.Run("Run with a Task run error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		taskRunner := newMockTaskRunner(t)

		task := NewFakeTask(t).WithTaksName("some-task").Build()

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		// First loop
		storage.On("GetNext", mock.Anything).Return(task, nil).Once()

		// Try and fail
		taskRunner.On("Run", mock.Anything, task.args).Return(errors.New("some-error")).Once()

		// Reschedule the task
		now1 := time.Now()
		tools.ClockMock.On("Now").Return(now1).Once()
		storage.On("Patch", mock.Anything, task.id, map[string]any{
			"registered_at": ptr.To(sqlstorage.SQLTime(now1.Add(defaultRetryDelay))),
			"retries":       1,
		}).Return(nil).Once()

		// Second loop
		storage.On("GetNext", mock.Anything).Return(nil, errNotFound).Once()

		err := svc.Run(context.Background())
		require.NoError(t, err)
	})

	t.Run("Run with a Task run too many error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		taskRunner := newMockTaskRunner(t)

		// This task doesn't have any remaining retry to do.
		task := NewFakeTask(t).WithTaksName("some-task").WithRetries(defaultMaxRetries).Build()

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		// First loop
		storage.On("GetNext", mock.Anything).Return(task, nil).Once()

		// Try and fail
		taskRunner.On("Run", mock.Anything, task.args).Return(errors.New("some-error")).Once()

		// Mark as failed
		storage.On("Patch", mock.Anything, task.id, map[string]any{"status": failed}).Return(nil).Once()

		// Second loop
		storage.On("GetNext", mock.Anything).Return(nil, errNotFound).Once()

		err := svc.Run(context.Background())
		require.NoError(t, err)
	})

	t.Run("Run with a status update error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		taskRunner := newMockTaskRunner(t)

		// This task doesn't have any remaining retry to do.
		task := NewFakeTask(t).WithTaksName("some-task").WithRetries(defaultMaxRetries).Build()

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		// First loop
		storage.On("GetNext", mock.Anything).Return(task, nil).Once()

		// Try and fail
		taskRunner.On("Run", mock.Anything, task.args).Return(errors.New("some-error")).Once()

		// Mark as failed
		storage.On("Patch", mock.Anything, task.id, map[string]any{"status": failed}).Return(errors.New("some-error")).Once()

		// No second loop

		err := svc.Run(context.Background())
		require.ErrorContains(t, err, "some-error")
	})
}
