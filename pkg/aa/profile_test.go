// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func readprofile(path string) string {
	file := paths.New("../../").Join(path)
	lines, err := file.ReadFileAsLines()
	if err != nil {
		panic(err)
	}
	res := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			continue
		}
		res += line + "\n"
	}
	return res[:len(res)-1]
}

func TestAppArmorProfile_String(t *testing.T) {
	tests := []struct {
		name string
		p    *AppArmorProfile
		want string
	}{
		{
			name: "empty",
			p:    &AppArmorProfile{},
			want: ``,
		},
		{
			name: "foo",
			p: &AppArmorProfile{
				Preamble: Preamble{
					Abi:      []Abi{{IsMagic: true, Path: "abi/4.0"}},
					Includes: []Include{{IsMagic: true, Path: "tunables/global"}},
					Aliases:  []Alias{{Path: "/mnt/usr", RewrittenPath: "/usr"}},
					Variables: []Variable{{
						Name:   "exec_path",
						Values: []string{"@{bin}/foo", "@{lib}/foo"},
					}},
				},
				Profile: Profile{
					Name:        "foo",
					Attachments: []string{"@{exec_path}"},
					Attributes:  map[string]string{"security.tagged": "allowed"},
					Flags:       []string{"complain", "attach_disconnected"},
					Rules: []ApparmorRule{
						&Include{IsMagic: true, Path: "abstractions/base"},
						&Include{IsMagic: true, Path: "abstractions/nameservice-strict"},
						rlimit1,
						&Capability{Name: "dac_read_search"},
						&Capability{Name: "dac_override"},
						&Network{Domain: "inet", Type: "stream"},
						&Network{Domain: "inet6", Type: "stream"},
						&Mount{
							MountConditions: MountConditions{
								FsType:  "fuse.portal",
								Options: []string{"rw", "rbind"},
							},
							Source:     "@{run}/user/@{uid}/ ",
							MountPoint: "/",
						},
						&Umount{
							MountConditions: MountConditions{},
							MountPoint:      "@{run}/user/@{uid}/",
						},
						&Signal{
							Access: "receive",
							Set:    "term",
							Peer:   "at-spi-bus-launcher",
						},
						&Ptrace{Access: "read", Peer: "nautilus"},
						&Unix{
							Access:   "send receive",
							Type:     "stream",
							Address:  "@/tmp/.ICE-unix/1995",
							Peer:     "gnome-shell",
							PeerAddr: "none",
						},
						&Dbus{
							Access: "bind",
							Bus:    "session",
							Name:   "org.gnome.*",
						},
						&Dbus{
							Access:    "receive",
							Bus:       "system",
							Name:      ":1.3",
							Path:      "/org/freedesktop/DBus",
							Interface: "org.freedesktop.DBus",
							Member:    "AddMatch",
							Label:     "power-profiles-daemon",
						},
						&File{Path: "/opt/intel/oneapi/compiler/*/linux/lib/*.so./*", Access: "rm"},
						&File{Path: "@{PROC}/@{pid}/task/@{tid}/comm", Access: "rw"},
						&File{Path: "@{sys}/devices/@{pci}/class", Access: "r"},
						includeLocal1,
					},
				},
			},
			want: readprofile("tests/string.aa"),
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
					Rules: []ApparmorRule{
						file2, network1, includeLocal1, dbus2, signal1, ptrace1,
						capability2, file1, dbus1, unix2, signal2, mount2,
					},
				},
			},
			want: &AppArmorProfile{
				Profile: Profile{
					Rules: []ApparmorRule{
						capability2, network1, mount2, signal1, signal2, ptrace1,
						unix2, dbus2, dbus1, file1, file2, includeLocal1,
					},
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

func TestAppArmorProfile_Integration(t *testing.T) {
	tests := []struct {
		name string
		p    *AppArmorProfile
		want string
	}{
		{
			name: "aa-status",
			p: &AppArmorProfile{
				Preamble: Preamble{
					Abi:      []Abi{{IsMagic: true, Path: "abi/3.0"}},
					Includes: []Include{{IsMagic: true, Path: "tunables/global"}},
					Variables: []Variable{{
						Name:   "exec_path",
						Values: []string{"@{bin}/aa-status", "@{bin}/apparmor_status"},
					}},
				},
				Profile: Profile{
					Name:        "aa-status",
					Attachments: []string{"@{exec_path}"},
					Rules: Rules{
						&Include{IfExists: true, IsMagic: true, Path: "local/aa-status"},
						&Capability{Name: "dac_read_search"},
						&File{Path: "@{exec_path}", Access: "mr"},
						&File{Path: "@{PROC}/@{pids}/attr/apparmor/current", Access: "r"},
						&File{Path: "@{PROC}/", Access: "r"},
						&File{Path: "@{sys}/module/apparmor/parameters/enabled", Access: "r"},
						&File{Path: "@{sys}/kernel/security/apparmor/profiles", Access: "r"},
						&File{Path: "@{PROC}/@{pids}/attr/current", Access: "r"},
						&Include{IsMagic: true, Path: "abstractions/consoles"},
						&File{Qualifier: Qualifier{Owner: true}, Path: "@{PROC}/@{pid}/mounts", Access: "r"},
						&Include{IsMagic: true, Path: "abstractions/base"},
						&File{Path: "/dev/tty@{int}", Access: "rw"},
						&Capability{Name: "sys_ptrace"},
						&Ptrace{Access: "read"},
					},
				},
			},
			want: readprofile("apparmor.d/profiles-a-f/aa-status"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Sort()
			tt.p.MergeRules()
			tt.p.Format()
			if got := tt.p.String(); "\n"+got != tt.want {
				t.Errorf("AppArmorProfile = |%v|, want |%v|", "\n"+got, tt.want)
			}
		})
	}
}
