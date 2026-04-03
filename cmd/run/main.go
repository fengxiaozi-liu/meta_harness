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
	// Parse flags
	backendType := flag.String("backend", "claude", "Backend type: codex or claude")
	cliPath := flag.String("cli-path", "", "Path to the CLI executable")
	goal := flag.String("goal", "", "The task goal")
	maxIterations := flag.Int("max-iterations", 5, "Maximum number of iterations")
	flag.Parse()

	// Validate required flags
	if *cliPath == "" {
		log.Fatal("--cli-path is required")
	}
	if *goal == "" {
		log.Fatal("--goal is required")
	}

	// Create task spec
	spec := schema.TaskSpec{
		ID:                 generateTaskID(),
		Goal:               *goal,
		Constraints:        []string{"follow repo conventions"},
		AcceptanceCriteria: []string{"tests pass", "files conform to schema"},
	}

	// Create backend
	backend, err := adapter.NewBackend(adapter.BackendType(*backendType), *cliPath)
	if err != nil {
		log.Fatalf("Failed to create backend: %v", err)
	}

	// Create controller
	ctrl := controller.NewController(controller.ControllerConfig{
		Backend:       backend,
		MaxIterations: *maxIterations,
	})

	// Run
	ctx := context.Background()
	result, err := ctrl.Run(ctx, spec)
	if err != nil {
		log.Fatalf("Execution failed: %v", err)
	}

	// Output result
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		log.Fatalf("Failed to encode result: %v", err)
	}

	// Print final state
	fmt.Printf("\nFinal state: %s\n", result.Status)
	fmt.Printf("Iterations: %d\n", result.Iteration)
}

func generateTaskID() string {
	// Simple task ID generation
	return fmt.Sprintf("task-%d", os.Getpid())
}
