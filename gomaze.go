package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

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
		cli.BoolFlag{
			Name:  "seed",
			Usage: "Set seed for generating specific maze",
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
		se := c.GlobalBool("seed")
		wi := c.GlobalString("format")

		m := NewMaze(th, tw, se, wi)
		m.Resize()
		m.Generate()
		m.printMaze(m.Height, m.Width, m.Format)

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
