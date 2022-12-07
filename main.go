package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/mplewis/figyr"
)

const desc = "EZVJ uses VLC to play random snippets of videos."

type Config struct {
	VLCHost     string `figyr:"default=localhost,description=The host running the VLC server"`
	VLCPort     int    `figyr:"default=8080,description=The port on which VLC's web interface is running"`
	VLCPassword string `figyr:"required,description=The password to the VLC server"`
	VideoDir    string `figyr:"required,description=The directory containing the videos to play"`
	// default exclude end period = 10%
	PlayDurationMin time.Duration `figyr:"default=1m,description=The minimum duration of a video to play"`
	PlayDurationMax time.Duration `figyr:"default=5m,description=The maximum duration of a video to play"`
	ExcludeStart    float64       `figyr:"default=0.1,description=Exclude this percentage of the video from playing at the start"`
	ExcludeEnd      float64       `figyr:"default=0.1,description=Exclude this percentage of the video from playing at the end"`
}

type PlaylistItem struct {
	Name     string
	ID       int
	Duration int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var cfg Config
	figyr.New(desc).MustParse(&cfg)
	log.Printf("VLC host: %s:%d\n", cfg.VLCHost, cfg.VLCPort)
	log.Printf("Video directory: %s\n", cfg.VideoDir)
	log.Printf("Play duration: %s to %s\n", cfg.PlayDurationMin, cfg.PlayDurationMax)
	log.Printf("Exclude: before %d%%, after %d%%\n", int(cfg.ExcludeStart*100), int((1-cfg.ExcludeEnd)*100))

	p := NewPlayer(cfg)
	check(p.VLC.EmptyPlaylist())
	files := listFiles(cfg.VideoDir)
	for _, f := range files {
		p.Add(f)
		log.Printf("Added to playlist: %s\n", f)
	}

	fs := false
	for {
		item := p.PlayRandomItem()
		dur := p.PickRandomPlayDuration()
		durT := fmt.Sprintf("%02d:%02d", int(dur.Minutes()), int(dur.Seconds())%60)
		start := p.SeekToRandomPosition(dur)
		startT := fmt.Sprintf("%02d:%02d", start/60, start%60)
		end := start + int(dur.Seconds())
		endT := fmt.Sprintf("%02d:%02d", end/60, end%60)
		log.Printf("Now playing: %s: %s-%s (%s)\n", item.Name, startT, endT, durT)

		if !fs {
			fs = true
			go func() {
				time.Sleep(5 * time.Second)     // HACK: we have to wait until the video is definitely playing to fs it
				check(p.VLC.ToggleFullscreen()) // HACK: right now we cannot check for fullscreen mode
			}()
		}
		time.Sleep(dur)
	}
}
