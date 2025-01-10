package timer

import (
	"time"
)

var running bool

// Sends 1 to the channel "receiver" each second.
// Sends 0 to the channel "receiver" the time provided as "max" is over.
func StartTimer(receiver chan int, max int) {
	timeout := time.After(time.Duration(max) * time.Second)
	running = true

	for {
		select {
		case <-timeout:
			receiver <- 0
			running = false
			return
		default:
		}
		time.Sleep(time.Second)
		receiver <- 1
	}
}

func IsRunning() bool {
	return running
}
