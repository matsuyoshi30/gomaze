package main

import (
	"errors"
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

func startGame(width, height int, seed bool, format bool) (Result, int, int, error) {
	s, err := initScreen()
	if err != nil {
		return STOPPED, 0, 0, err
	}
	defer s.Fini()

	w, h := s.Size()
	m := NewMaze(w/2, h, seed, format, false)

	game := Game{
		screen: s,
		maze:   m,
	}

	res, err := game.Loop()
	return res, w / 2, h, err
}

func startSearch(width, height int, seed bool, format bool, bfs, dfs bool) (Result, error) {
	s, err := initScreen()
	if err != nil {
		return STOPPED, err
	}
	defer s.Fini()

	w, h := s.Size()
	m := NewMaze(w/2, h, seed, format, true)

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	game := Game{
		screen: s,
		maze:   m,
		bfs:    bfs,
		dfs:    dfs,
		ticker: ticker,
	}

	res, err := game.Loop()
	return res, err
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
			Name:  "bfs",
			Usage: "Show breadth first search",
		},
		&cli.BoolFlag{
			Name:  "dfs",
			Usage: "Show depth first search",
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

		bfs := c.Bool("bfs")
		dfs := c.Bool("dfs")

		if sc {
			start := time.Now()
			res, w, h, err := startGame(tw, th, se, wi)
			if err != nil {
				return err
			}
			end := time.Now()

			if res == GOALED {
				fmt.Println("Congrats!")
				fmt.Printf("[Maze size] Width: %d / Height: %d\n", w, h)
				fmt.Printf("[Your time] %s\n", end.Sub(start))
			}

			return nil
		} else if bfs || dfs {
			if bfs && dfs {
				return errors.New("only choose one whether bfs or dfs")
			}

			start := time.Now()
			res, err := startSearch(tw, th, se, wi, bfs, dfs)
			if err != nil {
				return err
			}
			end := time.Now()

			if res == GOALED {
				fmt.Println("Search complete!")
				fmt.Printf("[Search time] %s\n", end.Sub(start))
			}

			return nil
		} else {
			m := NewMaze(tw, th, se, wi, false)
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
