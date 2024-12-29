package tasks

import (
	"context"
	"encoding/json"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
)

type Task interface {
	Priority() int
	Name() string
	Validate() error
	Args() json.RawMessage
}

type Service interface {
	RegisterTask(ctx context.Context, task Task) error
	Run(ctx context.Context) error
}

type TaskRunner interface {
	Run(ctx context.Context, args json.RawMessage) error
	Name() string
}

func Init(runners []TaskRunner, tools tools.Tools, db sqlstorage.Querier) Service {
	storage := newSqlStorage(db)

	return newService(tools, storage, runners)
}
