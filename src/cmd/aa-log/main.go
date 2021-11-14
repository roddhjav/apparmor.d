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

func removeDuplicateLog(logs []string) []string {
	list := []string{}
	keys := make(map[string]interface{})
	keys[""] = true
	for _, log := range logs {
		if _, v := keys[log]; !v {
			keys[log] = true
			list = append(list, log)
		}
	}
	return list
}

func NewApparmorLogs(file *os.File, profile string) AppArmorLogs {
	log := ""
	exp := fmt.Sprintf("^.*apparmor=(\"DENIED\"|\"ALLOWED\").* profile=\"%s.*\".*$", profile)
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
	cleanAppArmorLogs := []*regexp.Regexp{
		regexp.MustCompile(`type=AVC msg=audit(.*): `),
		regexp.MustCompile(` fsuid.*`),
	}
	for _, clean := range cleanAppArmorLogs {
		log = clean.ReplaceAllLiteralString(log, "")
	}
	replaceAppArmorLogs := regexp.MustCompile(`pid=.* comm`)
	log = replaceAppArmorLogs.ReplaceAllLiteralString(log, "comm")

	// Remove doublon in logs
	logs := strings.Split(log, "\n")
	logs = removeDuplicateLog(logs)

	// Parse log into ApparmorLog struct
	aaLogs := make(AppArmorLogs, 0)
	for _, log := range logs {
		tmp := strings.Split(log, " ")
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

func (aaLogs AppArmorLogs) String() string {
	res := ""
	state := map[string]string{
		"DENIED":  BoldRed + "DENIED " + Reset,
		"ALLOWED": BoldGreen + "ALLOWED" + Reset,
	}
	keys := []string{
		"profile", "operation", "name", "info", "comm", "laddr",
		"lport", "faddr", "fport", "family", "sock_type", "protocol",
		"requested_mask", "denied_mask", // "fsuid", "ouid", "FSUID", "OUID",
	}
	colors := map[string]string{
		"profile":        FgBlue,
		"operation":      FgYellow,
		"name":           FgMagenta,
		"requested_mask": "requested_mask=" + BoldRed,
		"denied_mask":    "denied_mask=" + BoldRed,
	}
	for _, log := range aaLogs {
		res += state[log["apparmor"]]
		delete(log, "apparmor")

		for _, key := range keys {
			if log[key] != "" {
				if colors[key] != "" {
					res += " " + colors[key] + log[key] + Reset
				} else {
					res += " " + key + "=" + log[key]
				}
				delete(log, key)
			}
		}

		for key, value := range log {
			res += " " + key + "=" + value
		}
		res += "\n"
	}
	return res
}
}

func main() {
	profile := ""
	if len(os.Args) >= 2 {
		profile = os.Args[1]
	}

	file, err := os.Open(LogFile)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	aaLogs := NewApparmorLogs(file, profile)
	fmt.Print(aaLogs.String())
}
