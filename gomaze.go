package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/urfave/cli"
)

func initScreen() tcell.Screen {
	encoding.Register()
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	s.EnableMouse()

	return s
}

func startGame(height int, width int, seed bool, format bool) error {
	s := initScreen()
	defer s.Fini()

	w, h := s.Size()
	m := NewMaze(h-4, (w-4)/2, seed, format)

	event := make(chan Event)

	game := Game{
		screen: s,
		maze:   m,
		event:  event,
	}

	go inputLoop(s, event)

	return game.Loop()
}

func main() {
	app := cli.NewApp()
	app.Name = "gomaze"
	app.Usage = "Generate maze"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "height",
			Usage: "Set the height of maze",
			Value: 30,
		},
		cli.IntFlag{
			Name:  "width",
			Usage: "Set the width of maze",
			Value: 30,
		},
		cli.BoolFlag{
			Name:  "seed",
			Usage: "Set seed for generating specific maze",
		},
		cli.BoolFlag{
			Name:  "screen",
			Usage: "TUI mode",
		},
		cli.BoolFlag{
			Name:  "format",
			Usage: "Format output bold",
		},
		cli.BoolFlag{
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
		th := c.Int("height")
		tw := c.Int("width")
		se := c.Bool("seed")
		sc := c.Bool("screen")
		wi := c.Bool("format")

		if sc {
			return startGame(th, tw, se, wi)
		} else {
			m := NewMaze(th, tw, se, wi)
			m.printMaze()

			return nil
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
