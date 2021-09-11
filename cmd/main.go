package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"the_game_the_life/game"
	"time"

	"github.com/nsf/termbox-go"
)

const cellsInLine = 1

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	gameLoop := game.NewLoop(
		game.WithView(func(m map[game.Pair]*game.Cell) {
			cellID := 1
			diedCount := 0
			livedCount := 0
			for pair, cell := range m {
				if cell.CurrentStatus() == game.CellLive {
					if cellID == cellsInLine {
						fmt.Printf("[ %d : %d ] -> %s\n", pair.X(), pair.Y(), cell.String())
						cellID = 0
					} else {
						fmt.Printf("[ %d : %d ] -> %s\t\t", pair.X(), pair.Y(), cell.String())
					}

					livedCount++
					cellID++
				} else {
					diedCount++
				}
				// if cellID == cellsInLine - 1 {
				// 	fmt.Printf("[ %d : %d ] -> %s\n", pair.X(), pair.Y(), cell.String())
				// 	cellID = 0
				// } else {
				// 	fmt.Printf("[ %d : %d ] -> %s\t\t", pair.X(), pair.Y(), cell.String())
				// }
				// cellID++
			}
			fmt.Println()
			fmt.Println()
			fmt.Printf("lived: %d\t\tdied: %d\n\n", livedCount, diedCount)
			time.Sleep(1000 * time.Millisecond)
		}),
		game.WithView(func(m map[game.Pair]*game.Cell) {
			for pair, cell := range m {
				if cell.CurrentStatus() == game.CellLive {
					termbox.SetCell(pair.X(), pair.Y(), '█', termbox.ColorDefault, termbox.ColorDefault)
				} else {
					termbox.SetCell(pair.X(), pair.Y(), '█', termbox.ColorDefault, termbox.ColorDefault)
				}
			}
			time.Sleep(1 * time.Second)
		}),
		game.WithLiveCell(game.NewPair(10, 10)),
		game.WithLiveCell(game.NewPair(11, 11)),
		game.WithLiveCell(game.NewPair(10, 11)),
		game.WithLiveCell(game.NewPair(9, 12)),

		game.WithLiveCell(game.NewPair(9, 12)),
		game.WithLiveCell(game.NewPair(9, 11)),
		game.WithLiveCell(game.NewPair(9, 10)),
	)
	go gameLoop.StartTheLife()
	c := make(chan os.Signal, 1)
	signal.Notify(c, []os.Signal{syscall.SIGINT, syscall.SIGTERM}...)

	// Block until we receive our signal.
	<-c
	log.Println("stop")
	gameLoop.Stop()
	termbox.Close()
}
