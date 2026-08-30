// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package flatpak

import (
	"os"
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

// testdata holds the metadata files of real flatpak applications.
var testdata = paths.New("../../tests/testdata/flatpak")

// setupTmp keeps the test temporary files under /tmp/tests.
func setupTmp(t *testing.T) {
	t.Helper()
	if err := os.MkdirAll("/tmp/tests", 0o755); err != nil {
		t.Fatalf("mkdir /tmp/tests: %v", err)
	}
	t.Setenv("TMPDIR", "/tmp/tests")
}

// loadTestMetadata loads the metadata of a reference application.
func loadTestMetadata(t *testing.T, name string) *FlatpakMetadata {
	t.Helper()
	meta, err := LoadFlatpakMetadata(testdata.Join(name))
	if err != nil {
		t.Fatalf("LoadFlatpakMetadata(%s) error = %v", name, err)
	}
	return meta
}

func TestLoadFlatpakMetadata(t *testing.T) {
	tests := []struct {
		name            string
		app             string
		wantRuntime     string
		wantSDK         string
		wantBase        string
		wantCommand     string
		wantShared      []string
		wantSockets     []string
		wantDevices     []string
		wantFeatures    []string
		wantPersistent  []string
		wantFilesystems []string
		wantSession     []Dbus
		wantSystem      []Dbus
		wantExtensions  int
	}{
		{
			name:        "chrome",
			app:         "com.google.Chrome",
			wantRuntime: "org.freedesktop.Platform/x86_64/25.08",
			wantSDK:     "org.freedesktop.Sdk/x86_64/25.08",
			wantBase:    "app/org.chromium.Chromium.BaseApp/x86_64/25.08",
			wantCommand: "chrome",
			wantShared:  []string{"ipc", "network"},
			wantSockets: []string{"cups", "pcsc", "pulseaudio", "wayland", "x11"},
			wantDevices: []string{"all"},
			wantFilesystems: []string{
				"host-etc", "~/.config/kioslaverc", "xdg-music", "xdg-pictures",
				"xdg-videos", "/run/.heim_org.h5l.kcm-socket", "~/.config/dconf:ro",
				"xdg-download", "xdg-run/dconf", "xdg-documents", "xdg-run/pipewire-0",
			},
			wantSession: []Dbus{
				{Interface: "org.freedesktop.Notifications", Action: "talk"},
				{Interface: "org.freedesktop.FileManager1", Action: "talk"},
				{Interface: "org.mpris.MediaPlayer2.chromium.*", Action: "own"},
				{Interface: "org.kde.StatusNotifierWatcher", Action: "talk"},
				{Interface: "org.freedesktop.ScreenSaver", Action: "talk"},
				{Interface: "org.freedesktop.secrets", Action: "talk"},
				{Interface: "ca.desrt.dconf", Action: "talk"},
				{Interface: "org.gnome.SessionManager", Action: "talk"},
			},
			wantSystem: []Dbus{
				{Interface: "org.freedesktop.Avahi", Action: "talk"},
				{Interface: "org.freedesktop.UPower", Action: "talk"},
				{Interface: "org.bluez", Action: "talk"},
			},
			wantExtensions: 1,
		},
		{
			name:           "steam has no base and uses features",
			app:            "com.valvesoftware.Steam",
			wantRuntime:    "org.freedesktop.Platform/x86_64/25.08",
			wantSDK:        "org.freedesktop.Sdk/x86_64/25.08",
			wantBase:       "",
			wantCommand:    "steam",
			wantShared:     []string{"network", "ipc"},
			wantSockets:    []string{"x11", "wayland", "pulseaudio"},
			wantDevices:    []string{"all"},
			wantFeatures:   []string{"devel", "multiarch", "bluetooth", "per-app-dev-shm"},
			wantPersistent: []string{"."},
			wantFilesystems: []string{
				"xdg-config/MangoHud:ro", "/mnt", "xdg-pictures:ro",
				"xdg-run/app/com.discordapp.Discord:create", "/run/udev:ro",
				"xdg-run/speech-dispatcher:ro", "/run/media", "xdg-music:ro",
				"/media", "xdg-run/pipewire-0:ro",
			},
			wantSession: []Dbus{
				{Interface: "com.steampowered.*", Action: "own"},
				{Interface: "org.kde.StatusNotifierWatcher", Action: "talk"},
				{Interface: "org.freedesktop.Notifications", Action: "talk"},
				{Interface: "org.freedesktop.ScreenSaver", Action: "talk"},
				{Interface: "org.freedesktop.PowerManagement", Action: "talk"},
				{Interface: "org.gnome.SessionManager", Action: "talk"},
			},
			wantSystem: []Dbus{
				{Interface: "org.freedesktop.UPower", Action: "talk"},
				{Interface: "org.freedesktop.UDisks2", Action: "talk"},
			},
			wantExtensions: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := loadTestMetadata(t, tt.app)

			if meta.Name != tt.app {
				t.Errorf("LoadFlatpakMetadata() Name = %v, want %v", meta.Name, tt.app)
			}
			if meta.Runtime != tt.wantRuntime {
				t.Errorf("LoadFlatpakMetadata() Runtime = %v, want %v", meta.Runtime, tt.wantRuntime)
			}
			if meta.SDK != tt.wantSDK {
				t.Errorf("LoadFlatpakMetadata() SDK = %v, want %v", meta.SDK, tt.wantSDK)
			}
			if meta.Base != tt.wantBase {
				t.Errorf("LoadFlatpakMetadata() Base = %v, want %v", meta.Base, tt.wantBase)
			}
			if meta.Command != tt.wantCommand {
				t.Errorf("LoadFlatpakMetadata() Command = %v, want %v", meta.Command, tt.wantCommand)
			}
			if !slices.Equal(meta.Shared, tt.wantShared) {
				t.Errorf("LoadFlatpakMetadata() Shared = %v, want %v", meta.Shared, tt.wantShared)
			}
			if !slices.Equal(meta.Sockets, tt.wantSockets) {
				t.Errorf("LoadFlatpakMetadata() Sockets = %v, want %v", meta.Sockets, tt.wantSockets)
			}
			if !slices.Equal(meta.Devices, tt.wantDevices) {
				t.Errorf("LoadFlatpakMetadata() Devices = %v, want %v", meta.Devices, tt.wantDevices)
			}
			if !slices.Equal(meta.Features, tt.wantFeatures) {
				t.Errorf("LoadFlatpakMetadata() Features = %v, want %v", meta.Features, tt.wantFeatures)
			}
			if !slices.Equal(meta.Persistent, tt.wantPersistent) {
				t.Errorf("LoadFlatpakMetadata() Persistent = %v, want %v", meta.Persistent, tt.wantPersistent)
			}
			if !slices.Equal(meta.Filesystems, tt.wantFilesystems) {
				t.Errorf("LoadFlatpakMetadata() Filesystems = %v, want %v", meta.Filesystems, tt.wantFilesystems)
			}
			if !slices.Equal(meta.DbusSession, tt.wantSession) {
				t.Errorf("LoadFlatpakMetadata() DbusSession = %v, want %v", meta.DbusSession, tt.wantSession)
			}
			if !slices.Equal(meta.DbusSystem, tt.wantSystem) {
				t.Errorf("LoadFlatpakMetadata() DbusSystem = %v, want %v", meta.DbusSystem, tt.wantSystem)
			}
			if len(meta.Extensions) != tt.wantExtensions {
				t.Errorf("LoadFlatpakMetadata() len(Extensions) = %v, want %v", len(meta.Extensions), tt.wantExtensions)
			}
		})
	}
}

// TestLoadFlatpakMetadata_Extension covers the manual mapping of the
// "Extension <name>" sections, whose names are not known in advance.
func TestLoadFlatpakMetadata_Extension(t *testing.T) {
	meta := loadTestMetadata(t, "com.valvesoftware.Steam")

	want := Extension{
		Name:           "com.valvesoftware.Steam.Utility",
		Directory:      "utils",
		AddLdPath:      "lib",
		Version:        "stable",
		Autodelete:     true,
		NoAutodownload: true,
	}
	idx := slices.IndexFunc(meta.Extensions, func(e Extension) bool { return e.Name == want.Name })
	if idx == -1 {
		t.Fatalf("LoadFlatpakMetadata() Extensions = %v, want one named %v", meta.Extensions, want.Name)
	}
	if meta.Extensions[idx] != want {
		t.Errorf("LoadFlatpakMetadata() Extension = %v, want %v", meta.Extensions[idx], want)
	}
}

func TestLoadFlatpakMetadata_Error(t *testing.T) {
	if _, err := LoadFlatpakMetadata(testdata.Join("does.not.Exist")); err == nil {
		t.Errorf("LoadFlatpakMetadata() error = %v, wantErr %v", err, true)
	}
}

func TestResolveNegations(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		want  []string
	}{
		{
			name:  "nothing negated",
			items: []string{"network", "ipc"},
			want:  []string{"network", "ipc"},
		},
		{
			name:  "negation removes both entries",
			items: []string{"network", "ipc", "!network"},
			want:  []string{"ipc"},
		},
		{
			name:  "negation without a match is dropped",
			items: []string{"ipc", "!network"},
			want:  []string{"ipc"},
		},
		{
			name:  "all negated",
			items: []string{"network", "!network"},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveNegations(tt.items)
			if !slices.Equal(got, tt.want) {
				t.Errorf("resolveNegations() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFlatpakMetadata_LoadOverride(t *testing.T) {
	tests := []struct {
		name            string
		override        string // content of the override file, "" means not created
		wantShared      []string
		wantSockets     []string
		wantFilesystems []string
		wantSession     int // number of session bus policies after the merge
	}{
		{
			name:        "no override file",
			wantShared:  []string{"ipc", "network"},
			wantSockets: []string{"cups", "pcsc", "pulseaudio", "wayland", "x11"},
			wantSession: 8,
		},
		{
			name:        "override adds a socket",
			override:    "[Context]\nsockets=gpg-agent;\n",
			wantShared:  []string{"ipc", "network"},
			wantSockets: []string{"cups", "pcsc", "pulseaudio", "wayland", "x11", "gpg-agent"},
			wantSession: 8,
		},
		{
			name:        "override negates a shared feature",
			override:    "[Context]\nshared=!network;\n",
			wantShared:  []string{"ipc"},
			wantSockets: []string{"cups", "pcsc", "pulseaudio", "wayland", "x11"},
			wantSession: 8,
		},
		{
			name:        "override adds a bus policy",
			override:    "[Session Bus Policy]\norg.freedesktop.portal.Desktop=talk\n",
			wantShared:  []string{"ipc", "network"},
			wantSockets: []string{"cups", "pcsc", "pulseaudio", "wayland", "x11"},
			wantSession: 9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupTmp(t)
			meta := loadTestMetadata(t, "com.google.Chrome")
			dir := paths.New(t.TempDir())
			if tt.override != "" {
				if err := dir.Join(meta.Name).WriteFile([]byte(tt.override)); err != nil {
					t.Fatalf("write override: %v", err)
				}
			}

			if err := meta.LoadOverride(dir); err != nil {
				t.Fatalf("LoadOverride() error = %v, wantErr %v", err, false)
			}
			if !slices.Equal(meta.Shared, tt.wantShared) {
				t.Errorf("LoadOverride() Shared = %v, want %v", meta.Shared, tt.wantShared)
			}
			if !slices.Equal(meta.Sockets, tt.wantSockets) {
				t.Errorf("LoadOverride() Sockets = %v, want %v", meta.Sockets, tt.wantSockets)
			}
			if len(meta.DbusSession) != tt.wantSession {
				t.Errorf("LoadOverride() len(DbusSession) = %v, want %v", len(meta.DbusSession), tt.wantSession)
			}
		})
	}
}

func TestFlatpakMetadata_Parts(t *testing.T) {
	tests := []struct {
		name        string
		appID       string
		wantTLD     string
		wantVendor  string
		wantProduct string
		wantName    string
	}{
		{
			name:        "three parts",
			appID:       "com.google.Chrome",
			wantTLD:     "com",
			wantVendor:  "google",
			wantProduct: "Chrome",
			wantName:    "chrome",
		},
		{
			name:        "four parts uses the last one",
			appID:       "org.freedesktop.Platform.GL32",
			wantTLD:     "org",
			wantVendor:  "freedesktop",
			wantProduct: "Platform",
			wantName:    "gl32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := &FlatpakMetadata{Application: Application{Name: tt.appID}}
			tld, vendor, product, name := meta.Parts()
			if tld != tt.wantTLD || vendor != tt.wantVendor ||
				product != tt.wantProduct || name != tt.wantName {
				t.Errorf("Parts() = %v %v %v %v, want %v %v %v %v",
					tld, vendor, product, name,
					tt.wantTLD, tt.wantVendor, tt.wantProduct, tt.wantName)
			}
		})
	}
}
