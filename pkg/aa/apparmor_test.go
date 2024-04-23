// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"strings"
	"testing"

	"github.com/arduino/go-paths-helper"
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
		f    *AppArmorProfileFile
		want string
	}{
		{
			name: "empty",
			f:    &AppArmorProfileFile{},
			want: ``,
		},
		{
			name: "foo",
			f: &AppArmorProfileFile{
				Preamble: Preamble{
					Abi:      []*Abi{{IsMagic: true, Path: "abi/4.0"}},
					Includes: []*Include{{IsMagic: true, Path: "tunables/global"}},
					Aliases:  []*Alias{{Path: "/mnt/usr", RewrittenPath: "/usr"}},
					Variables: []*Variable{{
						Name: "exec_path", Define: true,
						Values: []string{"@{bin}/foo", "@{lib}/foo"},
					}},
				},
				Profiles: []*Profile{{
					Header: Header{
						Name:        "foo",
						Attachments: []string{"@{exec_path}"},
						Attributes:  map[string]string{"security.tagged": "allowed"},
						Flags:       []string{"complain", "attach_disconnected"},
					},
					Rules: []Rule{
						&Include{IsMagic: true, Path: "abstractions/base"},
						&Include{IsMagic: true, Path: "abstractions/nameservice-strict"},
						rlimit1,
						&Capability{Names: []string{"dac_read_search"}},
						&Capability{Names: []string{"dac_override"}},
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
							Access: []string{"receive"},
							Set:    []string{"term"},
							Peer:   "at-spi-bus-launcher",
						},
						&Ptrace{Access: []string{"read"}, Peer: "nautilus"},
						&Unix{
							Access:    []string{"send", "receive"},
							Type:      "stream",
							Address:   "@/tmp/.ICE-unix/1995",
							PeerLabel: "gnome-shell",
							PeerAddr:  "none",
						},
						&Dbus{
							Access: []string{"bind"},
							Bus:    "session",
							Name:   "org.gnome.*",
						},
						&Dbus{
							Access:    []string{"receive"},
							Bus:       "system",
							Path:      "/org/freedesktop/DBus",
							Interface: "org.freedesktop.DBus",
							Member:    "AddMatch",
							PeerName:  ":1.3",
							PeerLabel: "power-profiles-daemon",
						},
						&File{Path: "/opt/intel/oneapi/compiler/*/linux/lib/*.so./*", Access: []string{"r", "m"}},
						&File{Path: "@{PROC}/@{pid}/task/@{tid}/comm", Access: []string{"r", "w"}},
						&File{Path: "@{sys}/devices/@{pci}/class", Access: []string{"r"}},
						includeLocal1,
					},
				}},
			},
			want: readprofile("tests/string.aa"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.String(); got != tt.want {
				t.Errorf("AppArmorProfile.String() = |%v|, want |%v|", got, tt.want)
			}
		})
	}
}

func TestAppArmorProfile_AddRule(t *testing.T) {
	tests := []struct {
		name string
		log  map[string]string
		want *AppArmorProfileFile
	}{
		{
			name: "capability",
			log:  capability1Log,
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{capability1},
				}},
			},
		},
		{
			name: "network",
			log:  network1Log,
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{network1},
				}},
			},
		},
		{
			name: "mount",
			log:  mount2Log,
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{mount2},
				}},
			},
		},
		{
			name: "signal",
			log:  signal1Log,
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{signal1},
				}},
			},
		},
		{
			name: "ptrace",
			log:  ptrace2Log,
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{ptrace2},
				}},
			},
		},
		{
			name: "unix",
			log:  unix1Log,
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{unix1},
				}},
			},
		},
		{
			name: "dbus",
			log:  dbus2Log,
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{dbus2},
				}},
			},
		},
		{
			name: "file",
			log:  file2Log,
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{file2},
				}},
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
		origin *AppArmorProfileFile
		want   *AppArmorProfileFile
	}{
		{
			name: "all",
			origin: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{
						file2, network1, includeLocal1, dbus2, signal1, ptrace1,
						capability2, file1, dbus1, unix2, signal2, mount2,
					},
				}},
			},
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{
						capability2, network1, mount2, signal1, signal2, ptrace1,
						unix2, dbus2, dbus1, file1, file2, includeLocal1,
					},
				}},
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
		origin *AppArmorProfileFile
		want   *AppArmorProfileFile
	}{
		{
			name: "all",
			origin: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{capability1, capability1, network1, network1, file1, file1},
				}},
			},
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{capability1, network1, file1},
				}},
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
		f    *AppArmorProfileFile
		want string
	}{
		{
			name: "aa-status",
			f: &AppArmorProfileFile{
				Preamble: Preamble{
					Abi:      []*Abi{{IsMagic: true, Path: "abi/3.0"}},
					Includes: []*Include{{IsMagic: true, Path: "tunables/global"}},
					Variables: []*Variable{{
						Name:   "exec_path",
						Values: []string{"@{bin}/aa-status", "@{bin}/apparmor_status"},
					}},
				},
				Profiles: []*Profile{{
					Header: Header{
						Name:        "aa-status",
						Attachments: []string{"@{exec_path}"},
					},
					Rules: Rules{
						&Include{IfExists: true, IsMagic: true, Path: "local/aa-status"},
						&Capability{Names: []string{"dac_read_search"}},
						&File{Path: "@{exec_path}", Access: []string{"m", "r"}},
						&File{Path: "@{PROC}/@{pids}/attr/apparmor/current", Access: []string{"r"}},
						&File{Path: "@{PROC}/", Access: []string{"r"}},
						&File{Path: "@{sys}/module/apparmor/parameters/enabled", Access: []string{"r"}},
						&File{Path: "@{sys}/kernel/security/apparmor/profiles", Access: []string{"r"}},
						&File{Path: "@{PROC}/@{pids}/attr/current", Access: []string{"r"}},
						&Include{IsMagic: true, Path: "abstractions/consoles"},
						&File{Owner: true, Path: "@{PROC}/@{pid}/mounts", Access: []string{"r"}},
						&Include{IsMagic: true, Path: "abstractions/base"},
						&File{Path: "/dev/tty@{int}", Access: []string{"r", "w"}},
						&Capability{Names: []string{"sys_ptrace"}},
						&Ptrace{Access: []string{"read"}},
					},
				}},
			},
			want: readprofile("apparmor.d/profiles-a-f/aa-status"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Sort()
			tt.f.MergeRules()
			tt.f.Format()
			if got := tt.f.String(); "\n"+got != tt.want {
				t.Errorf("AppArmorProfile = |%v|, want |%v|", "\n"+got, tt.want)
			}
		})
	}
}
