package roles

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/Peltoche/onlyfun/internal/tools/sqlstorage"
)

const (
	tableName     = "roles"
	permSeparator = ","
)

var errNotFound = errors.New("not found")

var allFields = []string{"name", "permissions"}

type sqlStorage struct {
	db sqlstorage.Querier
}

func newSqlStorage(db sqlstorage.Querier) *sqlStorage {
	return &sqlStorage{db}
}

func (s *sqlStorage) Save(ctx context.Context, r *Role) error {
	var rawPerms strings.Builder
	rawPerms.Grow(len(r.permissions) * 15)

	if len(r.permissions) > 0 {
		rawPerms.WriteString(string(r.permissions[0]))

		for _, p := range r.permissions[1:] {
			rawPerms.WriteString(permSeparator)
			rawPerms.WriteString(string(p))
		}
	}

	_, err := sq.
		Insert(tableName).
		Columns(allFields...).
		Values(r.name, rawPerms.String()).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetAll(ctx context.Context) ([]Role, error) {
	rows, err := sq.
		Select(allFields...).
		From(tableName).
		RunWith(s.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	roles := []Role{}
	for rows.Next() {
		var res Role

		var rawPerms string

		err := rows.Scan(
			&res.name,
			&rawPerms)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a row: %w", err)
		}

		for _, permStr := range strings.Split(rawPerms, permSeparator) {
			res.permissions = append(res.permissions, Permission(permStr))
		}

		roles = append(roles, res)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return roles, nil
}

func (s *sqlStorage) GetByName(ctx context.Context, roleName string) (*Role, error) {
	res := Role{}

	var rawPerms string

	err := sq.
		Select(allFields...).
		From(tableName).
		Where(sq.Eq{"name": roleName}).
		RunWith(s.db).
		ScanContext(ctx,
			&res.name,
			&rawPerms)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	for _, permStr := range strings.Split(rawPerms, permSeparator) {
		res.permissions = append(res.permissions, Permission(permStr))
	}

	return &res, nil
}
