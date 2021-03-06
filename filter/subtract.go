package filter

import (
	"image"

	"gocv.io/x/gocv"
)

//CreateSubtractFilter Creates subtract filter with settings
func CreateSubtractFilter(region *image.Rectangle) Subtract {
	return Subtract{
		region:     region,
		subtractor: gocv.NewBackgroundSubtractorKNN(),
	}
}

//Subtract Is a background subtractor, used to find movement
type Subtract struct {
	region     *image.Rectangle
	subtractor gocv.BackgroundSubtractorKNN
}

//Produce Produces a new frame with background subtracted
func (s Subtract) Produce(frame gocv.Mat) (gocv.Mat, error) {
	resultingFrame := frame.Region(*s.region)
	err := s.Apply(&resultingFrame)
	return resultingFrame, err
}

//Apply Applies background subtraction to frame
func (s Subtract) Apply(frame *gocv.Mat) error {
	s.subtractor.Apply(*frame, frame)
	return nil
}
