package alignutil

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"

	"imgai/pkg/common"
)

type landmarkResult struct {
	Error     string   `json:"error"`
	Landmarks [][2]int `json:"landmarks"`
}

// DetectLandmarks runs the Python script and extracts eye landmarks
func DetectLandmarks(imagePath, modelPath, scriptPath string) ([2]image.Point, error) {
	log := common.GetLogger()

	script, err := resolveScript(scriptPath)
	if err != nil {
		log.Errorf("failed to resolve script: %v", err)
		return [2]image.Point{}, err
	}

	cmd := exec.Command(script, imagePath, modelPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("script error output: %s", string(output))
		return [2]image.Point{}, fmt.Errorf("failed to execute landmark extractor: %v", err)
	}

	var result landmarkResult
	if err := json.Unmarshal(output, &result); err != nil {
		log.Errorf("failed to parse JSON output: %v", err)
		return [2]image.Point{}, fmt.Errorf("invalid JSON from Python script: %w", err)
	}
	if result.Error != "" {
		log.Debugf("landmark extraction error: %s", result.Error)
		return [2]image.Point{}, fmt.Errorf("landmark extraction error: %s", result.Error)
	}

	// Expect 68 points
	if len(result.Landmarks) < 46 {
		log.Warnf("landmark output too short: got %d points", len(result.Landmarks))
		return [2]image.Point{}, fmt.Errorf("landmark output too short: got %d points", len(result.Landmarks))
	}

	left := result.Landmarks[36]
	right := result.Landmarks[45]

	log.Debugf("Image: %s", imagePath)
	log.Debugf("Left eye: %v, Right eye: %v", left, right)

	return [2]image.Point{
		{X: left[0], Y: left[1]},
		{X: right[0], Y: right[1]},
	}, nil
}

// resolveScript ensures script exists and returns its full path
func resolveScript(path string) (string, error) {
	log := common.GetLogger()

	pathsToTry := []string{
		path,
		"./scripts/landmark_extractor.py",
		"landmark_extractor.py",
	}
	for _, p := range pathsToTry {
		if _, err := os.Stat(p); err == nil {
			fullPath, _ := filepath.Abs(p)
			log.Tracef("\nResolved script path: %s", fullPath)
			return fullPath, nil
		}
	}
	log.Errorf("landmark script not found in known locations")
	return "", fmt.Errorf("landmark script not found in known locations")
}
