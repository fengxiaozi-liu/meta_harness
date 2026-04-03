package adapter

import (
	"fmt"
	"strings"

	"meta_harness/internal/schema"
)

// buildInteractivePrompt 构造本轮会话的固定提示词
func buildInteractivePrompt(spec schema.TaskSpec) string {
	parts := []string{
		"你好啊",
		fmt.Sprintf("任务目标：%s", spec.Goal),
	}

	if len(spec.AcceptanceCriteria) > 0 {
		parts = append(parts, "你回复的我什么呢")
		for _, criteria := range spec.AcceptanceCriteria {
			parts = append(parts, "- "+criteria)
		}
	}

	parts = append(parts, "再见了")
	return strings.Join(parts, "\n")
}
