package tasks

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

type fakeTaskBuilder struct {
	t    testing.TB
	task *taskData
}

func newFakeTask(t testing.TB) *fakeTaskBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()
	registeredAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	return &fakeTaskBuilder{
		t: t,
		task: &taskData{
			RegisteredAt: registeredAt,
			ID:           uuidProvider.New(),
			Name:         gofakeit.Name(),
			Status:       queuing,
			Args:         json.RawMessage(`{"foo":"bar"}`),
			Priority:     2,
			Retries:      0,
		},
	}
}

func (f *fakeTaskBuilder) WithStatus(status Status) *fakeTaskBuilder {
	f.task.Status = status

	return f
}

func (f *fakeTaskBuilder) RegisteredAt(date time.Time) *fakeTaskBuilder {
	f.task.RegisteredAt = date

	return f
}

func (f *fakeTaskBuilder) WithPriority(priority int) *fakeTaskBuilder {
	f.task.Priority = priority

	return f
}

func (f *fakeTaskBuilder) WithRetries(retries int) *fakeTaskBuilder {
	f.task.Retries = retries

	return f
}

func (f *fakeTaskBuilder) WithTaksName(name string) *fakeTaskBuilder {
	f.task.Name = name

	return f
}

func (f *fakeTaskBuilder) Build() *taskData {
	return f.task
}

func (f *fakeTaskBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *taskData {
	f.t.Helper()

	storage := newSqlStorage(db)

	task := f.Build()

	err := storage.Save(ctx, task)
	require.NoError(f.t, err)

	return task
}
