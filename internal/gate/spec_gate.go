package gate

import (
	"meta_harness/internal/schema"
)

// SpecGate checks if a task spec is executable
type SpecGate interface {
	Check(spec schema.TaskSpec) schema.SpecGateResult
}

// DefaultSpecGate is a basic implementation of SpecGate
type DefaultSpecGate struct {
	MinGoalLength    int
	RequireCriteria bool
}

// NewDefaultSpecGate creates a new DefaultSpecGate
func NewDefaultSpecGate() *DefaultSpecGate {
	return &DefaultSpecGate{
		MinGoalLength:    10,
		RequireCriteria:  true,
	}
}

// Check validates the task spec
func (g *DefaultSpecGate) Check(spec schema.TaskSpec) schema.SpecGateResult {
	result := schema.SpecGateResult{Pass: true}

	if spec.ID == "" {
		result.Pass = false
		result.Issues = append(result.Issues, "task ID is required")
	}

	if len(spec.Goal) < g.MinGoalLength {
		result.Pass = false
		result.Issues = append(result.Issues, "goal is too short")
	}

	if g.RequireCriteria && len(spec.AcceptanceCriteria) == 0 {
		result.Pass = false
		result.Issues = append(result.Issues, "acceptance criteria is required")
	}

	if len(result.Issues) > 0 {
		result.Reason = "spec validation failed"
	}

	return result
}

// Ensure DefaultSpecGate implements SpecGate
var _ SpecGate = (*DefaultSpecGate)(nil)
