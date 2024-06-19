// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
)

const DBUS Kind = "dbus"

func init() {
	requirements[DBUS] = requirement{
		"access": []string{
			"send", "receive", "bind", "eavesdrop", "r", "read",
			"w", "write", "rw",
		},
		"bus": []string{"system", "session", "accessibility"},
	}
}

type Dbus struct {
	RuleBase
	Qualifier
	Access    []string
	Bus       string
	Name      string
	Path      string
	Interface string
	Member    string
	PeerName  string
	PeerLabel string
}

func newDbusFromLog(log map[string]string) Rule {
	name := ""
	peerName := ""
	if log["mask"] == "bind" {
		name = log["name"]
	} else {
		peerName = log["name"]
	}
	return &Dbus{
		RuleBase:  newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    []string{log["mask"]},
		Bus:       log["bus"],
		Name:      name,
		Path:      log["path"],
		Interface: log["interface"],
		Member:    log["member"],
		PeerName:  peerName,
		PeerLabel: log["peer_label"],
	}
}

func (r *Dbus) Validate() error {
	if err := validateValues(r.Kind(), "access", r.Access); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return validateValues(r.Kind(), "bus", []string{r.Bus})
}

func (r *Dbus) Compare(other Rule) int {
	o, _ := other.(*Dbus)
	if res := compare(r.Access, o.Access); res != 0 {
		return res
	}
	if res := compare(r.Bus, o.Bus); res != 0 {
		return res
	}
	if res := compare(r.Name, o.Name); res != 0 {
		return res
	}
	if res := compare(r.Path, o.Path); res != 0 {
		return res
	}
	if res := compare(r.Interface, o.Interface); res != 0 {
		return res
	}
	if res := compare(r.Member, o.Member); res != 0 {
		return res
	}
	if res := compare(r.PeerName, o.PeerName); res != 0 {
		return res
	}
	if res := compare(r.PeerLabel, o.PeerLabel); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Dbus) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Dbus) Constraint() constraint {
	return blockKind
}

func (r *Dbus) Kind() Kind {
	return DBUS
}
