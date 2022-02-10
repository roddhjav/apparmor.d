// aa-log - Review AppArmor generated messages
// Copyright (C) 2021 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"os"
	"reflect"
	"testing"
)

var refDnsmasq = AppArmorLogs{
	{
		"apparmor":       "DENIED",
		"profile":        "dnsmasq",
		"operation":      "open",
		"name":           "/proc/sys/kernel/osrelease",
		"comm":           "dnsmasq",
		"requested_mask": "r",
		"denied_mask":    "r",
	},
	{
		"apparmor":       "DENIED",
		"profile":        "dnsmasq",
		"operation":      "open",
		"name":           "/proc/1/environ",
		"comm":           "dnsmasq",
		"requested_mask": "r",
		"denied_mask":    "r",
	},
	{
		"apparmor":       "DENIED",
		"profile":        "dnsmasq",
		"operation":      "open",
		"name":           "/proc/cmdline",
		"comm":           "dnsmasq",
		"requested_mask": "r",
		"denied_mask":    "r",
	},
}

var refKmod = AppArmorLogs{
	{
		"apparmor":       "ALLOWED",
		"profile":        "kmod",
		"operation":      "file_inherit",
		"comm":           "modprobe",
		"family":         "unix",
		"sock_type":      "stream",
		"protocol":       "0",
		"requested_mask": "send receive",
	},
}

var refMan = AppArmorLogs{
	{
		"apparmor":       "ALLOWED",
		"profile":        "man",
		"operation":      "exec",
		"name":           "/usr/bin/preconv",
		"info":           "no new privs",
		"comm":           "man",
		"requested_mask": "x",
		"denied_mask":    "x",
		"error":          "-1",
	},
}

var refStringKmod = "\033[1;32mALLOWED\033[0m \033[34mkmod\033[0m \033[33mfile_inherit\033[0m comm=modprobe family=unix sock_type=stream protocol=0 requested_mask=\033[1;31m\"send receive\"\033[0m\n"
var refStringMan = "\033[1;32mALLOWED\033[0m \033[34mman\033[0m \033[33mexec\033[0m \033[35m/usr/bin/preconv\033[0m info=\"no new privs\" comm=man requested_mask=\033[1;31mx\033[0m denied_mask=\033[1;31mx\033[0m error=-1\n"

func TestNewApparmorLogs(t *testing.T) {
	tests := []struct {
		name string
		path string
		want AppArmorLogs
	}{
		{
			name: "dnsmasq",
			path: "../../tests/audit.log",
			want: refDnsmasq,
		},
		{
			name: "kmod",
			path: "../../tests/audit.log",
			want: refKmod,
		},
		{
			name: "man",
			path: "../../tests/audit.log",
			want: refMan,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, _ := os.Open(tt.path)
			if got := NewApparmorLogs(file, tt.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewApparmorLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorLogs_String(t *testing.T) {
	tests := []struct {
		name   string
		aaLogs AppArmorLogs
		want   string
	}{
		{
			name:   "kmod",
			aaLogs: refKmod,
			want:   refStringKmod,
		},
		{
			name:   "man",
			aaLogs: refMan,
			want:   refStringMan,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.aaLogs.String(); got != tt.want {
				t.Errorf("AppArmorLogs.String() = %v, want %v len: %d - %d", got, tt.want, len(got), len(tt.want))
			}
		})
	}
}

func Test_app(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		profile string
		wantErr bool
	}{
		{
			name:    "OK",
			path:    "../../tests/audit.log",
			profile: "",
			wantErr: false,
		},
		{
			name:    "No logfile",
			path:    "../../tests/log",
			profile: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := aaLog(tt.path, tt.profile); (err != nil) != tt.wantErr {
				t.Errorf("aaLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
