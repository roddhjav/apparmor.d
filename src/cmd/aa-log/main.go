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
type AppArmorLog struct {
	State     string
	Profile   string
	Operation string
	Name      string
	Content   string
}

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

func parseApparmorLogs(file *os.File, profile string) []AppArmorLog {
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
	aaLogs := make([]AppArmorLog, 0)
	getState := regexp.MustCompile(`apparmor=\"([A-Z]*)\"`)
	getProfile := regexp.MustCompile(`profile=\"([A-Za-z0-9_.\-/]*)\" `)
	getOperation := regexp.MustCompile(`operation="(\w+)"`)
	getName := regexp.MustCompile(` name=([A-Za-z0-9_.\-/#"]*)`)
	cleanContent := []*regexp.Regexp{getState, getProfile, getOperation, getName}
	for _, log := range logs {
		name := ""
		names := getName.FindStringSubmatch(log)
		if len(names) >= 2 {
			name = strings.Trim(names[1], `"`)
		}

		content := log
		for _, clean := range cleanContent {
			content = clean.ReplaceAllLiteralString(content, "")
		}
		aaLogs = append(aaLogs,
			AppArmorLog{
				State:     getState.FindStringSubmatch(log)[1],
				Profile:   getProfile.FindStringSubmatch(log)[1],
				Operation: getOperation.FindStringSubmatch(log)[1],
				Name:      name,
				Content:   content[2:],
			})
	}

	return aaLogs
}

func printFormatedApparmorLogs(aaLogs []AppArmorLog) {
	state := map[string]string{
		"DENIED":  BoldRed + "DENIED " + Reset,
		"ALLOWED": BoldGreen + "ALLOWED" + Reset,
	}
	for _, log := range aaLogs {
		fmt.Println(state[log.State],
			FgBlue+log.Profile,
			FgYellow+log.Operation,
			FgMagenta+log.Name+Reset,
			log.Content)
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

	aaLogs := parseApparmorLogs(file, profile)
	printFormatedApparmorLogs(aaLogs)
}
