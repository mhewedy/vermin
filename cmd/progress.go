package cmd

import (
	"fmt"
	"time"
)

func PrintProgress(title string) *chan bool {
	quit := make(chan bool, 1)
	i := 0
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
	return &quit
}
