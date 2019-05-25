package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"time"
)

type Game struct {
	screen tcell.Screen
	maze   *Maze
	event  chan Event
	ticker *time.Ticker
}

type Event struct {
	Type string
}

func (g *Game) display() {
	g.screen.Clear()
	for i, row := range g.maze.Points {
		for j := range row {
			st := tcell.StyleDefault.Foreground(tcell.ColorWhite)
			sts := g.maze.Points[i][j].status

			if sts == START {
				g.screen.SetContent(j*2, i, 'S', nil, st)
				g.screen.SetContent(j*2+1, i, ' ', nil, st)
			} else if sts == GOAL {
				g.screen.SetContent(j*2, i, ' ', nil, st)
				g.screen.SetContent(j*2+1, i, 'G', nil, st)
			} else if sts == WALL {
				g.screen.SetContent(j*2, i, tcell.RuneVLine, nil, st)
				g.screen.SetContent(j*2+1, i, tcell.RuneVLine, nil, st)
			} else {
				g.screen.SetContent(j*2, i, ' ', nil, st)
				g.screen.SetContent(j*2+1, i, ' ', nil, st)
			}
		}
	}
}

func (g *Game) Loop() error {
	for {
		g.display()
		g.screen.Show()

		select {
		case ev := <-g.event:
			switch ev.Type {
			case "done":
				return nil
			default:
				return fmt.Errorf(ev.Type)
			}
		case <-g.ticker.C:
			g.screen.Show()
		}
	}
}
