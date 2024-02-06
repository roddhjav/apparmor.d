// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
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

func DbusFromLog(log map[string]string) ApparmorRule {
	return &Dbus{
		Qualifier: NewQualifierFromLog(log),
		Access:    log["mask"],
		Bus:       log["bus"],
		Name:      log["name"],
		Path:      log["path"],
		Interface: log["interface"],
		Member:    log["member"],
		Label:     log["peer_label"],
	}
}

func (r *Dbus) Less(other any) bool {
	o, _ := other.(*Dbus)
	if r.Qualifier.Equals(o.Qualifier) {
		if r.Access == o.Access {
			if r.Bus == o.Bus {
				if r.Name == o.Name {
					if r.Path == o.Path {
						if r.Interface == o.Interface {
							if r.Member == o.Member {
								return r.Label < o.Label
							}
							return r.Member < o.Member
						}
						return r.Interface < o.Interface
					}
					return r.Path < o.Path
				}
				return r.Name < o.Name
			}
			return r.Bus < o.Bus
		}
		return r.Access < o.Access
	}
	return r.Qualifier.Less(o.Qualifier)
}

func (r *Dbus) Equals(other any) bool {
	o, _ := other.(*Dbus)
	return r.Access == o.Access && r.Bus == o.Bus && r.Name == o.Name &&
		r.Path == o.Path && r.Interface == o.Interface &&
		r.Member == o.Member && r.Label == o.Label && r.Qualifier.Equals(o.Qualifier)
}
