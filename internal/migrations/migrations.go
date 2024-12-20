package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Peltoche/onlyfun/internal/tools"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed *.sql
var fs embed.FS

func Run(db *sql.DB, tools tools.Tools) error {
	// Error not possible
	d, _ := iofs.New(fs, ".")

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("failed to setup the sqlite3 instance: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", d, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("failed to create a migrate manager: %w", err)
	}

	if tools != nil {
		m.Log = &migrateLogger{tools.Logger()}
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("database migration error: %w", err)
	}

	return nil
}

type migrateLogger struct {
	Logger *slog.Logger
}

func (t *migrateLogger) Printf(format string, v ...any) {
	t.Logger.Debug(fmt.Sprintf(format, v...))
}

func (t *migrateLogger) Verbose() bool {
	return true
}
