package testutil

import (
	"context"
	"time"

	"meta_harness/internal/schema"
)

// MockBackend implements adapter.ExecutorBackend for testing
type MockBackend struct {
	ResultToReturn schema.ExecutionResult
	ErrToReturn    error
	CallCount      int
}

// NewMockBackend creates a mock backend with default success response
func NewMockBackend() *MockBackend {
	return &MockBackend{
		ResultToReturn: schema.ExecutionResult{
			Backend:      "mock",
			Status:       "success",
			Summary:      "mock execution complete",
			FilesChanged: []string{"mock.go"},
			StartedAt:    time.Now(),
			FinishedAt:   time.Now(),
		},
	}
}

// NewFailingMockBackend creates a mock backend that always fails
func NewFailingMockBackend(err error) *MockBackend {
	return &MockBackend{
		ErrToReturn: err,
	}
}

// Execute runs the mock backend
func (m *MockBackend) Execute(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
	m.CallCount++
	if m.ErrToReturn != nil {
		return schema.ExecutionResult{}, m.ErrToReturn
	}
	return m.ResultToReturn, nil
}

// Reset resets the call count
func (m *MockBackend) Reset() {
	m.CallCount = 0
}
