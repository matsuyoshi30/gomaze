package main

import (
	"time"

	"github.com/gdamore/tcell"
)

type Game struct {
	screen tcell.Screen
	maze   *Maze
	bfs    bool
	dfs    bool
	queue  []*Point
	stack  []*Point
	ticker *time.Ticker
}

var st = tcell.StyleDefault.Foreground(tcell.ColorWhite)
var ans = tcell.StyleDefault.Foreground(tcell.ColorLightGreen)

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
			} else if p.status == VISITED {
				g.screen.SetContent(j*2, i, '#', nil, st)
				g.screen.SetContent(j*2+1, i, '#', nil, st)
			} else if p.status == ROUTE {
				g.screen.SetContent(j*2, i, '#', nil, ans)
				g.screen.SetContent(j*2+1, i, '#', nil, ans)
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
	NOTGOALED
)

func (g *Game) Loop() (Result, error) {
	e := make(chan Event)
	go input(g.screen, e)

	if g.bfs {
		g.queue = append(g.queue, g.maze.Points[1][1])
		g.maze.Points[1][1].status = VISITED
		g.maze.Points[1][1].cost = 1
	} else {
		g.stack = append(g.stack, g.maze.Points[1][1])
		g.maze.Points[1][1].status = VISITED
		g.maze.Points[1][1].cost = 1
	}

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
		case <-g.ticker.C:
			if g.next() == GOALED {
				return GOALED, nil
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

func (g *Game) next() Result {
	if g.bfs {
		return g.bfsearch()
	} else {
		return g.dfsearch()
	}
}

var dx = [4]int{1, 0, 0, -1}
var dy = [4]int{0, 1, -1, 0}

func (g *Game) bfsearch() Result {
	n, queue := g.queue[0], g.queue[1:]
	for i := 0; i < 4; i++ {
		p := g.maze.Points[n.y+dy[i]][n.x+dx[i]]

		if p.status == GOAL {
			g.shortest()
			return GOALED
		}
		if p.status == PATH {
			p.status = VISITED
			p.cost = g.maze.Points[n.y][n.x].cost + 1
			g.queue = append(g.queue, p)
			return NOTGOALED
		}
	}
	g.queue = queue
	return NOTGOALED
}

func (g *Game) dfsearch() Result {
	n, stack := g.stack[0], g.stack[1:]
	for i := 0; i < 4; i++ {
		p := g.maze.Points[n.y+dy[i]][n.x+dx[i]]

		if p.status == GOAL {
			g.shortest()
			return GOALED
		}
		if p.status == PATH {
			p.status = VISITED
			p.cost = g.maze.Points[n.y][n.x].cost + 1
			g.stack = append([]*Point{p}, g.stack...)
			return NOTGOALED
		}
	}
	g.stack = stack
	return NOTGOALED
}

func (g *Game) shortest() {
	p := g.maze.Points[g.maze.Height-2][g.maze.Width-2] // from goal
	for {
		for i := 0; i < 4; i++ {
			n := g.maze.Points[p.y+dy[i]][p.x+dx[i]]
			if n.status == START {
				p.status = ROUTE
				g.display()
				return
			}

			if n.status == VISITED && n.cost == p.cost-1 {
				p.status = ROUTE
				p = n
				break
			}
		}
	}
}
