// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logs

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/util"
)

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

var (
	quoted                bool
	isAppArmorLogTemplate = regexp.MustCompile(`apparmor=("DENIED"|"ALLOWED"|"AUDIT")`)
	regAALogs             = []struct {
		regex *regexp.Regexp
		repl  string
	}{
		{regexp.MustCompile(`.*apparmor="`), `apparmor="`},
		{regexp.MustCompile(`(peer_|)pid=[0-9]* `), ""},
		{regexp.MustCompile(` fsuid.*`), ""},
		{regexp.MustCompile(` exe=.*`), ""},
	}
	// Apparmor log states
	state = map[string]string{
		"DENIED":  BoldRed + "DENIED " + Reset,
		"ALLOWED": BoldGreen + "ALLOWED" + Reset,
		"AUDIT":   BoldYellow + "AUDIT  " + Reset,
	}
	// Print order of impression
	keys = []string{
		"profile", "label", // Profile name
		"operation", "name",
		"mask", "bus", "path", "interface", "member", // dbus
		"info", "comm",
		"laddr", "lport", "faddr", "fport", "family", "sock_type", "protocol",
		"requested_mask", "denied_mask", "signal", "peer", // "fsuid", "ouid", "FSUID", "OUID",
	}
	// Color template to use
	colors = map[string]string{
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
)

type AppArmorLog map[string]string

// AppArmorLogs describes all apparmor log entries
type AppArmorLogs []AppArmorLog

// SystemdLog is a simplified systemd json log representation.
type SystemdLog struct {
	Message string `json:"MESSAGE"`
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

// NewApparmorLogs return a new ApparmorLogs list of map from a log file
func NewApparmorLogs(file io.Reader, profile string) AppArmorLogs {
	log := ""
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
			log += line + "\n"
		}
	}

	// Clean logs
	for _, aa := range regAALogs {
		log = aa.regex.ReplaceAllLiteralString(log, aa.repl)
	}

	// Remove doublon in logs
	logs := strings.Split(log, "\n")
	logs = util.RemoveDuplicate(logs)

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
		aa["profile"] = util.DecodeHex(aa["profile"])
		toDecode := []string{"name", "comm"}
		for _, name := range toDecode {
			if value, ok := aa[name]; ok {
				aa[name] = util.DecodeHex(value)
			}
		}

		aaLogs = append(aaLogs, aa)
	}

	return aaLogs
}

// String returns a formatted AppArmor logs string
func (aaLogs AppArmorLogs) String() string {
	res := ""
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
