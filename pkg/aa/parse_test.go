// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"reflect"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/util"
)

func Test_tokenizeRule(t *testing.T) {
	for _, tt := range testRules {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenize(tt.raw); !reflect.DeepEqual(got, tt.tokens) {
				t.Errorf("tokenize() = %v, want %v", got, tt.tokens)
			}
		})
	}
}

func Test_AppArmorProfileFile_Parse(t *testing.T) {
	for _, tt := range testBlocks {
		t.Run(tt.name, func(t *testing.T) {
			got := &AppArmorProfileFile{}
			if err := got.Parse(tt.raw); (err != nil) != tt.wParseErr {
				t.Errorf("AppArmorProfileFile.Parse() error = %v, wantErr %v", err, tt.wParseErr)
			}
			if !reflect.DeepEqual(got, tt.apparmor) {
				t.Errorf("AppArmorProfileFile.Parse() = |%v|, want |%v|", got, tt.apparmor)
			}
		})
	}
}

var (
	// Test cases for tokenize
	testRules = []struct {
		name   string
		raw    string
		tokens []string
	}{
		{
			name:   "empty",
			raw:    "",
			tokens: []string{},
		},
		{
			name:   "abi",
			raw:    `abi <abi/4.0>`,
			tokens: []string{"abi", "<abi/4.0>"},
		},
		{
			name:   "alias",
			raw:    `alias /mnt/usr -> /usr`,
			tokens: []string{"alias", "/mnt/usr", "->", "/usr"},
		},
		{
			name:   "variable",
			raw:    `@{name} = torbrowser "tor browser"`,
			tokens: []string{"@{name}", "=", "torbrowser", `"tor browser"`},
		},
		{
			name:   "variable-2",
			raw:    `@{exec_path} += @{bin}/@{name}`,
			tokens: []string{"@{exec_path}", "+", "=", "@{bin}/@{name}"},
		},
		{
			name:   "variable-3",
			raw:    `@{empty}="dummy"`,
			tokens: []string{"@{empty}", "=", `"dummy"`},
		},
		{
			name:   "variable-4",
			raw:    `@{XDG_PROJECTS_DIR}+="Git"`,
			tokens: []string{"@{XDG_PROJECTS_DIR}", "+", "=", `"Git"`},
		},
		{
			name:   "header",
			raw:    `profile foo @{exec_path} xattrs=(security.tagged=allowed) flags=(complain attach_disconnected)`,
			tokens: []string{"profile", "foo", "@{exec_path}", "xattrs=(security.tagged=allowed)", "flags=(complain attach_disconnected)"},
		},
		{
			name:   "include",
			raw:    `include <tunables/global>`,
			tokens: []string{"include", "<tunables/global>"},
		},
		{
			name:   "include-if-exists",
			raw:    `include if exists "/etc/apparmor.d/dummy"`,
			tokens: []string{"include", "if", "exists", `"/etc/apparmor.d/dummy"`},
		},
		{
			name:   "rlimit",
			raw:    `set rlimit nproc <= 200`,
			tokens: []string{"set", "rlimit", "nproc", "<=", "200"},
		},
		{
			name:   "userns",
			raw:    `userns`,
			tokens: []string{"userns"},
		},
		{
			name:   "capability",
			raw:    `capability dac_read_search`,
			tokens: []string{"capability", "dac_read_search"},
		},
		{
			name:   "network",
			raw:    `network netlink raw`,
			tokens: []string{"network", "netlink", "raw"},
		},
		{
			name:   "mount",
			raw:    `mount /{,**}`,
			tokens: []string{"mount", "/{,**}"},
		},
		{
			name:   "mount-2",
			raw:    `mount               options=(rw rbind)                 /tmp/newroot/ -> /tmp/newroot/`,
			tokens: []string{"mount", "options=(rw rbind)", "/tmp/newroot/", "->", "/tmp/newroot/"},
		},
		{
			name:   "mount-3",
			raw:    `mount               options=(rw silent rprivate)                     -> /oldroot/`,
			tokens: []string{"mount", "options=(rw silent rprivate)", "->", "/oldroot/"},
		},
		{
			name:   "mount-4",
			raw:    `mount fstype=devpts options=(rw nosuid noexec)                devpts -> /newroot/dev/pts/`,
			tokens: []string{"mount", "fstype=devpts", "options=(rw nosuid noexec)", "devpts", "->", "/newroot/dev/pts/"},
		},
		{
			name:   "signal",
			raw:    `signal (receive) set=(cont, term,winch) peer=at-spi-bus-launcher`,
			tokens: []string{"signal", "(receive)", "set=(cont, term,winch)", "peer=at-spi-bus-launcher"},
		},
		{
			name:   "unix",
			raw:    `unix (send receive) type=stream addr="@/tmp/.ICE[0-9]*-unix/19 5" peer=(label="@{p_systemd}", addr=none)`,
			tokens: []string{"unix", "(send receive)", "type=stream", "addr=\"@/tmp/.ICE[0-9]*-unix/19 5\"", "peer=(label=\"@{p_systemd}\", addr=none)"},
		},
		{
			name: "unix-2",
			raw: `  unix (connect, receive, send)
			type=stream
			peer=(addr="@/tmp/ibus/dbus-????????")`,
			tokens: []string{"unix", "(connect, receive, send)\n", "type=stream\n", `peer=(addr="@/tmp/ibus/dbus-????????")`},
		},
		{
			name: "dbus",
			raw:  `dbus receive bus=system path=/org/freedesktop/DBus interface=org.freedesktop.DBus member=AddMatch peer=(name=:1.3, label=power-profiles-daemon)`,
			tokens: []string{
				"dbus", "receive", "bus=system",
				"path=/org/freedesktop/DBus", "interface=org.freedesktop.DBus",
				"member=AddMatch", "peer=(name=:1.3, label=power-profiles-daemon)",
			},
		},
		{
			name:   "file-1",
			raw:    `owner @{user_config_dirs}/powerdevilrc{,.@{rand6}} rwl -> @{user_config_dirs}/#@{int}`,
			tokens: []string{"owner", "@{user_config_dirs}/powerdevilrc{,.@{rand6}}", "rwl", "->", "@{user_config_dirs}/#@{int}"},
		},
		{
			name:   "file-2",
			raw:    `@{sys}/devices/@{pci}/class r`,
			tokens: []string{"@{sys}/devices/@{pci}/class", "r"},
		},
		{
			name:   "file-3",
			raw:    `owner @{PROC}/@{pid}/task/@{tid}/comm rw`,
			tokens: []string{"owner", "@{PROC}/@{pid}/task/@{tid}/comm", "rw"},
		},
		{
			name:   "file-4",
			raw:    `owner /{var/,}tmp/#@{int} rw`,
			tokens: []string{"owner", "/{var/,}tmp/#@{int}", "rw"},
		},
	}

	// Test cases for Parse
	testBlocks = []struct {
		name      string
		raw       string
		apparmor  *AppArmorProfileFile
		wParseErr bool
	}{
		{
			name:      "empty",
			raw:       "",
			apparmor:  &AppArmorProfileFile{},
			wParseErr: false,
		},
		{
			name: "comment",
			raw: `
			# IsLineRule comment
			include <tunables/global> # comment included
			@{lib_dirs} = @{lib}/@{name} /opt/@{name} # comment in variable`,
			apparmor: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{RuleBase: RuleBase{IsLineRule: true, Comment: " IsLineRule comment"}},
					&Include{
						RuleBase: RuleBase{Comment: " comment included"},
						Path:     "tunables/global", IsMagic: true,
					},
					&Variable{
						RuleBase: RuleBase{Comment: " comment in variable"},
						Name:     "lib_dirs", Define: true,
						Values: []string{"@{lib}/@{name}", "/opt/@{name}"},
					},
				},
			},
			wParseErr: false,
		},
		{
			name: "cornercases",
			raw: `# Simple test
			include <tunables/global>

			# { commented block }
			@{name} = {D,d}ummy
			@{exec_path} = @{bin}/@{name}
			alias /mnt/{,usr.sbin.}mount.cifs -> /sbin/mount.cifs,
			@{coreutils} += gawk {,e,f}grep head
			profile @{exec_path} {
			`,
			apparmor: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{RuleBase: RuleBase{IsLineRule: true, Comment: " Simple test"}},
					&Comment{RuleBase: RuleBase{IsLineRule: true, Comment: " { commented block }"}},
					&Include{IsMagic: true, Path: "tunables/global"},
					&Variable{Name: "name", Values: []string{"{D,d}ummy"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"@{bin}/@{name}"}, Define: true},
					&Alias{Path: "/mnt/{,usr.sbin.}mount.cifs", RewrittenPath: "/sbin/mount.cifs"},
					&Variable{Name: "coreutils", Values: []string{"gawk", "{,e,f}grep", "head"}, Define: false},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "@{exec_path}",
							Attachments: []string{},
							Attributes:  map[string]string{},
							Flags:       []string{},
						},
					},
				},
			},
			wParseErr: false,
		},
		{
			name: "string.aa",
			raw:  util.MustReadFile(testData.Join("string.aa")),
			apparmor: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{RuleBase: RuleBase{Comment: " Simple test profile for the AppArmorProfileFile.String() method", IsLineRule: true}},
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
					&Include{IsMagic: true, Path: "tunables/global"},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{"@{bin}/foo", "@{lib}/foo"},
					},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "foo",
							Attachments: []string{"@{exec_path}"},
							Attributes:  map[string]string{"security.tagged": "allowed"},
							Flags:       []string{"complain", "attach_disconnected"},
						},
					},
				},
			},
			wParseErr: false,
		},
	}
)
