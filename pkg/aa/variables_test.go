// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func TestDefaultTunables(t *testing.T) {
	tests := []struct {
		name string
		want *AppArmorProfile
	}{
		{
			name: "aa",
			want: &AppArmorProfile{
				Preamble: Preamble{
					Variables: []Variable{
						{"bin", []string{"/{,usr/}{,s}bin"}},
						{"lib", []string{"/{,usr/}lib{,exec,32,64}"}},
						{"multiarch", []string{"*-linux-gnu*"}},
						{"user_share_dirs", []string{"/home/*/.local/share"}},
						{"etc_ro", []string{"/{,usr/}etc/"}},
						{"int", []string{"[0-9]{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}"}},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultTunables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppArmorProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_ParseVariables(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []Variable
	}{
		{
			name: "firefox",
			content: `@{firefox_name} = firefox{,-esr,-bin}
			@{firefox_lib_dirs} = /{usr/,}lib{,32,64}/@{firefox_name} /opt/@{firefox_name}
			@{firefox_config_dirs} = @{HOME}/.mozilla/
			@{firefox_cache_dirs} = @{user_cache_dirs}/mozilla/
			@{exec_path} = /{usr/,}bin/@{firefox_name} @{firefox_lib_dirs}/@{firefox_name}
			`,
			want: []Variable{
				{"firefox_name", []string{"firefox{,-esr,-bin}"}},
				{"firefox_lib_dirs", []string{"/{usr/,}lib{,32,64}/@{firefox_name}", "/opt/@{firefox_name}"}},
				{"firefox_config_dirs", []string{"@{HOME}/.mozilla/"}},
				{"firefox_cache_dirs", []string{"@{user_cache_dirs}/mozilla/"}},
				{"exec_path", []string{"/{usr/,}bin/@{firefox_name}", "@{firefox_lib_dirs}/@{firefox_name}"}},
			},
		},
		{
			name: "xorg",
			content: `@{exec_path}  = /{usr/,}bin/X
			@{exec_path} += /{usr/,}bin/Xorg{,.bin}
			@{exec_path} += /{usr/,}lib/Xorg{,.wrap}
			@{exec_path} += /{usr/,}lib/xorg/Xorg{,.wrap}`,
			want: []Variable{
				{"exec_path", []string{
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
			want: []Variable{
				{"lib_dirs", []string{"@{lib}/", "/snap/snapd/@{int}@{lib}"}},
				{"exec_path", []string{"@{lib_dirs}/snapd/snapd"}},
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
			name:  "empty",
			input: "@{}",
			want:  []string{"@{}"},
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
		variables []Variable
		want      []string
	}{
		{
			name: "firefox",
			variables: []Variable{
				{"firefox_name", []string{"firefox{,-esr,-bin}"}},
				{"firefox_lib_dirs", []string{"/{usr/,}/lib{,32,64}/@{firefox_name}", "/opt/@{firefox_name}"}},
				{"exec_path", []string{"/{usr/,}bin/@{firefox_name}", "@{firefox_lib_dirs}/@{firefox_name}"}},
			},
			want: []string{
				"/{usr/,}bin/firefox{,-esr,-bin}",
				"/{usr/,}/lib{,32,64}/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
				"/opt/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
			},
		},
		{
			name: "chromium",
			variables: []Variable{
				{"name", []string{"chromium"}},
				{"lib_dirs", []string{"/{usr/,}lib/@{name}"}},
				{"exec_path", []string{"@{lib_dirs}/@{name}"}},
			},
			want: []string{
				"/{usr/,}lib/chromium/chromium",
			},
		},
		{
			name: "geoclue",
			variables: []Variable{
				{"libexec", []string{"/{usr/,}libexec"}},
				{"exec_path", []string{"@{libexec}/geoclue", "@{libexec}/geoclue-2.0/demos/agent"}},
			},
			want: []string{
				"/{usr/,}libexec/geoclue",
				"/{usr/,}libexec/geoclue-2.0/demos/agent",
			},
		},
		{
			name: "opera",
			variables: []Variable{
				{"multiarch", []string{"*-linux-gnu*"}},
				{"name", []string{"opera{,-beta,-developer}"}},
				{"lib_dirs", []string{"/{usr/,}lib/@{multiarch}/@{name}"}},
				{"exec_path", []string{"@{lib_dirs}/@{name}"}},
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
