package tasks

import (
	"context"
	"encoding/json"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
)

type Service interface {
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
