package executor

import (
	"context"
	"fmt"

	"meta_harness/internal/adapter"
	"meta_harness/internal/schema"
)

// Executor orchestrates the execution of tasks
type Executor struct {
	backend adapter.ExecutorBackend
}

// NewExecutor creates a new Executor
func NewExecutor(backend adapter.ExecutorBackend) *Executor {
	return &Executor{backend: backend}
}

// Execute runs the task using the configured backend
func (e *Executor) Execute(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
	if e.backend == nil {
		return schema.ExecutionResult{}, fmt.Errorf("no backend configured")
	}

	result, err := e.backend.Execute(ctx, spec)
	if err != nil {
		return result, fmt.Errorf("executor failed: %w", err)
	}

	return result, nil
}
