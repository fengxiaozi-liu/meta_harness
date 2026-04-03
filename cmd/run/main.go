package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"meta_harness/internal/adapter"
	"meta_harness/internal/controller"
	"meta_harness/internal/schema"
)

func main() {
	// 解析命令行参数
	backendType := flag.String("backend", "codex", "Backend type: codex or claude")
	goal := flag.String("goal", "我是一个策划师", "The task goal")
	maxIterations := flag.Int("max-iterations", 5, "Maximum number of iterations")
	flag.Parse()

	if *goal == "" {
		log.Fatal("--goal is required")
	}

	// 构造任务规格
	spec := schema.TaskSpec{
		ID:                 generateTaskID(),
		Goal:               *goal,
		Constraints:        []string{"follow repo conventions"},
		AcceptanceCriteria: []string{"tests pass", "files conform to schema"},
	}

	// 创建后端
	backend, err := adapter.NewBackend(adapter.BackendType(*backendType))
	if err != nil {
		log.Fatalf("Failed to create backend: %v", err)
	}

	// 创建控制器
	ctrl := controller.NewController(controller.ControllerConfig{
		Backend:       backend,
		MaxIterations: *maxIterations,
	})

	// 运行流程
	ctx := context.Background()
	result, err := ctrl.Run(ctx, spec)
	if err != nil {
		log.Fatalf("Execution failed: %v", err)
	}

	// 输出结果
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		log.Fatalf("Failed to encode result: %v", err)
	}

	// 打印最终状态
	fmt.Printf("\nFinal state: %s\n", result.Status)
	fmt.Printf("Iterations: %d\n", result.Iteration)
}

func generateTaskID() string {
	// 生成简单任务 ID
	return fmt.Sprintf("task-%d", os.Getpid())
}
