// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"reflect"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

var (
	apparmorDDir = paths.New("../../../apparmor.d")
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
			file: paths.New("dbus"),
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
				File:    paths.New("dbus"),
				Raw:     "  #aa:dbus own bus=system name=org.gnome.DisplayManager",
			},
		},
		{
			name: "only",
			file: paths.New("only"),
			match: []string{
				"    #aa:only opensuse",
				"only",
				"opensuse",
			},
			want: &Option{
				Name:    "only",
				ArgMap:  map[string]string{"opensuse": ""},
				ArgList: []string{"opensuse"},
				File:    paths.New("only"),
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
		wantErr bool
	}{
		{
			name:    "none",
			file:    paths.New("dummy"),
			profile: `  `,
			want:    `  `,
		},
		{
			name:    "present",
			file:    paths.New("fake-own"),
			profile: `  #aa:dbus own bus=system name=org.freedesktop.systemd1`,
			want:    dbusOwnSystemd1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.file, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}
