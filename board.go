package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	PATH = iota
	WALL
	CURRENT
	START
	GOAL
)

type Point struct {
	x, y   int
	status int
}

type Maze struct {
	Points [][]*Point
	Width  int
	Height int
	Seed   bool
	Format bool
}

func NewMaze(w int, h int, s bool, f bool) *Maze {
	m := Maze{
		Width:  w,
		Height: h,
		Seed:   s,
		Format: f,
	}
	m.Resize()
	m.Generate()

	return &m
}

func (m *Maze) Resize() {
	if m.Height < 5 {
		m.Height = 5
	}
	if m.Width < 5 {
		m.Width = 5
	}

	if m.Height%2 == 0 {
		m.Height--
	}
	if m.Width%2 == 0 {
		m.Width--
	}
}

func (m *Maze) Generate() {
	wall := make([]*Point, 0) // wall
	cand := make([]*Point, 0) // wall candidate

	h := m.Height
	w := m.Width
	m.Points = make([][]*Point, h)
	for i := 0; i < h; i++ {
		m.Points[i] = make([]*Point, w)
		for j := 0; j < w; j++ {
			p := &Point{
				x: j,
				y: i,
			}

			if i == 1 && j == 0 {
				p.status = START
				wall = append(wall, p)
			} else if i == h-2 && j == w-1 {
				p.status = GOAL
				wall = append(wall, p)
			} else if i == 0 || i == h-1 || j == 0 || j == w-1 {
				p.status = WALL
				wall = append(wall, p)
			} else {
				p.status = PATH
				if i%2 == 0 && j%2 == 0 {
					cand = append(cand, p)
				}
			}

			m.Points[i][j] = p
		}
	}

	for len(cand) > 0 {
		if m.Seed {
			rand.NewSource(1) // for test
		} else {
			rand.Seed(time.Now().UnixNano())
		}
		r := rand.Intn(len(cand))
		cp := cand[r]
		cand = append(cand[:r], cand[r+1:]...) // remove cp

		if cp.status != WALL {
			cp.status = CURRENT
			current := make([]*Point, 0)
			current = append(current, cp)

			// Up Down Right Left
			dx := [4]int{1, -1, 0, 0}
			dy := [4]int{0, 0, 1, -1}

			for {
				kw := make([]*Point, 0)
				kkw := make([]*Point, 0)
				for i := 0; i < 4; i++ {
					_y := cp.y + dy[i]
					_x := cp.x + dx[i]
					if 0 <= _y && _y < h && 0 <= _x && _x < w {
						np := m.Points[_y][_x]
						__y := cp.y + dy[i]*2
						__x := cp.x + dx[i]*2
						if 0 <= __y && __y < h && 0 <= __x && __x < w {
							nnp := m.Points[__y][__x]

							if np.status == PATH && nnp.status != CURRENT {
								kw = append(kw, np)
								kkw = append(kkw, nnp)
							}
						}
					}
				}

				if len(kw) > 0 {
					// 候補から進む方向をランダムに選定
					dr := rand.Intn(len(kw))
					dp := kw[dr]
					ddp := kkw[dr]

					if ddp.status == WALL {
						for _, c := range current {
							c.status = WALL
						}
						dp.status = WALL
						break
					} else {
						dp.status = CURRENT
						ddp.status = CURRENT
						current = append(current, dp, ddp)
						cp = ddp
					}
				} else {
					ddp := current[len(current)-1]
					ddp.status = PATH
					dp := current[len(current)-2]
					dp.status = PATH
					current = current[:len(current)-2]

					cp = ddp
				}
			}
		}
	}
}

func (m *Maze) printMaze() {
	format := m.Format

	for _, row := range m.Points {
		for _, p := range row {
			cell := "  "

			sts := p.status
			if sts == START {
				cell = "S "
			} else if sts == GOAL {
				cell = " G"
			} else if sts == WALL {
				if format {
					cell = "\033[07m  \033[00m"
				} else {
					cell = "||"
				}
			}

			fmt.Printf("%s", cell)
		}
		fmt.Println()
	}
}
