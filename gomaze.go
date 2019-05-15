package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/urfave/cli"
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
}

func Resize(th, tw int) (int, int) {
	if th%2 == 0 {
		th++
	} else if th < 5 {
		th = 5
	}
	if tw%2 == 0 {
		tw++
	} else if tw < 5 {
		tw = 5
	}

	return th, tw
}

func NewMaze(h, w int) *Maze {
	m := &Maze{
		Width:  w,
		Height: h,
	}

	m.Generate()

	return m
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
		cli.StringFlag{
			Name:  "format",
			Usage: "Format output, normal or bold",
			Value: "normal",
		},
	}

	app.Action = func(c *cli.Context) error {
		th := c.GlobalInt("height")
		tw := c.GlobalInt("width")
		wi := c.GlobalString("format")

		h, w := Resize(th, tw)

		m := NewMaze(h, w)
		m.printMaze(h, w, wi)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
