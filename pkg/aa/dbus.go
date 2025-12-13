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
	Base
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

func newDbus(q Qualifier, rule rule) (Rule, error) {
	accesses, err := toAccess(DBUS, rule.GetString())
	if err != nil {
		return nil, err
	}
	if err := rule.ValidateMapKeys([]string{"bus", "name", "path", "interface", "member", "peer"}); err != nil {
		return nil, err
	}
	if err := rule.GetValues("peer").ValidateMapKeys([]string{"name", "label"}); err != nil {
		return nil, err
	}
	return &Dbus{
		Base:      newBase(rule),
		Qualifier: q,
		Access:    accesses,
		Bus:       rule.GetValuesAsString("bus"),
		Name:      rule.GetValuesAsString("name"),
		Path:      rule.GetValuesAsString("path"),
		Interface: rule.GetValuesAsString("interface"),
		Member:    rule.GetValuesAsString("member"),
		PeerName:  rule.GetValues("peer").GetValuesAsString("name"),
		PeerLabel: rule.GetValues("peer").GetValuesAsString("label"),
	}, nil
}

func newDbusFromLog(log map[string]string) Rule {
	name := ""
	peerName := ""
	if log["mask"] == "bind" {
		name = log["name"]
	} else {
		peerName = log["name"]
	}
	member, present := log["member"]
	if !present {
		member = log["method"]
	}
	return &Dbus{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Access:    []string{log["mask"]},
		Bus:       log["bus"],
		Name:      name,
		Path:      log["path"],
		Interface: log["interface"],
		Member:    member,
		PeerName:  peerName,
		PeerLabel: log["peer_label"],
	}
}

func (r *Dbus) Kind() Kind {
	return DBUS
}

func (r *Dbus) Constraint() Constraint {
	return BlockRule
}

func (r *Dbus) String() string {
	return renderTemplate(r.Kind(), r)
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

func (r *Dbus) Merge(other Rule) bool {
	o, _ := other.(*Dbus)

	if !r.Equal(o.Qualifier) {
		return false
	}
	if r.Bus == o.Bus && r.Name == o.Name && r.Path == o.Path &&
		r.Interface == o.Interface && r.Member == o.Member &&
		r.PeerName == o.PeerName && r.PeerLabel == o.PeerLabel {
		r.Access = merge(r.Kind(), "access", r.Access, o.Access)
		b := &r.Base
		return b.merge(o.Base)
	}
	return false
}

func (r *Dbus) Lengths() []int {
	return []int{} // No len for dbus
}

func (r *Dbus) setPaddings(max []int) {} // No paddings for dbus
