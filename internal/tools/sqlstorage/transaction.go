package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/Peltoche/onlyfun/internal/tools/errs"
	"github.com/Peltoche/onlyfun/internal/tools/logger"
)

var txKey = contextKey{"txKey"}

type contextKey struct {
	key string
}

type TransacService struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewTransacGenerator(db *sql.DB, tools tools.Tools) *TransacService {
	return &TransacService{
		db:     db,
		logger: tools.Logger(),
	}
}

// WithinTransaction runs function within transaction
//
// The transaction commits when function were finished without error
func (t *TransacService) WithinTransaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	// begin transaction
	tx, err := t.db.Begin()
	if err != nil {
		return errs.Internal(fmt.Errorf("failed to start transaction: %w", err))
	}

	defer func() {
		// finalize transaction on panic, etc.
		errTx := tx.Commit()
		if errTx != nil && !errors.Is(errTx, sql.ErrTxDone) {
			logger.LogEntrySetAttrs(ctx, slog.String("commit-error", err.Error()))
		}
	}()

	// run callback
	err = tFunc(context.WithValue(ctx, txKey, tx))
	if err != nil {
		// if error, rollback
		if errRollback := tx.Rollback(); errRollback != nil {
			logger.LogEntrySetAttrs(ctx, slog.String("rollback-error", err.Error()))
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit(); errCommit != nil {
		logger.LogEntrySetAttrs(ctx, slog.String("commit-error", err.Error()))
	}
	return nil
}
