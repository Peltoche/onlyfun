package taskrunner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestTasksService(t *testing.T) {
	ctx := context.Background()

	t.Run("Run success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		taskRunner := newMockTaskRunner(t)

		task := newFakeTask(t).WithTaksName("some-task").Build()

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		// First loop
		storage.On("GetNext", mock.Anything).Return(task, nil).Once()
		taskRunner.On("Run", mock.Anything, task.Args).Return(nil).Once()
		storage.On("Delete", mock.Anything, task.ID).Return(nil).Once()

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

		task := newFakeTask(t).WithTaksName("some-task").Build()

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		// First loop
		storage.On("GetNext", mock.Anything).Return(task, nil).Once()

		// Try and fail
		taskRunner.On("Run", mock.Anything, task.Args).Return(errors.New("some-error")).Once()

		// Reschedule the task
		now1 := time.Now()
		tools.ClockMock.On("Now").Return(now1).Once()

		updatedTask := *task
		updatedTask.RegisteredAt = now1.Add(defaultRetryDelay)
		updatedTask.Retries = 1

		storage.On("Update", mock.Anything, &updatedTask).Return(nil).Once()

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
		task := newFakeTask(t).WithTaksName("some-task").WithRetries(defaultMaxRetries).Build()

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		// First loop
		storage.On("GetNext", mock.Anything).Return(task, nil).Once()

		// Try and fail
		taskRunner.On("Run", mock.Anything, task.Args).Return(errors.New("some-error")).Once()

		// Mark as failed
		updatedTask := *task
		updatedTask.Status = failed
		storage.On("Update", mock.Anything, &updatedTask).Return(nil).Once()

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
		task := newFakeTask(t).WithTaksName("some-task").WithRetries(defaultMaxRetries).Build()

		taskRunner.On("Name").Return("some-task").Once()
		svc := newService(tools, storage, []TaskRunner{taskRunner})

		// First loop
		storage.On("GetNext", mock.Anything).Return(task, nil).Once()

		// Try and fail
		taskRunner.On("Run", mock.Anything, task.Args).Return(errors.New("some-error")).Once()

		// Mark as failed
		updatedTask := *task
		updatedTask.Status = failed
		storage.On("Update", mock.Anything, &updatedTask).Return(errors.New("some-error")).Once()

		// No second loop

		err := svc.Run(context.Background())
		require.ErrorContains(t, err, "some-error")
	})

	t.Run("RegisterTask success", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage, []TaskRunner{})

		task := taskStub{
			name:          "test",
			priority:      1,
			validateError: nil,
			args: map[string]string{
				"key": "value",
			},
		}

		taskID := uuid.UUID("some-uuid")

		now := time.Now()
		tools.UUIDMock.On("New").Return(taskID).Once()
		tools.ClockMock.On("Now").Return(now).Once()
		storage.On("Save", ctx, &taskData{
			ID:           taskID,
			Priority:     1,
			Status:       queuing,
			Name:         "test",
			RegisteredAt: now,
			Args:         json.RawMessage(`{"key":"value"}`),
		}).Return(nil).Once()

		err := svc.RegisterTask(ctx, &task)
		require.NoError(t, err)
	})

	t.Run("RegisterTask with a validation error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage, []TaskRunner{})

		task := taskStub{
			name:          "test",
			priority:      1,
			validateError: errors.New("some-error"), // return a validation error
			args:          map[string]string{},
		}

		err := svc.RegisterTask(ctx, &task)
		require.ErrorIs(t, err, errs.ErrValidation)
		require.ErrorContains(t, err, "some-error")
	})

	t.Run("RegisterTask with a storage error", func(t *testing.T) {
		tools := tools.NewMock(t)
		storage := newMockStorage(t)
		svc := newService(tools, storage, []TaskRunner{})

		task := taskStub{
			name:          "test",
			priority:      1,
			validateError: nil,
			args: map[string]string{
				"key": "value",
			},
		}

		taskID := uuid.UUID("some-uuid")

		now := time.Now()
		tools.UUIDMock.On("New").Return(taskID).Once()
		tools.ClockMock.On("Now").Return(now).Once()
		storage.On("Save", ctx, &taskData{
			ID:           taskID,
			Priority:     1,
			Status:       queuing,
			Name:         "test",
			RegisteredAt: now,
			Args:         json.RawMessage(`{"key":"value"}`),
		}).Return(errors.New("some-error")).Once()

		err := svc.RegisterTask(ctx, &task)
		require.ErrorIs(t, err, errs.ErrInternal)
		require.ErrorContains(t, err, "some-error")
	})
}

type taskStub struct {
	name          string
	priority      int
	validateError error
	args          any
}

func (s *taskStub) Name() string {
	return s.name
}

func (s *taskStub) Priority() int {
	return s.priority
}

func (s *taskStub) Validate() error {
	return s.validateError
}

func (s *taskStub) Args() json.RawMessage {
	res, _ := json.Marshal(s.args)

	return res
}
