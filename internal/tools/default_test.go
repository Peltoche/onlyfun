package tools

import (
	"log/slog"
	"testing"

	"github.com/Peltoche/onlyfun/internal/tools/clock"
	"github.com/Peltoche/onlyfun/internal/tools/password"
	"github.com/Peltoche/onlyfun/internal/tools/response"
	"github.com/Peltoche/onlyfun/internal/tools/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDefaultToolbox(t *testing.T) {
	tools := NewToolbox(Config{})

	assert.IsType(t, new(clock.Default), tools.Clock())
	assert.IsType(t, new(uuid.Default), tools.UUID())
	assert.IsType(t, new(response.Default), tools.ResWriter())
	assert.IsType(t, new(slog.Logger), tools.Logger())
	assert.IsType(t, new(password.Argon2IDPassword), tools.Password())
}

func TestToolboxForTest(t *testing.T) {
	tools := NewToolboxForTest(t)

	assert.IsType(t, new(clock.Default), tools.Clock())
	assert.IsType(t, new(uuid.Default), tools.UUID())
	assert.IsType(t, new(response.Default), tools.ResWriter())
	assert.IsType(t, new(slog.Logger), tools.Logger())
	assert.IsType(t, new(password.Argon2IDPassword), tools.Password())
}
