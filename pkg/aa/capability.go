// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
)

const CAPABILITY Kind = "capability"

func init() {
	requirements[CAPABILITY] = requirement{
		"name": {
			"audit_control", "audit_read", "audit_write", "block_suspend", "bpf",
			"checkpoint_restore", "chown", "dac_override", "dac_read_search",
			"fowner", "fsetid", "ipc_lock", "ipc_owner", "kill", "lease",
			"linux_immutable", "mac_admin", "mac_override", "mknod", "net_admin",
			"net_bind_service", "net_broadcast", "net_raw", "perfmon", "setfcap",
			"setgid", "setpcap", "setuid", "sys_admin", "sys_boot", "sys_chroot",
			"sys_module", "sys_nice", "sys_pacct", "sys_ptrace", "sys_rawio",
			"sys_resource", "sys_time", "sys_tty_config", "syslog", "wake_alarm",
		},
	}
}

type Capability struct {
	Base
	Qualifier
	Names []string
}

func newCapability(q Qualifier, rule rule) (Rule, error) {
	names, err := toValues(CAPABILITY, "name", rule.GetString())
	if err != nil {
		return nil, err
	}
	return &Capability{
		Base:      newBase(rule),
		Qualifier: q,
		Names:     names,
	}, rule.ValidateMapKeys([]string{})
}

func newCapabilityFromLog(log map[string]string) Rule {
	return &Capability{
		Base:      newBaseFromLog(log),
		Qualifier: newQualifierFromLog(log),
		Names:     Must(toValues(CAPABILITY, "name", log["capname"])),
	}
}

func (r *Capability) Kind() Kind {
	return CAPABILITY
}

func (r *Capability) Constraint() Constraint {
	return BlockRule
}

func (r *Capability) String() string {
	return renderTemplate(r.Kind(), r)
}

func (r *Capability) Validate() error {
	if err := validateValues(r.Kind(), "name", r.Names); err != nil {
		return fmt.Errorf("%s: %w", r, err)
	}
	return nil
}

func (r *Capability) Compare(other Rule) int {
	o, _ := other.(*Capability)
	if res := compare(r.Names, o.Names); res != 0 {
		return res
	}
	return r.Qualifier.Compare(o.Qualifier)
}

func (r *Capability) Merge(other Rule) bool {
	return false // Never merge capabilities
}

func (r *Capability) Lengths() []int {
	return []int{
		r.getLenAudit(),
		r.getLenAccess(),
		length("", r.Names),
	}
}

func (r *Capability) setPaddings(max []int) {
	r.Paddings = append(r.Qualifier.setPaddings(max[:2]), setPaddings(
		max[2:], []string{""}, []any{r.Names})...,
	)
}
