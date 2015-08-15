package main

import (
	vlc "./.."
	"flag"
	"fmt"
	"time"
)

func main() {
	doHelp := false
	flag.BoolVar(&doHelp, "help", false, "Display usage info.")
	flag.Parse()

	if doHelp || flag.NArg() != 1 {
		fmt.Println("Usage: info <file>")
		return
	}

	str := flag.Arg(1)
	v, err := vlc.New()
	if err != nil {
		fmt.Println("vlc.New():", err)
		return
	}
	m, err := v.MediaNewPath(str)
	if err != nil {
		fmt.Println("vlc.MediaNewPath:", err)
		return
	}
	m.Parse()
	err = vlc.LastError()
	if err != nil {
		fmt.Println("media.Parse():", err)
		return
	}

	player, err := v.NewPlayer()
	if err != nil {
		fmt.Println("vlc.NewPlayer():", err)
		return
	}

	player.SetMedia(m)

	err = player.Play()
	if err != nil {
		fmt.Println("player.Play():", err)
		return
	}

	for k, t := range vlc.MetaTags {
		v := m.GetMeta(t)
		err = vlc.LastError()
		if err != nil {
			fmt.Println("media.GetMeta():", err)
			return
		}
		fmt.Println(k, ":", v)
	}

	time.Sleep(time.Second * 20)
	fmt.Println("end.")
}
