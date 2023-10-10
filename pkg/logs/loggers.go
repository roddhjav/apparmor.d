// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
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
func GetApparmorLogs(file io.Reader, profile string) []string {
	res := ""
	isAppArmorLog := isAppArmorLogTemplate.Copy()
	if profile != "" {
		exp := `apparmor=("DENIED"|"ALLOWED"|"AUDIT")`
		exp = fmt.Sprintf(exp+`.* (profile="%s.*"|label="%s.*")`, profile, profile)
		isAppArmorLog = regexp.MustCompile(exp)
	}

	// Select Apparmor logs
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if isAppArmorLog.MatchString(line) {
			res += line + "\n"
		}
	}

	// Clean & remove doublon in logs
	res = util.DecodeHexInString(res)
	for _, aa := range regCleanLogs {
		res = aa.Regex.ReplaceAllLiteralString(res, aa.Repl)
	}
	logs := strings.Split(res, "\n")
	return util.RemoveDuplicate(logs)
}

// GetAuditLogs return a reader with the logs entries from Auditd
func GetAuditLogs(path string) (io.Reader, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return file, err
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
