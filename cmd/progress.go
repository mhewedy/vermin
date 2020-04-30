package cmd

import (
	"fmt"
	"time"
)

func PrintProgress() *chan bool {
	quit := make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				fmt.Println()
				return
			default:
				fmt.Print(".")
				time.Sleep(3 * time.Second)
			}
		}
	}()
	return &quit
}
