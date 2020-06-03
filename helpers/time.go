package helpers

import "time"

// SubtractTimeInSeconds returns a float64 of seconds from time1 to time2
func SubtractTimeInSeconds(time1, time2 time.Time) float64 {
	diff := time2.Sub(time1).Seconds()
	return diff
}
