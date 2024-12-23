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

type FakeTaskBuilder struct {
	t    testing.TB
	task *Task
}

func NewFakeTask(t testing.TB) *FakeTaskBuilder {
	t.Helper()

	uuidProvider := uuid.NewProvider()
	registeredAt := gofakeit.DateRange(time.Now().Add(-time.Hour*1000), time.Now())

	return &FakeTaskBuilder{
		t: t,
		task: &Task{
			registeredAt: registeredAt,
			id:           uuidProvider.New(),
			name:         gofakeit.Name(),
			status:       queuing,
			args:         json.RawMessage(`{"foo":"bar"}`),
			priority:     2,
			retries:      0,
		},
	}
}

func (f *FakeTaskBuilder) WithRetries(retries int) *FakeTaskBuilder {
	f.task.retries = retries

	return f
}

func (f *FakeTaskBuilder) WithTaksName(name string) *FakeTaskBuilder {
	f.task.name = name

	return f
}

func (f *FakeTaskBuilder) Build() *Task {
	return f.task
}

func (f *FakeTaskBuilder) BuildAndStore(ctx context.Context, db sqlstorage.Querier) *Task {
	f.t.Helper()

	storage := newSqlStorage(db)

	task := f.Build()

	err := storage.Save(ctx, task)
	require.NoError(f.t, err)

	return task
}
