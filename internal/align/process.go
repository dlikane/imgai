package align

import (
	"fmt"
	"imgai/internal/align/alignutil"
	"imgai/pkg/common"
	"sort"

	"github.com/schollz/progressbar/v3"
)

func Run(cfg alignutil.Config) error {
	log := common.GetLogger()

	log.Info("🔍 Detecting base image landmarks...")
	baseEyes, err := alignutil.DetectLandmarks(cfg.BaseImagePath, cfg.ModelPath, cfg.ScriptPath)
	if err != nil {
		return fmt.Errorf("failed to detect base image landmarks: %w", err)
	}

	log.Info("🎞️  Opening video writer...")
	writer, err := alignutil.NewVideoWriter(cfg.VideoPath, cfg.FrameSize, cfg.FPS)
	if err != nil {
		return fmt.Errorf("failed to create video writer: %w", err)
	}
	defer writer.Close()

	log.Info("📂 Reading input directory...")
	imageFiles, err := alignutil.GetImageFiles(cfg.InputDir)
	if err != nil {
		return err
	}

	sort.Strings(imageFiles)
	log.Infof("📸 Found %d image(s) to process", len(imageFiles))

	bar := progressbar.NewOptions(len(imageFiles),
		progressbar.OptionSetDescription("📷 Aligning"),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(20),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionShowIts(),
		progressbar.OptionThrottle(65),
	)

	var skipped int
	var written int

	for _, name := range imageFiles {
		bar.Add(1)

		if alignutil.ProcessImage(name, cfg, baseEyes, writer) {
			written++
		} else {
			skipped++
		}
	}

	log.Infof("\n✅ Video saved to: %s", cfg.VideoPath)
	log.Infof("📊 Processing summary:")
	log.Infof("   Total files:       %d", len(imageFiles))
	log.Infof("   Skipped:           %d", skipped)
	log.Infof("   Written to video:  %d", written)

	return nil
}
