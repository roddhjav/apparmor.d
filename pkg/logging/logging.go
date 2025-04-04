// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logging

import (
	"fmt"
	"os"
)

// Colors
const (
	reset      = "\033[0m"
	bold       = "\033[1m"
	boldRed    = "\033[1;31m"
	boldGreen  = "\033[1;32m"
	boldYellow = "\033[1;33m"
)

// Logging messages prefix
const (
	bulletText  = bold + " ⋅ " + reset
	fatalText   = boldRed + " ✗ Error: " + reset
	errorText   = boldRed + " ✗ " + reset
	successText = boldGreen + " ✓ " + reset
	warningText = boldYellow + " ‼ " + reset
)

var Indent = ""

// Print prints a formatted message. Arguments are handled in the manner of fmt.Print.
func Print(msg string, a ...interface{}) int {
	n, _ := fmt.Fprintf(os.Stdout, msg, a...)
	return n
}

// Println prints a formatted message. Arguments are handled in the manner of fmt.Println.
func Println(msg string) int {
	n, _ := fmt.Fprintf(os.Stdout, "%s\n", msg)
	return n
}

// Bulletf returns a formatted bullet point string
func Bulletf(msg string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s\n", Indent, bulletText, fmt.Sprintf(msg, a...))
}

// Bullet prints a formatted bullet point string
func Bullet(msg string, a ...interface{}) int {
	return Print("%s", Bulletf(msg, a...))
}

// Stepf returns a formatted step string
func Stepf(msg string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s\033[0m\n", Indent, boldGreen, fmt.Sprintf(msg, a...))
}

// Step prints a step title
func Step(msg string, a ...interface{}) int {
	return Print("%s", Stepf(msg, a...))
}

// Successf returns a formatted success string
func Successf(msg string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s\n", Indent, successText, fmt.Sprintf(msg, a...))
}

// Success prints a formatted success message to stdout
func Success(msg string, a ...interface{}) int {
	return Print("%s", Successf(msg, a...))
}

// Warningf returns a formatted warning string
func Warningf(msg string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s\n", Indent, warningText, fmt.Sprintf(msg, a...))
}

// Warning prints a formatted warning message to stdout
func Warning(msg string, a ...interface{}) int {
	return Print("%s", Warningf(msg, a...))
}

// Error returns a formatted error message
func Error(msg string, a ...interface{}) int {
	return Print("%s", fmt.Sprintf("%s%s%s\n", Indent, errorText, fmt.Sprintf(msg, a...)))
}

// Fatalf returns a formatted error message
func Fatalf(msg string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s\n", Indent, fatalText, fmt.Sprintf(msg, a...))
}

// Fatal is equivalent to Print() followed by a call to os.Exit(1).
func Fatal(msg string, a ...interface{}) {
	fmt.Fprint(os.Stderr, Fatalf(msg, a...))
	os.Exit(1)
}
