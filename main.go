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
		start := p.SeekToRandomPosition(item, dur)
		startT := fmt.Sprintf("%02d:%02d", start/60, start%60)
		end := start + int(dur.Seconds())
		endT := fmt.Sprintf("%02d:%02d", end/60, end%60)
		log.Printf("Now playing: %s: %s-%s (%s)\n", item.Name, startT, endT, durT)
		if !fs {
			check(p.VLC.ToggleFullscreen()) // HACK: right now we cannot check for fullscreen mode
			fs = true
		}
		time.Sleep(dur)
	}
}
