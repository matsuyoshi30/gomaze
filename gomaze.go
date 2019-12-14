package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/urfave/cli"
)

func initScreen() (tcell.Screen, error) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err = s.Init(); err != nil {
		return nil, err
	}

	return s, nil
}

func startGame(width, height int, seed bool, format bool) error {
	s, err := initScreen()
	if err != nil {
		return err
	}
	defer s.Fini()

	w, h := s.Size()
	m := NewMaze(w/2, h, seed, format)

	game := Game{
		screen: s,
		maze:   m,
	}

	return game.Loop()
}

func main() {
	app := cli.NewApp()
	app.Name = "gomaze"
	app.Usage = "Generate maze"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:  "height",
			Usage: "Set the height of maze",
			Value: 30,
		},
		&cli.IntFlag{
			Name:  "width",
			Usage: "Set the width of maze",
			Value: 30,
		},
		&cli.BoolFlag{
			Name:  "seed",
			Usage: "Set seed for generating specific maze",
		},
		&cli.BoolFlag{
			Name:  "screen",
			Usage: "TUI mode",
		},
		&cli.BoolFlag{
			Name:  "format",
			Usage: "Format output bold",
		},
		&cli.BoolFlag{
			Name: "debug",
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetOutput(os.Stderr)
		} else {
			file, err := os.Open(os.DevNull)
			if err != nil {
				return err
			}
			log.SetOutput(file)
		}
		return nil
	}

	app.Action = func(c *cli.Context) error {
		tw := c.Int("width")
		th := c.Int("height")
		se := c.Bool("seed")
		sc := c.Bool("screen")
		wi := c.Bool("format")

		if sc {
			return startGame(tw, th, se, wi)
		} else {
			m := NewMaze(tw, th, se, wi)
			m.printMaze()

			return nil
		}
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

func wait() {
	time.Sleep(time.Second * 2)
}
