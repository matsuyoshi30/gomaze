package main

import (
	"github.com/gdamore/tcell"
)

type Game struct {
	screen tcell.Screen
	maze   *Maze
}

var st = tcell.StyleDefault.Foreground(tcell.ColorWhite)

func (g *Game) display() {
	g.screen.Clear()

	for i, row := range g.maze.Points {
		for j, p := range row {
			if p.status == START {
				g.screen.SetContent(j*2, i, 'S', nil, st)
				g.screen.SetContent(j*2+1, i, ' ', nil, st)
			} else if p.status == GOAL {
				g.screen.SetContent(j*2, i, ' ', nil, st)
				g.screen.SetContent(j*2+1, i, 'G', nil, st)
			} else if p.status == WALL {
				g.screen.SetContent(j*2, i, tcell.RuneVLine, nil, st)
				g.screen.SetContent(j*2+1, i, tcell.RuneVLine, nil, st)
			} else {
				g.screen.SetContent(j*2, i, ' ', nil, st)
				g.screen.SetContent(j*2+1, i, ' ', nil, st)
			}
		}
	}

	g.screen.Show()
}

func (g *Game) Loop() error {
	quit := make(chan struct{})

	g.display()
	go func() {
		for {
			ev := g.screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyEsc, tcell.KeyCtrlC:
					close(quit)
					return
				}
			case *tcell.EventResize:
				g.screen.Sync()
			}
		}
	}()

	<-quit

	g.screen.Fini()

	return nil
}
