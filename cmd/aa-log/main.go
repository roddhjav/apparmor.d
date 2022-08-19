// aa-log - Review AppArmor generated messages
// Copyright (C) 2021-2022 Alexandre Pujol <alexandre@pujol.io>
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
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Command line options
var (
	dbus  bool
	help  bool
	path  string
)

// LogFile is the default path to the file to query
const LogFile = "/var/log/audit/audit.log"

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

// getJournalctlDbusSessionLogs return a reader with the logs entries
func getJournalctlDbusSessionLogs(file io.Reader, useFile bool) (io.Reader, error) {
	var logs []SystemdLog
	var stdout bytes.Buffer
	var value string

	if useFile {
		content, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		value = string(content)
	} else {
		cmd := exec.Command("journalctl", "--user", "-b", "-u", "dbus.service", "-o", "json")
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
	exp := "apparmor=(\"DENIED\"|\"ALLOWED\"|\"AUDIT\")"
	if profile != "" {
		exp = fmt.Sprintf(exp+".* (profile=\"%s.*\"|label=\"%s.*\")", profile, profile)
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
	regex := regexp.MustCompile(`type=(USER_|)AVC msg=audit(.*): (pid=.*msg='|)apparmor`)
	log = regex.ReplaceAllLiteralString(log, "apparmor")
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
		if name, ok := aa["name"]; ok {
			aa["name"] = decodeHex(name)
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

func aaLog(path string, profile string, dbus bool) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	/* #nosec G307 */
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if dbus {
		file, err := getJournalctlDbusSessionLogs(file, path != LogFile)
		if err != nil {
			return err
		}
		aaLogs := NewApparmorLogs(file, profile)
		fmt.Print(aaLogs.String())
	} else {
		aaLogs := NewApparmorLogs(file, profile)
		fmt.Print(aaLogs.String())
	}
	return nil
}

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.StringVar(&path, "f", LogFile,
		"Set a log`file` or a suffix to the default log file.")
	flag.BoolVar(&dbus, "d", false, "Show dbus session event.")
}

func main() {
	flag.Parse()
	if help {
		fmt.Printf(`aa-log [-h] [-d] [-f file] [profile]

  Review AppArmor generated messages in a colorful way.
  It can be given an optional profile name to filter the output with.

`)
		flag.PrintDefaults()
		os.Exit(0)
	}

	profile := ""
	if len(flag.Args()) >= 1 {
		profile = flag.Args()[0]
	}

	logfile := filepath.Clean(LogFile + "." + path)
	if _, err := os.Stat(logfile); err != nil {
		logfile = path
	}

	err := aaLog(logfile, profile, dbus)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
