package progress

import (
	"fmt"
	"time"
)

type StopFunc func()

func Show(title string) StopFunc {
	quit := make(chan bool, 1)
	i, appendln := 0, false

	go func() {
		for {
			select {
			case <-quit:
				close(quit)
				return
			default:
				const d = 3 * time.Second
				if i == 0 {
					time.Sleep(1 * time.Second)
				} else if i == 1 {
					appendln = true
					fmt.Print(title + " ")
					time.Sleep(d)
				} else {
					fmt.Print(".")
					time.Sleep(d)
				}
			}
			i++
		}
	}()

	return func() {
		quit <- true
		if appendln {
			fmt.Println()
		}
	}
}
