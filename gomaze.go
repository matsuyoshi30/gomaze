package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/urfave/cli"
)

func initScreen() tcell.Screen {
	s, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if err := s.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	return s
}

func startGame(height int, width int, seed bool, format string) error {
	s := initScreen()
	defer s.Fini()

	w, h := s.Size()
	m := NewMaze(h-4, (w-2)/2, seed, format)
	m.Generate()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	event := make(chan Event)

	game := Game{
		screen: s,
		maze:   m,
		event:  event,
		ticker: ticker,
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
		cli.StringFlag{
			Name:  "format",
			Usage: "Format output, normal or bold",
			Value: "normal",
		},
	}

	app.Action = func(c *cli.Context) error {
		th := c.GlobalInt("height")
		tw := c.GlobalInt("width")
		se := c.GlobalBool("seed")
		sc := c.GlobalBool("screen")
		wi := c.GlobalString("format")

		if sc {
			return startGame(th, tw, se, wi)
		} else {
			m := NewMaze(th, tw, se, wi)
			m.Generate()
			m.printMaze(m.Height, m.Width, m.Format)

			return nil
		}
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
