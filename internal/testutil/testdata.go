package testutil

import (
	"meta_harness/internal/schema"
)

// ValidTaskSpec returns a valid task spec for testing
func ValidTaskSpec() schema.TaskSpec {
	return schema.TaskSpec{
		ID:                 "test-001",
		Goal:               "Implement a hello world function in Python",
		Constraints:        []string{"Follow PEP 8", "Use type hints"},
		AcceptanceCriteria: []string{"Function named hello", "Returns a string", "Has unit tests"},
	}
}

// InvalidTaskSpecEmptyGoal returns a task spec with empty goal
func InvalidTaskSpecEmptyGoal() schema.TaskSpec {
	return schema.TaskSpec{
		ID:                 "test-002",
		Goal:               "",
		AcceptanceCriteria: []string{"test"},
	}
}

// InvalidTaskSpecNoCriteria returns a task spec without acceptance criteria
func InvalidTaskSpecNoCriteria() schema.TaskSpec {
	return schema.TaskSpec{
		ID:     "test-003",
		Goal:   "Some task",
	}
}

// SampleExecutionResult returns a sample execution result
func SampleExecutionResult() schema.ExecutionResult {
	return schema.ExecutionResult{
		Backend:      "claude",
		Status:       "success",
		Summary:      "Successfully implemented hello function",
		FilesChanged: []string{"hello.py", "hello_test.py"},
		Stdout:       "Implementation complete",
		Stderr:       "",
	}
}

// SampleReviewFindings returns sample review findings
func SampleReviewFindings() []schema.ReviewFinding {
	return []schema.ReviewFinding{
		{
			Severity: "high",
			Type:     "missing_test",
			File:     "hello.py",
			Message:  "No test coverage for edge cases",
		},
		{
			Severity: "medium",
			Type:     "convention_violation",
			File:     "hello.py",
			Message:  "Line too long (100 characters)",
		},
	}
}

// NoBlockerFindings returns review findings with no high severity blockers
func NoBlockerFindings() []schema.ReviewFinding {
	return []schema.ReviewFinding{
		{
			Severity: "low",
			Type:     "style",
			File:     "hello.py",
			Message:  "Consider using f-string",
		},
	}
}

// SampleValidationResult returns a passing validation result
func SampleValidationResult() schema.ValidationResult {
	return schema.ValidationResult{
		Pass:         true,
		Errors:       []string{},
		Warnings:     []string{},
		FilesChecked: []string{"hello.py", "hello_test.py"},
	}
}

// FailingValidationResult returns a failing validation result
func FailingValidationResult() schema.ValidationResult {
	return schema.ValidationResult{
		Pass:    false,
		Errors:  []string{"lint error: undefined variable 'x'"},
		Warnings: []string{"unused import 'os'"},
		FilesChecked: []string{"hello.py"},
	}
}
