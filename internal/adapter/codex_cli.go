package adapter

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"meta_harness/internal/schema"
)

// CodexBackend implements ExecutorBackend for Codex CLI
type CodexBackend struct {
	CLIPath string
}

// NewCodexBackend creates a new Codex backend
func NewCodexBackend(cliPath string) *CodexBackend {
	return &CodexBackend{CLIPath: cliPath}
}

// Execute runs the Codex CLI with the given task spec
func (c *CodexBackend) Execute(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
	start := time.Now()

	// Build the command - Codex CLI interface
	cmd := exec.CommandContext(ctx, c.CLIPath, "--prompt", spec.Goal)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	finished := time.Now()

	result := schema.ExecutionResult{
		Backend:    "codex",
		Status:     "success",
		StartedAt:  start,
		FinishedAt: finished,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
	}

	if err != nil {
		result.Status = "error"
		result.Summary = fmt.Sprintf("Codex CLI error: %v", err)
		return result, fmt.Errorf("codex execution failed: %w", err)
	}

	result.Summary = parseCodexSummary(stdout.String())
	return result, nil
}

// parseCodexSummary extracts a summary from Codex output
func parseCodexSummary(output string) string {
	// TODO: Implement proper parsing based on actual Codex output format
	if len(output) > 200 {
		return output[:200] + "..."
	}
	return output
}

// Ensure CodexBackend implements ExecutorBackend
var _ ExecutorBackend = (*CodexBackend)(nil)
