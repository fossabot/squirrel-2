package util

import "fmt"

// SecondsToHuman returns human readable time format.
func SecondsToHuman(duration uint64) string {
	var hours, minutes, seconds uint64 = 0, 0, 0
	if duration >= 3600 {
		hours = duration / 3600
		duration -= hours * 3600
	}
	if duration >= 60 {
		minutes = duration / 60
		duration -= minutes * 60
	}
	seconds = duration

	str := ""

	if hours > 0 {
		str = fmt.Sprintf("%02dh %02dm %02ds", hours, minutes, seconds)
	} else if minutes > 0 {
		str = fmt.Sprintf("%02dm %02ds", minutes, seconds)
	} else {
		str = fmt.Sprintf("%02ds", seconds)
	}
	return str
}
