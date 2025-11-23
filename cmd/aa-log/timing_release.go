//go:build !dev

package main

import "time"

const withTime = false

func timeNow() time.Time {
	return time.Time{}
}

func printTiming(start, endRead, endParse, end time.Time) {}
