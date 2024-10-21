// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

func TestFilterOnly_Apply(t *testing.T) {
	tests := []struct {
		name    string
		dist    string
		family  string
		opt     *Option
		profile string
		want    string
		wantErr bool
	}{
		{
			name:   "inline",
			dist:   "debian",
			family: "apt",
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{"apt": ""},
				ArgList: []string{"apt"},
				File:    nil,
				Raw:     "  @{bin}/arch-audit rPx, #aa:only apt",
			},
			profile: "  @{bin}/arch-audit rPx, #aa:only apt",
			want:    "  @{bin}/arch-audit rPx,",
		},
		{
			name:   "paragraph",
			dist:   "arch",
			family: "pacman",
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{"zypper": ""},
				ArgList: []string{"zypper"},
				File:    nil,
				Raw:     "  #aa:only zypper",
			},
			profile: `
        /tmp/apt-changelog-@{rand6}/ w,
        /tmp/apt-changelog-@{rand6}/*.changelog rw,
  owner /tmp/alpm_*/{,**} rw,
  owner /tmp/apt-changelog-@{rand6}/.apt-acquire-privs-test.@{rand6} rw,
  owner /tmp/packagekit* rw,

        @{run}/systemd/inhibit/*.ref rw,
  owner @{run}/systemd/users/@{uid} r,

  #aa:only zypper
        @{run}/zypp.pid rwk,
  owner @{run}/zypp-rpm.pid rwk,
  owner @{run}/zypp/packages/ r,

  owner /dev/shm/AP_0x@{rand6}/{,**} rw,
  owner /dev/shm/ r,`,
			want: `
        /tmp/apt-changelog-@{rand6}/ w,
        /tmp/apt-changelog-@{rand6}/*.changelog rw,
  owner /tmp/alpm_*/{,**} rw,
  owner /tmp/apt-changelog-@{rand6}/.apt-acquire-privs-test.@{rand6} rw,
  owner /tmp/packagekit* rw,

        @{run}/systemd/inhibit/*.ref rw,
  owner @{run}/systemd/users/@{uid} r,

  owner /dev/shm/AP_0x@{rand6}/{,**} rw,
  owner /dev/shm/ r,`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prebuild.Distribution = tt.dist
			prebuild.Family = tt.family
			got, err := Directives["only"].Apply(tt.opt, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterOnly.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FilterOnly.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterExclude_Apply(t *testing.T) {
	tests := []struct {
		name    string
		dist    string
		family  string
		opt     *Option
		profile string
		want    string
		wantErr bool
	}{
		{
			name:   "inline",
			dist:   "debian",
			family: "apt",
			opt: &Option{
				Name:    "exclude",
				ArgMap:  map[string]string{"debian": ""},
				ArgList: []string{"debian"},
				File:    nil,
				Raw:     "  @{bin}/dpkg rPx -> child-dpkg, #aa:exclude debian",
			},
			profile: "  @{bin}/dpkg rPx -> child-dpkg, #aa:exclude debian",
			want:    "",
		},
		{
			name:   "inline-keep",
			dist:   "whonix",
			family: "apt",
			opt: &Option{
				Name:    "exclude",
				ArgMap:  map[string]string{"debian": ""},
				ArgList: []string{"debian"},
				File:    nil,
				Raw:     "  @{bin}/dpkg rPx -> child-dpkg, #aa:exclude debian",
			},
			profile: "  @{bin}/dpkg rPx -> child-dpkg, #aa:exclude debian",
			want:    "  @{bin}/dpkg rPx -> child-dpkg,",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prebuild.Distribution = tt.dist
			prebuild.Family = tt.family
			got, err := Directives["exclude"].Apply(tt.opt, tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterExclude.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FilterExclude.Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}
