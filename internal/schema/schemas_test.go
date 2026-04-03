package schema

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTaskSpec_JSON(t *testing.T) {
	spec := TaskSpec{
		ID:                 "test-001",
		Goal:               "implement hello world",
		Constraints:        []string{"follow PEP 8"},
		AcceptanceCriteria: []string{"returns string"},
	}

	data, err := json.Marshal(spec)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled TaskSpec
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.ID != spec.ID {
		t.Errorf("ID = %v, want %v", unmarshaled.ID, spec.ID)
	}
	if unmarshaled.Goal != spec.Goal {
		t.Errorf("Goal = %v, want %v", unmarshaled.Goal, spec.Goal)
	}
}

func TestExecutionResult_JSON(t *testing.T) {
	now := time.Now()
	result := ExecutionResult{
		Backend:      "codex",
		Status:       "success",
		Summary:      "done",
		FilesChanged: []string{"a.go", "b.go"},
		Stdout:       "output",
		Stderr:       "",
		PatchRef:     "artifacts/patch.diff",
		StartedAt:    now,
		FinishedAt:   now,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled ExecutionResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.Backend != result.Backend {
		t.Errorf("Backend = %v, want %v", unmarshaled.Backend, result.Backend)
	}
	if unmarshaled.Status != result.Status {
		t.Errorf("Status = %v, want %v", unmarshaled.Status, result.Status)
	}
	if len(unmarshaled.FilesChanged) != len(result.FilesChanged) {
		t.Errorf("FilesChanged length = %v, want %v", len(unmarshaled.FilesChanged), len(result.FilesChanged))
	}
}

func TestReviewFinding_JSON(t *testing.T) {
	finding := ReviewFinding{
		Severity: "high",
		Type:     "missing_test",
		File:     "a.go",
		Message:  "test coverage insufficient",
	}

	data, err := json.Marshal(finding)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled ReviewFinding
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.Severity != finding.Severity {
		t.Errorf("Severity = %v, want %v", unmarshaled.Severity, finding.Severity)
	}
	if unmarshaled.Type != finding.Type {
		t.Errorf("Type = %v, want %v", unmarshaled.Type, finding.Type)
	}
}

func TestValidationResult_JSON(t *testing.T) {
	result := ValidationResult{
		Pass:         true,
		Errors:       []string{},
		Warnings:     []string{"unused variable"},
		FilesChecked: []string{"a.go", "b.go"},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled ValidationResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.Pass != result.Pass {
		t.Errorf("Pass = %v, want %v", unmarshaled.Pass, result.Pass)
	}
	if len(unmarshaled.Errors) != len(result.Errors) {
		t.Errorf("Errors length = %v, want %v", len(unmarshaled.Errors), len(result.Errors))
	}
}

func TestSpecGateResult_JSON(t *testing.T) {
	result := SpecGateResult{
		Pass:   false,
		Reason: "validation failed",
		Issues: []string{"goal too short", "missing criteria"},
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	var unmarshaled SpecGateResult
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	if unmarshaled.Pass != result.Pass {
		t.Errorf("Pass = %v, want %v", unmarshaled.Pass, result.Pass)
	}
	if len(unmarshaled.Issues) != len(result.Issues) {
		t.Errorf("Issues length = %v, want %v", len(unmarshaled.Issues), len(result.Issues))
	}
}

func TestStopCriteria_Defaults(t *testing.T) {
	criteria := StopCriteria{}

	if criteria.SpecGatePass != false {
		t.Errorf("SpecGatePass default = %v, want false", criteria.SpecGatePass)
	}
	if criteria.ReviewBlockers != 0 {
		t.Errorf("ReviewBlockers default = %v, want 0", criteria.ReviewBlockers)
	}
	if criteria.ValidationPass != false {
		t.Errorf("ValidationPass default = %v, want false", criteria.ValidationPass)
	}
	if criteria.TestsPass != false {
		t.Errorf("TestsPass default = %v, want false", criteria.TestsPass)
	}
}

func TestExecutionState_InitialState(t *testing.T) {
	spec := TaskSpec{ID: "test-001", Goal: "test goal"}
	state := ExecutionState{
		TaskSpec:  spec,
		Iteration: 0,
	}

	if state.Iteration != 0 {
		t.Errorf("Iteration = %v, want 0", state.Iteration)
	}
	if len(state.Results) != 0 {
		t.Errorf("Results length = %v, want 0", len(state.Results))
	}
	if len(state.Findings) != 0 {
		t.Errorf("Findings length = %v, want 0", len(state.Findings))
	}
}
