// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

// TODO: space in variable need to be tested.
// @{name} = "Mullvad VPN"
// profile mullvad-gui /{opt/"Mullvad/mullvad-gui,opt/VPN"/mullvad-gui,mullvad-gui}  flags=(attach_disconnected,complain) {

func TestDefaultTunables(t *testing.T) {
	tests := []struct {
		name string
		want *AppArmorProfile
	}{
		{
			name: "aa",
			want: &AppArmorProfile{
				Preamble: Preamble{
					Variables: []*Variable{
						{Name: "bin", Values: []string{"/{,usr/}{,s}bin"}},
						{Name: "lib", Values: []string{"/{,usr/}lib{,exec,32,64}"}},
						{Name: "multiarch", Values: []string{"*-linux-gnu*"}},
						{Name: "HOME", Values: []string{"/home/*"}},
						{Name: "user_share_dirs", Values: []string{"/home/*/.local/share"}},
						{Name: "etc_ro", Values: []string{"/{,usr/}etc/"}},
						{Name: "int", Values: []string{"[0-9]{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}"}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultTunables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultTunables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_ParseVariables(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []*Variable
	}{
		{
			name: "firefox",
			content: `@{firefox_name} = firefox{,-esr,-bin}
			@{firefox_lib_dirs} = /{usr/,}lib{,32,64}/@{firefox_name} /opt/@{firefox_name}
			@{firefox_config_dirs} = @{HOME}/.mozilla/
			@{firefox_cache_dirs} = @{user_cache_dirs}/mozilla/
			@{exec_path} = /{usr/,}bin/@{firefox_name} @{firefox_lib_dirs}/@{firefox_name}
			`,
			want: []*Variable{
				{Name: "firefox_name", Values: []string{"firefox{,-esr,-bin}"}},
				{Name: "firefox_lib_dirs", Values: []string{"/{usr/,}lib{,32,64}/@{firefox_name}", "/opt/@{firefox_name}"}},
				{Name: "firefox_config_dirs", Values: []string{"@{HOME}/.mozilla/"}},
				{Name: "firefox_cache_dirs", Values: []string{"@{user_cache_dirs}/mozilla/"}},
				{Name: "exec_path", Values: []string{"/{usr/,}bin/@{firefox_name}", "@{firefox_lib_dirs}/@{firefox_name}"}},
			},
		},
		{
			name: "xorg",
			content: `@{exec_path}  = /{usr/,}bin/X
			@{exec_path} += /{usr/,}bin/Xorg{,.bin}
			@{exec_path} += /{usr/,}lib/Xorg{,.wrap}
			@{exec_path} += /{usr/,}lib/xorg/Xorg{,.wrap}`,
			want: []*Variable{
				{Name: "exec_path", Values: []string{
					"/{usr/,}bin/X",
					"/{usr/,}bin/Xorg{,.bin}",
					"/{usr/,}lib/Xorg{,.wrap}",
					"/{usr/,}lib/xorg/Xorg{,.wrap}"},
				},
			},
		},
		{
			name: "snapd",
			content: `@{lib_dirs} = @{lib}/ /snap/snapd/@{int}@{lib}
			@{exec_path} = @{lib_dirs}/snapd/snapd`,
			want: []*Variable{
				{Name: "lib_dirs", Values: []string{"@{lib}/", "/snap/snapd/@{int}@{lib}"}},
				{Name: "exec_path", Values: []string{"@{lib_dirs}/snapd/snapd"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAppArmorProfile()
			p.ParseVariables(tt.content)
			if !reflect.DeepEqual(p.Variables, tt.want) {
				t.Errorf("AppArmorProfile.ParseVariables() = %v, want %v", p.Variables, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_resolve(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "default",
			input: "@{etc_ro}",
			want:  []string{"/{,usr/}etc/"},
		},
		{
			name:  "empty",
			input: "@{}",
			want:  []string{"@{}"},
		},
		{
			name:  "nil",
			input: "@{foo}",
			want:  []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := DefaultTunables()
			if got := p.resolve(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfile.resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_ResolveAttachments(t *testing.T) {
	tests := []struct {
		name      string
		variables []*Variable
		want      []string
	}{
		{
			name: "firefox",
			variables: []*Variable{
				{Name: "firefox_name", Values: []string{"firefox{,-esr,-bin}"}},
				{Name: "firefox_lib_dirs", Values: []string{"/{usr/,}/lib{,32,64}/@{firefox_name}", "/opt/@{firefox_name}"}},
				{Name: "exec_path", Values: []string{"/{usr/,}bin/@{firefox_name}", "@{firefox_lib_dirs}/@{firefox_name}"}},
			},
			want: []string{
				"/{usr/,}bin/firefox{,-esr,-bin}",
				"/{usr/,}/lib{,32,64}/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
				"/opt/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
			},
		},
		{
			name: "chromium",
			variables: []*Variable{
				{Name: "name", Values: []string{"chromium"}},
				{Name: "lib_dirs", Values: []string{"/{usr/,}lib/@{name}"}},
				{Name: "exec_path", Values: []string{"@{lib_dirs}/@{name}"}},
			},
			want: []string{
				"/{usr/,}lib/chromium/chromium",
			},
		},
		{
			name: "geoclue",
			variables: []*Variable{
				{Name: "libexec", Values: []string{"/{usr/,}libexec"}},
				{Name: "exec_path", Values: []string{"@{libexec}/geoclue", "@{libexec}/geoclue-2.0/demos/agent"}},
			},
			want: []string{
				"/{usr/,}libexec/geoclue",
				"/{usr/,}libexec/geoclue-2.0/demos/agent",
			},
		},
		{
			name: "opera",
			variables: []*Variable{
				{Name: "multiarch", Values: []string{"*-linux-gnu*"}},
				{Name: "name", Values: []string{"opera{,-beta,-developer}"}},
				{Name: "lib_dirs", Values: []string{"/{usr/,}lib/@{multiarch}/@{name}"}},
				{Name: "exec_path", Values: []string{"@{lib_dirs}/@{name}"}},
			},
			want: []string{
				"/{usr/,}lib/*-linux-gnu*/opera{,-beta,-developer}/opera{,-beta,-developer}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewAppArmorProfile()
			p.Variables = tt.variables
			p.ResolveAttachments()
			if !reflect.DeepEqual(p.Attachments, tt.want) {
				t.Errorf("AppArmorProfile.ResolveAttachments() = %v, want %v", p.Attachments, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_NestAttachments(t *testing.T) {
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
			p := NewAppArmorProfile()
			p.Attachments = tt.Attachments
			if got := p.NestAttachments(); got != tt.want {
				t.Errorf("AppArmorProfile.NestAttachments() = %v, want %v", got, tt.want)
			}
		})
	}
}
