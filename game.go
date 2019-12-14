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

	wall := tcell.RuneVLine
	path := ' '
	if g.maze.Format {
		wall = tcell.RuneBlock
	}

	for i, row := range g.maze.Points {
		for j, p := range row {
			if p.status == START {
				g.screen.SetContent(j*2, i, 'S', nil, st)
				g.screen.SetContent(j*2+1, i, ' ', nil, st)
			} else if p.status == GOAL {
				g.screen.SetContent(j*2, i, ' ', nil, st)
				g.screen.SetContent(j*2+1, i, 'G', nil, st)
			} else if p.status == WALL {
				g.screen.SetContent(j*2, i, wall, nil, st)
				g.screen.SetContent(j*2+1, i, wall, nil, st)
			} else if p.status == CURRENT {
				g.screen.SetContent(j*2, i, '@', nil, st)
				g.screen.SetContent(j*2+1, i, '@', nil, st)
			} else {
				g.screen.SetContent(j*2, i, path, nil, st)
				g.screen.SetContent(j*2+1, i, path, nil, st)
			}
		}
	}

	g.screen.Show()
}

type Event int

const (
	EXIT Event = iota
	RIGHT
	LEFT
	UP
	DOWN
)

type Result int

const (
	GOALED Result = iota
	STOPPED
)

func (g *Game) Loop() (Result, error) {
	e := make(chan Event)
	go input(g.screen, e)

	for {
		g.display()

		select {
		case ev := <-e:
			switch ev {
			case EXIT:
				return STOPPED, nil
			case RIGHT:
				if g.maze.CheckGoal() { // gaol
					return GOALED, nil
				}
				if g.maze.CheckMaze(RIGHT) {
					g.maze.MoveCurrent(RIGHT)
				}
			case LEFT:
				if g.maze.CheckMaze(LEFT) {
					g.maze.MoveCurrent(LEFT)
				}
			case UP:
				if g.maze.CheckMaze(UP) {
					g.maze.MoveCurrent(UP)
				}
			case DOWN:
				if g.maze.CheckMaze(DOWN) {
					g.maze.MoveCurrent(DOWN)
				}
			}
		}
	}

	return STOPPED, nil
}

func input(s tcell.Screen, e chan<- Event) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEsc, tcell.KeyCtrlC:
				e <- EXIT
			case tcell.KeyRight:
				e <- RIGHT
			case tcell.KeyLeft:
				e <- LEFT
			case tcell.KeyUp:
				e <- UP
			case tcell.KeyDown:
				e <- DOWN
			}
		case *tcell.EventResize:
			s.Sync()
		default:
			continue
		}
	}
}
