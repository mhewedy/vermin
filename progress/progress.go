package progress

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"time"
)

const format = "\r✔ %s\n"

type StopFunc func()

func Immediate(msg ...string) {
	var msgs string
	for _, m := range msg {
		msgs += m + " "
	}

	fmt.Printf(format, msgs)
}

func Show(msg string, initialWait bool) StopFunc {
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
				if initialWait && i == 0 {
					time.Sleep(600 * time.Millisecond)
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
			fmt.Printf(format, msg)
		}
	}
}
