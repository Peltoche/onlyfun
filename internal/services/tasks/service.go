package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/ptr"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

const (
	defaultRetryDelay = 5 * time.Second
	defaultMaxRetries = 5
)

type storage interface {
	Save(ctx context.Context, task *Task) error
	GetNext(ctx context.Context) (*Task, error)
	Patch(ctx context.Context, taskID uuid.UUID, fields map[string]any) error
	GetLastRegisteredTask(ctx context.Context, name string) (*Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	Delete(ctx context.Context, taskID uuid.UUID) error
}

type service struct {
	storage storage
	runners map[string]TaskRunner
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
		clock:   tools.Clock(),
		log:     tools.Logger(),
	}
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

		runner, ok := s.runners[task.name]
		if !ok {
			logger.Error(fmt.Sprintf("unhandled task name: %s", task.name))

			err := s.storage.Patch(ctx, task.id, map[string]any{"status": failed})
			if err != nil {
				return fmt.Errorf("failed to Patch task: %w", err)
			}

			continue
		}

		var updateErr error
		err = runner.Run(ctx, task.args)
		switch {
		case err == nil:
			logger.DebugContext(ctx, "task succeed")

			updateErr = s.storage.Delete(ctx, task.id)

		case task.retries < defaultMaxRetries:
			task.retries++
			logger.With(slog.String("error", err.Error())).
				ErrorContext(ctx, fmt.Sprintf("task failed (#%d), retry later", task.retries))

			updateErr = s.storage.Patch(ctx, task.id, map[string]any{
				"retries":       task.retries,
				"registered_at": ptr.To(sqlstorage.SQLTime(s.clock.Now().Add(defaultRetryDelay))),
			})

		default:
			task.status = failed
			updateErr = s.storage.Patch(ctx, task.id, map[string]any{"status": failed})
			logger.With(slog.String("error", err.Error())).
				ErrorContext(ctx, "task failed, too many retries")
		}

		if updateErr != nil {
			return fmt.Errorf("failed to Patch the task status: %w", err)
		}
	}
}
