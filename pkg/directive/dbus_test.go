// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"testing"
)

const (
	dbusOwnSystemd1 = `  #aa/dbus own bus=system name=org.freedesktop.systemd1
  include <abstractions/bus/system/own>
  dbus bind bus=system name=org.freedesktop.systemd1{,.*},
  dbus receive bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.systemd1{,.*}
       peer=(name="@{busname}"),
  dbus send bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.systemd1{,.*}
       peer=(name="{@{busname},org.freedesktop.DBus}"),
  # DBus.Properties: reply to properties request from anyone
  dbus (send receive) bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.Properties
       member={Get,GetAll,Set,PropertiesChanged}
       peer=(name="{@{busname},org.freedesktop.DBus}"),
  # DBus.Introspectable: allow clients to introspect the service
  dbus receive bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.Introspectable
       member=Introspect
       peer=(name="@{busname}"),
  # DBus.ObjectManager: allow clients to enumerate sources
  dbus receive bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       member=GetManagedObjects
       peer=(name="{@{busname},org.freedesktop.systemd1{,.*}}"),
  dbus send bus=system path=/org/freedesktop/systemd1{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       member={InterfacesAdded,InterfacesRemoved}
       peer=(name="{@{busname},org.freedesktop.DBus}"),
`

	dbusOwnInterface = `  #aa/dbus own bus=session name=com.rastersoft.ding interface+=org.gtk.Actions
  include <abstractions/bus/session/own>
  dbus bind bus=session name=com.rastersoft.ding{,.*},
  dbus receive bus=session path=/com/rastersoft/ding{,/**}
       interface=com.rastersoft.ding{,.*}
       peer=(name="@{busname}"),
  dbus send bus=session path=/com/rastersoft/ding{,/**}
       interface=com.rastersoft.ding{,.*}
       peer=(name="{@{busname},org.freedesktop.DBus}"),
  dbus receive bus=session path=/com/rastersoft/ding{,/**}
       interface=org.gtk.Actions
       peer=(name="@{busname}"),
  dbus send bus=session path=/com/rastersoft/ding{,/**}
       interface=org.gtk.Actions
       peer=(name="{@{busname},org.freedesktop.DBus}"),
  # DBus.Properties: reply to properties request from anyone
  dbus (send receive) bus=session path=/com/rastersoft/ding{,/**}
       interface=org.freedesktop.DBus.Properties
       member={Get,GetAll,Set,PropertiesChanged}
       peer=(name="{@{busname},org.freedesktop.DBus}"),
  # DBus.Introspectable: allow clients to introspect the service
  dbus receive bus=session path=/com/rastersoft/ding{,/**}
       interface=org.freedesktop.DBus.Introspectable
       member=Introspect
       peer=(name="@{busname}"),
  # DBus.ObjectManager: allow clients to enumerate sources
  dbus receive bus=session path=/com/rastersoft/ding{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       member=GetManagedObjects
       peer=(name="{@{busname},com.rastersoft.ding{,.*}}"),
  dbus send bus=session path=/com/rastersoft/ding{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       member={InterfacesAdded,InterfacesRemoved}
       peer=(name="{@{busname},org.freedesktop.DBus}"),
`

	dbusTalkAccounts = `  #aa/dbus talk bus=system name=org.freedesktop.Accounts label=accounts-daemon
  # Unix: allow connection to the profile
  unix type=stream peer=(label=accounts-daemon),
  # org.freedesktop.Accounts: send and receive anything to the interface on the specific peer label
  dbus (send receive) bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.Accounts{,.*}
       peer=(name="{@{busname},org.freedesktop.Accounts{,.*},org.freedesktop.DBus}", label=accounts-daemon),
  dbus send bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.Accounts{,.*}
       peer=(name="org.freedesktop.Accounts{,.*}"),
  # DBus.Properties: read and send properties
  dbus (send receive) bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.DBus.Properties
       member={Get,GetAll,Set,PropertiesChanged}
       peer=(name="{@{busname},org.freedesktop.Accounts{,.*},org.freedesktop.DBus}", label=accounts-daemon),
  # DBus.Introspectable: allow service introspection
  dbus send bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.DBus.Introspectable
       member=Introspect
       peer=(name="{@{busname},org.freedesktop.Accounts{,.*},org.freedesktop.DBus}", label=accounts-daemon),
  # DBus.ObjectManager: allow clients to enumerate sources
  dbus send bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       member=GetManagedObjects
       peer=(name="{@{busname},org.freedesktop.Accounts{,.*},org.freedesktop.DBus}", label=accounts-daemon),
  dbus receive bus=system path=/org/freedesktop/Accounts{,/**}
       interface=org.freedesktop.DBus.ObjectManager
       member={InterfacesAdded,InterfacesRemoved}
       peer=(name="{@{busname},org.freedesktop.Accounts{,.*},org.freedesktop.DBus}", label=accounts-daemon),
`

	dbusSeePowerProfiles = `  #aa/dbus see bus=system name=net.hadess.PowerProfiles label=power-profiles-daemon
  # Unix: allow connection to the profile
  unix type=stream peer=(label=power-profiles-daemon),
  # DBus.Properties: read all properties from the interface
  dbus send bus=system path=/net/hadess/PowerProfiles{,/**}
       interface=org.freedesktop.DBus.Properties
       member={Get,GetAll}
       peer=(name="{@{busname},net.hadess.PowerProfiles{,.*}}", label=power-profiles-daemon),
  # DBus.Properties: receive property changed events
  dbus receive bus=system path=/net/hadess/PowerProfiles{,/**}
       interface=org.freedesktop.DBus.Properties
       member=PropertiesChanged
       peer=(name="{@{busname},net.hadess.PowerProfiles{,.*}}", label=power-profiles-daemon),
  # DBus.Introspectable: allow service introspection
  dbus send bus=system path=/net/hadess/PowerProfiles{,/**}
       interface=org.freedesktop.DBus.Introspectable
       member=Introspect
       peer=(name="{@{busname},net.hadess.PowerProfiles{,.*}}", label=power-profiles-daemon),
`
)

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
					"bus":        "session",
					"name":       "com.rastersoft.ding",
					"interface+": "org.gtk.Actions",
					"own":        "",
				},
				ArgList: []string{"own", "bus=session", "name=com.rastersoft.ding", "interface+=org.gtk.Actions"},
				File:    nil,
				Raw:     "  #aa:dbus own bus=session name=com.rastersoft.ding interface+=org.gtk.Actions",
			},
			profile: "  #aa:dbus own bus=session name=com.rastersoft.ding interface+=org.gtk.Actions",
			want:    dbusOwnInterface,
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
			want:    dbusTalkAccounts,
		},
		{
			name: "see",
			opt: &Option{
				Name: "dbus",
				ArgMap: map[string]string{
					"bus":   "system",
					"name":  "net.hadess.PowerProfiles",
					"label": "power-profiles-daemon",
					"talk":  "",
				},
				ArgList: []string{"see", "bus=system", "name=net.hadess.PowerProfiles", "power-profiles-daemon"},
				File:    nil,
				Raw:     "  #aa:dbus see bus=system name=net.hadess.PowerProfiles label=power-profiles-daemon",
			},
			profile: "  #aa:dbus see bus=system name=net.hadess.PowerProfiles label=power-profiles-daemon",
			want:    dbusSeePowerProfiles,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			drctv := NewDbus()
			drctv.SetConfig(cfg)
			got, err := drctv.Apply(tt.opt, tt.profile)
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
