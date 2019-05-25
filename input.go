package main

import (
	"github.com/gdamore/tcell"
)

func inputLoop(s tcell.Screen, event chan<- Event) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEsc || ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
				event <- Event{Type: "done"}
			}
		default:
			continue
		}
	}
}
