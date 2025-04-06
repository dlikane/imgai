package alignutil

import (
	"image"
	"math"

	"imgai/pkg/common"
)

// CalculateTransform calculates rotation, scaling, and translation needed to align two images.
func CalculateTransform(baseEyes, targetEyes [2]image.Point) (angle, scale, dx, dy float64) {
	log := common.GetLogger()

	// Calculate the distance between eyes for base and target images
	baseDist := math.Hypot(float64(baseEyes[1].X-baseEyes[0].X), float64(baseEyes[1].Y-baseEyes[0].Y))
	targetDist := math.Hypot(float64(targetEyes[1].X-targetEyes[0].X), float64(targetEyes[1].Y-targetEyes[0].Y))

	// Calculate the angle between the eyes for base and target images
	baseAngle := math.Atan2(float64(baseEyes[1].Y-baseEyes[0].Y), float64(baseEyes[1].X-baseEyes[0].X))
	targetAngle := math.Atan2(float64(targetEyes[1].Y-targetEyes[0].Y), float64(targetEyes[1].X-targetEyes[0].X))
	angle = -(baseAngle - targetAngle)

	// Calculate the scaling factor
	scale = baseDist / targetDist

	// Translation based on eye position difference
	dx = float64(baseEyes[0].X - targetEyes[0].X)
	dy = float64(baseEyes[0].Y - targetEyes[0].Y)

	// Only these logs interfere with progress bar, so they are set to debug
	log.Debugf("Scale : base=%.2f, target=%.2f, result=%.4f", baseDist, targetDist, scale)
	log.Debugf("MoveX : base=%d, target=%d, dx=%.2f", baseEyes[0].X, targetEyes[0].X, dx)
	log.Debugf("MoveY : base=%d, target=%d, dy=%.2f", baseEyes[0].Y, targetEyes[0].Y, dy)

	return angle, scale, dx, dy
}
