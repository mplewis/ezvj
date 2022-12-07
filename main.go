package main

import (
	"github.com/k0kubun/pp/v3"
	"github.com/mplewis/figyr"
)

const desc = "EZVJ uses VLC to play random snippets of videos."

type Config struct {
	VLCHost     string `figyr:"default=localhost,description=The host running the VLC server"`
	VLCPort     int    `figyr:"default=8080,description=The port on which VLC's web interface is running"`
	VLCPassword string `figyr:"required,description=The password to the VLC server"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var cfg Config
	figyr.New(desc).MustParse(&cfg)
	pp.Println(cfg)
}
