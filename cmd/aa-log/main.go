// aa-log - Review AppArmor generated messages
// Copyright (C) 2021 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// LogFile is the path to the file to query
const LogFile = "/var/log/audit/audit.log"

// Colors
const (
	Reset     = "\033[0m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	BoldRed   = "\033[1;31m"
	BoldGreen = "\033[1;32m"
)

// AppArmorLog describes a apparmor log entry
type AppArmorLog map[string]string

// AppArmorLogs describes all apparmor log entries
type AppArmorLogs []AppArmorLog

var quoted bool

func splitQuoted(r rune) bool {
	if r == '"' {
		quoted = !quoted
	}
	return !quoted && r == ' '
}

func toQuote(str string) string {
	if strings.Contains(str, " ") {
		return `"` + str + `"`
	}
	return str
}

func removeDuplicateLog(logs []string) []string {
	list := []string{}
	keys := map[string]interface{}{"": true}
	for _, log := range logs {
		if _, v := keys[log]; !v {
			keys[log] = true
			list = append(list, log)
		}
	}
	return list
}

// NewApparmorLogs return a new ApparmorLogs list of map from a log file
func NewApparmorLogs(file *os.File, profile string) AppArmorLogs {
	log := ""
	exp := "apparmor=(\"DENIED\"|\"ALLOWED\")"
	if profile != "" {
		exp = fmt.Sprintf(exp+".* profile=\"%s.*\"", profile)
	}
	isAppArmorLog := regexp.MustCompile(exp)

	// Select Apparmor logs
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if isAppArmorLog.MatchString(line) {
			log += line + "\n"
		}
	}

	// Clean logs
	regexAppArmorLogs := map[*regexp.Regexp]string{
		regexp.MustCompile(`type=AVC msg=audit(.*): apparmor`): "apparmor",
		regexp.MustCompile(` fsuid.*`):                         "",
		regexp.MustCompile(`pid=.* comm`):                      "comm",
	}
	for regex, value := range regexAppArmorLogs {
		log = regex.ReplaceAllLiteralString(log, value)
	}

	// Remove doublon in logs
	logs := strings.Split(log, "\n")
	logs = removeDuplicateLog(logs)

	// Parse log into ApparmorLog struct
	aaLogs := make(AppArmorLogs, 0)
	for _, log := range logs {
		quoted = false
		tmp := strings.FieldsFunc(log, splitQuoted)

		aa := make(AppArmorLog)
		for _, item := range tmp {
			kv := strings.Split(item, "=")
			if len(kv) >= 2 {
				aa[kv[0]] = strings.Trim(kv[1], `"`)
			}
		}
		aaLogs = append(aaLogs, aa)
	}

	return aaLogs
}

// String returns a formatted AppArmor logs string
func (aaLogs AppArmorLogs) String() string {
	res := ""
	state := map[string]string{
		"DENIED":  BoldRed + "DENIED " + Reset,
		"ALLOWED": BoldGreen + "ALLOWED" + Reset,
	}
	// Order of impression
	keys := []string{
		"profile", "operation", "name", "info", "comm", "laddr",
		"lport", "faddr", "fport", "family", "sock_type", "protocol",
		"requested_mask", "denied_mask", "signal", "peer", // "fsuid", "ouid", "FSUID", "OUID",
	}
	// Optional colors template to use
	colors := map[string]string{
		"profile":        FgBlue,
		"operation":      FgYellow,
		"name":           FgMagenta,
		"requested_mask": "requested_mask=" + BoldRed,
		"denied_mask":    "denied_mask=" + BoldRed,
	}
	for _, log := range aaLogs {
		seen := map[string]bool{"apparmor": true}
		res += state[log["apparmor"]]

		for _, key := range keys {
			if log[key] != "" {
				if colors[key] != "" {
					res += " " + colors[key] + toQuote(log[key]) + Reset
				} else {
					res += " " + key + "=" + toQuote(log[key])
				}
				seen[key] = true
			}
		}

		for key, value := range log {
			if !seen[key] {
				res += " " + key + "=" + toQuote(value)
			}
		}
		res += "\n"
	}
	return res
}

func main() {
	profile := ""
	if len(os.Args) >= 2 {
		profile = os.Args[1]
	}

	file, err := os.Open(LogFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
			os.Exit(1)
		}
	}()

	aaLogs := NewApparmorLogs(file, profile)
	fmt.Print(aaLogs.String())
}
