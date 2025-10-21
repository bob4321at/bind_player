package main

import (
	"os"
	"strings"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
)

const targetSampleRate = beep.SampleRate(44100)

func main() {
	speaker.Init(targetSampleRate, targetSampleRate.N(time.Second/30))
	for true {
		home_path, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}
		song_path_bytes, err := os.ReadFile(home_path + "/Documents/current_song")
		song_path := ""
		if err != nil {
			f, newerr := os.Create(home_path + "/Documents/current_song")
			if newerr != nil {
				panic(newerr)
			}
			f.WriteString("non")
			var newnewerr error
			song_path_bytes, newnewerr = os.ReadFile(home_path + "/Documents/current_song")
			if newnewerr != nil {
				panic(newnewerr)
			}
			f.Close()
		} else {
			if len(song_path_bytes) != 0 {
				song_path = string(song_path_bytes)
				if i := strings.Index(song_path, "^"); i != 1 {
					song_path = song_path[:i]
				}
				song_path = home_path + "/Music/" + song_path
			}
		}
		if song_path == "" {
			time.Sleep(time.Second)
		} else {
			song_f, err := os.Open(song_path)
			if err != nil {
				panic(err)
			}

			streamer, format, err := mp3.Decode(song_f)
			if err != nil {
				panic(err)
			}
			defer streamer.Close()

			resampled := beep.Resample(4, format.SampleRate, targetSampleRate, streamer)

			done := make(chan bool)
			speaker.Play(beep.Seq(resampled, beep.Callback(func() {
				done <- true
			})))

			<-done
		}
	}
}
