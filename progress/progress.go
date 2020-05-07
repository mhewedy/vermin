package progress

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"time"
)

type StopFunc func()

func Show(title string) StopFunc {
	quit := make(chan bool, 1)
	i, appendln := 0, false

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(title),
		progressbar.OptionSpinnerType(11),
	)

	go func() {
		for {
			select {
			case <-quit:
				close(quit)
				return
			default:
				if i == 0 {
					time.Sleep(1 * time.Second)
				} else {
					appendln = true
					_ = bar.Add(1)
					time.Sleep(100 * time.Millisecond)
				}
			}
			i++
		}
	}()

	return func() {
		quit <- true
		if appendln {
			fmt.Printf("\r✔ %s", title)
			fmt.Println()
		}
	}
}
