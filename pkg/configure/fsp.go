// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"regexp"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	tunables = map[string]string{
		// Set systemd profiles name
		"sd":           "sd",
		"sdu":          "sdu",
		"systemd_user": "systemd-user",
		"systemd":      "systemd",

		// With FSP on apparmor 4.1+, the dbus profiles don't get stacked as they
		"dbus_system":  "dbus-system",
		"dbus_session": "dbus-session",

		// Update name of stacked profiles
		"apt_news":               "",
		"colord":                 "",
		"e2scrub_all":            "",
		"e2scrub":                "",
		"fprintd":                "",
		"fwupd":                  "",
		"fwupdmgr":               "",
		"geoclue":                "",
		"irqbalance":             "",
		"logrotate":              "",
		"ModemManager":           "",
		"nm_priv_helper":         "",
		"pcscd":                  "",
		"polkitd":                "",
		"power_profiles_daemon":  "",
		"rsyslogd":               "",
		"systemd_coredump":       "",
		"systemd_homed":          "",
		"systemd_hostnamed":      "",
		"systemd_importd":        "",
		"systemd_initctl":        "",
		"systemd_journal_remote": "",
		"systemd_journald":       "",
		"systemd_localed":        "",
		"systemd_logind":         "",
		"systemd_machined":       "",
		"systemd_networkd":       "",
		"systemd_oomd":           "",
		"systemd_resolved":       "",
		"systemd_rfkill":         "",
		"systemd_timedated":      "",
		"systemd_timesyncd":      "",
		"systemd_userdbd":        "",
		"upowerd":                "",
	}
)

type FullSystemPolicy struct {
	tasks.BaseTask
}

// NewFullSystemPolicy creates a new FullSystemPolicy task.
func NewFullSystemPolicy() *FullSystemPolicy {
	return &FullSystemPolicy{
		BaseTask: tasks.BaseTask{
			Keyword: "fsp",
			Msg:     "Configure AppArmor for full system policy",
		},
	}
}

func (p FullSystemPolicy) Apply() ([]string, error) {
	res := []string{}

	// Install full system policy profiles
	if err := paths.New("apparmor.d/groups/_full/").CopyFS(p.RootApparmor); err != nil {
		return res, err
	}

	// Set profile name for FSP
	path := p.RootApparmor.Join("tunables/multiarch.d/profiles")
	out, err := path.ReadFileAsString()
	if err != nil {
		return res, err
	}
	for varname, profile := range tunables {
		pattern := regexp.MustCompile(`(@\{p_` + varname + `}=)([^\s]+)`)
		if profile == "" {
			out = pattern.ReplaceAllString(out, `@{p_`+varname+`}={$2,sd//&$2,$2//&sd}`)
		} else {
			out = pattern.ReplaceAllString(out, `@{p_`+varname+`}=`+profile)
		}
	}
	if err := path.WriteFile([]byte(out)); err != nil {
		return res, err
	}

	// Set systemd unit drop-in files
	return res, paths.CopyTo(prebuild.SystemdDir.Join("full"), p.Root.Join("systemd"))
}
