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
	Format string
}

func NewMaze(h, w int, s bool, f string) *Maze {
	m := &Maze{
		Width:  w,
		Height: h,
		Seed:   s,
		Format: f,
	}

	return m
}

func (m *Maze) Resize() {
	if m.Height%2 == 0 {
		m.Height++
	} else if m.Height < 5 {
		m.Height = 5
	}
	if m.Width%2 == 0 {
		m.Width++
	} else if m.Width < 5 {
		m.Width = 5
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
			if i == 0 || i == h-1 || j == 0 || j == w-1 {
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
			rand.NewSource(1)
		} else {
			rand.Seed(time.Now().UnixNano())
		}
		r := rand.Intn(len(cand))
		cp := cand[r]
		cand = append(cand[:r], cand[r+1:]...)

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
					np := m.Points[cp.y+dy[i]][cp.x+dx[i]]
					nnp := m.Points[cp.y+dy[i]*2][cp.x+dx[i]*2]
					if np.status == 0 && nnp.status != 2 {
						kw = append(kw, np)
						kkw = append(kkw, nnp)
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

func (m *Maze) printMaze(h, w int, format string) {
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			var cell string
			if m.Points[i][j].status == WALL {
				if i == 1 && j == 0 {
					cell = "S "
				} else if i == h-2 && j == w-1 {
					cell = " G"
				} else {
					if format == "bold" {
						cell = "\033[07m  \033[00m"
					} else if format == "normal" {
						cell = "||"
					}
				}
			} else {
				cell = "  "
			}
			fmt.Printf("%s", cell)
		}
		fmt.Println()
	}
}