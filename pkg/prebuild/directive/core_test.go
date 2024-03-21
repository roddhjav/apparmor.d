// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"reflect"
	"testing"

	"github.com/arduino/go-paths-helper"
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
				Args: map[string]string{
					"bus":  "system",
					"name": "org.gnome.DisplayManager",
					"own":  "",
				},
				File: paths.New(""),
				Raw:  "  #aa:dbus own bus=system name=org.gnome.DisplayManager",
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
				Name: "only",
				Args: map[string]string{"opensuse": ""},
				File: paths.New(""),
				Raw:  "    #aa:only opensuse",
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
