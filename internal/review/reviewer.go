package review

import (
	"meta_harness/internal/schema"
)

// Reviewer generates structured review findings
type Reviewer interface {
	Review(spec schema.TaskSpec, result schema.ExecutionResult) ([]schema.ReviewFinding, error)
}

// DefaultReviewer is a basic implementation of Reviewer
type DefaultReviewer struct{}

// NewDefaultReviewer creates a new DefaultReviewer
func NewDefaultReviewer() *DefaultReviewer {
	return &DefaultReviewer{}
}

// Review analyzes the execution result and generates findings
func (r *DefaultReviewer) Review(spec schema.TaskSpec, result schema.ExecutionResult) ([]schema.ReviewFinding, error) {
	var findings []schema.ReviewFinding

	// Check if execution was successful
	if result.Status == "error" {
		findings = append(findings, schema.ReviewFinding{
			Severity: "high",
			Type:     "execution_error",
			File:     "",
			Message:  result.Summary,
		})
		return findings, nil
	}

	// Check for missing tests (example rule)
	hasTestFiles := false
	for _, f := range result.FilesChanged {
		if len(f) >= 5 && (f[len(f)-5:] == "_test" || f[len(f)-3:] == ".py" && contains(f, "test")) {
			hasTestFiles = true
			break
		}
	}

	if !hasTestFiles && len(spec.AcceptanceCriteria) > 0 {
		findings = append(findings, schema.ReviewFinding{
			Severity: "medium",
			Type:     "missing_test",
			File:     "",
			Message:  "no test files detected in changes",
		})
	}

	return findings, nil
}

// contains checks if s contains substr
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Ensure DefaultReviewer implements Reviewer
var _ Reviewer = (*DefaultReviewer)(nil)
