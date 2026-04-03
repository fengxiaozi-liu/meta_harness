package controller

import (
	"context"
	"fmt"
	"time"

	"meta_harness/internal/adapter"
	"meta_harness/internal/executor"
	"meta_harness/internal/gate"
	"meta_harness/internal/review"
	"meta_harness/internal/schema"
	"meta_harness/internal/state"
	"meta_harness/internal/validate"
)

// Controller orchestrates the entire workflow
type Controller struct {
	executor      *executor.Executor
	specGate      gate.SpecGate
	reviewer      review.Reviewer
	validator     validate.Validator
	stateMachine  *state.StateMachine
	maxIterations int
}

// ControllerConfig holds controller configuration
type ControllerConfig struct {
	Backend        adapter.ExecutorBackend
	SpecGate       gate.SpecGate
	Reviewer       review.Reviewer
	Validator      validate.Validator
	MaxIterations  int
}

// NewController creates a new Controller
func NewController(cfg ControllerConfig) *Controller {
	if cfg.MaxIterations <= 0 {
		cfg.MaxIterations = 5
	}
	if cfg.SpecGate == nil {
		cfg.SpecGate = gate.NewDefaultSpecGate()
	}
	if cfg.Reviewer == nil {
		cfg.Reviewer = review.NewDefaultReviewer()
	}
	if cfg.Validator == nil {
		cfg.Validator = validate.NewCompositeValidator(
			validate.NewFileValidator(nil),
			validate.NewTestRunner(),
		)
	}

	return &Controller{
		executor:      executor.NewExecutor(cfg.Backend),
		specGate:      cfg.SpecGate,
		reviewer:      cfg.Reviewer,
		validator:     cfg.Validator,
		stateMachine:  state.NewStateMachine(),
		maxIterations: cfg.MaxIterations,
	}
}

// Run executes the complete workflow
func (c *Controller) Run(ctx context.Context, spec schema.TaskSpec) (*schema.ExecutionState, error) {
	execState := &schema.ExecutionState{
		TaskSpec:  spec,
		Iteration: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Step 1: Spec Gate
	c.stateMachine.Transition(state.StateReady)
	gateResult := c.specGate.Check(spec)
	execState.StopCriteria.SpecGatePass = gateResult.Pass

	if !gateResult.Pass {
		c.stateMachine.Transition(state.StateFailed)
		execState.Status = string(state.StateFailed)
		return execState, fmt.Errorf("spec gate failed: %v", gateResult.Issues)
	}

	// Main loop
	for execState.Iteration < c.maxIterations {
		execState.Iteration++
		execState.UpdatedAt = time.Now()

		// Step 2: Execute
		c.stateMachine.Transition(state.StateExecuting)
		result, err := c.executor.Execute(ctx, spec)
		if err != nil {
			c.stateMachine.Transition(state.StateFailed)
			execState.Status = string(state.StateFailed)
			return execState, fmt.Errorf("execution failed: %w", err)
		}
		execState.Results = append(execState.Results, result)

		// Step 3: Review
		c.stateMachine.Transition(state.StateReviewing)
		findings, err := c.reviewer.Review(spec, result)
		if err != nil {
			c.stateMachine.Transition(state.StateFailed)
			execState.Status = string(state.StateFailed)
			return execState, fmt.Errorf("review failed: %w", err)
		}
		execState.Findings = append(execState.Findings, findings...)

		// Count high severity blockers
		blockers := countBlockers(findings)
		execState.StopCriteria.ReviewBlockers = blockers

		if blockers > 0 {
			c.stateMachine.Transition(state.StateReworkRequired)
			execState.Status = string(state.StateReworkRequired)
			// Continue to next iteration for rework
			continue
		}

		// Step 4: Validate
		c.stateMachine.Transition(state.StateValidating)
		validation := c.validator.Validate(result.FilesChanged)
		execState.Validation = validation
		execState.StopCriteria.ValidationPass = validation.Pass

		if !validation.Pass {
			c.stateMachine.Transition(state.StateReworkRequired)
			execState.Status = string(state.StateReworkRequired)
			continue
		}

		// Success
		execState.StopCriteria.TestsPass = true
		c.stateMachine.Transition(state.StateAccepted)
		execState.Status = string(state.StateAccepted)
		return execState, nil
	}

	// Max iterations reached
	execState.StopCriteria.IterationCount = execState.Iteration
	c.stateMachine.Transition(state.StateFailed)
	execState.Status = string(state.StateFailed)
	return execState, fmt.Errorf("max iterations (%d) reached", c.maxIterations)
}

// countBlockers counts high severity findings
func countBlockers(findings []schema.ReviewFinding) int {
	count := 0
	for _, f := range findings {
		if f.Severity == "high" {
			count++
		}
	}
	return count
}

// GetState returns the current state
func (c *Controller) GetState() state.State {
	return c.stateMachine.CurrentState
}
