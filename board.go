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
	h := m.Height
	w := m.Width

	wall := make([]*Point, 0) // wall
	cand := make([]*Point, 0) // wall candidate

	m.Points = make([][]*Point, m.Height)
	for i := 0; i < m.Height; i++ {
		m.Points[i] = make([]*Point, m.Width)
		for j := 0; j < m.Width; j++ {
			p := &Point{x: j, y: i}

			if i == 1 && j == 0 {
				p.status = START
				wall = append(wall, p)
			} else if i == m.Height-2 && j == m.Width-1 {
				p.status = GOAL
				wall = append(wall, p)
			} else if i == 0 || i == m.Height-1 || j == 0 || j == m.Width-1 {
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
	if m.Seed {
		rand.NewSource(1) // for test
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	for len(cand) > 0 {
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
				kw := make([]*Point, 0)  // 1つ隣
				kkw := make([]*Point, 0) // 2つ隣
				for i := 0; i < 4; i++ {
					_y := cp.y + dy[i]
					_x := cp.x + dx[i]
					if 0 <= _y && _y < h && 0 <= _x && _x < w {
						np := m.Points[_y][_x]
						__y := cp.y + dy[i]*2
						__x := cp.x + dx[i]*2
						if 0 <= __y && __y < h && 0 <= __x && __x < w {
							nnp := m.Points[__y][__x]

							if np.status == PATH && nnp.status != CURRENT { // 1つ隣が道で、2つ隣が現在地じゃない
								kw = append(kw, np)
								kkw = append(kkw, nnp)
							}
						}
					}
				}

				if len(kw) > 0 {
					// 候補から進む方向をランダムに選定
					dp := kw[rand.Intn(len(kw))]
					ddp := kkw[rand.Intn(len(kw))]

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
					if len(current)-2 > 0 {
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
