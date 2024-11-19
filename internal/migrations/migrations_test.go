package migrations

import (
	"database/sql"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestStorage(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:")
	require.NoError(t, err)

	err = db.Ping()
	require.NoError(t, err)

	return db
}

func TestRunMigration(t *testing.T) {
	tools := tools.NewMock(t)
	db := newTestStorage(t)

	err := Run(db, tools)
	require.NoError(t, err)

	row := db.QueryRow(`SELECT COUNT(*) FROM sqlite_schema 
  where type='table' AND name NOT LIKE 'sqlite_%'`)

	require.NoError(t, row.Err())
	var res int
	row.Scan(&res)

	// There is more than 3 tables
	assert.Positive(t, res)
}

func TestRunMigrationTwice(t *testing.T) {
	tools := tools.NewMock(t)
	db := newTestStorage(t)

	err := Run(db, tools)
	require.NoError(t, err)

	err = Run(db, tools)
	require.NoError(t, err)
}
