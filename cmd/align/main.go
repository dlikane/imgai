package main

import (
	"log"

	"imgai/internal/align"
	"imgai/pkg/common"
)

func main() {
	logger := common.GetLogger()
	logger.SetLevel(common.GetLogger().Level) // ensure it's at info level (already default, but explicit)
	logger.Info("Starting align tool")

	if err := align.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}
