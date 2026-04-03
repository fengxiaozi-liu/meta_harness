package gate

import (
	"testing"

	"meta_harness/internal/schema"
)

func TestDefaultSpecGate_Check(t *testing.T) {
	gate := NewDefaultSpecGate()

	tests := []struct {
		name       string
		spec       schema.TaskSpec
		wantPass   bool
		wantIssues int
	}{
		{
			name: "valid spec with all fields",
			spec: schema.TaskSpec{
				ID:                 "task-001",
				Goal:               "Implement a hello world function",
				Constraints:        []string{"follow PEP 8"},
				AcceptanceCriteria: []string{"returns string", "has tests"},
			},
			wantPass:   true,
			wantIssues: 0,
		},
		{
			name: "missing ID",
			spec: schema.TaskSpec{
				ID:                 "",
				Goal:               "Implement a hello world function",
				AcceptanceCriteria: []string{"returns string"},
			},
			wantPass:   false,
			wantIssues: 1,
		},
		{
			name: "goal too short",
			spec: schema.TaskSpec{
				ID:                 "task-002",
				Goal:               "Hi",
				AcceptanceCriteria: []string{"returns string"},
			},
			wantPass:   false,
			wantIssues: 1,
		},
		{
			name: "missing acceptance criteria",
			spec: schema.TaskSpec{
				ID:                 "task-003",
				Goal:               "Implement a hello world function",
				Constraints:        []string{"follow PEP 8"},
				AcceptanceCriteria: []string{},
			},
			wantPass:   false,
			wantIssues: 1,
		},
		{
			name: "multiple issues",
			spec: schema.TaskSpec{
				ID:                 "",
				Goal:               "Hi",
				AcceptanceCriteria: []string{},
			},
			wantPass:   false,
			wantIssues: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gate.Check(tt.spec)
			if result.Pass != tt.wantPass {
				t.Errorf("Check() pass = %v, want %v", result.Pass, tt.wantPass)
			}
			if len(result.Issues) != tt.wantIssues {
				t.Errorf("Check() issues count = %d, want %d", len(result.Issues), tt.wantIssues)
			}
		})
	}
}

func TestDefaultSpecGate_Reason(t *testing.T) {
	gate := NewDefaultSpecGate()

	// Test that reason is set when pass is false
	spec := schema.TaskSpec{
		ID:                 "",
		Goal:               "Hi",
		AcceptanceCriteria: []string{},
	}

	result := gate.Check(spec)
	if result.Pass && result.Reason != "" {
		t.Errorf("Check() reason should be empty when pass is true")
	}
	if !result.Pass && result.Reason == "" {
		t.Errorf("Check() reason should not be empty when pass is false")
	}
}

func TestNewDefaultSpecGate(t *testing.T) {
	gate := NewDefaultSpecGate()
	if gate.MinGoalLength != 10 {
		t.Errorf("DefaultSpecGate MinGoalLength = %d, want 10", gate.MinGoalLength)
	}
	if !gate.RequireCriteria {
		t.Errorf("DefaultSpecGate RequireCriteria = %v, want true", gate.RequireCriteria)
	}
}
