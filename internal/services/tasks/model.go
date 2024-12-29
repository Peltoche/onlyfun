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

type taskData struct {
	RegisteredAt time.Time
	ID           uuid.UUID
	Name         string
	Status       Status
	Args         json.RawMessage
	Priority     int
	Retries      int
}
