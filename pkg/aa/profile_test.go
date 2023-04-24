// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"
)

func TestNewAppArmorProfile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    *AppArmorProfile
	}{
		{
			name:    "aa",
			content: "",
			want: &AppArmorProfile{
				Content: "",
				Variables: map[string][]string{
					"libexec":         {},
					"etc_ro":          {"/{usr/,}etc/"},
					"multiarch":       {"*-linux-gnu*"},
					"user_share_dirs": {"/home/*/.local/share"},
				},
				Attachments: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAppArmorProfile(tt.content); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppArmorProfile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_ParseVariables(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    map[string][]string
	}{
		{
			name: "firefox",
			content: `@{firefox_name} = firefox{,-esr,-bin}
			@{firefox_lib_dirs} = /{usr/,}lib{,32,64}/@{firefox_name} /opt/@{firefox_name}
			@{firefox_config_dirs} = @{HOME}/.mozilla/
			@{firefox_cache_dirs} = @{user_cache_dirs}/mozilla/
			@{exec_path} = /{usr/,}bin/@{firefox_name} @{firefox_lib_dirs}/@{firefox_name}
			`,
			want: map[string][]string{
				"firefox_name":        {"firefox{,-esr,-bin}"},
				"firefox_config_dirs": {"@{HOME}/.mozilla/"},
				"firefox_lib_dirs":    {"/{usr/,}lib{,32,64}/@{firefox_name}", "/opt/@{firefox_name}"},
				"firefox_cache_dirs":  {"@{user_cache_dirs}/mozilla/"},
				"exec_path":           {"/{usr/,}bin/@{firefox_name}", "@{firefox_lib_dirs}/@{firefox_name}"},
			},
		},
		{
			name: "xorg",
			content: `@{exec_path}  = /{usr/,}bin/X
			@{exec_path} += /{usr/,}bin/Xorg{,.bin}
			@{exec_path} += /{usr/,}lib/Xorg{,.wrap}
			@{exec_path} += /{usr/,}lib/xorg/Xorg{,.wrap}`,
			want: map[string][]string{
				"exec_path": {
					"/{usr/,}bin/X",
					"/{usr/,}bin/Xorg{,.bin}",
					"/{usr/,}lib/Xorg{,.wrap}",
					"/{usr/,}lib/xorg/Xorg{,.wrap}",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &AppArmorProfile{
				Content:     tt.content,
				Variables:   map[string][]string{},
				Attachments: []string{},
			}

			p.ParseVariables()
			if !reflect.DeepEqual(p.Variables, tt.want) {
				t.Errorf("AppArmorProfile.ParseVariables() = %v, want %v", p.Variables, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_resolve(t *testing.T) {
	tests := []struct {
		name      string
		variables map[string][]string
		input     string
		want      []string
	}{
		{
			name:      "empty",
			variables: Tunables,
			input:     "@{}",
			want:      []string{"@{}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &AppArmorProfile{
				Content:     "",
				Variables:   tt.variables,
				Attachments: []string{},
			}
			if got := p.resolve(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfile.resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_ResolveAttachments(t *testing.T) {
	tests := []struct {
		name      string
		variables map[string][]string
		want      []string
	}{
		{
			name: "firefox",
			variables: map[string][]string{
				"firefox_name":     {"firefox{,-esr,-bin}"},
				"firefox_lib_dirs": {"/{usr/,}/lib{,32,64}/@{firefox_name}", "/opt/@{firefox_name}"},
				"exec_path":        {"/{usr/,}bin/@{firefox_name}", "@{firefox_lib_dirs}/@{firefox_name}"},
			},
			want: []string{
				"/{usr/,}bin/firefox{,-esr,-bin}",
				"/{usr/,}/lib{,32,64}/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
				"/opt/firefox{,-esr,-bin}/firefox{,-esr,-bin}",
			},
		},
		{
			name: "chromium",
			variables: map[string][]string{
				"chromium_name":     {"chromium"},
				"chromium_lib_dirs": {"/{usr/,}lib/@{chromium_name}"},
				"exec_path":         {"@{chromium_lib_dirs}/@{chromium_name}"},
			},
			want: []string{
				"/{usr/,}lib/chromium/chromium",
			},
		},
		{
			name: "geoclue",
			variables: map[string][]string{
				"libexec":   {"/{usr/,}libexec"},
				"exec_path": {"@{libexec}/geoclue", "@{libexec}/geoclue-2.0/demos/agent"},
			},
			want: []string{
				"/{usr/,}libexec/geoclue",
				"/{usr/,}libexec/geoclue-2.0/demos/agent",
			},
		},
		{
			name: "opera",
			variables: map[string][]string{
				"multiarch":         {"*-linux-gnu*"},
				"chromium_name":     {"opera{,-beta,-developer}"},
				"chromium_lib_dirs": {"/{usr/,}lib/@{multiarch}/@{chromium_name}"},
				"exec_path":         {"@{chromium_lib_dirs}/@{chromium_name}"},
			},
			want: []string{
				"/{usr/,}lib/*-linux-gnu*/opera{,-beta,-developer}/opera{,-beta,-developer}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &AppArmorProfile{
				Content:     "",
				Variables:   tt.variables,
				Attachments: []string{},
			}
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
			p := &AppArmorProfile{
				Content:     "",
				Variables:   map[string][]string{},
				Attachments: tt.Attachments,
			}
			if got := p.NestAttachments(); got != tt.want {
				t.Errorf("AppArmorProfile.NestAttachments() = %v, want %v", got, tt.want)
			}
		})
	}
}
