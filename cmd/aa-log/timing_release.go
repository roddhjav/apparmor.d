//go:build !dev
// +build !dev

package main

import "time"

const withTime = false

func printTiming(start, endRead, endParse, end time.Time) {}
