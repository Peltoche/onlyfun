package tasks

import (
	"encoding/json"
	"time"

	"github.com/Peltoche/onlyfun/internal/tools/uuid"
)

type Status string

const (
	queuing Status = "queuing"
	failed  Status = "failed"
)

type Task struct {
	registeredAt time.Time
	id           uuid.UUID
	name         string
	status       Status
	args         json.RawMessage
	priority     int
	retries      int
}

func (t Task) RegisteredAt() time.Time { return t.registeredAt }
func (t Task) ID() uuid.UUID           { return t.id }
func (t Task) Name() string            { return t.name }
func (t Task) Status() Status          { return t.status }
func (t Task) Args() json.RawMessage   { return t.args }
func (t Task) Priority() int           { return t.priority }
func (t Task) Retries() int            { return t.retries }
