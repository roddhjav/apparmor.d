// aa-log - Review AppArmor generated messages
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const usage = `aa-log [-h] [--systemd] [--file file] [profile]

    Review AppArmor generated messages in a colorful way. Supports logs from
    auditd, systemd, syslog as well as dbus session events.

    It can be given an optional profile name to filter the output with.

    Default logs are read from '/var/log/audit/audit.log'. Other files in 
    '/var/log/audit/' can easily be checked: 'aa-log -f 1' parses 'audit.log.1' 

Options:
    -h, --help         Show this help message and exit.
    -f, --file FILE    Set a logfile or a suffix to the default log file.
    -s, --systemd      Parse systemd logs from journalctl.

`

// Command line options
var (
	help    bool
	path    string
	systemd bool
)

// LogFiles is the list of default path to query
var LogFiles = []string{
	"/var/log/audit/audit.log",
	"/var/log/syslog",
}

// Colors
const (
	Reset      = "\033[0m"
	FgGreen    = "\033[32m"
	FgYellow   = "\033[33m"
	FgBlue     = "\033[34m"
	FgMagenta  = "\033[35m"
	FgCian     = "\033[36m"
	FgWhite    = "\033[37m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
)

// AppArmorLog describes a apparmor log entry
type AppArmorLog map[string]string

// AppArmorLogs describes all apparmor log entries
type AppArmorLogs []AppArmorLog

// SystemdLog is a simplified systemd json log representation.
type SystemdLog struct {
	Message string `json:"MESSAGE"`
}

var (
	quoted bool
	isHexa = regexp.MustCompile("^[0-9A-Fa-f]+$")
)

func inSlice(item string, slice []string) bool {
	sort.Strings(slice)
	i := sort.SearchStrings(slice, item)
	return i < len(slice) && slice[i] == item
}

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

func decodeHex(str string) string {
	if isHexa.MatchString(str) {
		bs, _ := hex.DecodeString(str)
		return string(bs)
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

// getAuditLogs return a reader with the logs entries from Auditd
func getAuditLogs(path string) (io.Reader, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	return file, err
}

// getJournalctlLogs return a reader with the logs entries from Systemd
func getJournalctlLogs(path string, useFile bool) (io.Reader, error) {
	var logs []SystemdLog
	var stdout bytes.Buffer
	var value string

	if useFile {
		// content, err := os.ReadFile(filepath.Clean(path))
		content, err := ioutil.ReadFile(filepath.Clean(path))
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

// NewApparmorLogs return a new ApparmorLogs list of map from a log file
func NewApparmorLogs(file io.Reader, profile string) AppArmorLogs {
	log := ""
	exp := `apparmor=("DENIED"|"ALLOWED"|"AUDIT")`
	if profile != "" {
		exp = fmt.Sprintf(exp+`.* (profile="%s.*"|label="%s.*")`, profile, profile)
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
	regex := regexp.MustCompile(`.*apparmor="`)
	log = regex.ReplaceAllLiteralString(log, `apparmor="`)
	regexAppArmorLogs := map[*regexp.Regexp]string{
		regexp.MustCompile(`(peer_|)pid=[0-9]* `): "",
		regexp.MustCompile(` fsuid.*`):            "",
		regexp.MustCompile(` exe=.*`):             "",
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
		aa["profile"] = decodeHex(aa["profile"])
		toDecode := []string{"name", "comm"}
		for _, name := range toDecode {
			if value, ok := aa[name]; ok {
				aa[name] = decodeHex(value)
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
		"AUDIT":   BoldYellow + "AUDIT  " + Reset,
	}
	// Order of impression
	keys := []string{
		"profile", "label", // Profile name
		"operation", "name",
		"mask", "bus", "path", "interface", "member", // dbus
		"info", "comm",
		"laddr", "lport", "faddr", "fport", "family", "sock_type", "protocol",
		"requested_mask", "denied_mask", "signal", "peer", // "fsuid", "ouid", "FSUID", "OUID",
	}
	// Optional colors template to use
	colors := map[string]string{
		"profile":        FgBlue,
		"label":          FgBlue,
		"operation":      FgYellow,
		"name":           FgMagenta,
		"mask":           BoldRed,
		"bus":            FgCian + "bus=",
		"path":           "path=" + FgWhite,
		"requested_mask": "requested_mask=" + BoldRed,
		"denied_mask":    "denied_mask=" + BoldRed,
		"interface":      "interface=" + FgWhite,
		"member":         "member=" + FgGreen,
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
			if !seen[key] && value != "" {
				res += " " + key + "=" + toQuote(value)
			}
		}
		res += "\n"
	}
	return res
}

func getLogFile(path string) string {
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

func aaLog(logger string, path string, profile string) error {
	var err error
	var file io.Reader

	switch logger {
	case "auditd":
		file, err = getAuditLogs(path)
	case "systemd":
		file, err = getJournalctlLogs(path, !inSlice(path, LogFiles))
	default:
		err = fmt.Errorf("Logger %s not supported.", logger)
	}
	if err != nil {
		return err
	}
	aaLogs := NewApparmorLogs(file, profile)
	fmt.Print(aaLogs.String())
	return nil
}

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.StringVar(&path, "f", "", "Set a logfile or a suffix to the default log file.")
	flag.StringVar(&path, "file", "", "Set a logfile or a suffix to the default log file.")
	flag.BoolVar(&systemd, "s", false, "Parse systemd logs from journalctl.")
	flag.BoolVar(&systemd, "systemd", false, "Parse systemd logs from journalctl.")
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	profile := ""
	if len(flag.Args()) >= 1 {
		profile = flag.Args()[0]
	}

	logger := "auditd"
	if systemd {
		logger = "systemd"
	}

	logfile := getLogFile(path)
	err := aaLog(logger, logfile, profile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
