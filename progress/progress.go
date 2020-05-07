package progress

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"time"
)

type StopFunc func()

func Immediate(msg ...string) {
	var msgs string
	for _, m := range msg {
		msgs += m + " "
	}
	fmt.Printf("\r✔ %s\n", msgs)
}

func Show(msg string) StopFunc {
	quit := make(chan bool, 1)
	i, isWritten := 0, false

	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(msg),
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
					isWritten = true
					_ = bar.Add(1)
					time.Sleep(100 * time.Millisecond)
				}
			}
			i++
		}
	}()

	return func() {
		quit <- true
		if isWritten {
			fmt.Printf("\r✔ %s\n", msg)
		}
	}
}
