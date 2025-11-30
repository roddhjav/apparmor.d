// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

func Test_cmp(t *testing.T) {
	tests := []struct {
		name     string
		operator string
		refValue float64
		value    float64
		want     bool
	}{
		{
			name:     "3.2 < 5.0",
			operator: "<",
			refValue: 3.2,
			value:    5.0,
			want:     true,
		},
		{
			name:     "5.0 == 5.0",
			operator: "==",
			refValue: 5.0,
			value:    5.0,
			want:     true,
		},
		{
			name:     "5.0 >= 4.1",
			operator: ">=",
			refValue: 5.0,
			value:    4.1,
			want:     true,
		},
		{
			name:     "3.2 < 5.0",
			operator: "==",
			refValue: 3.2,
			value:    5.0,
			want:     false,
		},
		{
			name:     "3.2 <= 5.0",
			operator: "<=",
			refValue: 3.2,
			value:    5.0,
			want:     true,
		},
		{
			name:     "4.1 >= 4.1",
			operator: ">=",
			refValue: 4.1,
			value:    4.1,
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cmp(tt.refValue, tt.operator, tt.value)
			if got != tt.want {
				t.Errorf("cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_compare(t *testing.T) {
	tests := []struct {
		name     string
		refValue float64
		arg      string
		want     bool
	}{
		{
			name:     "3.1 < 4.0",
			refValue: 3.1,
			arg:      "apparmor<4.0",
			want:     true,
		},
		{
			name:     "3.2 < 5.0",
			refValue: 3.2,
			arg:      "apparmor<5.0",
			want:     true,
		},
		{
			name:     "5.0 == 5.0",
			refValue: 5.0,
			arg:      "apparmor==5.0",
			want:     true,
		},
		{
			name:     "5.0 >= 4.1",
			refValue: 5.0,
			arg:      "apparmor>=4.1",
			want:     true,
		},
		{
			name:     "3.2 == 5.0",
			refValue: 3.2,
			arg:      "apparmor==5.0",
			want:     false,
		},
		{
			name:     "3.2 <= 5.0",
			refValue: 3.2,
			arg:      "apparmor<=5.0",
			want:     true,
		},
		{
			name:     "4.1 >= 4.1",
			refValue: 4.1,
			arg:      "apparmor>=4.1",
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := compare(tt.refValue, "apparmor", tt.arg)
			if got != tt.want {
				t.Errorf("compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{
			name:   "inline-exclude",
			dist:   "ubuntu",
			family: "apt",
			opt: &Option{
				Name:    "exclude",
				ArgMap:  map[string]string{"ubuntu": ""},
				ArgList: []string{"ubuntu"},
				File:    nil,
				Raw:     "  @{bin}/dpkg rPx -> child-dpkg, #aa:exclude ubuntu",
			},
			profile: "  @{bin}/dpkg rPx -> child-dpkg, #aa:exclude ubuntu",
			want:    "",
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

func TestFilterCmp_Apply(t *testing.T) {
	tests := []struct {
		name    string
		abi     int
		version float64
		opt     *Option
		want    string
		wantErr bool
	}{
		{
			name:    "apparmor3.1<4.0",
			version: 3.1,
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{},
				ArgList: []string{"apparmor<4.0"},
				File:    nil,
				Raw:     "  /dev/shm/ r, #aa:only apparmor>=4.1",
			},
			want: "  /dev/shm/ r,",
		},
		{
			name:    "apparmor5.0>=4.1",
			version: 5.0,
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{},
				ArgList: []string{"apparmor>=4.1"},
				File:    nil,
				Raw:     "  priority=100 @{bin}/bwrap Px, #aa:only apparmor>=4.1",
			},
			want: "  priority=100 @{bin}/bwrap Px,",
		},
		{
			name:    "apparmor4.1>=4.1",
			version: 4.1,
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{},
				ArgList: []string{"apparmor>=4.1"},
				File:    nil,
				Raw:     "  priority=100 @{bin}/bwrap Px, #aa:only apparmor>=4.1",
			},
			want: "  priority=100 @{bin}/bwrap Px,",
		},
		{
			name:    "apparmor4.0>=4.1",
			version: 4.0,
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{},
				ArgList: []string{"apparmor>=4.1"},
				File:    nil,
				Raw:     "  priority=100 @{bin}/bwrap Px, #aa:only apparmor>=4.1",
			},
			want: "",
		},
		{
			name: "abi3=3",
			abi:  3,
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{},
				ArgList: []string{"abi==3"},
				File:    nil,
				Raw:     "  /dev/shm/ r, #aa:only abi==3",
			},
			want: "  /dev/shm/ r,",
		},
		{
			name: "abi 3<4",
			abi:  3,
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{},
				ArgList: []string{"abi<4"},
				File:    nil,
				Raw:     "  /efi/ r, #aa:only abi<4",
			},
			want: "  /efi/ r,",
		},
		{
			name: "abi 3>=5",
			abi:  3,
			opt: &Option{
				Name:    "only",
				ArgMap:  map[string]string{},
				ArgList: []string{"abi>=5"},
				File:    nil,
				Raw:     "  /dev/shm/ r, #aa:only abi>=5",
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prebuild.Version = tt.version
			prebuild.ABI = tt.abi
			got, err := Directives["only"].Apply(tt.opt, tt.opt.Raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterOnly.Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FilterOnly.Apply() = |%v|, want |%v|", got, tt.want)
			}
		})
	}
}
