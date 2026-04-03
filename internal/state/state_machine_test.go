package state

import (
	"testing"
)

func TestNewStateMachine(t *testing.T) {
	sm := NewStateMachine()
	if sm.CurrentState != StateSpecPending {
		t.Errorf("expected initial state to be SPEC_PENDING, got %v", sm.CurrentState)
	}
}

func TestValidTransitions(t *testing.T) {
	sm := NewStateMachine()

	tests := []struct {
		name    string
		from    State
		to      State
		wantErr bool
	}{
		// Valid transitions from SPEC_PENDING
		{"SPEC_PENDING to READY", StateSpecPending, StateReady, false},
		{"SPEC_PENDING to FAILED", StateSpecPending, StateFailed, false},

		// Valid transitions from READY
		{"READY to EXECUTING", StateReady, StateExecuting, false},
		{"READY to FAILED", StateReady, StateFailed, false},

		// Valid transitions from EXECUTING
		{"EXECUTING to REVIEWING", StateExecuting, StateReviewing, false},
		{"EXECUTING to FAILED", StateExecuting, StateFailed, false},

		// Valid transitions from REVIEWING
		{"REVIEWING to VALIDATING", StateReviewing, StateValidating, false},
		{"REVIEWING to REWORK_REQUIRED", StateReviewing, StateReworkRequired, false},
		{"REVIEWING to FAILED", StateReviewing, StateFailed, false},

		// Valid transitions from VALIDATING
		{"VALIDATING to ACCEPTED", StateValidating, StateAccepted, false},
		{"VALIDATING to REWORK_REQUIRED", StateValidating, StateReworkRequired, false},
		{"VALIDATING to FAILED", StateValidating, StateFailed, false},

		// Valid transitions from REWORK_REQUIRED
		{"REWORK_REQUIRED to EXECUTING", StateReworkRequired, StateExecuting, false},
		{"REWORK_REQUIRED to FAILED", StateReworkRequired, StateFailed, false},

		// Invalid transitions from terminal states
		{"ACCEPTED to EXECUTING (invalid)", StateAccepted, StateExecuting, true},
		{"FAILED to READY (invalid)", StateFailed, StateReady, true},
		{"ACCEPTED to FAILED (invalid)", StateAccepted, StateFailed, true},

		// Invalid skip transitions
		{"SPEC_PENDING to EXECUTING (invalid skip)", StateSpecPending, StateExecuting, true},
		{"READY to REVIEWING (invalid skip)", StateReady, StateReviewing, true},
		{"EXECUTING to ACCEPTED (invalid skip)", StateExecuting, StateAccepted, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm.CurrentState = tt.from
			err := sm.Transition(tt.to)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transition() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && sm.CurrentState != tt.to {
				t.Errorf("Transition() state = %v, want %v", sm.CurrentState, tt.to)
			}
		})
	}
}

func TestCanTransition(t *testing.T) {
	sm := NewStateMachine()
	sm.CurrentState = StateReady

	tests := []struct {
		name string
		to   State
		want bool
	}{
		{"to EXECUTING", StateExecuting, true},
		{"to FAILED", StateFailed, true},
		{"to REVIEWING (invalid)", StateReviewing, false},
		{"to VALIDATING (invalid)", StateValidating, false},
		{"to ACCEPTED (invalid)", StateAccepted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sm.CanTransition(tt.to); got != tt.want {
				t.Errorf("CanTransition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidTransitionError(t *testing.T) {
	err := &InvalidTransitionError{From: StateReady, To: StateAccepted}
	expected := "invalid transition from READY to ACCEPTED"
	if err.Error() != expected {
		t.Errorf("InvalidTransitionError.Error() = %v, want %v", err.Error(), expected)
	}
}

func TestStateConstants(t *testing.T) {
	// Verify state constants are defined correctly
	states := []State{
		StateSpecPending,
		StateReady,
		StateExecuting,
		StateReviewing,
		StateValidating,
		StateReworkRequired,
		StateAccepted,
		StateFailed,
	}

	expected := []string{
		"SPEC_PENDING",
		"READY",
		"EXECUTING",
		"REVIEWING",
		"VALIDATING",
		"REWORK_REQUIRED",
		"ACCEPTED",
		"FAILED",
	}

	for i, state := range states {
		if string(state) != expected[i] {
			t.Errorf("State constant mismatch: got %v, want %v", state, expected[i])
		}
	}
}
