package main

import (
	"bufio"
	"fmt"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/farukkandemir/freq/internal/jamendo"
	"github.com/fatih/color"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

type App struct {
	view               string
	currentTag         string
	musicCtrl          *beep.Ctrl
	streamer           beep.StreamCloser
	speakerInitialized bool
}

var (
	Accent = color.New(color.FgHiCyan).Add(color.Bold)
	Text   = color.New(color.FgWhite)
	Dim    = color.New(color.FgHiBlack)
)

func printHeader() {
	fmt.Println()
	fmt.Print("  ")
	Accent.Print("f r e q")
	Dim.Println("  |  v1.0")
	fmt.Println("  " + Dim.Sprint(strings.Repeat("─", 32)))

	Dim.Println("  search for e.g. chill, lofi, ambient…")
	fmt.Println()
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func (app *App) stopMusic() {
	if app.musicCtrl != nil {
		speaker.Clear()
	}
	if app.streamer != nil {
		app.streamer.Close()
		app.streamer = nil
	}
	app.musicCtrl = nil
}

func (app *App) startMusic(tag string) error {

	app.stopMusic()

	jamClient := jamendo.NewJamendoClient()

	tracks, err := jamClient.GetTrack(tag)
	if err != nil {
		return fmt.Errorf("API error: %v", err)
	}
	if len(tracks.Results) == 0 {
		return fmt.Errorf("no tracks found for tag: %s", tag)
	}

	track := tracks.Results[rand.IntN(len(tracks.Results))]

	resp, err := http.Get(track.Audio)
	if err != nil {
		return fmt.Errorf("network error: %v", err)
	}

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		resp.Body.Close()
		return fmt.Errorf("audio decode error: %v", err)
	}

	if !app.speakerInitialized {
		err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		if err != nil {
			streamer.Close()
			return fmt.Errorf("speaker error: %v", err)
		}
		app.speakerInitialized = true
	}

	app.musicCtrl = &beep.Ctrl{Streamer: streamer, Paused: false}
	app.streamer = streamer
	speaker.Play(app.musicCtrl)

	return nil
}

func main() {
	app := App{
		view: "home",
	}

	reader := bufio.NewScanner(os.Stdin)

	for {

		clearScreen()

		switch app.view {
		case "home":
			printHeader()

			fmt.Print("search...>")

			reader.Scan()
			input := reader.Text()

			app.currentTag = input
			app.view = "loading"

		case "loading":
			printHeader()
			fmt.Printf("  %s Fetching some %s for you...\n", Accent.Sprint("∞"), app.currentTag)

			err := app.startMusic(app.currentTag)

			if err != nil {
				fmt.Printf("\n  %s Error: %v\n", color.RedString("!!"), err)
				fmt.Println("  Returning home in 3 seconds...")
				time.Sleep(3 * time.Second)
				app.view = "home"
			} else {
				app.view = "play"
			}

		case "play":
			printHeader()

			fmt.Print("  ")
			Accent.Print("ON AIR")
			Dim.Printf("  Playing: #%s\n", app.currentTag)
			fmt.Println("  " + Dim.Sprint(strings.Repeat("━", 32)))

			fmt.Println("\n  Controls:")
			fmt.Printf("  %s %s\n", Accent.Sprint("⌙"), "Type 'b' to Stop & Search")
			fmt.Printf("  %s %s\n", Accent.Sprint("⌙"), "Type 'p' to Pause/Resume")
			fmt.Printf("  %s %s\n", Accent.Sprint("⌙"), "Type 'q' to Exit")

			fmt.Println()

			fmt.Print("  ❯ ")

			reader.Scan()
			input := strings.ToLower(reader.Text())

			switch input {
			case "p":
				if app.musicCtrl != nil {
					app.musicCtrl.Paused = !app.musicCtrl.Paused
				}
			case "b":
				app.stopMusic()
				app.view = "home"

			case "q", "exit", "quit":
				app.stopMusic()
				Accent.Print("Shutting down..")
				time.Sleep(time.Second)
				os.Exit(0)
			default:
				fmt.Println("Invalid command, try again.")
			}
		}

	}

}
