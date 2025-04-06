package alignutil

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"imgai/pkg/common"
	"os"
	"path/filepath"
	"strings"
)

const maxScaleFactor = 2.0

// GetImageFiles returns a sorted list of .jpg image filenames in the directory
func GetImageFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read input dir: %w", err)
	}

	var imageFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".jpg") {
			imageFiles = append(imageFiles, entry.Name())
		}
	}
	return imageFiles, nil
}

// ProcessImage handles detection, transform, and writing of a single image.
// Returns true if written to video, false if skipped.
func ProcessImage(name string, cfg Config, baseEyes [2]image.Point, writer *gocv.VideoWriter) (written bool) {
	log := common.GetLogger()
	defer func() {
		if r := recover(); r != nil {
			log.Warnf("❌ Recovered from panic in %s: %v", name, r)
			written = false
		}
	}()

	imgPath := filepath.Join(cfg.InputDir, name)
	log.Debugf("Image: %s", imgPath)

	targetEyes, err := DetectLandmarks(imgPath, cfg.ModelPath, cfg.ScriptPath)
	if err != nil {
		log.Debugf("⚠️  Skipping %s: %v", name, err)
		return false
	}

	mat := gocv.IMRead(imgPath, gocv.IMReadColor)
	if mat.Empty() {
		log.Warnf("⚠️  Could not read %s", name)
		return false
	}
	defer mat.Close()

	angle, scale, dx, dy := CalculateTransform(baseEyes, targetEyes)
	if scale > maxScaleFactor {
		log.Warnf("⚠️  Skipping %s: scale %.2f too large", name, scale)
		return false
	}

	aligned := AlignImage(targetEyes, mat, angle, scale, dx, dy)
	defer aligned.Close()

	final := ResizeAndCropFill(aligned, cfg.FrameSize)
	defer final.Close()

	if err := writer.Write(final); err != nil {
		log.Warnf("⚠️  Failed to write %s to video: %v", name, err)
		return false
	}

	return true
}
