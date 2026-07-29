// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package detector

import (
	"errors"
	"os"
	"slices"
	"strings"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

// newTestSystem returns a System over temp trees, with a Run stub serving
// the given canned outputs, keyed by the full command line. Commands without
// a canned output fail.
func newTestSystem(t *testing.T, commands map[string]string) *System {
	t.Helper()
	return &System{
		Root:     paths.New(t.TempDir()),
		Apparmor: paths.New(t.TempDir()),
		Run: func(name string, arg ...string) (string, error) {
			cmdline := strings.Join(append([]string{name}, arg...), " ")
			out, ok := commands[cmdline]
			if !ok {
				return "", errors.New("stubbed: " + cmdline)
			}
			return out, nil
		},
	}
}

// loginctlCommands returns the Run outputs for a set of logind sessions,
// given per session id property values.
func loginctlCommands(sessions map[string]map[string]string) map[string]string {
	commands := map[string]string{}
	var list []string
	ids := make([]string, 0, len(sessions))
	for id := range sessions {
		ids = append(ids, id)
	}
	slices.Sort(ids)
	for _, id := range ids {
		list = append(list, id+" 1000 alex seat0 42 user tty2 no -")
		for prop, value := range sessions[id] {
			commands["loginctl show-session "+id+" -p "+prop+" --value"] = value + "\n"
		}
	}
	commands["loginctl list-sessions --no-legend"] = strings.Join(list, "\n") + "\n"
	return commands
}

// setupOSRelease stubs the tasks os-release globals.
func setupOSRelease(t *testing.T, family string, release map[string]string) {
	t.Helper()
	oldFamily, oldRelease := tasks.Family, tasks.Release
	tasks.Family, tasks.Release = family, release
	t.Cleanup(func() { tasks.Family, tasks.Release = oldFamily, oldRelease })
}

func symlink(t *testing.T, target string, link *paths.Path) {
	t.Helper()
	if err := link.Parent().MkdirAll(); err != nil {
		t.Fatalf("mkdir %s: %v", link.Parent(), err)
	}
	if err := os.Symlink(target, link.String()); err != nil {
		t.Fatalf("symlink %s: %v", link, err)
	}
}

func writeFile(t *testing.T, p *paths.Path, content string) {
	t.Helper()
	if err := p.Parent().MkdirAll(); err != nil {
		t.Fatalf("mkdir %s: %v", p.Parent(), err)
	}
	if err := p.WriteFile([]byte(content)); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
}

func TestRegister(t *testing.T) {
	for _, name := range []string{
		"ABI", "VERSION", "OS_FAMILY", "OS_ID", "OS_VERSION_ID",
		"DE", "DM", "DS", "GPU",
	} {
		t.Run(name, func(t *testing.T) {
			if _, ok := Registry[name]; !ok {
				t.Errorf("Registry[%s] = %v, want %v", name, "missing", "registered")
			}
		})
	}
}

func TestVersion_Detect(t *testing.T) {
	tests := []struct {
		name string
		out  string // apparmor_parser -V output; "" = command error
		want []string
	}{
		{name: "release", out: "AppArmor parser version 4.0.2\n...", want: []string{"4"}},
		{name: "beta", out: "AppArmor parser version 5.0.0~beta1\n", want: []string{"5"}},
		{name: "command error", want: nil},
		{name: "garbage output", out: "no digits here", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commands := map[string]string{}
			if tt.out != "" {
				commands["apparmor_parser -V"] = tt.out
			}
			sys := newTestSystem(t, commands)
			if got := (version{}).Detect(sys); !slices.Equal(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAbi_Detect(t *testing.T) {
	tests := []struct {
		name    string
		content string
		missing bool
		want    []string
	}{
		{name: "abi 5.0", content: "# header\n  abi <abi/5.0>,\n\n  /usr/** r,\n", want: []string{"5"}},
		{name: "abi 5", content: "abi <abi/5>,\n", want: []string{"5"}},
		{name: "no abi line", content: "# nothing\n", want: nil},
		{name: "missing file", missing: true, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := newTestSystem(t, nil)
			if !tt.missing {
				writeFile(t, sys.Apparmor.Join("abstractions/base-strict"), tt.content)
			}
			if got := (abi{}).Detect(sys); !slices.Equal(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOsRelease_Detect(t *testing.T) {
	tests := []struct {
		name     string
		detector osRelease
		family   string
		release  map[string]string
		want     []string
	}{
		{name: "family", detector: osRelease{name: "OS_FAMILY"}, family: "pacman", want: []string{"pacman"}},
		{name: "id", detector: osRelease{name: "OS_ID", key: "ID"}, release: map[string]string{"ID": "arch"}, want: []string{"arch"}},
		{name: "version id", detector: osRelease{name: "OS_VERSION_ID", key: "VERSION_ID"}, release: map[string]string{"VERSION_ID": "rolling"}, want: []string{"rolling"}},
		{name: "unknown", detector: osRelease{name: "OS_ID", key: "ID"}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupOSRelease(t, tt.family, tt.release)
			if got := tt.detector.Detect(newTestSystem(t, nil)); !slices.Equal(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisplayManager_Detect(t *testing.T) {
	tests := []struct {
		name   string
		target string // display-manager.service symlink target; "" = no symlink
		want   []string
	}{
		{name: "gdm enabled", target: "/usr/lib/systemd/system/gdm.service", want: []string{"gdm"}},
		{name: "sddm enabled", target: "/usr/lib/systemd/system/sddm.service", want: []string{"sddm"}},
		{name: "gdm3 unit name normalized", target: "/lib/systemd/system/gdm3.service", want: []string{"gdm"}},
		{name: "no display manager", want: []string{"none"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := newTestSystem(t, nil)
			if tt.target != "" {
				symlink(t, tt.target, sys.Root.Join("etc/systemd/system/display-manager.service"))
			}
			if got := (displayManager{}).Detect(sys); !slices.Equal(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDesktop_Detect(t *testing.T) {
	tests := []struct {
		name      string
		sessions  map[string]map[string]string // session id -> property -> value
		processes []string                     // comm of running processes
		want      []string
	}{
		{
			name:     "session desktop",
			sessions: map[string]map[string]string{"2": {"Desktop": "KDE"}},
			want:     []string{"kde"},
		},
		{
			name:     "colon separated desktop list",
			sessions: map[string]map[string]string{"2": {"Desktop": "ubuntu:GNOME"}},
			want:     []string{"gnome"},
		},
		{
			name: "multiple sessions deduped",
			sessions: map[string]map[string]string{
				"2": {"Desktop": "GNOME"},
				"3": {"Desktop": "gnome"},
				"4": {"Desktop": "plasma"},
			},
			want: []string{"gnome", "kde"},
		},
		{
			name:      "empty session desktop falls back to running shell",
			sessions:  map[string]map[string]string{"2": {"Desktop": ""}},
			processes: []string{"systemd", "gnome-shell"},
			want:      []string{"gnome"},
		},
		{
			name:      "no logind falls back to running shell",
			processes: []string{"cosmic-session"},
			want:      []string{"cosmic"},
		},
		{name: "nothing detected", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := newTestSystem(t, loginctlCommands(tt.sessions))
			for i, comm := range tt.processes {
				writeFile(t, sys.Root.Join("proc", string(rune('1'+i)), "comm"), comm+"\n")
			}
			if got := (desktop{}).Detect(sys); !slices.Equal(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDisplayServer_Detect(t *testing.T) {
	tests := []struct {
		name     string
		sessions map[string]map[string]string // session id -> property -> value
		want     []string
	}{
		{
			name:     "wayland session",
			sessions: map[string]map[string]string{"2": {"Type": "wayland"}},
			want:     []string{"wayland"},
		},
		{
			name: "tty sessions filtered, list deduped",
			sessions: map[string]map[string]string{
				"1": {"Type": "tty"},
				"2": {"Type": "wayland"},
				"3": {"Type": "x11"},
				"4": {"Type": "wayland"},
			},
			want: []string{"wayland", "x11"},
		},
		{name: "no session", want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := newTestSystem(t, loginctlCommands(tt.sessions))
			if got := (displayServer{}).Detect(sys); !slices.Equal(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGpu_Detect(t *testing.T) {
	tests := []struct {
		name    string
		vendors map[string]string // card name -> vendor file content
		want    []string
	}{
		{name: "nvidia", vendors: map[string]string{"card0": "0x10de\n"}, want: []string{"nvidia"}},
		{name: "hybrid intel nvidia", vendors: map[string]string{"card0": "0x8086\n", "card1": "0x10de\n"}, want: []string{"intel", "nvidia"}},
		{name: "amd", vendors: map[string]string{"card0": "0x1002\n"}, want: []string{"amd"}},
		{name: "duplicated vendor", vendors: map[string]string{"card0": "0x1002\n", "card1": "0x1002\n"}, want: []string{"amd"}},
		{name: "unknown vendor skipped", vendors: map[string]string{"card0": "0xdead\n"}, want: []string{"none"}},
		{name: "no drm", want: []string{"none"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sys := newTestSystem(t, nil)
			for card, vendor := range tt.vendors {
				writeFile(t, sys.Root.Join("sys/class/drm/"+card+"/device/vendor"), vendor)
			}
			if got := (gpu{}).Detect(sys); !slices.Equal(got, tt.want) {
				t.Errorf("Detect() = %v, want %v", got, tt.want)
			}
		})
	}
}
