package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"path"
	"strconv"
	"time"

	vlcctrl "github.com/CedArctic/go-vlc-ctrl"
)

type Player struct {
	vlcctrl.VLC
	Config
}

func NewPlayer(cfg Config) Player {
	v, err := vlcctrl.NewVLC(cfg.VLCHost, cfg.VLCPort, cfg.VLCPassword)
	check(err)
	return Player{VLC: v, Config: cfg}
}

func (p Player) Add(f string) {
	path := fmt.Sprintf("file://%s", url.PathEscape(path.Join(p.Config.VideoDir, f)))
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

func (p Player) SeekToRandomPosition(item PlaylistItem, playDuration time.Duration) int {
	start := int(float64(item.Duration) * p.Config.ExcludeStart)
	end := item.Duration - int(float64(item.Duration)*p.Config.ExcludeEnd) - int(playDuration.Seconds())
	pos := rand.Intn(end-start) + start
	check(p.VLC.Seek(fmt.Sprintf("%ds", pos)))

	// HACK: ensure VLC actually seeks to the position we want
	go func() {
		// VLC can take a few seconds to load a video - keep checking for up to 15 seconds
		attemptFor := time.Duration(math.Min(15, playDuration.Seconds())) * time.Second
		startChecking := time.Now()
		for time.Since(startChecking) < attemptFor {
			time.Sleep(100 * time.Millisecond) // if we don't wait, VLC will lie and say it seeked when the file wasn't yet open
			s, _ := p.VLC.GetStatus()          // HACK: this throws an error every time, but s.Time is still accurate
			if s.Time < uint(attemptFor.Seconds()) {
				check(p.VLC.Seek(fmt.Sprintf("%ds", pos)))
			}
		}
	}()

	return pos
}
