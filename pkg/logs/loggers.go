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

const (
	// dbusBrokerFields is the list of dbus-broker journal fields needed to rebuild
	// an apparmor log entry.
	dbusBrokerFields = "DBUS_BROKER_POLICY_TYPE,DBUS_BROKER_TRANSMIT_ACTION," +
		"DBUS_BROKER_SENDER_SECURITY_LABEL,DBUS_BROKER_RECEIVER_SECURITY_LABEL," +
		"DBUS_BROKER_MESSAGE_TYPE,DBUS_BROKER_MESSAGE_DESTINATION," +
		"DBUS_BROKER_MESSAGE_PATH,DBUS_BROKER_MESSAGE_INTERFACE,DBUS_BROKER_MESSAGE_MEMBER"

	// ownNameDenied is the error returned to a client when the apparmor policy
	// denies it a dbus name. The bus does not log anything on such denial: the
	// client message is the only trace of it.
	ownNameDenied = "Request to own name refused by policy"
)

// LogFiles is the list of default path to query
var LogFiles = []string{
	"/var/log/audit/audit.log",
	"/var/log/syslog",
}

// SystemdLog is a simplified systemd json log representation.
//
// dbus-broker mediates the dbus apparmor policy itself: its denials are only
// reported in the journal, in its own format, and never reach the audit
// subsystem. They are rebuilt as regular apparmor log entries.
type systemdLog struct {
	Message     string `json:"MESSAGE"`
	Identifier  string `json:"SYSLOG_IDENTIFIER"`
	UID         string `json:"_UID"`
	PolicyType  string `json:"DBUS_BROKER_POLICY_TYPE"`
	Action      string `json:"DBUS_BROKER_TRANSMIT_ACTION"`
	Sender      string `json:"DBUS_BROKER_SENDER_SECURITY_LABEL"`
	Receiver    string `json:"DBUS_BROKER_RECEIVER_SECURITY_LABEL"`
	Type        string `json:"DBUS_BROKER_MESSAGE_TYPE"`
	Destination string `json:"DBUS_BROKER_MESSAGE_DESTINATION"`
	Path        string `json:"DBUS_BROKER_MESSAGE_PATH"`
	Interface   string `json:"DBUS_BROKER_MESSAGE_INTERFACE"`
	Member      string `json:"DBUS_BROKER_MESSAGE_MEMBER"`
}

// String returns the log entry as an apparmor log line. Only the dbus denials
// reported by the bus or its clients are rebuilt, all other entries are left
// untouched.
func (l systemdLog) String() string {
	switch {
	case l.PolicyType == "apparmor":
		return fmt.Sprintf(
			`apparmor="DENIED" operation="dbus_%s" bus="%s" path="%s" interface="%s" member="%s" mask="%s" name="%s" label="%s" peer_label="%s"`,
			l.Type, l.bus(), l.Path, l.Interface, l.Member, l.Action, l.Destination,
			securityLabel(l.Sender), securityLabel(l.Receiver),
		)

	case strings.Contains(l.Message, ownNameDenied):
		// Neither the profile nor the requested name are reported. The syslog
		// identifier is the application id when it is dbus activated, its
		// binary name otherwise: both are only a guess.
		return fmt.Sprintf(
			`apparmor="DENIED" operation="dbus_bind" bus="%s" mask="bind" name="%s" label="%s" info="%s"`,
			l.bus(), l.Identifier, l.Identifier, l.Message,
		)
	}
	return l.Message
}

// bus returns the dbus bus the log entry comes from. The bus is not reported,
// only the system bus runs as root.
func (l systemdLog) bus() string {
	if l.UID == "0" {
		return "system"
	}
	return "session"
}

// securityLabel returns the profile name of a dbus-broker security label:
// 'gnome-clocks (enforce)' -> 'gnome-clocks'
func securityLabel(label string) string {
	name, _, _ := strings.Cut(label, " ")
	return name
}

// GetApparmorLogs return a list of cleaned apparmor logs from a file
func GetApparmorLogs(file io.Reader, profile string, namespace string) []string {
	var logs []string

	exp := `apparmor=("DENIED"|"ALLOWED"|"AUDIT")`
	if profile != "" {
		exp += fmt.Sprintf(`.* (profile="%s.*"|label="%s.*")`, profile, profile)
	}
	if namespace != "" {
		exp += fmt.Sprintf(`.* namespace="root//%s.*"`, namespace)
	}
	isAppArmorLog := regexp.MustCompile(exp)

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
		// The tests logs are generated with the same arguments:
		// journalctl -b --grep=… --output=json --output-fields=… > systemd.log
		//
		// Dbus denials are reported by the bus, but also by the clients
		// themselves, under their own identifier: the logs cannot be filtered
		// on a list of identifiers. A single grep is also faster than one
		// indexed query per identifier.
		args := []string{
			"--grep=apparmor|security policy denied|" + ownNameDenied,
			"--output=json",
			"--output-fields=MESSAGE,SYSLOG_IDENTIFIER,_UID," + dbusBrokerFields,
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
		if strings.Contains(line, "apparmor") || strings.Contains(line, ownNameDenied) {
			jctlRaw = append(jctlRaw, line)
		}
	}

	jctlStr := "[" + strings.Join(jctlRaw, ",\n") + "]"
	if err := json.Unmarshal([]byte(jctlStr), &logs); err != nil {
		return nil, err
	}

	var res strings.Builder
	for _, log := range logs {
		res.WriteString(log.String())
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
	mode := info.Mode()
	if mode&(os.ModeNamedPipe|os.ModeCharDevice) != 0 {
		return nil
	}
	if !mode.IsRegular() {
		return fmt.Errorf("not a regular file: %s", filename)
	}
	if info.Size() == 0 {
		return fmt.Errorf("file is empty: %s", filename)
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
