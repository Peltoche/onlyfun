package perms

import (
	"context"
	"fmt"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
)

type WithRole interface {
	Role() *Role
}

type Service interface {
	IsAuthorized(withRole WithRole, askedPerm Permission) bool
}

func Init(ctx context.Context, db sqlstorage.Querier, tools tools.Tools) (Service, error) {
	storage := newSqlStorage(db)

	svc := newService(tools, storage)

	err := svc.bootstrap(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to bootstrap the permissions: %w", err)
	}

	return svc, nil
}
