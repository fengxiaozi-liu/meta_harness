package schema

import "time"

// TaskSpec represents the input task specification
type TaskSpec struct {
	ID                 string   `json:"id"`
	Goal               string   `json:"goal"`
	Constraints        []string `json:"constraints"`
	AcceptanceCriteria []string `json:"acceptance_criteria"`
}

// ExecutionResult represents the result from an executor backend
type ExecutionResult struct {
	Backend      string    `json:"backend"`
	Status       string    `json:"status"` // "success", "failure", "error"
	Summary      string    `json:"summary"`
	FilesChanged []string  `json:"files_changed"`
	Stdout       string    `json:"stdout"`
	Stderr       string    `json:"stderr"`
	PatchRef     string    `json:"patch_ref"`     // path to artifact
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
}

// ReviewFinding represents a single review finding
type ReviewFinding struct {
	Severity string `json:"severity"` // "high", "medium", "low"
	Type     string `json:"type"`      // "missing_test", "convention_violation", etc.
	File     string `json:"file"`
	Message  string `json:"message"`
}

// ValidationResult represents the result of validation checks
type ValidationResult struct {
	Pass         bool     `json:"pass"`
	Errors       []string `json:"errors"`
	Warnings     []string `json:"warnings"`
	FilesChecked []string `json:"files_checked"`
}

// SpecGateResult represents the result of spec gate check
type SpecGateResult struct {
	Pass   bool     `json:"pass"`
	Reason string   `json:"reason,omitempty"`
	Issues []string `json:"issues,omitempty"`
}

// StopCriteria defines when the system should stop
type StopCriteria struct {
	SpecGatePass       bool `json:"spec_gate_pass"`
	ReviewBlockers     int  `json:"review_blockers"`
	ValidationPass     bool `json:"validation_pass"`
	TestsPass          bool `json:"tests_pass"`
	IterationCount     int  `json:"iteration_count"`
	MaxIterations      int  `json:"max_iterations"`
}

// ExecutionState represents the current state of execution
type ExecutionState struct {
	TaskSpec     TaskSpec          `json:"task_spec"`
	Status       string            `json:"status"` // State from state machine as string
	Iteration    int               `json:"iteration"`
	Results      []ExecutionResult  `json:"results"`
	Findings     []ReviewFinding    `json:"findings"`
	Validation   ValidationResult   `json:"validation"`
	StopCriteria StopCriteria      `json:"stop_criteria"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}
