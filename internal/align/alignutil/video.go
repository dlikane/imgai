package alignutil

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"

	"imgai/pkg/common"
)

// NewVideoWriter creates a new video writer for output.
func NewVideoWriter(path string, size image.Point, fps int) (*gocv.VideoWriter, error) {
	log := common.GetLogger()

	writer, err := gocv.VideoWriterFile(path, "avc1", float64(fps), size.X, size.Y, true)
	if err != nil {
		log.Errorf("Failed to open video writer: %v", err)
		return nil, fmt.Errorf("failed to open video writer: %w", err)
	}

	log.Infof("Created video writer: %s (%dx%d @ %dfps)", path, size.X, size.Y, fps)
	return writer, nil
}
