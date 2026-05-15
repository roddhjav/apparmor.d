// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"reflect"
	"testing"
)

func TestGetFlags(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		want    []string
	}{
		{
			name:    "no flags",
			profile: "profile foo /usr/bin/foo {\n",
			want:    nil,
		},
		{
			name:    "single flag",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n",
			want:    []string{"complain"},
		},
		{
			name:    "multiple flags",
			profile: "profile foo /usr/bin/foo flags=(attach_disconnected,complain) {\n",
			want:    []string{"attach_disconnected", "complain"},
		},
		{
			name:    "empty profile",
			profile: "",
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetFlags(tt.profile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetFlags(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		flags   []string
		want    string
	}{
		{
			name:    "add flags to profile without flags",
			profile: "profile foo /usr/bin/foo {\n",
			flags:   []string{"complain"},
			want:    "profile foo /usr/bin/foo flags=(complain) {\n",
		},
		{
			name:    "add flags to profile with if statement",
			profile: "profile foo /usr/bin/foo {\n  if true {\n    /bin/true rix,\n  }\n}\n",
			flags:   []string{"complain"},
			want:    "profile foo /usr/bin/foo flags=(complain) {\n  if true {\n    /bin/true rix,\n  }\n}\n",
		},
		{
			name:    "add multiple flags",
			profile: "profile foo /usr/bin/foo {\n",
			flags:   []string{"attach_disconnected", "complain"},
			want:    "profile foo /usr/bin/foo flags=(attach_disconnected,complain) {\n",
		},
		{
			name:    "replace existing flags",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n",
			flags:   []string{"enforce"},
			want:    "profile foo /usr/bin/foo flags=(enforce) {\n",
		},
		{
			name:    "remove flags with empty slice",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n",
			flags:   []string{},
			want:    "profile foo /usr/bin/foo {\n",
		},
		{
			name:    "remove flags with nil",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n",
			flags:   nil,
			want:    "profile foo /usr/bin/foo {\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetFlags(tt.profile, tt.flags); got != tt.want {
				t.Errorf("SetFlags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsUnconfined(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		want    bool
	}{
		{
			name:    "no flags",
			profile: "profile foo /usr/bin/foo {\n}\n",
			want:    false,
		},
		{
			name:    "unconfined only",
			profile: "profile foo /usr/bin/foo flags=(unconfined) {\n}\n",
			want:    true,
		},
		{
			name:    "unconfined with other flags",
			profile: "profile foo /usr/bin/foo flags=(attach_disconnected, unconfined) {\n}\n",
			want:    true,
		},
		{
			name:    "complain only",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n}\n",
			want:    false,
		},
		{
			name:    "substring should not match",
			profile: "profile foo /usr/bin/foo flags=(unconfinedx) {\n}\n",
			want:    false,
		},
		{
			name:    "second profile is unconfined",
			profile: "profile a /usr/bin/a {\n}\nprofile b /usr/bin/b flags=(unconfined) {\n}\n",
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsUnconfined(tt.profile); got != tt.want {
				t.Errorf("IsUnconfined() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetMode(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		mode    string
		want    string
		wantErr bool
	}{
		{
			name:    "set complain mode",
			profile: "profile foo /usr/bin/foo {\n",
			mode:    "complain",
			want:    "profile foo /usr/bin/foo flags=(complain) {\n",
		},
		{
			name:    "set enforce mode removes mode flags",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n",
			mode:    "enforce",
			want:    "profile foo /usr/bin/foo {\n",
		},
		{
			name:    "replace complain with kill",
			profile: "profile foo /usr/bin/foo flags=(complain) {\n",
			mode:    "kill",
			want:    "profile foo /usr/bin/foo flags=(kill) {\n",
		},
		{
			name:    "preserve non-mode flags when setting mode",
			profile: "profile foo /usr/bin/foo flags=(attach_disconnected,complain) {\n",
			mode:    "enforce",
			want:    "profile foo /usr/bin/foo flags=(attach_disconnected) {\n",
		},
		{
			name:    "preserve non-mode flags when changing mode",
			profile: "profile foo /usr/bin/foo flags=(attach_disconnected,complain) {\n",
			mode:    "kill",
			want:    "profile foo /usr/bin/foo flags=(attach_disconnected,kill) {\n",
		},
		{
			name:    "unknown mode returns error",
			profile: "profile foo /usr/bin/foo {\n",
			mode:    "invalid",
			want:    "profile foo /usr/bin/foo {\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SetMode(tt.profile, tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetMode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SetMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
