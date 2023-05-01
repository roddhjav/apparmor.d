// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logs

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// GetAuditLogs return a reader with the logs entries from Auditd
func GetAuditLogs(path string) (io.Reader, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	return file, err
}

// GetJournalctlLogs return a reader with the logs entries from Systemd
func GetJournalctlLogs(path string, useFile bool) (io.Reader, error) {
	var logs []systemdLog
	var stdout bytes.Buffer
	var value string

	if useFile {
		content, err := os.ReadFile(filepath.Clean(path))
		if err != nil {
			return nil, err
		}
		value = string(content)
	} else {
		// journalctl -b -o json > systemd.log
		cmd := exec.Command("journalctl", "--boot", "--output=json")
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			return nil, err
		}
		value = stdout.String()
	}

	value = strings.Replace(value, "\n", ",\n", -1)
	value = strings.TrimSuffix(value, ",\n")
	value = `[` + value + `]`
	if err := json.Unmarshal([]byte(value), &logs); err != nil {
		return nil, err
	}
	res := ""
	for _, log := range logs {
		res += log.Message + "\n"
	}
	return strings.NewReader(res), nil
}

func GetLogFile(path string) string {
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
