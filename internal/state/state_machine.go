package state

// State represents the state machine states
type State string

const (
	StateSpecPending  State = "SPEC_PENDING"
	StateReady        State = "READY"
	StateExecuting    State = "EXECUTING"
	StateReviewing    State = "REVIEWING"
	StateValidating   State = "VALIDATING"
	StateReworkRequired State = "REWORK_REQUIRED"
	StateAccepted     State = "ACCEPTED"
	StateFailed        State = "FAILED"
)

// ValidTransitions defines valid state transitions
var ValidTransitions = map[State][]State{
	StateSpecPending:  {StateReady, StateFailed},
	StateReady:       {StateExecuting, StateFailed},
	StateExecuting:   {StateReviewing, StateFailed},
	StateReviewing:   {StateValidating, StateReworkRequired, StateFailed},
	StateValidating:  {StateReworkRequired, StateAccepted, StateFailed},
	StateReworkRequired: {StateExecuting, StateFailed},
	StateAccepted:    {},
	StateFailed:      {},
}

// StateMachine manages state transitions
type StateMachine struct {
	CurrentState State
}

// NewStateMachine creates a new state machine in SPEC_PENDING state
func NewStateMachine() *StateMachine {
	return &StateMachine{CurrentState: StateSpecPending}
}

// CanTransition checks if a transition is valid
func (sm *StateMachine) CanTransition(to State) bool {
	allowed, ok := ValidTransitions[sm.CurrentState]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

// Transition moves to a new state
func (sm *StateMachine) Transition(to State) error {
	if !sm.CanTransition(to) {
		return &InvalidTransitionError{From: sm.CurrentState, To: to}
	}
	sm.CurrentState = to
	return nil
}

// InvalidTransitionError represents an invalid state transition
type InvalidTransitionError struct {
	From, To State
}

func (e *InvalidTransitionError) Error() string {
	return "invalid transition from " + string(e.From) + " to " + string(e.To)
}
