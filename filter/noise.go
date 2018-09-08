package filter

import "gocv.io/x/gocv"

//CreateNoiseFilter Creates noise filter
func CreateNoiseFilter(referenceFrame *gocv.Mat, frameNumber *int) Noise {
	return Noise{frameNumber: frameNumber, referenceFrame: referenceFrame}
}

//NoiseFilter Removes moving noise pixels,
//Should ease detection
type Noise struct {
	frameNumber    *int
	referenceFrame *gocv.Mat
}

func (nf Noise) Apply(frame *gocv.Mat) error {
	for i := 0; i < frame.Rows(); i++ {
		for j := 0; j < frame.Cols(); j++ {
			if *nf.frameNumber > 0 && nf.referenceFrame.GetUCharAt(i, j) != frame.GetUCharAt(i, j) {
				frame.SetUCharAt(i, j, 0)
			}
		}
	}
	return nil
}

//Produce Comparison with previous frame for moving pixels detection and removal
func (nf Noise) Produce(frame gocv.Mat) (gocv.Mat, error) {
	resultingFrame := frame.Clone()
	for i := 0; i < frame.Rows(); i++ {
		for j := 0; j < frame.Cols(); j++ {
			if *nf.frameNumber > 0 && nf.referenceFrame.GetUCharAt(i, j) != frame.GetUCharAt(i, j) {
				resultingFrame.SetUCharAt(i, j, 0)
			}
		}
	}
	return resultingFrame, nil
}
