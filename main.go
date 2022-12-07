package main

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
	"github.com/k0kubun/pp/v3"
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

type Player struct {
	vlcctrl.VLC
	Config
}

func NewPlayer(cfg Config) Player {
	v, err := vlcctrl.NewVLC(cfg.VLCHost, cfg.VLCPort, cfg.VLCPassword)
	check(err)
	return Player{VLC: v, Config: cfg}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func listFiles(dir string) []string {
	fs, err := os.ReadDir(dir)
	check(err)
	var files []string
	for _, f := range fs {
		if f.IsDir() {
			continue
		}
		files = append(files, f.Name())
	}
	return files
}

func (p Player) Add(f string) {
	path := fmt.Sprintf("file://%s", url.PathEscape(path.Join(p.Config.VideoDir, f)))
	fmt.Println(path)
	check(p.VLC.Add(path))
}

func (p Player) playlist() []PlaylistItem {
	root, err := p.VLC.Playlist()
	check(err)

	var pl vlcctrl.Node
	for _, n := range root.Children {
		if n.Name == "Playlist" {
			pl = n
			break
		}
	}
	if pl.Name == "" {
		panic("no playlist")
	}

	items := make([]PlaylistItem, len(pl.Children))
	for i, n := range pl.Children {
		if n.Type != "leaf" {
			continue
		}
		pp.Println(n)
		id, err := strconv.Atoi(n.ID)
		check(err)
		item := PlaylistItem{Name: n.Name, ID: id, Duration: n.Duration}
		items[i] = item
	}
	return items
}

func (p Player) PlayRandomItem() PlaylistItem {
	pl := p.playlist()
	i := rand.Intn(len(pl))
	item := pl[i]
	check(p.VLC.Play(item.ID))
	check(p.VLC.SelectSubtitleTrack(2)) // this is usually english subtitles
	return item
}

func (p Player) PickRandomPlayDuration() time.Duration {
	min := int(p.Config.PlayDurationMin.Nanoseconds())
	max := int(p.Config.PlayDurationMax.Nanoseconds())
	dur := rand.Intn(max-min) + min
	return time.Duration(dur)
}

func (p Player) SeekToRandomPosition(item PlaylistItem, playDuration time.Duration) {
	start := int(float64(item.Duration) * p.Config.ExcludeStart)
	end := item.Duration - int(float64(item.Duration)*p.Config.ExcludeEnd) - int(playDuration.Seconds())
	pos := rand.Intn(end-start) + start
	check(p.VLC.Seek(fmt.Sprintf("%ds", pos)))
}

func main() {
	var cfg Config
	figyr.New(desc).MustParse(&cfg)
	p := NewPlayer(cfg)

	rand.Seed(time.Now().UnixNano())

	check(p.VLC.EmptyPlaylist())
	files := listFiles(cfg.VideoDir)
	for _, f := range files {
		p.Add(f)
	}

	item := p.PlayRandomItem()
	dur := p.PickRandomPlayDuration()
	pp.Println(dur.Seconds())
	p.SeekToRandomPosition(item, dur)
}
