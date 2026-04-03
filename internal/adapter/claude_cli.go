package adapter

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"meta_harness/internal/schema"
)

// ClaudeBackend implements ExecutorBackend for Claude CLI
type ClaudeBackend struct {
	CLIPath string
}

// NewClaudeBackend creates a new Claude backend
func NewClaudeBackend(cliPath string) *ClaudeBackend {
	return &ClaudeBackend{CLIPath: cliPath}
}

// Execute runs the Claude CLI with the given task spec
func (c *ClaudeBackend) Execute(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
	start := time.Now()

	// Build the command - Claude CLI interface
	cmd := exec.CommandContext(ctx, c.CLIPath, "--prompt", spec.Goal)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	finished := time.Now()

	result := schema.ExecutionResult{
		Backend:    "claude",
		Status:     "success",
		StartedAt:  start,
		FinishedAt: finished,
		Stdout:     stdout.String(),
		Stderr:     stderr.String(),
	}

	if err != nil {
		result.Status = "error"
		result.Summary = fmt.Sprintf("Claude CLI error: %v", err)
		return result, fmt.Errorf("claude execution failed: %w", err)
	}

	result.Summary = parseClaudeSummary(stdout.String())
	return result, nil
}

// parseClaudeSummary extracts a summary from Claude output
func parseClaudeSummary(output string) string {
	// TODO: Implement proper parsing based on actual Claude output format
	if len(output) > 200 {
		return output[:200] + "..."
	}
	return output
}

// Ensure ClaudeBackend implements ExecutorBackend
var _ ExecutorBackend = (*ClaudeBackend)(nil)
