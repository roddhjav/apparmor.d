//go:build dev
// +build dev

package main

import (
	"fmt"
	"time"
)

const withTime = true

func printTiming(start, endRead, endParse, end time.Time) {
	printDuration := func(d time.Duration) string {
		if d >= time.Minute {
			return fmt.Sprintf("%.2fmin", d.Minutes())
		} else if d >= time.Second {
			return fmt.Sprintf("%.2fs", d.Seconds())
		} else if d >= time.Millisecond {
			return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
		}
		return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1000)
	}
	readDur := endRead.Sub(start)
	parseDur := endParse.Sub(endRead)
	printDur := end.Sub(endParse)
	totalDur := end.Sub(start)
	fmt.Printf("\x1b[3;97m( Read %s | Parse %s | Print %s | Total %s )\x1b[0m\n",
		printDuration(readDur), printDuration(parseDur),
		printDuration(printDur), printDuration(totalDur),
	)
}
