package tasks

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Peltoche/onlyfun/internal/tools/ptr"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

const tableName = "tasks"

var errNotFound = errors.New("not found")

var allFields = []string{"id", "priority", "name", "status", "retries", "registered_at", "args"}

type sqlStorage struct {
	db sqlstorage.Querier
}

func newSqlStorage(db sqlstorage.Querier) *sqlStorage {
	return &sqlStorage{db}
}

func (s *sqlStorage) Save(ctx context.Context, task *Task) error {
	rawArgs, err := json.Marshal(task.args)
	if err != nil {
		return fmt.Errorf("failed to marshal the args: %w", err)
	}

	_, err = sq.
		Insert(tableName).
		Columns(allFields...).
		Values(
			task.id,
			task.priority,
			task.name,
			task.status,
			task.retries,
			ptr.To(sqlstorage.SQLTime(task.registeredAt)),
			rawArgs,
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetLastRegisteredTask(ctx context.Context, name string) (*Task, error) {
	var res Task
	var rawArgs json.RawMessage
	var sqlRegisteredAt sqlstorage.SQLTime

	err := sq.
		Select(allFields...).
		From(tableName).
		Where(sq.Eq{"name": name}).
		OrderBy("registered_at DESC").
		Limit(1).
		RunWith(s.db).
		ScanContext(ctx,
			&res.id,
			&res.priority,
			&res.name,
			&res.status,
			&res.retries,
			&sqlRegisteredAt,
			&rawArgs,
		)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	res.registeredAt = sqlRegisteredAt.Time()
	err = json.Unmarshal(rawArgs, &res.args)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the arg: %w", err)
	}

	return &res, nil
}

func (s *sqlStorage) GetByID(ctx context.Context, id uuid.UUID) (*Task, error) {
	row := sq.
		Select(allFields...).
		From(tableName).
		Where(sq.Eq{"id": id}).
		OrderBy("priority", "registered_at").
		RunWith(s.db).
		QueryRowContext(ctx)

	return s.scanRow(row)
}

func (s *sqlStorage) GetNext(ctx context.Context) (*Task, error) {
	row := sq.
		Select(allFields...).
		From(tableName).
		Where(sq.Eq{"status": queuing}).
		OrderBy("priority", "registered_at").
		RunWith(s.db).
		QueryRowContext(ctx)

	return s.scanRow(row)
}

func (s *sqlStorage) Patch(ctx context.Context, taskID uuid.UUID, fields map[string]any) error {
	_, err := sq.Update(tableName).
		SetMap(fields).
		Where(sq.Eq{"id": taskID}).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) Delete(ctx context.Context, taskID uuid.UUID) error {
	_, err := sq.Delete(tableName).
		Where(sq.Eq{"id": taskID}).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) scanRow(row sq.RowScanner) (*Task, error) {
	var res Task
	var rawArgs json.RawMessage
	var sqlRegisteredAt sqlstorage.SQLTime

	err := row.Scan(
		&res.id,
		&res.priority,
		&res.name,
		&res.status,
		&res.retries,
		&sqlRegisteredAt,
		&rawArgs,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan the sql result: %w", err)
	}

	res.registeredAt = sqlRegisteredAt.Time()
	err = json.Unmarshal(rawArgs, &res.args)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the args: %w", err)
	}

	return &res, nil
}
