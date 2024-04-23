// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func TestProfile_AddRule(t *testing.T) {
	tests := []struct {
		name string
		log  map[string]string
		want *Profile
	}{
		{
			name: "capability",
			log:  capability1Log,
			want: &Profile{
				Rules: Rules{capability1},
			},
		},
		{
			name: "network",
			log:  network1Log,
			want: &Profile{
				Rules: Rules{network1},
			},
		},
		{
			name: "mount",
			log:  mount2Log,
			want: &Profile{
				Rules: Rules{mount2},
			},
		},
		{
			name: "signal",
			log:  signal1Log,
			want: &Profile{
				Rules: Rules{signal1},
			},
		},
		{
			name: "ptrace",
			log:  ptrace2Log,
			want: &Profile{
				Rules: Rules{ptrace2},
			},
		},
		{
			name: "unix",
			log:  unix1Log,
			want: &Profile{
				Rules: Rules{unix1},
			},
		},
		{
			name: "dbus",
			log:  dbus2Log,
			want: &Profile{
				Rules: Rules{dbus2},
			},
		},
		{
			name: "file",
			log:  file2Log,
			want: &Profile{
				Rules: Rules{file2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := &Profile{}
			got.AddRule(tt.log)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Profile.AddRule() = |%v|, want |%v|", got, tt.want)
			}
		})
	}
}
