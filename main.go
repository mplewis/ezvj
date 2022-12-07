package main

import (
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

	for {
		item := p.PlayRandomItem()
		dur := p.PickRandomPlayDuration()
		pos := p.SeekToRandomPosition(item, dur)
		log.Printf("Playing %s for %d seconds at pos %d\n", item.Name, int(dur.Seconds()), pos)
		time.Sleep(dur)
	}
}
