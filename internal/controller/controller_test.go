package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	"meta_harness/internal/adapter"
	"meta_harness/internal/schema"
)

// MockExecutorBackend is a mock implementation for testing
type MockExecutorBackend struct {
	ExecuteFunc func(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error)
}

func (m *MockExecutorBackend) Execute(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, spec)
	}
	return schema.ExecutionResult{}, nil
}

// MockReviewer is a mock implementation for testing
type MockReviewer struct {
	ReviewFunc func(spec schema.TaskSpec, result schema.ExecutionResult) ([]schema.ReviewFinding, error)
}

func (m *MockReviewer) Review(spec schema.TaskSpec, result schema.ExecutionResult) ([]schema.ReviewFinding, error) {
	if m.ReviewFunc != nil {
		return m.ReviewFunc(spec, result)
	}
	return nil, nil
}

// MockValidator is a mock implementation for testing
type MockValidator struct {
	ValidateFunc func(files []string) schema.ValidationResult
}

func (m *MockValidator) Validate(files []string) schema.ValidationResult {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(files)
	}
	return schema.ValidationResult{Pass: true}
}

// MockSpecGate is a mock implementation for testing
type MockSpecGate struct {
	CheckFunc func(spec schema.TaskSpec) schema.SpecGateResult
}

func (m *MockSpecGate) Check(spec schema.TaskSpec) schema.SpecGateResult {
	if m.CheckFunc != nil {
		return m.CheckFunc(spec)
	}
	return schema.SpecGateResult{Pass: true}
}

func TestController_NewController(t *testing.T) {
	backend := &MockExecutorBackend{}
	cfg := ControllerConfig{
		Backend:       backend,
		MaxIterations: 3,
	}

	ctrl := NewController(cfg)
	if ctrl == nil {
		t.Fatal("NewController returned nil")
	}
	if ctrl.maxIterations != 3 {
		t.Errorf("maxIterations = %d, want 3", ctrl.maxIterations)
	}
}

func TestController_NewController_Defaults(t *testing.T) {
	backend := &MockExecutorBackend{}
	cfg := ControllerConfig{
		Backend: backend,
	}

	ctrl := NewController(cfg)
	// Should default to 5 iterations
	if ctrl.maxIterations != 5 {
		t.Errorf("default maxIterations = %d, want 5", ctrl.maxIterations)
	}
}

func TestController_Run_SpecGateFails(t *testing.T) {
	backend := &MockExecutorBackend{}
	specGate := &MockSpecGate{
		CheckFunc: func(spec schema.TaskSpec) schema.SpecGateResult {
			return schema.SpecGateResult{
				Pass:   false,
				Reason: "spec too vague",
				Issues: []string{"goal is empty"},
			}
		},
	}

	ctrl := NewController(ControllerConfig{
		Backend:  backend,
		SpecGate: specGate,
	})

	spec := schema.TaskSpec{
		ID:                 "test-001",
		Goal:               "",
		AcceptanceCriteria: []string{"test"},
	}

	ctx := context.Background()
	state, err := ctrl.Run(ctx, spec)

	if err == nil {
		t.Error("Run() expected error, got nil")
	}
	if state.Status != "FAILED" {
		t.Errorf("state.Status = %v, want FAILED", state.Status)
	}
}

func TestController_Run_SuccessfulExecution(t *testing.T) {
	backend := &MockExecutorBackend{
		ExecuteFunc: func(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
			return schema.ExecutionResult{
				Backend:      "mock",
				Status:       "success",
				Summary:      "done",
				FilesChanged: []string{"a.go", "b.go"},
				StartedAt:    time.Now(),
				FinishedAt:   time.Now(),
			}, nil
		},
	}

	specGate := &MockSpecGate{
		CheckFunc: func(spec schema.TaskSpec) schema.SpecGateResult {
			return schema.SpecGateResult{Pass: true}
		},
	}

	reviewer := &MockReviewer{
		ReviewFunc: func(spec schema.TaskSpec, result schema.ExecutionResult) ([]schema.ReviewFinding, error) {
			return []schema.ReviewFinding{}, nil // No blockers
		},
	}

	validator := &MockValidator{
		ValidateFunc: func(files []string) schema.ValidationResult {
			return schema.ValidationResult{Pass: true}
		},
	}

	ctrl := NewController(ControllerConfig{
		Backend:       backend,
		SpecGate:      specGate,
		Reviewer:      reviewer,
		Validator:     validator,
		MaxIterations: 3,
	})

	spec := schema.TaskSpec{
		ID:                 "test-001",
		Goal:               "implement feature",
		AcceptanceCriteria: []string{"works"},
	}

	ctx := context.Background()
	state, err := ctrl.Run(ctx, spec)

	if err != nil {
		t.Errorf("Run() unexpected error: %v", err)
	}
	if state.Status != "ACCEPTED" {
		t.Errorf("state.Status = %v, want ACCEPTED", state.Status)
	}
	if state.Iteration != 1 {
		t.Errorf("state.Iteration = %d, want 1", state.Iteration)
	}
}

func TestController_Run_ExecutionError(t *testing.T) {
	backend := &MockExecutorBackend{
		ExecuteFunc: func(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
			return schema.ExecutionResult{}, errors.New("cli timeout")
		},
	}

	specGate := &MockSpecGate{
		CheckFunc: func(spec schema.TaskSpec) schema.SpecGateResult {
			return schema.SpecGateResult{Pass: true}
		},
	}

	ctrl := NewController(ControllerConfig{
		Backend:  backend,
		SpecGate: specGate,
	})

	spec := schema.TaskSpec{
		ID:                 "test-001",
		Goal:               "implement feature",
		AcceptanceCriteria: []string{"works"},
	}

	ctx := context.Background()
	state, err := ctrl.Run(ctx, spec)

	if err == nil {
		t.Error("Run() expected error, got nil")
	}
	if state.Status != "FAILED" {
		t.Errorf("state.Status = %v, want FAILED", state.Status)
	}
}

func TestController_Run_ReviewBlockers(t *testing.T) {
	backend := &MockExecutorBackend{
		ExecuteFunc: func(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
			return schema.ExecutionResult{
				Backend:      "mock",
				Status:       "success",
				FilesChanged: []string{"a.go"},
				StartedAt:    time.Now(),
				FinishedAt:   time.Now(),
			}, nil
		},
	}

	specGate := &MockSpecGate{
		CheckFunc: func(spec schema.TaskSpec) schema.SpecGateResult {
			return schema.SpecGateResult{Pass: true}
		},
	}

	reviewer := &MockReviewer{
		ReviewFunc: func(spec schema.TaskSpec, result schema.ExecutionResult) ([]schema.ReviewFinding, error) {
			return []schema.ReviewFinding{
				{Severity: "high", Type: "missing_test", File: "a.go", Message: "test missing"},
			}, nil
		},
	}

	ctrl := NewController(ControllerConfig{
		Backend:  backend,
		SpecGate: specGate,
		Reviewer: reviewer,
		MaxIterations: 2,
	})

	spec := schema.TaskSpec{
		ID:                 "test-001",
		Goal:               "implement feature",
		AcceptanceCriteria: []string{"works"},
	}

	ctx := context.Background()
	state, err := ctrl.Run(ctx, spec)

	if err == nil {
		t.Error("Run() expected error due to max iterations, got nil")
	}
	if state.Status != "FAILED" {
		t.Errorf("state.Status = %v, want FAILED", state.Status)
	}
	if state.Iteration != 2 {
		t.Errorf("state.Iteration = %v, want 2", state.Iteration)
	}
}

func TestController_Run_ValidationFails(t *testing.T) {
	backend := &MockExecutorBackend{
		ExecuteFunc: func(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
			return schema.ExecutionResult{
				Backend:      "mock",
				Status:       "success",
				FilesChanged: []string{"a.go"},
				StartedAt:    time.Now(),
				FinishedAt:   time.Now(),
			}, nil
		},
	}

	specGate := &MockSpecGate{
		CheckFunc: func(spec schema.TaskSpec) schema.SpecGateResult {
			return schema.SpecGateResult{Pass: true}
		},
	}

	reviewer := &MockReviewer{
		ReviewFunc: func(spec schema.TaskSpec, result schema.ExecutionResult) ([]schema.ReviewFinding, error) {
			return []schema.ReviewFinding{}, nil
		},
	}

	validator := &MockValidator{
		ValidateFunc: func(files []string) schema.ValidationResult {
			return schema.ValidationResult{
				Pass:   false,
				Errors: []string{"lint error"},
			}
		},
	}

	ctrl := NewController(ControllerConfig{
		Backend:   backend,
		SpecGate:  specGate,
		Reviewer:  reviewer,
		Validator: validator,
		MaxIterations: 1,
	})

	spec := schema.TaskSpec{
		ID:                 "test-001",
		Goal:               "implement feature",
		AcceptanceCriteria: []string{"works"},
	}

	ctx := context.Background()
	state, err := ctrl.Run(ctx, spec)

	if err == nil {
		t.Error("Run() expected error, got nil")
	}
	if state.Status != "FAILED" {
		t.Errorf("state.Status = %v, want FAILED", state.Status)
	}
}

func TestCountBlockers(t *testing.T) {
	tests := []struct {
		name     string
		findings []schema.ReviewFinding
		want     int
	}{
		{
			name:     "no findings",
			findings: []schema.ReviewFinding{},
			want:     0,
		},
		{
			name: "no high severity",
			findings: []schema.ReviewFinding{
				{Severity: "medium", Type: "warning"},
				{Severity: "low", Type: "info"},
			},
			want: 0,
		},
		{
			name: "one high severity",
			findings: []schema.ReviewFinding{
				{Severity: "high", Type: "error"},
				{Severity: "low", Type: "info"},
			},
			want: 1,
		},
		{
			name: "multiple high severity",
			findings: []schema.ReviewFinding{
				{Severity: "high", Type: "error1"},
				{Severity: "high", Type: "error2"},
				{Severity: "medium", Type: "warning"},
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := countBlockers(tt.findings)
			if got != tt.want {
				t.Errorf("countBlockers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_GetState(t *testing.T) {
	backend := &MockExecutorBackend{}
	ctrl := NewController(ControllerConfig{Backend: backend})

	state := ctrl.GetState()
	if state != "SPEC_PENDING" {
		t.Errorf("GetState() = %v, want SPEC_PENDING", state)
	}
}

// Ensure mock types implement the required interfaces
var _ adapter.ExecutorBackend = (*MockExecutorBackend)(nil)
