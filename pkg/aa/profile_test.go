// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"testing"
)

func TestAppArmorProfile_String(t *testing.T) {
	tests := []struct {
		name string
		p    AppArmorProfile
		want string
	}{
		{
			name: "empty",
			p:    AppArmorProfile{},
			want: `profile  {
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.String(); got != tt.want {
				t.Errorf("AppArmorProfile.String() = |%v|, want |%v|", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_AddRule(t *testing.T) {
	tests := []struct {
		name string
		log  map[string]string
		want *AppArmorProfile
	}{
		{
			name: "capability",
			log:  capability1Log,
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{capability1},
				},
			},
		},
		{
			name: "network",
			log:  network1Log,
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{network1},
				},
			},
		},
		{
			name: "mount",
			log:  mount2Log,
			want: &AppArmorProfile{
				Profile: Profile{
					Flags: []string{"attach_disconnected"},
					Rules: []ApparmorRule{mount2},
				},
			},
		},
		{
			name: "signal",
			log:  signal1Log,
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{signal1},
				},
			},
		},
		{
			name: "ptrace",
			log:  ptrace2Log,
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{ptrace2},
				},
			},
		},
		{
			name: "unix",
			log:  unix1Log,
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{unix1},
				},
			},
		},
		{
			name: "dbus",
			log:  dbus2Log,
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{dbus2},
				},
			},
		},
		{
			name: "file",
			log:  file2Log,
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{file2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAppArmorProfile()
			got.AddRule(tt.log)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfile.AddRule() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_Sort(t *testing.T) {
	tests := []struct {
		name   string
		origin *AppArmorProfile
		want   *AppArmorProfile
	}{
		{
			name: "all",
			origin: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{file2, network1, dbus2, signal1, ptrace1, capability2, file1, dbus1, unix2, signal2, mount2},
				},
			},
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{capability2, network1, mount2, signal1, signal2, ptrace1, unix2, dbus2, dbus1, file2, file1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.origin
			got.Sort()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfile.Sort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_MergeRules(t *testing.T) {
	tests := []struct {
		name   string
		origin *AppArmorProfile
		want   *AppArmorProfile
	}{
		{
			name: "all",
			origin: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{capability1, capability1, network1, network1, file1, file1},
				},
			},
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{capability1, network1, file1},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.origin
			got.MergeRules()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfile.MergeRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
