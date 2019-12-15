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
	VISITED // only used for search
	START
	GOAL
)

type Point struct {
	x, y   int
	status int
}

type Maze struct {
	Points  [][]*Point
	Width   int
	Height  int
	Current Point
	Seed    bool
	Format  bool
}

func NewMaze(w int, h int, s bool, f bool, search bool) *Maze {
	m := Maze{
		Width:  w,
		Height: h,
		Seed:   s,
		Format: f,
	}
	m.Resize()
	m.Generate()

	m.SetCurrent(1, 1)

	return &m
}

func (m *Maze) SetCurrent(x, y int) {
	m.Points[y][x].status = CURRENT
	m.Current.x = x
	m.Current.y = y
}

func (m *Maze) UnsetCurrent(x, y int) {
	m.Points[y][x].status = PATH
}

func (m *Maze) GetCurrent() (int, int) {
	return m.Current.x, m.Current.y
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
	cand := m.initialize()

	if m.Seed {
		rand.NewSource(1) // for test
	} else {
		rand.Seed(time.Now().UnixNano())
	}

	// Right Left Up Down
	dx := [4]int{1, -1, 0, 0}
	dy := [4]int{0, 0, 1, -1}

	for len(cand) > 0 { // loop wall candidate
		r := rand.Intn(len(cand))
		cp := cand[r]                          // select one candidate point
		cand = append(cand[:r], cand[r+1:]...) // remove cp

		if cp.status != WALL {
			cp.status = CURRENT
			current := make([]*Point, 0)
			current = append(current, cp)

			for {
				kw := make([]*Point, 0)  // 1つ隣
				kkw := make([]*Point, 0) // 2つ隣

				for i := 0; i < 4; i++ {
					_y := cp.y + dy[i]
					_x := cp.x + dx[i]
					if 0 <= _y && _y < m.Height && 0 <= _x && _x < m.Width {
						np := m.Points[_y][_x]

						__y := cp.y + dy[i]*2
						__x := cp.x + dx[i]*2
						if 0 <= __y && __y < m.Height && 0 <= __x && __x < m.Width {
							nnp := m.Points[__y][__x]

							if np.status == PATH && nnp.status != CURRENT { // 1つ隣が道で、2つ隣が今まで掘っていた線地じゃない
								kw = append(kw, np)
								kkw = append(kkw, nnp)
							}
						}
					}
				}

				if len(kw) > 0 {
					// 候補から進む方向をランダムに選定
					_r := rand.Intn(len(kw))
					dp := kw[_r]
					ddp := kkw[_r]

					if ddp.status == WALL { // 進行方向の2つ先が壁
						for _, c := range current { // 今まで掘っていた線を wall に
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
				} else { // 候補がない
					for _, c := range current {
						c.status = WALL
					}
				}
			}
		}
	}
}

// initialize maze point (start, goal and wall)
func (m *Maze) initialize() []*Point {
	cand := make([]*Point, 0) // wall candidate

	m.Points = make([][]*Point, m.Height)
	for i := 0; i < m.Height; i++ {
		m.Points[i] = make([]*Point, m.Width)
		for j := 0; j < m.Width; j++ {
			p := &Point{x: j, y: i}

			if i == 1 && j == 0 { // (1, 0)
				p.status = START
			} else if i == m.Height-2 && j == m.Width-1 { // (width, height-1)
				p.status = GOAL
			} else if i == 0 || i == m.Height-1 || j == 0 || j == m.Width-1 { // (0, _), (width, _), (_, 0), (_, height)
				p.status = WALL
			} else {
				p.status = PATH
				if i%2 == 0 && j%2 == 0 {
					cand = append(cand, p)
				}
			}

			m.Points[i][j] = p
		}
	}

	return cand
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
			} else if sts == CURRENT {
				cell = "@@"
			}

			fmt.Printf("%s", cell)
		}
		fmt.Println()
	}
}

func (m *Maze) CheckMaze(e Event) bool {
	x, y := m.GetCurrent()

	switch e {
	case RIGHT:
		if m.Points[y][x+1].status == GOAL { // goal
			return true
		}
		if x+1 > m.Width {
			return false
		}
		if m.Points[y][x+1].status == PATH {
			return true
		}
		return false
	case LEFT:
		if x-1 < 0 {
			return false
		}
		if m.Points[y][x-1].status == PATH {
			return true
		}
		return false
	case UP:
		if y-1 < 0 {
			return false
		}
		if m.Points[y-1][x].status == PATH {
			return true
		}
		return false
	case DOWN:
		if y+1 > m.Height {
			return false
		}
		if m.Points[y+1][x].status == PATH {
			return true
		}
		return false
	}

	return false
}

func (m *Maze) MoveCurrent(e Event) {
	x, y := m.GetCurrent()
	m.UnsetCurrent(x, y)

	switch e {
	case RIGHT:
		m.SetCurrent(x+1, y)
	case LEFT:
		m.SetCurrent(x-1, y)
	case UP:
		m.SetCurrent(x, y-1)
	case DOWN:
		m.SetCurrent(x, y+1)
	}
}

func (m *Maze) CheckGoal() bool {
	x, y := m.GetCurrent()
	if m.Points[y][x+1].status == GOAL {
		return true
	}
	return false
}
