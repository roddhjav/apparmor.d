// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"reflect"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func TestNewOption(t *testing.T) {
	tests := []struct {
		name  string
		file  *paths.Path
		match []string
		want  *Option
	}{
		{
			name: "dbus",
			file: nil,
			match: []string{
				"  #aa:dbus own bus=system name=org.gnome.DisplayManager",
				"dbus",
				"own bus=system name=org.gnome.DisplayManager",
			},
			want: &Option{
				Name: "dbus",
				ArgMap: map[string]string{
					"bus":  "system",
					"name": "org.gnome.DisplayManager",
					"own":  "",
				},
				ArgList: []string{"own", "bus=system", "name=org.gnome.DisplayManager"},
				File:    nil,
				Raw:     "  #aa:dbus own bus=system name=org.gnome.DisplayManager",
			},
		},
		{
			name: "only",
			file: nil,
			match: []string{
				"    #aa:only opensuse",
				"only",
				"opensuse",
			},
			want: &Option{
				Name:    "only",
				ArgMap:  map[string]string{"opensuse": ""},
				ArgList: []string{"opensuse"},
				File:    nil,
				Raw:     "    #aa:only opensuse",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewOption(tt.file, tt.match); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOption() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		file    *paths.Path
		profile string
		want    string
	}{
		{
			name:    "none",
			file:    nil,
			profile: `  `,
			want:    `  `,
		},
		{
			name:    "present",
			file:    nil,
			profile: `  #aa:dbus own bus=system name=org.freedesktop.systemd1`,
			want:    dbusOwnSystemd1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Run(tt.file, tt.profile); got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
