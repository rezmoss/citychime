package main

import (
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/getlantern/systray"
)

var (
	audioBuffer *beep.Buffer
)

func main() {
	initAudio()
	systray.Run(onReady, onExit)
}

func initAudio() {
	// Get the path to the Resources directory in the app bundle
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	resourcesPath := filepath.Join(filepath.Dir(execPath), "..", "Resources")
	chimeFile := filepath.Join(resourcesPath, "chime.mp3")

	// Open and decode the MP3 file
	f, err := os.Open(chimeFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	audioBuffer = beep.NewBuffer(format)
	audioBuffer.Append(streamer)
	
	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		log.Fatal(err)
	}
}

func onReady() {
	systray.SetIcon(getIcon())
	systray.SetTitle("City Hall Clock")
	systray.SetTooltip("City Hall Clock")

	mQuit := systray.AddMenuItem("Quit", "Quit the app")

	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()

	go runClock()
}

func onExit() {
	// Cleanup tasks go here
}

func runClock() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			if t.Minute() == 0 || t.Minute() == 15 || t.Minute() == 30 || t.Minute() == 45 {
				go chime(t)
			}
		}
	}
}

func chime(t time.Time) {
	hour := t.Hour()
	minute := t.Minute()

	var chimeTimes int
	if minute == 0 {
		chimeTimes = hour % 12
		if chimeTimes == 0 {
			chimeTimes = 12
		}
	} else {
		chimeTimes = 1
	}

	for i := 0; i < chimeTimes; i++ {
		streamer := audioBuffer.Streamer(0, audioBuffer.Len())
		speaker.Play(streamer)
		time.Sleep(time.Duration(audioBuffer.Len()) * time.Second / time.Duration(audioBuffer.Format().SampleRate))
		
		if i < chimeTimes-1 {
			time.Sleep(500 * time.Millisecond)  // Wait between chimes
		}
	}
}

func getIcon() []byte {
    execPath, err := os.Executable()
    if err != nil {
        log.Fatal(err)
    }
    iconPath := filepath.Join(filepath.Dir(execPath), "..", "Resources", "icon.png")
    
    // Read the icon file
    icon, err := os.ReadFile(iconPath)
    if err != nil {
        log.Fatal(err)
    }
    
    return icon
}

func drawLine(img *image.RGBA, start, end image.Point, col color.Color) {
	dx := abs(end.X - start.X)
	dy := abs(end.Y - start.Y)
	sx, sy := 1, 1
	if start.X >= end.X {
		sx = -1
	}
	if start.Y >= end.Y {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(start.X, start.Y, col)
		if start.X == end.X && start.Y == end.Y {
			return
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			start.X += sx
		}
		if e2 < dx {
			err += dx
			start.Y += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}