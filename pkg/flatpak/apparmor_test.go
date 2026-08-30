// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package flatpak

import (
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/aa"
)

// stubDbusLabels installs known dbus labels and marks the lazy scan as done so
// the tests never read the system profile tree.
func stubDbusLabels(t *testing.T) {
	t.Helper()
	origExact, origPrefix := dbusLabelMap, dbusLabelPrefixMap
	t.Cleanup(func() {
		dbusLabelMap, dbusLabelPrefixMap = origExact, origPrefix
	})
	dbusLabelOnce.Do(func() {})
	dbusLabelMap = map[string]string{
		"org.bluez":                     "bluetoothd",
		"org.freedesktop.UPower":        "upowerd",
		"org.kde.StatusNotifierWatcher": "plasmashell",
	}
	dbusLabelPrefixMap = map[string]string{}
}

func TestProfileName(t *testing.T) {
	tests := []struct {
		name  string
		appID string
		want  string
	}{
		{
			name:  "chrome",
			appID: "com.google.Chrome",
			want:  "flatpak.com.google.Chrome",
		},
		{
			name:  "steam",
			appID: "com.valvesoftware.Steam",
			want:  "flatpak.com.valvesoftware.Steam",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProfileName(tt.appID); got != tt.want {
				t.Errorf("ProfileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFilesystemSpec(t *testing.T) {
	tests := []struct {
		name       string
		spec       string
		want       string
		wantAccess aa.AccessMask
	}{
		{
			name:       "no access suffix is read write",
			spec:       "xdg-music",
			want:       "xdg-music",
			wantAccess: aa.MustAccess(aa.FILE, "r", "w", "l", "k"),
		},
		{
			name:       "read only suffix",
			spec:       "~/.config/dconf:ro",
			want:       "~/.config/dconf",
			wantAccess: aa.MustAccess(aa.FILE, "r"),
		},
		{
			name:       "create suffix is read write",
			spec:       "xdg-run/app/com.discordapp.Discord:create",
			want:       "xdg-run/app/com.discordapp.Discord",
			wantAccess: aa.MustAccess(aa.FILE, "r", "w", "l", "k"),
		},
		{
			name:       "absolute path read only",
			spec:       "/run/udev:ro",
			want:       "/run/udev",
			wantAccess: aa.MustAccess(aa.FILE, "r"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, access := parseFilesystemSpec(tt.spec)
			if got != tt.want {
				t.Errorf("parseFilesystemSpec() = %v, want %v", got, tt.want)
			}
			if access != tt.wantAccess {
				t.Errorf("parseFilesystemSpec() access = %v, want %v", access, tt.wantAccess)
			}
		})
	}
}

func TestExtractPathComponents(t *testing.T) {
	tests := []struct {
		name     string
		fs       string
		wantBase string
		wantSub  string
	}{
		{
			name:     "no subdirectory",
			fs:       "xdg-music",
			wantBase: "xdg-music",
			wantSub:  "",
		},
		{
			name:     "xdg prefix with subdirectory",
			fs:       "xdg-run/dconf",
			wantBase: "xdg-run",
			wantSub:  "dconf",
		},
		{
			name:     "home prefix with subdirectory",
			fs:       "~/.config/kioslaverc",
			wantBase: "~",
			wantSub:  ".config/kioslaverc",
		},
		{
			name:     "absolute path is kept whole",
			fs:       "/run/media",
			wantBase: "/run/media",
			wantSub:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base, sub := extractPathComponents(tt.fs)
			if base != tt.wantBase || sub != tt.wantSub {
				t.Errorf("extractPathComponents() = %v %v, want %v %v",
					base, sub, tt.wantBase, tt.wantSub)
			}
		})
	}
}

func TestAddIncludes(t *testing.T) {
	tests := []struct {
		name    string
		items   []string
		allowed []string
		prefix  string
		want    []string
	}{
		{
			name:    "known items",
			items:   []string{"x11", "wayland"},
			allowed: allowedSockets,
			prefix:  "sockets",
			want: []string{
				"abstractions/flatpak/sockets/x11",
				"abstractions/flatpak/sockets/wayland",
			},
		},
		{
			name:    "unknown item is skipped",
			items:   []string{"x11", "not-a-socket"},
			allowed: allowedSockets,
			prefix:  "sockets",
			want:    []string{"abstractions/flatpak/sockets/x11"},
		},
		{
			name:    "no item",
			items:   nil,
			allowed: allowedDevices,
			prefix:  "devices",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rules := addIncludes(tt.items, tt.allowed, tt.prefix)
			if len(rules) != len(tt.want) {
				t.Fatalf("addIncludes() = %v, want %v", rules, tt.want)
			}
			for i, rule := range rules {
				include, ok := rule.(*aa.Include)
				if !ok {
					t.Fatalf("addIncludes()[%d] = %T, want *aa.Include", i, rule)
				}
				if include.Path != tt.want[i] {
					t.Errorf("addIncludes()[%d] = %v, want %v", i, include.Path, tt.want[i])
				}
			}
		})
	}
}

func TestFlatpakAppArmorProfile_addRuntime(t *testing.T) {
	tests := []struct {
		name    string
		runtime string
		want    string
	}{
		{
			name:    "freedesktop platform",
			runtime: "org.freedesktop.Platform/x86_64/25.08",
			want:    "abstractions/flatpak/platform/org.freedesktop",
		},
		{
			name:    "gnome platform",
			runtime: "org.gnome.Platform/x86_64/49",
			want:    "abstractions/flatpak/platform/org.gnome",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FlatpakAppArmorProfile{
				Metadata: &FlatpakMetadata{
					Application: Application{Runtime: tt.runtime},
				},
			}
			rules := p.addRuntime()
			if len(rules) != 1 {
				t.Fatalf("addRuntime() = %v, want 1 rule", rules)
			}
			if got := rules[0].(*aa.Include).Path; got != tt.want {
				t.Errorf("addRuntime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlatpakAppArmorProfile_addBaseApp(t *testing.T) {
	tests := []struct {
		name string
		base string
		want string // "" means no rule is expected
	}{
		{
			name: "chromium base app",
			base: "app/org.chromium.Chromium.BaseApp/x86_64/25.08",
			want: "abstractions/flatpak/baseapp/org.chromium.Chromium",
		},
		{
			name: "no base app",
			base: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &FlatpakAppArmorProfile{
				Metadata: &FlatpakMetadata{
					Application: Application{Base: tt.base},
				},
			}
			rules := p.addBaseApp()
			if tt.want == "" {
				if len(rules) != 0 {
					t.Fatalf("addBaseApp() = %v, want no rule", rules)
				}
				return
			}
			if len(rules) != 1 {
				t.Fatalf("addBaseApp() = %v, want 1 rule", rules)
			}
			if got := rules[0].(*aa.Include).Path; got != tt.want {
				t.Errorf("addBaseApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlatpakAppArmorProfile_Generate(t *testing.T) {
	tests := []struct {
		name         string
		app          string
		mode         string
		wantContains []string
		wantAbsent   []string
	}{
		{
			name: "chrome in complain mode",
			app:  "com.google.Chrome",
			mode: "complain",
			wantContains: []string{
				"profile flatpak.com.google.Chrome flags=(attach_disconnected.path=/att/flatpak.com.google.Chrome complain mediate_deleted) {",
				"include <abstractions/flatpak/platform/org.freedesktop>",
				"include <abstractions/flatpak/baseapp/org.chromium.Chromium>",
				"include <abstractions/flatpak/shared/network>",
				"include <abstractions/flatpak/sockets/x11>",
				"include <abstractions/flatpak/devices/all>",
				"owner @{user_music_dirs}/** rwlk,",
				"owner @{HOME}/.config/dconf/{,**} r,",
				"include if exists <local/flatpak.com.google.Chrome>",
				"profile flatpak.dbus.com.google.Chrome",
				"dbus bind bus=session name=com.google.Chrome{,.*},",
				"include <abstractions/notifications> # org.freedesktop.Notifications",
			},
			wantAbsent: []string{
				"abstractions/flatpak/features/",
			},
		},
		{
			name: "steam in enforce mode",
			app:  "com.valvesoftware.Steam",
			mode: "enforce",
			wantContains: []string{
				"profile flatpak.com.valvesoftware.Steam",
				"include <abstractions/flatpak/features/devel>",
				"include <abstractions/flatpak/features/multiarch>",
				"include <abstractions/flatpak/features/bluetooth>",
				"include <abstractions/deny-sensitive-home>",
				"profile flatpak.dbus.com.valvesoftware.Steam",
			},
			wantAbsent: []string{
				"complain",
				"abstractions/flatpak/baseapp/",
			},
		},
	}

	stubDbusLabels(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := loadTestMetadata(t, tt.app)
			p := NewFlatpakAppArmorProfile(meta, tt.mode)
			if err := p.Generate(); err != nil {
				t.Fatalf("Generate() error = %v, wantErr %v", err, false)
			}
			p.Format()
			got := p.String()

			if p.FileName != ProfileName(tt.app) {
				t.Errorf("Generate() FileName = %v, want %v", p.FileName, ProfileName(tt.app))
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(got, want) {
					t.Errorf("Generate() = missing %q in:\n%s", want, got)
				}
			}
			for _, absent := range tt.wantAbsent {
				if strings.Contains(got, absent) {
					t.Errorf("Generate() = unexpected %q in:\n%s", absent, got)
				}
			}
		})
	}
}
