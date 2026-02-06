package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/farukkandemir/freq/internal/library"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

func main() {

	defaultMusicDir := "C:\\Users\\faruk\\Desktop\\Musics"

	if len(os.Args) < 2 {
		fmt.Println("FREQ >> Digital Audio Utility")
		fmt.Println("Usage: freq <command> [arguments]")
		return
	}

	command := os.Args[1]

	switch command {
	case "local":

		if len(os.Args) < 3 {
			fmt.Println("Please provide a song name")
			return
		}

		songName := os.Args[2]

		if filepath.Ext(songName) == "" {
			songName += ".mp3"
		}

		var fullPath string
		if filepath.IsAbs(songName) {
			fullPath = songName
		} else {
			fullPath = filepath.Join(defaultMusicDir, songName)
		}

		if _, err := os.Stat(fullPath); err != nil {
			fmt.Println("File not found:", fullPath)
			return
		}

		f, err := os.Open(fullPath)
		defer f.Close()

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		streamer, format, err := mp3.Decode(f)

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

		ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}
		speaker.Play(ctrl)

		for {
			fmt.Print("Press [ENTER] to pause/resume. ")
			fmt.Scanln()

			speaker.Lock()

			ctrl.Paused = !ctrl.Paused
			speaker.Unlock()
		}

	case "scan":

		if len(os.Args) < 3 {
			fmt.Println("Please provide a path")
			return
		}

		folderName := os.Args[2]

		files, err := library.Scan(folderName)

		if err != nil {
			fmt.Println(err)
		}

		for _, file := range files {

			fmt.Println(file)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)

	}

}
