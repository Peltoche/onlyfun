package taskrunner

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

func (s *sqlStorage) Save(ctx context.Context, task *taskData) error {
	rawArgs, err := json.Marshal(task.Args)
	if err != nil {
		return fmt.Errorf("failed to marshal the args: %w", err)
	}

	_, err = sq.
		Insert(tableName).
		Columns(allFields...).
		Values(
			task.ID,
			task.Priority,
			task.Name,
			task.Status,
			task.Retries,
			ptr.To(sqlstorage.SQLTime(task.RegisteredAt)),
			rawArgs,
		).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetByID(ctx context.Context, id uuid.UUID) (*taskData, error) {
	row := sq.
		Select(allFields...).
		From(tableName).
		Where(sq.Eq{"id": id}).
		OrderBy("priority", "registered_at").
		RunWith(s.db).
		QueryRowContext(ctx)

	return s.scanRow(row)
}

func (s *sqlStorage) GetNext(ctx context.Context) (*taskData, error) {
	row := sq.
		Select(allFields...).
		From(tableName).
		Where(sq.Eq{"status": queuing}).
		OrderBy("priority", "registered_at ASC").
		RunWith(s.db).
		QueryRowContext(ctx)

	return s.scanRow(row)
}

func (s *sqlStorage) Update(ctx context.Context, task *taskData) error {
	rawArgs, err := json.Marshal(task.Args)
	if err != nil {
		return fmt.Errorf("failed to marshal the args: %w", err)
	}

	_, err = sq.Update(tableName).
		SetMap(map[string]any{
			"priority":      task.Priority,
			"name":          task.Name,
			"status":        task.Status,
			"retries":       task.Retries,
			"registered_at": ptr.To(sqlstorage.SQLTime(task.RegisteredAt)),
			"args":          rawArgs,
		}).
		Where(sq.Eq{"id": task.ID}).
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

func (s *sqlStorage) scanRow(row sq.RowScanner) (*taskData, error) {
	var res taskData
	var rawArgs json.RawMessage
	var sqlRegisteredAt sqlstorage.SQLTime

	err := row.Scan(
		&res.ID,
		&res.Priority,
		&res.Name,
		&res.Status,
		&res.Retries,
		&sqlRegisteredAt,
		&rawArgs,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan the sql result: %w", err)
	}

	res.RegisteredAt = sqlRegisteredAt.Time()
	err = json.Unmarshal(rawArgs, &res.Args)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal the args: %w", err)
	}

	return &res, nil
}
