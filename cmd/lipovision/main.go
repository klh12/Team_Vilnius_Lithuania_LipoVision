package main

import (
	"context"
	"image"

	"github.com/Vilnius-Lithuania-iGEM-2018/lipovision/device"
	"github.com/Vilnius-Lithuania-iGEM-2018/lipovision/device/dropletgenomics"
	"github.com/Vilnius-Lithuania-iGEM-2018/lipovision/device/video"
	"github.com/Vilnius-Lithuania-iGEM-2018/lipovision/gui"
	"github.com/gotk3/gotk3/gtk"
	log "github.com/sirupsen/logrus"
)

var (
	mainCtx      context.Context
	mainCancel   context.CancelFunc
	activeDevice device.Device
)

func chooseFileCreateDevice(win *gtk.Window) device.Device {
	chooser, err := gtk.FileChooserDialogNewWith1Button(
		"Select video file", win, gtk.FILE_CHOOSER_ACTION_OPEN,
		"Open", gtk.RESPONSE_ACCEPT)
	if err != nil {
		log.Fatal("File chooser failed: ", err)
	}
	defer chooser.Destroy()

	filter, _ := gtk.FileFilterNew()
	filter.AddPattern("*.mp4")
	filter.SetName(".mp4")
	chooser.AddFilter(filter)

	resp := chooser.Run()
	log.Info(resp)

	videoFile := chooser.GetFilename()
	return video.Create(videoFile, 24)
}

func main() {
	mainCtx, mainCancel = context.WithCancel(context.Background())
	defer mainCancel()

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window")
	}
	win.SetTitle("LipoVision")
	win.Connect("destroy", func() {
		mainCancel()
		gtk.MainQuit()
	})
	win.SetDefaultSize(890, 500)

	content, err := gui.NewMainControl()
	if err != nil {
		panic(err)
	}
	win.Add(content.Root())
	win.ShowAll()

	registerDeviceChange(content, win)

	gtk.Main()
}

func registerDeviceChange(content *gui.MainControl, win *gtk.Window) {
	content.StreamControl.ComboBox.Connect("changed", func(combo *gtk.ComboBoxText) {
		mainCancel()
		mainCtx, mainCancel = context.WithCancel(context.Background())
		selection := combo.GetActiveText()
		switch selection {
		case "Video file...":
			activeDevice = chooseFileCreateDevice(win)
		case "DropletGenomics":
			activeDevice = dropletgenomics.Create(4)
		default:
			errDialog := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL,
				gtk.MESSAGE_ERROR, gtk.BUTTONS_OK,
				"Chosen device %s, does not exist", selection)
			errDialog.Run()
		}

		imageStream := make(chan image.Image, 10)
		go content.StreamControl.ShowStream(imageStream)

		go func() {
			log.WithFields(log.Fields{
				"device": selection,
			}).Info("Stream processor started")
			streamCtx, streamCancel := context.WithCancel(mainCtx)
			deviceStream := activeDevice.Stream(streamCtx)
			defer streamCancel()
			defer close(imageStream)
		Process:
			for {
				select {
				case <-streamCtx.Done():
					break Process
				case frame, ok := <-deviceStream:
					if ok {
						imageStream <- frame.Frame()
					} else {
						break Process
					}
				}
			}
			log.WithFields(log.Fields{
				"device": selection,
			}).Info("Stream processor exited")
		}()
	})
}
