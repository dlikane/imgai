package align

import (
	"fmt"
	"image"
	"imgai/internal/align/alignutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"imgai/pkg/common"
)

var baseImage string
var inputDir string
var outputDir string
var modelPath string
var scriptPath string
var videoPath string
var frameWidth, frameHeight int
var fps int

func init() {
	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Align images from input directory and produce output video",
		RunE:  runAlign,
	}

	runCmd.Flags().StringVar(&baseImage, "base", "", "Base image filename (must be inside input directory) [required]")
	runCmd.Flags().StringVar(&inputDir, "input", "data/input", "Directory of source images")
	runCmd.Flags().StringVar(&outputDir, "output", "data/output", "Directory for aligned output and video")
	runCmd.Flags().StringVar(&modelPath, "model", "models/shape_predictor_68_face_landmarks.dat", "Path to Dlib model file")
	runCmd.Flags().StringVar(&scriptPath, "script", "scripts/landmark_extractor.py", "Path to Python landmark script")
	runCmd.Flags().StringVar(&videoPath, "video", "output.mp4", "Output video file path (inside outputDir)")
	runCmd.Flags().IntVar(&frameWidth, "width", 1000, "Output frame width")
	runCmd.Flags().IntVar(&frameHeight, "height", 1500, "Output frame height")
	runCmd.Flags().IntVar(&fps, "fps", 50, "Frames per second in output video")

	rootCmd.AddCommand(runCmd)
}

func runAlign(cmd *cobra.Command, args []string) error {
	log := common.GetLogger()

	log.Infof("Base image: %s", baseImage)
	log.Infof("Input dir: %s", inputDir)
	log.Infof("Output dir: %s", outputDir)
	log.Infof("Model: %s", modelPath)
	log.Infof("Script: %s", scriptPath)

	basePath := filepath.Join(inputDir, baseImage)
	videoOutputPath := filepath.Join(outputDir, videoPath)

	cfg := alignutil.Config{
		BaseImagePath: basePath,
		InputDir:      inputDir,
		OutputDir:     outputDir,
		ModelPath:     modelPath,
		ScriptPath:    scriptPath,
		VideoPath:     videoOutputPath,
		FrameSize:     image.Pt(frameWidth, frameHeight),
		FPS:           fps,
	}

	if err := cfg.Validate(); err != nil {
		log.Errorf("Validation error: %v", err)
		return err
	}
	if err := cfg.EnsureOutputDir(); err != nil {
		log.Errorf("Failed to create output dir: %v", err)
		return fmt.Errorf("failed to create output dir: %v", err)
	}

	return Run(cfg)
}
