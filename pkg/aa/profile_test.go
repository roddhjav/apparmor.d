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

func TestProfile_GetAttachments(t *testing.T) {
	tests := []struct {
		name        string
		Attachments []string
		want        string
	}{
		{
			name: "firefox",
			Attachments: []string{
				"/{usr/,}bin/firefox{,-esr,-bin}",
				"/{usr/,}lib{,32,64}/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
				"/opt/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
			},
			want: "/{{usr/,}bin/firefox{,-esr,-bin},{usr/,}lib{,32,64}/firefox{,-esr,-bin}/firefox{,-esr,-bin},opt/firefox{,-esr,-bin}/firefox{,-esr,-bin}}",
		},
		{
			name: "geoclue",
			Attachments: []string{
				"/{usr/,}libexec/geoclue",
				"/{usr/,}libexec/geoclue-2.0/demos/agent",
			},
			want: "/{{usr/,}libexec/geoclue,{usr/,}libexec/geoclue-2.0/demos/agent}",
		},
		{
			name:        "null",
			Attachments: []string{},
			want:        "",
		},
		{
			name:        "empty",
			Attachments: []string{""},
			want:        "",
		},
		{
			name:        "not valid aare",
			Attachments: []string{"/file", "relative"},
			want:        "/{file,relative}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Profile{}
			p.Attachments = tt.Attachments
			if got := p.GetAttachments(); got != tt.want {
				t.Errorf("Profile.GetAttachments() = %v, want %v", got, tt.want)
			}
		})
	}
}
