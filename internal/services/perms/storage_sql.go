package perms

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
	tableName     = "permissions"
	permSeparator = ","
)

var errNotFound = errors.New("not found")

var allFields = []string{"role", "permissions"}

type sqlStorage struct {
	db sqlstorage.Querier
}

func newSqlStorage(db sqlstorage.Querier) *sqlStorage {
	return &sqlStorage{db}
}

func (s *sqlStorage) Save(ctx context.Context, role *Role, perms []Permission) error {
	var rawPerms strings.Builder
	rawPerms.Grow(len(perms) * 15)

	if len(perms) > 0 {
		rawPerms.WriteString(string(perms[0]))

		for _, p := range perms[1:] {
			rawPerms.WriteString(permSeparator)
			rawPerms.WriteString(string(p))
		}
	}

	_, err := sq.
		Insert(tableName).
		Columns(allFields...).
		Values(role, rawPerms.String()).
		RunWith(s.db).
		ExecContext(ctx)
	if err != nil {
		return fmt.Errorf("sql error: %w", err)
	}

	return nil
}

func (s *sqlStorage) GetAll(ctx context.Context) (map[Role][]Permission, error) {
	rows, err := sq.
		Select(allFields...).
		From(tableName).
		RunWith(s.db).
		QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	res := map[Role][]Permission{}
	for rows.Next() {
		var role Role
		var rawPerms string

		err := rows.Scan(
			&role,
			&rawPerms)
		if err != nil {
			return nil, fmt.Errorf("failed to scan a row: %w", err)
		}

		var permissions []Permission
		for _, permStr := range strings.Split(rawPerms, permSeparator) {
			permissions = append(permissions, Permission(permStr))
		}

		res[role] = permissions
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("scan error: %w", err)
	}

	return res, nil
}

func (s *sqlStorage) GetPermissions(ctx context.Context, role *Role) ([]Permission, error) {
	var rawPerms string

	err := sq.
		Select(allFields[1:]...). // skip the role, it's already given
		From(tableName).
		Where(sq.Eq{"role": role}).
		RunWith(s.db).
		ScanContext(ctx, &rawPerms)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("sql error: %w", err)
	}

	var permissions []Permission
	for _, permStr := range strings.Split(rawPerms, permSeparator) {
		permissions = append(permissions, Permission(permStr))
	}

	return permissions, nil
}
