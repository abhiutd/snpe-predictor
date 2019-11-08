package utils

import "time"

func Every(duration time.Duration, f func()) {
	f()
	ticker := time.NewTicker(duration)
	for {
		select {
		case <-ticker.C:
			f()
		}
	}
}
