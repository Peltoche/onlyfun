package tools

import (
	"log/slog"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/password"
	"github.com/Peltoche/onlyfun/internal/tools/response"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/neilotoole/slogt"
)

type Mock struct {
	ClockMock     *clock.Mock
	UUIDMock      *uuid.Mock
	LogTest       *slog.Logger
	PasswordMock  *password.Mock
	ResWriterMock *response.Mock
}

func NewMock(t *testing.T) *Mock {
	t.Helper()

	return &Mock{
		ClockMock:     clock.NewMock(t),
		UUIDMock:      uuid.NewMock(t),
		LogTest:       slogt.New(t),
		PasswordMock:  password.NewMock(t),
		ResWriterMock: response.NewMock(t),
	}
}

// Clock implements App.
func (m *Mock) Clock() clock.Clock {
	return m.ClockMock
}

// UUID implements App.
func (m *Mock) UUID() uuid.Service {
	return m.UUIDMock
}

// Logger implements App.
func (m *Mock) Logger() *slog.Logger {
	return m.LogTest
}

func (m *Mock) ResWriter() response.Writer {
	return m.ResWriterMock
}

func (m *Mock) Password() password.Password {
	return m.PasswordMock
}
