package adapter

import (
	"context"

	"meta_harness/internal/schema"
)

// ExecutorBackend defines the interface for execution backends
type ExecutorBackend interface {
	Execute(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error)
}

// BackendType represents the type of backend
type BackendType string

const (
	BackendCodex  BackendType = "codex"
	BackendClaude BackendType = "claude"
)

// NewBackend creates a new executor backend by type
func NewBackend(backendType BackendType) (ExecutorBackend, error) {
	switch backendType {
	case BackendCodex:
		return NewCodexBackend()
	case BackendClaude:
		return NewClaudeBackend()
	default:
		return nil, &UnknownBackendError{Type: backendType}
	}
}

// UnknownBackendError represents an unknown backend type
type UnknownBackendError struct {
	Type BackendType
}

func (e *UnknownBackendError) Error() string {
	return "unknown backend type: " + string(e.Type)
}
