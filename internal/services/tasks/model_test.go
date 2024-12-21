package tasks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Task_Getters(t *testing.T) {
	p := NewFakeTask(t).Build()

	assert.Equal(t, p.registeredAt, p.RegisteredAt())
	assert.Equal(t, p.id, p.ID())
	assert.Equal(t, p.name, p.Name())
	assert.Equal(t, p.status, p.Status())
	assert.Equal(t, p.args, p.Args())
	assert.Equal(t, p.priority, p.Priority())
	assert.Equal(t, p.retries, p.Retries())
}
