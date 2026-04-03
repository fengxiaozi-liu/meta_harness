package adapter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"meta_harness/internal/schema"
)

// ClaudeBackend implements ExecutorBackend for Claude CLI
type ClaudeBackend struct {
	CLIPath string
}

// NewClaudeBackend creates a new Claude backend
func NewClaudeBackend() (*ClaudeBackend, error) {
	path, err := exec.LookPath("claude")
	if err != nil {
		return nil, fmt.Errorf("未找到 claude CLI: %w", err)
	}
	return &ClaudeBackend{CLIPath: path}, nil
}

// Execute 启动交互式 Claude 会话
func (c *ClaudeBackend) Execute(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
	start := time.Now()
	cmd := exec.CommandContext(ctx, c.CLIPath, buildInteractivePrompt(spec))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	finished := time.Now()

	result := schema.ExecutionResult{
		Backend:    "claude",
		Status:     "success",
		StartedAt:  start,
		FinishedAt: finished,
		Summary:    "interactive claude session completed",
	}

	if err != nil {
		result.Status = "error"
		result.Summary = fmt.Sprintf("Claude CLI error: %v", err)
		return result, fmt.Errorf("claude execution failed: %w", err)
	}

	return result, nil
}

// Ensure ClaudeBackend implements ExecutorBackend
var _ ExecutorBackend = (*ClaudeBackend)(nil)
