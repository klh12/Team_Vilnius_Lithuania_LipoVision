package dropletgenomics

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"net/http"
	"strconv"
	"time"
)

const streamEndpoint string = "http://example.com/"
const frameRate int64 = 30

type Device struct {
	IPAddress         string
	HTTPPort          int
	PumpDataPort      int
	RecordingDataPort int
	PumpExperiment    int
	pumps             []pump
	camera            camera
}

func (d Device) Stream(ctx context.Context) <-chan Frame {
	stream := make(chan Frame, 10)
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			default:
				response, err := http.Get(streamEndpoint)
				if err != nil {
					fmt.Printf("Failed to connect to stream")
					continue
				}
				byteStream := response.Body
				var decodeError error = nil
				for decodeError == nil {
					var img image.Image
					img, decodeError = png.Decode(byteStream)
					frameLifetime, cancel := context.WithTimeout(ctx, time.Second/(time.Duration)(frameRate))
					stream <- Frame{frame: img, ctx: frameLifetime, cancel: cancel}
				}
			}
		}
	}()
	return stream
}

// Available determines whether connection to
// DropletGenomics device is established
func (device *Device) Available() bool {
	url := setupDeviceURL(device)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return false
	}
	return true
}

//Camera returns the device's camera data
func (device Device) Camera() camera {
	return device.camera
}

//Pump returns device's pump by it's id
func (device Device) Pump(index int) pump {
	return device.pumps[index]
}

func (device *Device) DefinePumpExperiment(numberOfPumps int) {
	//TODO : get endpoints from device in GMC
}

func (device *Device) EstablishPumpsWithID() {
	device.pumps = make([]pump, device.PumpExperiment, device.PumpExperiment)
	for i := 0; i < device.PumpExperiment; i++ {
		device.pumps[i] = NewPump(i)
	}
}
