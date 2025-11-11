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
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	return file, nil
}

// GetJournalctlLogs return a reader with the logs entries from Systemd
func GetJournalctlLogs(path string, since string, useFile bool) (io.Reader, error) {
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
		if since == "" {
			args = append(args, "--boot")
		} else {
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

// SelectLogFile return the path of the available log file to parse (audit, syslog, .1, .2)
func SelectLogFile(path string) string {
	info, err := os.Stat(filepath.Clean(path))
	if err == nil && !info.IsDir() {
		return path
	}
	for _, logfile := range LogFiles {
		if _, err := os.Stat(logfile); err == nil {
			oldLogfile := filepath.Clean(logfile + "." + path)
			if _, err := os.Stat(oldLogfile); err == nil {
				return oldLogfile
			} else {
				return logfile
			}
		}
	}
	return ""
}
