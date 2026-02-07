package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/farukkandemir/freq/internal/jamendo"
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

func main() {

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

		tag := os.Args[2]

		jamClient := jamendo.NewJamendoClient()

		tracks, err := jamClient.GetTrack(tag)

		if err != nil {
			fmt.Println("Something went wrong", err)
			return
		}

		selection := rand.Intn(len(tracks.Results))
		track := tracks.Results[selection]

		resp, err := http.Get(track.Audio)
		if err != nil {
			fmt.Println("Error fetching audio stream:", err)
			return
		}

		streamer, format, err := mp3.Decode(resp.Body)

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

	default:
		fmt.Printf("Unknown command: %s\n", command)

	}

}
