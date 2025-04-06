package alignutil

import (
	"fmt"
	"image"
	"os"

	"imgai/pkg/common"
)

type Config struct {
	BaseImagePath string
	InputDir      string
	OutputDir     string
	ModelPath     string
	ScriptPath    string
	VideoPath     string
	FrameSize     image.Point
	FPS           int
}

func (c Config) Validate() error {
	log := common.GetLogger()

	if _, err := os.Stat(c.BaseImagePath); err != nil {
		log.Errorf("Base image not found: %s", c.BaseImagePath)
		return fmt.Errorf("base image not found: %s", c.BaseImagePath)
	}
	if _, err := os.Stat(c.InputDir); err != nil {
		log.Errorf("Input directory not found: %s", c.InputDir)
		return fmt.Errorf("input directory not found: %s", c.InputDir)
	}
	if _, err := os.Stat(c.ModelPath); err != nil {
		log.Errorf("Model file not found: %s", c.ModelPath)
		return fmt.Errorf("model file not found: %s", c.ModelPath)
	}
	if _, err := os.Stat(c.ScriptPath); err != nil {
		log.Errorf("Landmark script not found: %s", c.ScriptPath)
		return fmt.Errorf("landmark script not found: %s", c.ScriptPath)
	}
	return nil
}

func (c Config) EnsureOutputDir() error {
	log := common.GetLogger()
	if err := os.MkdirAll(c.OutputDir, 0755); err != nil {
		log.Errorf("Failed to create output directory: %s", c.OutputDir)
		return err
	}
	log.Infof("Output directory ensured: %s", c.OutputDir)
	return nil
}
