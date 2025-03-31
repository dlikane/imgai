package imgalign

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
)

func findPythonScript(scriptName string, searchPaths []string) (string, error) {
	for _, path := range searchPaths {
		candidate := filepath.Join(path, scriptName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("Python script %s not found in search paths: %v", scriptName, searchPaths)
}

// DetectLandmarks detects facial landmarks using Dlib.
func DetectLandmarks(isLocal bool, imagePath string, modelPath string) ([2]image.Point, error) {
	// Run the Python script (copy it to /usr/local/bin)
	scriptPaths := []string{"", ".", "./scripts"}
	py, err := findPythonScript("landmark_extractor.py", scriptPaths)
	if err != nil {
		return [2]image.Point{}, err
	}
	cmd := exec.Command(py, imagePath, modelPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		fmt.Printf("Output: %s\n", string(output))
		return [2]image.Point{}, fmt.Errorf("failed to execute landmark extractor: %v for: %s", err, imagePath)
	}
	fmt.Printf("here is output: %s\n", string(output))
	// Parse the JSON output
	var result struct {
		Error     string   `json:"error"`
		Landmarks [][2]int `json:"landmarks"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return [2]image.Point{}, fmt.Errorf("failed to parse output: %v", err)
	}

	// Handle errors from the Python script
	if result.Error != "" {
		return [2]image.Point{}, fmt.Errorf("landmark extraction error: %s", result.Error)
	}

	// Map left and right eye positions
	leftEye := image.Point{X: result.Landmarks[36][0], Y: result.Landmarks[36][1]}
	rightEye := image.Point{X: result.Landmarks[45][0], Y: result.Landmarks[45][1]}
	fmt.Printf("Left eye: %v, Right eye: %v\n", leftEye, rightEye)

	return [2]image.Point{leftEye, rightEye}, nil
}
