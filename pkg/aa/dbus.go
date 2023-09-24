// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

type Dbus struct {
	Qualifier
	Access    string
	Bus       string
	Name      string
	Path      string
	Interface string
	Member    string
	Label     string
}

func DbusFromLog(log map[string]string, noNewPrivs, fileInherit bool) ApparmorRule {
	return &Dbus{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    log["mask"],
		Bus:       log["bus"],
		Name:      log["name"],
		Path:      log["path"],
		Interface: log["interface"],
		Member:    log["member"],
		Label:     log["peer_label"],
	}
}
