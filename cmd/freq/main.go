package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/farukkandemir/freq/internal/library"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

type JamendoTrack struct {
	Name       string `json:"name"`
	ArtistName string `json:"artist_name"`
	Audio      string `json:"audio"`
}

type JamendoResponse struct {
	Results []JamendoTrack `json:"results"`
}

var baseUrl string = "https://api.jamendo.com/v3.0/tracks/?client_id=6e1ba05b&format=jsonpretty&limit=1"

func main() {

	defaultMusicDir := "C:\\Users\\faruk\\Desktop\\Musics"

	if len(os.Args) < 2 {
		fmt.Println("FREQ >> Digital Audio Utility")
		fmt.Println("Usage: freq <command> [arguments]")
		return
	}

	command := os.Args[1]

	switch command {
	case "jam":

		if len(os.Args) < 3 {
			fmt.Println("Please provide a command for music(e.g freq jam chill)")
			return
		}

		songType := os.Args[2]

		fullPath := baseUrl + "&" + songType

		resp, err := http.Get(fullPath)

		if err != nil {
			fmt.Println(err)
			return
		}

		body, err := io.ReadAll(resp.Body)

		if err != nil {
			fmt.Println(err)
			return
		}

		var tracks JamendoResponse

		err = json.Unmarshal(body, &tracks)

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, track := range tracks.Results {
			fmt.Println(track)
		}

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
