// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package flatpak

import (
	"testing"
)

func TestExpandBraces(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "no braces",
			input: "org.freedesktop.UPower",
			want:  []string{"org.freedesktop.UPower"},
		},
		{
			name:  "trailing alternation",
			input: "org.kde.kwalletd{,5,6}",
			want:  []string{"org.kde.kwalletd", "org.kde.kwalletd5", "org.kde.kwalletd6"},
		},
		{
			name:  "mid-name alternation",
			input: "org.a11y.{B,b}us",
			want:  []string{"org.a11y.Bus", "org.a11y.bus"},
		},
		{
			name:  "multiple brace groups",
			input: "org.freedesktop.{S,s}ecret{,s}",
			want: []string{
				"org.freedesktop.Secret", "org.freedesktop.Secrets",
				"org.freedesktop.secret", "org.freedesktop.secrets",
			},
		},
		{
			name:  "trailing alternation with values",
			input: "org.freedesktop.portal.{Documents,FileTransfer}",
			want:  []string{"org.freedesktop.portal.Documents", "org.freedesktop.portal.FileTransfer"},
		},
		{
			name:  "unclosed brace",
			input: "org.broken.{test",
			want:  []string{"org.broken.{test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandBraces(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("expandBraces(%q) = %v, want %v", tt.input, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("expandBraces(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestAddDbusBind(t *testing.T) {
	tests := []struct {
		name        string
		dbusName    string
		profileName string
		wantExact   map[string][]string
		wantPrefix  map[string][]string
	}{
		{
			name:        "exact name",
			dbusName:    "org.bluez",
			profileName: "bluetoothd",
			wantExact:   map[string][]string{"org.bluez": {"bluetoothd"}},
			wantPrefix:  map[string][]string{},
		},
		{
			name:        "trailing @{var} becomes prefix",
			dbusName:    "org.gnome.evolution.dataserver.Sources@{int}",
			profileName: "evolution-source-registry",
			wantExact:   map[string][]string{},
			wantPrefix:  map[string][]string{"org.gnome.evolution.dataserver.Sources": {"evolution-source-registry"}},
		},
		{
			name:        "brace alternation expands to exact",
			dbusName:    "org.kde.kwalletd{,5,6}",
			profileName: "kwalletd",
			wantExact: map[string][]string{
				"org.kde.kwalletd":  {"kwalletd"},
				"org.kde.kwalletd5": {"kwalletd"},
				"org.kde.kwalletd6": {"kwalletd"},
			},
			wantPrefix: map[string][]string{},
		},
		{
			name:        "mid-name alternation expands to exact",
			dbusName:    "org.a11y.{B,b}us",
			profileName: "dbus-accessibility",
			wantExact: map[string][]string{
				"org.a11y.Bus": {"dbus-accessibility"},
				"org.a11y.bus": {"dbus-accessibility"},
			},
			wantPrefix: map[string][]string{},
		},
		{
			name:        "mid-name @{var} skipped entirely",
			dbusName:    "org.mpris.MediaPlayer2.@{profile_name}",
			profileName: "mpris-proxy",
			wantExact:   map[string][]string{},
			wantPrefix:  map[string][]string{"org.mpris.MediaPlayer2.": {"mpris-proxy"}},
		},
		{
			name:        "trailing rule terminator stripped",
			dbusName:    "org.bluez,",
			profileName: "bluetoothd",
			wantExact:   map[string][]string{"org.bluez": {"bluetoothd"}},
			wantPrefix:  map[string][]string{},
		},
		{
			name:        "{,.*} expands to exact and dotted prefix",
			dbusName:    "org.freedesktop.Flatpak{,.*},",
			profileName: "flatpak-session-helper",
			wantExact:   map[string][]string{"org.freedesktop.Flatpak": {"flatpak-session-helper"}},
			wantPrefix:  map[string][]string{"org.freedesktop.Flatpak.": {"flatpak-session-helper"}},
		},
		{
			name:        "nested braces with {,.*} suffix",
			dbusName:    "org.kde.kwalletd{,5,6}{,.*},",
			profileName: "kwalletd",
			wantExact: map[string][]string{
				"org.kde.kwalletd":  {"kwalletd"},
				"org.kde.kwalletd5": {"kwalletd"},
				"org.kde.kwalletd6": {"kwalletd"},
			},
			wantPrefix: map[string][]string{
				"org.kde.kwalletd.":  {"kwalletd"},
				"org.kde.kwalletd5.": {"kwalletd"},
				"org.kde.kwalletd6.": {"kwalletd"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exact := map[string][]string{}
			prefix := map[string][]string{}
			addDbusBind(tt.dbusName, tt.profileName, exact, prefix)

			for k, v := range tt.wantExact {
				if got, ok := exact[k]; !ok {
					t.Errorf("missing exact key %q", k)
				} else if len(got) != len(v) || got[0] != v[0] {
					t.Errorf("exact[%q] = %v, want %v", k, got, v)
				}
			}
			for k := range exact {
				if _, ok := tt.wantExact[k]; !ok {
					t.Errorf("unexpected exact key %q = %v", k, exact[k])
				}
			}

			for k, v := range tt.wantPrefix {
				if got, ok := prefix[k]; !ok {
					t.Errorf("missing prefix key %q", k)
				} else if len(got) != len(v) || got[0] != v[0] {
					t.Errorf("prefix[%q] = %v, want %v", k, got, v)
				}
			}
			for k := range prefix {
				if _, ok := tt.wantPrefix[k]; !ok {
					t.Errorf("unexpected prefix key %q = %v", k, prefix[k])
				}
			}
		})
	}
}

func TestLookupDbusLabel(t *testing.T) {
	// Save and restore global state
	origExact := dbusLabelMap
	origPrefix := dbusLabelPrefixMap
	origDefault := dbusLabelDefault
	defer func() {
		dbusLabelMap = origExact
		dbusLabelPrefixMap = origPrefix
		dbusLabelDefault = origDefault
	}()

	// Mark the lazy scan as done so it does not overwrite the test maps.
	dbusLabelOnce.Do(func() {})

	dbusLabelMap = map[string]string{
		"org.bluez":          "bluetoothd",
		"org.kde.kwalletd":   "kwalletd",
		"org.kde.kwalletd5":  "kwalletd",
		"org.gnome.Software": "gnome-software",
	}
	dbusLabelPrefixMap = map[string]string{
		"org.gnome.evolution.dataserver.Sources": "evolution-source-registry",
		"org.mpris.MediaPlayer2.":                "mpris-proxy",
	}
	dbusLabelDefault = map[string]string{
		"org.kde.KGlobalSettings": "kded",
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"exact match", "org.bluez", "bluetoothd"},
		{"exact match kwalletd", "org.kde.kwalletd", "kwalletd"},
		{"prefix match Sources5", "org.gnome.evolution.dataserver.Sources5", "evolution-source-registry"},
		{"prefix match mpris", "org.mpris.MediaPlayer2.firefox", "mpris-proxy"},
		{"default fallback", "org.kde.KGlobalSettings", "kded"},
		{"unconfined fallback", "org.nonexistent.Service", "unconfined"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := lookupDbusLabel(tt.input)
			if got != tt.want {
				t.Errorf("lookupDbusLabel(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestOwnersToLabelMap(t *testing.T) {
	tests := []struct {
		name   string
		owners map[string][]string
		want   map[string]string
	}{
		{
			name:   "single owner",
			owners: map[string][]string{"org.bluez": {"bluetoothd"}},
			want:   map[string]string{"org.bluez": "bluetoothd"},
		},
		{
			name:   "multiple owners sorted",
			owners: map[string][]string{"org.kde.StatusNotifierWatcher": {"xfce-panel", "gnome-shell", "kded"}},
			want:   map[string]string{"org.kde.StatusNotifierWatcher": `"{gnome-shell,kded,xfce-panel}"`},
		},
		{
			name:   "empty",
			owners: map[string][]string{},
			want:   map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ownersToLabelMap(tt.owners)
			if len(got) != len(tt.want) {
				t.Fatalf("ownersToLabelMap() = %v, want %v", got, tt.want)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("ownersToLabelMap()[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}
