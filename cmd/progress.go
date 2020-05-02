package cmd

import (
	"fmt"
	"time"
)

func PrintProgress(title string) *chan bool {
	fmt.Print(title + " ")
	quit := make(chan bool, 10)
	go func() {
		for {
			select {
			case <-quit:
				close(quit)
				return
			default:
				fmt.Print(".")
				time.Sleep(3 * time.Second)
			}
		}
	}()
	return &quit
}
