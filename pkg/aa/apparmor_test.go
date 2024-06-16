// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	testData = paths.New("../../tests/testdata/")
	intData  = paths.New("../../apparmor.d")
)

// mustReadProfileFile read a file and return its content as a slice of string.
// It panics if an error occurs. It removes the last comment line.
func mustReadProfileFile(path *paths.Path) string {
	res := strings.Split(util.MustReadFile(path), "\n")
	return strings.Join(res[:len(res)-2], "\n")
}

func TestAppArmorProfileFile_String(t *testing.T) {
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
				Preamble: Rules{
					&Comment{RuleBase: RuleBase{Comment: " Simple test profile for the AppArmorProfileFile.String() method", IsLineRule: true}},
					nil,
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
					&Include{IsMagic: true, Path: "tunables/global"},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{"@{bin}/foo", "@{lib}/foo"},
					},
				},
				Profiles: []*Profile{{
					Header: Header{
						Name:        "foo",
						Attachments: []string{"@{exec_path}"},
						Attributes:  map[string]string{"security.tagged": "allowed"},
						Flags:       []string{"complain", "attach_disconnected"},
					},
					Rules: Rules{
						&Include{IsMagic: true, Path: "abstractions/base"},
						&Include{IsMagic: true, Path: "abstractions/nameservice-strict"},
						rlimit1,
						&Capability{Names: []string{"dac_read_search"}},
						&Capability{Names: []string{"dac_override"}},
						&Network{Domain: "inet", Type: "stream"},
						&Network{Domain: "inet6", Type: "stream"},
						&Mount{
							RuleBase: RuleBase{Comment: " failed perms check"},
							MountConditions: MountConditions{
								FsType:  "fuse.portal",
								Options: []string{"rw", "rbind"},
							},
							Source:     "@{run}/user/@{uid}/",
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
						&Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
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
			want: util.MustReadFile(testData.Join("string.aa")),
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

func TestAppArmorProfileFile_Sort(t *testing.T) {
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
						file2, network1, userns1, include1, dbus2, signal1,
						ptrace1, includeLocal1, rlimit3, capability1, network2,
						mqueue2, iouring2, dbus1, link2, capability2, file1,
						unix2, signal2, mount2, all1, umount2, mount1, remount2,
						pivotroot1, changeprofile2,
					},
				}},
			},
			want: &AppArmorProfileFile{
				Profiles: []*Profile{{
					Rules: []Rule{
						include1, all1, rlimit3, userns1, capability1, capability2,
						network2, network1, mount2, mount1, remount2, umount2,
						pivotroot1, changeprofile2, mqueue2, iouring2, signal2,
						signal1, ptrace1, unix2, dbus2, dbus1, file1, file2,
						link2, includeLocal1,
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

func TestAppArmorProfileFile_MergeRules(t *testing.T) {
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

func TestAppArmorProfileFile_Integration(t *testing.T) {
	tests := []struct {
		name string
		f    *AppArmorProfileFile
		want string
	}{
		{
			name: "aa-status",
			f: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{RuleBase: RuleBase{Comment: " apparmor.d - Full set of apparmor profiles", IsLineRule: true}},
					&Comment{RuleBase: RuleBase{Comment: " Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>", IsLineRule: true}},
					&Comment{RuleBase: RuleBase{Comment: " SPDX-License-Identifier: GPL-2.0-only", IsLineRule: true}},
					nil,
					&Abi{IsMagic: true, Path: "abi/3.0"},
					&Include{IsMagic: true, Path: "tunables/global"},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{"@{bin}/aa-status", "@{bin}/apparmor_status"},
					},
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
			want: mustReadProfileFile(intData.Join("profiles-a-f/aa-status")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.f.Sort()
			tt.f.MergeRules()
			tt.f.Format()
			err := tt.f.Validate()
			if err != nil {
				t.Errorf("AppArmorProfile.Validate() = %v", err)
			}
			if got := tt.f.String(); got != tt.want {
				t.Errorf("AppArmorProfile = |%v|, want |%v|", got, tt.want)
			}
		})
	}
}
