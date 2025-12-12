// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logs

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/util"
)

// LogFiles is the list of default path to query
var LogFiles = []string{
	"/var/log/audit/audit.log",
	"/var/log/syslog",
}

// SystemdLog is a simplified systemd json log representation.
type systemdLog struct {
	Message string `json:"MESSAGE"`
}

// GetApparmorLogs return a list of cleaned apparmor logs from a file
func GetApparmorLogs(file io.Reader, profile string, namespace string) []string {
	var logs []string

	isAppArmorLog := isAppArmorLogTemplate.Copy()
	exp := `apparmor=("DENIED"|"ALLOWED"|"AUDIT")`
	if profile != "" {
		exp = fmt.Sprintf(exp+`.* (profile="%s.*"|label="%s.*")`, profile, profile)
		isAppArmorLog = regexp.MustCompile(exp)
	}
	if namespace != "" {
		exp = fmt.Sprintf(exp+`.* namespace="root//%s.*"`, namespace)
		isAppArmorLog = regexp.MustCompile(exp)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if isAppArmorLog.MatchString(line) {
			logs = append(logs,
				regCleanLogs.Replace(util.DecodeHexInString(line)),
			)
		}
	}
	return util.RemoveDuplicate(logs)
}

// GetAuditLogs return a reader with the logs entries from Auditd
func GetAuditLogs(path string) (io.Reader, error) {
	if path == "/dev/stdin" || path == "-" {
		return os.Stdin, nil
	}
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GetJournalctlLogs return a reader with the logs entries from Systemd
func GetJournalctlLogs(path string, boot string, since string, useFile bool) (io.Reader, error) {
	var logs []systemdLog
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	var scanner *bufio.Scanner

	if useFile {
		file, err := os.Open(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
		scanner = bufio.NewScanner(file)
	} else {
		// journalctl -b -o json -g apparmor -t kernel -t audit -t dbus-daemon --output-fields=MESSAGE > systemd.log
		args := []string{
			"--grep=apparmor", "--identifier=kernel",
			"--identifier=audit", "--identifier=dbus-daemon",
			"--output=json", "--output-fields=MESSAGE",
		}
		if boot != "" {
			args = append(args, "--boot="+boot)
		} else if since == "" {
			args = append(args, "--boot")
		}
		if since != "" {
			args = append(args, "--since="+since)
		}
		cmd := exec.Command("journalctl", args...)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil && stderr.Len() != 0 {
			return nil, fmt.Errorf("journalctl: %s", stderr.String())
		}
		scanner = bufio.NewScanner(&stdout)
	}

	var jctlRaw []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "apparmor") {
			jctlRaw = append(jctlRaw, line)
		}
	}

	jctlStr := "[" + strings.Join(jctlRaw, ",\n") + "]"
	if err := json.Unmarshal([]byte(jctlStr), &logs); err != nil {
		return nil, err
	}

	var res strings.Builder
	for _, log := range logs {
		res.WriteString(log.Message)
		res.WriteString("\n")
	}
	return strings.NewReader(res.String()), nil
}

// validateLogFile checks if a file exists, is readable, and is not empty.
func validateLogFile(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("not a regular file: %s", filename)
	}
	if info.Size() == 0 {
		return fmt.Errorf("file is empty: %s", filename)
	}
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to read: %s", filename)
	}
	if cerr := file.Close(); cerr != nil {
		return fmt.Errorf("unable to close file %s: %w", filename, cerr)
	}
	return nil
}

// SelectLogFile return the path of the available log file to parse (audit, syslog, .1, .2)
func SelectLogFile(input string) (string, error) {
	if input == "/dev/stdin" || input == "-" {
		return input, nil
	}

	// If a specific file path is provided
	if input != "" {
		path := filepath.Clean(input)

		// Check if it's a full path that exists
		if _, err := os.Stat(path); err == nil {
			if err := validateLogFile(path); err != nil {
				return "", err
			}
			return path, nil
		}

		// Try as a suffix to default log files (e.g., "1" -> audit.log.1)
		for _, logfile := range LogFiles {
			suffixedFile := logfile + "." + input
			if _, err := os.Stat(suffixedFile); err == nil {
				if err := validateLogFile(suffixedFile); err != nil {
					return "", err
				}
				return suffixedFile, nil
			}
		}

		return "", fmt.Errorf("log file not found: %s", input)
	}

	// No input provided, find first available default log file
	for _, logfile := range LogFiles {
		if _, err := os.Stat(logfile); err == nil {
			if err := validateLogFile(logfile); err != nil {
				return "", err
			}
			return logfile, nil
		}
	}
	return "", fmt.Errorf("no log file found")
}
