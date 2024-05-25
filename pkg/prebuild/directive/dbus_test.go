// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"testing"
)

const dbusOwnSystemd1 = `  dbus bind bus=system name=org.freedesktop.systemd1{,.*},
  dbus receive bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.systemd1{,.*}
       peer=(name=":1.@{int}"),
  dbus receive bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.Properties
       peer=(name=":1.@{int}"),
  dbus receive bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       peer=(name=":1.@{int}"),
  dbus send bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.systemd1{,.*}
       peer=(name="{:1.@{int},org.freedesktop.DBus}"),
  dbus send bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.Properties
       peer=(name="{:1.@{int},org.freedesktop.DBus}"),
  dbus send bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       peer=(name="{:1.@{int},org.freedesktop.DBus}"),
  dbus receive bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.Introspectable
       member=Introspect
       peer=(name=":1.@{int}"),`

func TestDbus_Apply(t *testing.T) {
	tests := []struct {
		name    string
		opt     *Option
		profile string
		want    string
		wantErr bool
	}{
		{
			name: "own",
			opt: &Option{
				Name: "dbus",
				ArgMap: map[string]string{
					"bus":  "system",
					"name": "org.freedesktop.systemd1",
					"own":  "",
				},
				ArgList: []string{"own", "bus=system", "name=org.freedesktop.systemd1"},
				File:    nil,
				Raw:     "  #aa:dbus own bus=system name=org.freedesktop.systemd1",
			},
			profile: "  #aa:dbus own bus=system name=org.freedesktop.systemd1",
			want:    dbusOwnSystemd1,
		},
		{
			name: "own-interface",
			opt: &Option{
				Name: "dbus",
				ArgMap: map[string]string{
					"bus":       "session",
					"name":      "com.rastersoft.dingextension",
					"interface": "org.gtk.Actions",
					"own":       "",
				},
				ArgList: []string{"own", "bus=session", "name=com.rastersoft.dingextension", "interface=org.gtk.Actions"},
				File:    nil,
				Raw:     "  #aa:dbus own bus=session name=com.rastersoft.dingextension interface=org.gtk.Actions",
			},
			profile: "  #aa:dbus own bus=session name=com.rastersoft.dingextension interface=org.gtk.Actions",
			want: `  dbus bind bus=session name=com.rastersoft.dingextension{,.*},
  dbus receive bus=session path=/com/rastersoft/dingextension{,/**}
       interface=com.rastersoft.dingextension{,.*}
       peer=(name=":1.@{int}"),
  dbus receive bus=session path=/com/rastersoft/dingextension{,/**}
       interface=org.gtk.Actions
       peer=(name=":1.@{int}"),
  dbus receive bus=session path=/com/rastersoft/dingextension{,/**}
       interface=org.freedesktop.DBus.Properties
       peer=(name=":1.@{int}"),
  dbus receive bus=session path=/com/rastersoft/dingextension{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       peer=(name=":1.@{int}"),
  dbus send bus=session path=/com/rastersoft/dingextension{,/**}
       interface=com.rastersoft.dingextension{,.*}
       peer=(name="{:1.@{int},org.freedesktop.DBus}"),
  dbus send bus=session path=/com/rastersoft/dingextension{,/**}
       interface=org.gtk.Actions
       peer=(name="{:1.@{int},org.freedesktop.DBus}"),
  dbus send bus=session path=/com/rastersoft/dingextension{,/**}
       interface=org.freedesktop.DBus.Properties
       peer=(name="{:1.@{int},org.freedesktop.DBus}"),
  dbus send bus=session path=/com/rastersoft/dingextension{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       peer=(name="{:1.@{int},org.freedesktop.DBus}"),
  dbus receive bus=session path=/com/rastersoft/dingextension{,/**}
       interface=org.freedesktop.DBus.Introspectable
       member=Introspect
       peer=(name=":1.@{int}"),`,
		},
		{
			name: "talk",
			opt: &Option{
				Name: "dbus",
				ArgMap: map[string]string{
					"bus":   "system",
					"name":  "org.freedesktop.Accounts",
					"label": "accounts-daemon",
					"talk":  "",
				},
				ArgList: []string{"talk", "bus=system", "name=org.freedesktop.Accounts", "label=accounts-daemon"},
				File:    nil,
				Raw:     "  #aa:dbus talk bus=system name=org.freedesktop.Accounts label=accounts-daemon",
			},
			profile: "  #aa:dbus talk bus=system name=org.freedesktop.Accounts label=accounts-daemon",
			want: `  dbus send bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.Accounts{,.*}
       peer=(name="{:1.@{int},org.freedesktop.Accounts{,.*}}", label=accounts-daemon),
  dbus send bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.DBus.Properties
       peer=(name="{:1.@{int},org.freedesktop.Accounts{,.*}}", label=accounts-daemon),
  dbus send bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       peer=(name="{:1.@{int},org.freedesktop.Accounts{,.*}}", label=accounts-daemon),
  dbus receive bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.Accounts{,.*}
       peer=(name="{:1.@{int},org.freedesktop.Accounts{,.*}}", label=accounts-daemon),
  dbus receive bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.DBus.Properties
       peer=(name="{:1.@{int},org.freedesktop.Accounts{,.*}}", label=accounts-daemon),
  dbus receive bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       peer=(name="{:1.@{int},org.freedesktop.Accounts{,.*}}", label=accounts-daemon),`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Directives["dbus"].Apply(tt.opt, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dbus.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Dbus.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
