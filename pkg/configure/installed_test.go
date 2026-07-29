// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"os"
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

func TestIsInstalled(t *testing.T) {
	tests := []struct {
		name        string
		attachments []string
		want        bool
	}{
		{
			name:        "existing binary",
			attachments: []string{"@{bin}/ls"},
			want:        true,
		},
		{
			name:        "nonexistent binary",
			attachments: []string{"@{bin}/this-program-does-not-exist-12345"},
			want:        false,
		},
		{
			name:        "unresolved variable kept",
			attachments: []string{"@{brave_path}"},
			want:        true,
		},
		{
			name:        "absolute existing path",
			attachments: []string{"/usr/bin/ls"},
			want:        true,
		},
		{
			name:        "post-prebuild resolved",
			attachments: []string{"/{,usr/}bin/ls"},
			want:        true,
		},
		{
			name:        "post-prebuild not installed",
			attachments: []string{"/{,usr/}bin/this-program-does-not-exist-12345"},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isInstalled(tt.attachments)
			if got != tt.want {
				t.Errorf("isInstalled(%v) = %v, want %v", tt.attachments, got, tt.want)
			}
		})
	}
}

func TestGetAttachments(t *testing.T) {
	tests := []struct {
		name    string
		profile string
		want    []string
		wantErr bool
	}{
		{
			name: "post-prebuild resolved",
			profile: `
abi <abi/5.0>,

profile dolphin /{,usr/}bin/dolphin {
  include <abstractions/base>
}`,
			want: []string{"/{,usr/}bin/dolphin"},
		},
		{
			name: "post-prebuild multiple",
			profile: `
abi <abi/5.0>,

profile claude {/{,usr/}bin/claude,/opt/claude-code/bin/claude} {
  include <abstractions/base>
}`,
			want: []string{"{/{,usr/}bin/claude,/opt/claude-code/bin/claude}"},
		},
		{
			name: "child profile without exec_path",
			profile: `
abi <abi/5.0>,

profile child-pager {
  include <abstractions/base>
}`,
			want: nil,
		},
		{
			name:    "parse error",
			profile: "profile broken {\n",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAttachments(tt.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAttachments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("getAttachments() = %v, want %v", got, tt.want)
				return
			}
			for i, w := range tt.want {
				if i < len(got) && got[i] != w {
					t.Errorf("getAttachments()[%d] = %q, want %q", i, got[i], w)
				}
			}
		})
	}
}

func TestSbinGlob(t *testing.T) {
	tests := []struct {
		name   string
		family string
		want   string
	}{
		{
			name:   "pacman merges sbin and bin",
			family: "pacman",
			want:   "/{,usr/}{,s}bin",
		},
		{
			name:   "other families keep sbin only",
			family: "apt",
			want:   "/{,usr/}sbin",
		},
	}

	original := tasks.Family
	t.Cleanup(func() { tasks.Family = original })
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tasks.Family = tt.family
			if got := sbinGlob(); got != tt.want {
				t.Errorf("sbinGlob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasIgnoredSuffix(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "shared library",
			filename: "libfoo.so",
			want:     true,
		},
		{
			name:     "systemd unit",
			filename: "foo.service",
			want:     true,
		},
		{
			name:     "python module",
			filename: "module.py",
			want:     true,
		},
		{
			name:     "versioned library not matched",
			filename: "libfoo.so.1",
			want:     false,
		},
		{
			name:     "regular program",
			filename: "bash",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasIgnoredSuffix(tt.filename); got != tt.want {
				t.Errorf("hasIgnoredSuffix(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestExpandPathBraces(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		want    []string
	}{
		{
			name:    "structural brace expanded",
			pattern: "/{,usr/}bin/foo",
			want:    []string{"/bin/foo", "/usr/bin/foo"},
		},
		{
			name:    "intra-segment brace kept verbatim",
			pattern: "/usr/lib{,exec,32,64}/foo",
			want:    []string{"/usr/lib{,exec,32,64}/foo"},
		},
		{
			name:    "version pattern kept verbatim",
			pattern: "/usr/lib/pg/[0-9]{[0-9],}/bin/pg_ctl",
			want:    []string{"/usr/lib/pg/[0-9]{[0-9],}/bin/pg_ctl"},
		},
		{
			name:    "nested structural alternation",
			pattern: "/{{,usr/}bin/x,opt/y}",
			want:    []string{"/bin/x", "/usr/bin/x", "/opt/y"},
		},
		{
			name:    "unbalanced brace returned as-is",
			pattern: "/usr/{bin/foo",
			want:    []string{"/usr/{bin/foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandPathBraces(tt.pattern)
			if len(got) != len(tt.want) {
				t.Errorf("expandPathBraces(%q) = %v, want %v", tt.pattern, got, tt.want)
				return
			}
			for _, w := range tt.want {
				if !slices.Contains(got, w) {
					t.Errorf("expandPathBraces(%q) missing %q, got %v", tt.pattern, w, got)
				}
			}
		})
	}
}

func TestGlobExists(t *testing.T) {
	setupTmp(t)
	base := paths.New(t.TempDir())
	// Mimic /usr/lib/postgresql/17/bin/pg_ctl and /etc/update-motd.d/10-foo.
	for _, rel := range []string{"lib/postgresql/17/bin", "motd.d"} {
		if err := base.Join(rel).MkdirAll(); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
	}
	for _, rel := range []string{"lib/postgresql/17/bin/pg_ctl", "motd.d/10-foo"} {
		if err := base.Join(rel).WriteFile([]byte("x")); err != nil {
			t.Fatalf("write: %v", err)
		}
	}
	root := base.String()

	tests := []struct {
		name    string
		pattern string
		want    bool
	}{
		{
			name:    "trailing wildcard segment matches",
			pattern: root + "/motd.d/*",
			want:    true,
		},
		{
			name:    "wildcard segment no match",
			pattern: root + "/nope.d/*",
			want:    false,
		},
		{
			name:    "version char-class matches two digits",
			pattern: root + "/lib{,exec}/postgresql/[0-9]{[0-9],}{[0-9],}/bin/pg_ctl",
			want:    true,
		},
		{
			name:    "version char-class no match",
			pattern: root + "/lib/postgresql/[0-9]/bin/pg_ctl",
			want:    false, // single [0-9] can't match "17"
		},
		{
			name:    "doublestar matches when parent exists",
			pattern: root + "/lib/postgresql/**",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := globExists(tt.pattern); got != tt.want {
				t.Errorf("globExists(%q) = %v, want %v", tt.pattern, got, tt.want)
			}
		})
	}
}

func TestPathIsInstalled(t *testing.T) {
	setupTmp(t)
	base := paths.New(t.TempDir())
	if err := base.Join("real").MkdirAll(); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := base.Join("real/prog").WriteFile([]byte("#!/bin/sh\n")); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.Symlink(base.Join("real").String(), base.Join("link").String()); err != nil {
		t.Fatalf("symlink: %v", err)
	}

	tests := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "unresolved variable kept",
			path: "not-absolute",
			want: true,
		},
		{
			name: "wildcard remnant skipped",
			path: "/usr/bin/*",
			want: false,
		},
		{
			name: "existing path",
			path: base.Join("real/prog").String(),
			want: true,
		},
		{
			name: "missing path",
			path: base.Join("real/absent").String(),
			want: false,
		},
		{
			name: "symlink component skipped",
			path: base.Join("link/prog").String(),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathIsInstalled(tt.path); got != tt.want {
				t.Errorf("pathIsInstalled(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestSelectInstalled_Apply(t *testing.T) {
	tests := []struct {
		name           string
		profiles       map[string]string // apparmor.d relative file -> content
		groups         map[string]string // profile basename -> group (as Merge records)
		include        []string          // include entries passed to NewSelectInstalled
		noReadProfiles []string          // apparmor.d relative files made unreadable
		roDirs         []string          // apparmor.d relative dirs made read-only
		missingRoot    bool              // build root does not exist
		wantErr        bool
		wantResult     []string // entries expected in the result
		wantKept       []string // apparmor.d relative files expected to remain
		wantRemoved    []string // apparmor.d relative files expected to be removed
	}{
		{
			name: "keep installed remove missing",
			profiles: map[string]string{
				"groups/misc/foo": "profile foo @{bin}/true {\n}\n",
				"groups/misc/bar": "profile bar /usr/bin/no-such-program-zz {\n}\n",
				"child":           "profile child {\n}\n",
				"garbage":         "profile broken {\n",
				"secret":          "profile secret /usr/bin/no-such-program-zz {\n}\n",
			},
			noReadProfiles: []string{"secret"},
			wantResult:     []string{"Kept 4 profiles", "Ignored 1"},
			wantKept:       []string{"groups/misc/foo", "child", "garbage", "secret"},
			wantRemoved:    []string{"groups/misc/bar"},
		},
		{
			name: "keep group helper when group installed",
			profiles: map[string]string{
				"code":        "profile code @{bin}/true {\n}\n",
				"code-shells": "profile code-shells flags=(attach_disconnected) {\n}\n",
			},
			groups:      map[string]string{"code": "code", "code-shells": "code"},
			wantKept:    []string{"code", "code-shells"},
			wantRemoved: []string{},
		},
		{
			name: "drop group helper when group not installed",
			profiles: map[string]string{
				"steam":         "profile steam /usr/bin/no-such-program-zz {\n}\n",
				"steam-runtime": "profile steam-runtime flags=(attach_disconnected) {\n}\n",
			},
			groups:      map[string]string{"steam": "steam", "steam-runtime": "steam"},
			wantRemoved: []string{"steam", "steam-runtime"},
		},
		{
			name: "always-keep profiles survive uninstalled group",
			profiles: map[string]string{
				"steam":         "profile steam /usr/bin/no-such-program-zz {\n}\n",
				"steam-runtime": "profile steam-runtime flags=(attach_disconnected) {\n}\n",
				"child-pager":   "profile child-pager flags=(attach_disconnected) {\n}\n",
				"dbus-session":  "profile dbus-session flags=(attach_disconnected) {\n}\n",
			},
			groups: map[string]string{
				"steam":         "steam",
				"steam-runtime": "steam",
				"child-pager":   "children",
				"dbus-session":  "bus",
			},
			wantKept:    []string{"child-pager", "dbus-session"},
			wantRemoved: []string{"steam", "steam-runtime"},
		},
		{
			name: "included profiles kept despite uninstalled attachment",
			profiles: map[string]string{
				"steam":         "profile steam /usr/bin/no-such-program-zz {\n}\n",
				"steam-runtime": "profile steam-runtime flags=(attach_disconnected) {\n}\n",
				"other":         "profile other /usr/bin/no-such-program-zz {\n}\n",
			},
			groups:      map[string]string{"steam": "steam", "steam-runtime": "steam"},
			include:     []string{"steam"},
			wantKept:    []string{"steam", "steam-runtime"},
			wantRemoved: []string{"other"},
		},
		{
			name: "namespaced profile kept despite uninstalled attachment",
			profiles: map[string]string{
				":podman:podman": "profile :podman:podman /usr/bin/no-such-program-zz {\n}\n",
			},
			wantKept: []string{":podman:podman"},
		},
		{
			name: "overwriting profile kept when installed",
			profiles: map[string]string{
				"systemd-detect-virt.apparmor.d": "profile systemd-detect-virt @{bin}/true {\n}\n",
				"absent.apparmor.d":              "profile absent /usr/bin/no-such-program-zz {\n}\n",
			},
			wantKept:    []string{"systemd-detect-virt.apparmor.d"},
			wantRemoved: []string{"absent.apparmor.d"},
		},
		{
			name:        "read profiles error",
			missingRoot: true,
			wantErr:     true,
		},
		{
			name: "remove profile error",
			profiles: map[string]string{
				"sub/bar": "profile bar /usr/bin/no-such-program-zz {\n}\n",
			},
			roDirs:  []string{"sub"},
			wantErr: true,
		},
	}

	setupTmp(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *tasks.TaskConfig
			if tt.missingRoot {
				c = tasks.NewTaskConfig(paths.New(t.TempDir()).Join("missing"))
			} else {
				c = newTaskConfigTmp(t)
			}
			seedFiles(t, c.RootApparmor, tt.profiles)
			if tt.groups != nil {
				c.Groups = tt.groups
			}
			for _, rel := range tt.noReadProfiles {
				chmodNoAccess(t, c.RootApparmor.Join(rel))
			}
			for _, rel := range tt.roDirs {
				chmodRO(t, c.RootApparmor.Join(rel))
			}

			task := NewSelectInstalled(tt.include...)
			task.SetConfig(c)
			got, err := task.Apply()
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			for _, want := range tt.wantResult {
				if !slices.Contains(got, want) {
					t.Errorf("Apply() = %v, want %v", got, want)
				}
			}
			for _, rel := range tt.wantKept {
				if c.RootApparmor.Join(rel).NotExist() {
					t.Errorf("Apply() kept %s = %v, want %v", rel, "removed", "exists")
				}
			}
			for _, rel := range tt.wantRemoved {
				if c.RootApparmor.Join(rel).Exist() {
					t.Errorf("Apply() removed %s = %v, want %v", rel, "exists", "removed")
				}
			}
		})
	}
}
