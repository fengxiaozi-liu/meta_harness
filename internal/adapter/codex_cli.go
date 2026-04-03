package adapter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"meta_harness/internal/schema"
)

// CodexBackend implements ExecutorBackend for Codex CLI
type CodexBackend struct {
	CLIPath string
}

// NewCodexBackend creates a new Codex backend
func NewCodexBackend() (*CodexBackend, error) {
	path, err := exec.LookPath("codex")
	if err != nil {
		return nil, fmt.Errorf("未找到 codex CLI: %w", err)
	}
	return &CodexBackend{CLIPath: path}, nil
}

// Execute 启动交互式 Codex 会话
func (c *CodexBackend) Execute(ctx context.Context, spec schema.TaskSpec) (schema.ExecutionResult, error) {
	start := time.Now()
	cmd := exec.CommandContext(ctx, c.CLIPath, buildInteractivePrompt(spec))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	finished := time.Now()

	result := schema.ExecutionResult{
		Backend:    "codex",
		Status:     "success",
		StartedAt:  start,
		FinishedAt: finished,
		Summary:    "interactive codex session completed",
	}

	if err != nil {
		result.Status = "error"
		result.Summary = fmt.Sprintf("Codex CLI error: %v", err)
		return result, fmt.Errorf("codex execution failed: %w", err)
	}

	return result, nil
}

// Ensure CodexBackend implements ExecutorBackend
var _ ExecutorBackend = (*CodexBackend)(nil)
