package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

const (
	defaultRetryDelay = 5 * time.Second
	defaultMaxRetries = 5
)

type storage interface {
	Save(ctx context.Context, task *taskData) error
	GetNext(ctx context.Context) (*taskData, error)
	GetByID(ctx context.Context, id uuid.UUID) (*taskData, error)
	Update(ctx context.Context, task *taskData) error
	Delete(ctx context.Context, taskID uuid.UUID) error
}

type service struct {
	storage storage
	runners map[string]TaskRunner
	uuid    uuid.Service
	clock   clock.Clock
	log     *slog.Logger
}

func newService(tools tools.Tools, storage storage, runners []TaskRunner) *service {
	runnerMap := make(map[string]TaskRunner, len(runners))

	for _, runner := range runners {
		runnerMap[runner.Name()] = runner
	}

	return &service{
		storage: storage,
		runners: runnerMap,
		uuid:    tools.UUID(),
		clock:   tools.Clock(),
		log:     tools.Logger(),
	}
}

func (s *service) RegisterTask(ctx context.Context, task Task) error {
	err := task.Validate()
	if err != nil {
		return errs.Validation(err)
	}

	err = s.storage.Save(ctx, &taskData{
		ID:           s.uuid.New(),
		Priority:     task.Priority(),
		Status:       queuing,
		Name:         task.Name(),
		RegisteredAt: s.clock.Now(),
		Args:         task.Args(),
	})
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to save the %q job : %w", task.Name(), err))
	}

	return nil
}

func (s *service) Run(ctx context.Context) error {
	for {
		task, err := s.storage.GetNext(ctx)
		if errors.Is(err, errNotFound) {
			// All the tasks have been processed
			return nil
		}

		if err != nil {
			return fmt.Errorf("failed to GetNext task: %w", err)
		}

		logger := s.log.With(slog.Any("task", task))

		runner, ok := s.runners[task.Name]
		if !ok {
			logger.Error(fmt.Sprintf("unhandled task name: %s", task.Name))

			task.Status = failed

			err := s.storage.Update(ctx, task)
			if err != nil {
				return fmt.Errorf("failed to Patch task: %w", err)
			}

			continue
		}

		var updateErr error
		err = runner.Run(ctx, task.Args)
		switch {
		case err == nil:
			logger.DebugContext(ctx, "task succeed")

			updateErr = s.storage.Delete(ctx, task.ID)

		case task.Retries < defaultMaxRetries:
			task.Retries++
			logger.With(slog.String("error", err.Error())).
				ErrorContext(ctx, fmt.Sprintf("task failed (#%d), retry later", task.Retries))

			task.RegisteredAt = s.clock.Now().Add(defaultRetryDelay)

			updateErr = s.storage.Update(ctx, task)

		default:
			task.Status = failed
			updateErr = s.storage.Update(ctx, task)
			logger.With(slog.String("error", err.Error())).
				ErrorContext(ctx, "task failed, too many retries")
		}

		if updateErr != nil {
			return fmt.Errorf("failed to Patch the task status: %w", err)
		}
	}
}
