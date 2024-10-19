package util

import "time"

func DoWithRetry(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts < 0 {
		if err := fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}
