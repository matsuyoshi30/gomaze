package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

// 0 ... 道
// 1 ... 壁
// 2 ... 伸ばし中
type Point struct {
	x, y   int
	status int
}

type Maze struct {
	Points [][]*Point
	Width  int
	Height int
}

func main() {

	// default value
	h := 30
	w := 30

	// 引数があればそうする
	if len(os.Args) > 1 {
		th, _ := strconv.Atoi(os.Args[1])
		if th < 5 {
			h = 5
		} else {
			h = th
		}
		if len(os.Args) > 2 {
			tw, _ := strconv.Atoi(os.Args[2])
			if tw < 5 {
				w = 5
			} else {
				w = tw
			}
		}
	}

	if h%2 == 0 {
		h++
	}
	if w%2 == 0 {
		w++
	}

	// 座標準備
	m := &Maze{
		Height: h,
		Width:  w,
	}

	// 壁
	wall := make([]*Point, 0)
	// 壁候補
	cand := make([]*Point, 0)

	m.Points = make([][]*Point, h)
	// 迷路生成
	for i := 0; i < h; i++ {
		m.Points[i] = make([]*Point, w)
		for j := 0; j < w; j++ {
			var s int

			// 外周設定
			if i == 0 || i == h-1 || j == 0 || j == w-1 {
				s = 1
			} else {
				s = 0
			}
			p := &Point{j, i, s}

			if i == 0 || i == h-1 || j == 0 || j == w-1 {
				wall = append(wall, p)
			} else {
				if i%2 == 0 && j%2 == 0 {
					cand = append(cand, p)
				}
			}

			m.Points[i][j] = p
		}
	}

	// 壁候補がなくなるまでループ
	for len(cand) > 0 {
		// 壁候補から、ランダムに壁伸ばし開始点を取得
		r := rand.Intn(len(cand))
		cp := cand[r]
		cand = append(cand[:r], cand[r+1:]...)
		// fmt.Println(cp)

		if cp.status == 1 {
		} else {
			// 拡張中の壁
			cp.status = 2
			current := make([]*Point, 0)
			current = append(current, cp)

			// 方向 Up Down Right Left
			dx := [4]int{1, -1, 0, 0}
			dy := [4]int{0, 0, 1, -1}

			// 壁にぶつかるまでループ
			for {
				// 点から進む方向の候補
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

					if ddp.status == 1 {
						for _, c := range current {
							c.status = 1
						}
						dp.status = 1
						break
					} else {
						dp.status = 2
						ddp.status = 2
						current = append(current, dp)
						current = append(current, ddp)
						cp = ddp
					}
				} else {
					ddp := current[len(current)-1]
					ddp.status = 0
					dp := current[len(current)-2]
					dp.status = 0
					current = current[:len(current)-2]

					cp = ddp
				}
			}
		}

	}

	// 描画
	fmt.Println()
	for i := 0; i < h; i++ {
		for j := 0; j < w; j++ {
			var cell string

			if m.Points[i][j].status == 1 {
				if i == 1 && j == 0 {
					cell = "S "
				} else if i == h-2 && j == w-1 {
					cell = " G"
				} else {
					cell = "\x1b[7m  \x1b[0m"
				}
			} else {
				cell = "  "
			}
			fmt.Printf("%s", cell)
		}
		fmt.Println()
	}
}
