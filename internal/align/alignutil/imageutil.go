package alignutil

import (
	"gocv.io/x/gocv"
	"image"
)

// ResizeAndCropFill resizes the image to fill the target height, then center-crops width
func ResizeAndCropFill(src gocv.Mat, targetSize image.Point) gocv.Mat {
	targetW := targetSize.X
	targetH := targetSize.Y

	aspectRatio := float64(src.Cols()) / float64(src.Rows())
	newW := int(float64(targetH) * aspectRatio)

	resized := gocv.NewMat()
	gocv.Resize(src, &resized, image.Point{X: newW, Y: targetH}, 0, 0, gocv.InterpolationLinear)

	// Center crop to target width
	xOffset := (resized.Cols() - targetW) / 2
	if xOffset < 0 {
		xOffset = 0
	}
	roi := image.Rect(xOffset, 0, xOffset+targetW, targetH)
	cropped := resized.Region(roi)
	resized.Close()
	return cropped
}
