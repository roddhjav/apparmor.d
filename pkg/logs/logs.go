// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package logs

import (
	"bufio"
	"fmt"
	"io"
	"os/user"
	"regexp"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/util"
)

// Colors
const (
	reset      = "\033[0m"
	fgGreen    = "\033[32m"
	fgYellow   = "\033[33m"
	fgBlue     = "\033[34m"
	fgMagenta  = "\033[35m"
	fgCian     = "\033[36m"
	fgWhite    = "\033[37m"
	boldRed    = "\033[1;31m"
	boldGreen  = "\033[1;32m"
	boldYellow = "\033[1;33m"
)

// Anonymized username
const Username = "AAD"

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
)

type AppArmorLog map[string]string

// AppArmorLogs describes all apparmor log entries
type AppArmorLogs []AppArmorLog

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

// Anonymize the logs before reporting
func (aaLogs AppArmorLogs) Anonymize() {
	user, _ := user.Current()
	keys := []string{"name", "comm"}
	regAnonymizeLogs := []struct {
		regex *regexp.Regexp
		repl  string
	}{
		{regexp.MustCompile(user.Username), Username},
		{regexp.MustCompile(`/home/[^/]+`), `/home/` + Username},
		{regexp.MustCompile(`[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*-[0-9a-fA-F]*`), `b08dfa60-83e7-567a-1921-a715000001fb`},
	}
	for _, log := range aaLogs {
		for _, key := range keys {
			if _, ok := log[key]; ok {
				for _, aa := range regAnonymizeLogs {
					log[key] = aa.regex.ReplaceAllLiteralString(log[key], aa.repl)
				}
			}
		}
	}
}

// String returns a formatted AppArmor logs string
func (aaLogs AppArmorLogs) String() string {
	// Apparmor log states
	state := map[string]string{
		"DENIED":  boldRed + "DENIED " + reset,
		"ALLOWED": boldGreen + "ALLOWED" + reset,
		"AUDIT":   boldYellow + "AUDIT  " + reset,
	}
	// Print order of impression
	keys := []string{
		"profile", "label", // Profile name
		"operation", "name",
		"mask", "bus", "path", "interface", "member", // dbus
		"info", "comm",
		"laddr", "lport", "faddr", "fport", "family", "sock_type", "protocol",
		"requested_mask", "denied_mask", "signal", "peer", // "fsuid", "ouid", "FSUID", "OUID",
	}
	// Color template to use
	colors := map[string]string{
		"profile":        fgBlue,
		"label":          fgBlue,
		"operation":      fgYellow,
		"name":           fgMagenta,
		"mask":           boldRed,
		"bus":            fgCian + "bus=",
		"path":           "path=" + fgWhite,
		"requested_mask": "requested_mask=" + boldRed,
		"denied_mask":    "denied_mask=" + boldRed,
		"interface":      "interface=" + fgWhite,
		"member":         "member=" + fgGreen,
	}
	res := ""
	for _, log := range aaLogs {
		seen := map[string]bool{"apparmor": true}
		res += state[log["apparmor"]]

		for _, key := range keys {
			if log[key] != "" {
				if colors[key] != "" {
					res += " " + colors[key] + toQuote(log[key]) + reset
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
