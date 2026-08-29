// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func tempPath(t *testing.T, name string) *paths.Path {
	t.Helper()
	return paths.New(t.TempDir()).Join(name)
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

func writeFiles(t *testing.T, dir *paths.Path, files map[string]string) {
	t.Helper()
	for name, content := range files {
		writeFile(t, dir.Join(name), content)
	}
}

func writeLinks(t *testing.T, dir *paths.Path, links map[string]string) {
	t.Helper()
	for name, target := range links {
		p := dir.Join(name)
		if err := p.Parent().MkdirAll(); err != nil {
			t.Fatalf("mkdir %s: %v", p.Parent(), err)
		}
		if err := os.Symlink(target, p.String()); err != nil {
			t.Fatalf("symlink %s: %v", p, err)
		}
	}
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// chmod changes a path's mode for the test and restores a writable mode on
// cleanup so t.TempDir removal succeeds. Permission-based error tests are
// skipped for root, which bypasses file modes.
func chmod(t *testing.T, p *paths.Path, mode os.FileMode) {
	t.Helper()
	if os.Geteuid() == 0 {
		t.Skip("running as root: permission errors cannot be provoked")
	}
	if err := os.Chmod(p.String(), mode); err != nil {
		t.Fatalf("chmod %s: %v", p, err)
	}
	t.Cleanup(func() { _ = os.Chmod(p.String(), 0o755) })
}

func TestHashFile(t *testing.T) {
	tests := []struct {
		name    string
		content string
		missing bool
		want    string
		wantErr bool
	}{
		{name: "hash content", content: "hello", want: sha256Hex("hello")},
		{name: "missing file", missing: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tempPath(t, "f")
			if !tt.missing {
				writeFile(t, p, tt.content)
			}
			got, err := hashFile(p)
			if (err != nil) != tt.wantErr {
				t.Fatalf("hashFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("hashFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileIdentity(t *testing.T) {
	tests := []struct {
		name    string
		content string // regular file content
		link    string // symlink target, takes precedence over content
		missing bool
		want    string
		wantErr bool
	}{
		{name: "regular file hash", content: "hello", want: sha256Hex("hello")},
		{name: "symlink identity", link: "../x", want: symlinkPrefix + "../x"},
		{name: "missing file", missing: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tempPath(t, "f")
			switch {
			case tt.missing:
			case tt.link != "":
				if err := os.Symlink(tt.link, p.String()); err != nil {
					t.Fatalf("symlink: %v", err)
				}
			default:
				writeFile(t, p, tt.content)
			}
			got, err := fileIdentity(p)
			if (err != nil) != tt.wantErr {
				t.Fatalf("fileIdentity() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("fileIdentity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWriteManifest(t *testing.T) {
	tests := []struct {
		name         string
		entries      map[string]string
		stateDirFile bool // the state dir path is an existing regular file
		want         string
		wantErr      bool
	}{
		{
			name: "sorted by path",
			entries: map[string]string{
				"a/b.conf": "hash-b",
				"a/a.conf": "hash-a",
				"c.conf":   "hash-c",
			},
			want: "hash-a a/a.conf\nhash-b a/b.conf\nhash-c c.conf\n",
		},
		{
			name:         "state dir is a file",
			entries:      map[string]string{"p": "h"},
			stateDirFile: true,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateDir := paths.New(t.TempDir())
			if tt.stateDirFile {
				stateDir = tempPath(t, "state")
				writeFile(t, stateDir, "not a directory")
			}
			err := writeManifest(stateDir, tt.entries)
			if (err != nil) != tt.wantErr {
				t.Fatalf("writeManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			data, err := stateDir.Join(manifestFile).ReadFile()
			if err != nil {
				t.Fatalf("read manifest: %v", err)
			}
			if string(data) != tt.want {
				t.Errorf("writeManifest() wrote %q, want %q", string(data), tt.want)
			}
		})
	}
}

func TestReadManifest(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		missing    bool
		unreadable bool
		want       map[string]string
	}{
		{
			name:    "missing manifest",
			missing: true,
			want:    map[string]string{},
		},
		{
			name:    "round trip",
			content: "hash-a a/a.conf\nhash-b b.conf\n",
			want:    map[string]string{"a/a.conf": "hash-a", "b.conf": "hash-b"},
		},
		{
			name:    "malformed lines ignored",
			content: "hash1 path1\nmalformed-no-space\nhash2 path2\n",
			want:    map[string]string{"path1": "hash1", "path2": "hash2"},
		},
		{
			name:       "unreadable manifest",
			content:    "hash path\n",
			unreadable: true,
			want:       map[string]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateDir := paths.New(t.TempDir())
			if !tt.missing {
				writeFile(t, stateDir.Join(manifestFile), tt.content)
			}
			if tt.unreadable {
				chmod(t, stateDir.Join(manifestFile), 0o000)
			}
			got := readManifest(stateDir)
			if len(got) != len(tt.want) {
				t.Fatalf("readManifest() = %v, want %v", got, tt.want)
			}
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("readManifest()[%q] = %q, want %q", k, got[k], v)
				}
			}
		})
	}
}

func TestInstall(t *testing.T) {
	tests := []struct {
		name         string
		build        map[string]string // build dir files for the first install
		buildLinks   map[string]string // build dir symlinks: relative path -> link target
		target       map[string]string // target dir files present before the first install
		targetLinks  map[string]string // target dir symlinks present before the first install
		twice        bool              // run install a second time after mutate
		mutate       func(t *testing.T, buildDir, targetDir *paths.Path)
		wantChanged  bool              // changed result of the last install run
		wantFiles    map[string]string // target relative path -> expected content
		wantLinks    map[string]string // target relative path -> expected link target
		wantGone     []string          // target relative paths that must not exist
		wantManifest []string          // exact sorted manifest keys
	}{
		{
			name:        "adds files",
			build:       map[string]string{"profile.a": "content-a", "sub/profile.b": "content-b"},
			wantChanged: true,
			wantFiles:   map[string]string{"profile.a": "content-a", "sub/profile.b": "content-b"},
			wantManifest: []string{
				"profile.a", "sub/profile.b",
			},
		},
		{
			name: "skips untracked target file",
			build: map[string]string{
				"foreign": "ours",
			},
			target: map[string]string{
				"foreign": "theirs",
			},
			wantChanged: false,
			wantFiles: map[string]string{
				"foreign": "theirs",
			},
			wantManifest: []string{},
		},
		{
			name: "installs alongside untracked target file",
			build: map[string]string{
				"foreign": "ours",
				"mine":    "m",
			},
			target: map[string]string{
				"foreign": "theirs",
			},
			wantChanged: true,
			wantFiles: map[string]string{
				"foreign": "theirs",
				"mine":    "m",
			},
			wantManifest: []string{
				"mine",
			},
		},
		{
			name:        "noop when unchanged",
			build:       map[string]string{"p": "same"},
			twice:       true,
			wantChanged: false,
		},
		{
			name:  "updates modified file",
			build: map[string]string{"p": "v1"},
			twice: true,
			mutate: func(t *testing.T, buildDir, targetDir *paths.Path) {
				writeFile(t, buildDir.Join("p"), "v2")
			},
			wantChanged: true,
			wantFiles:   map[string]string{"p": "v2"},
		},
		{
			name:  "repairs drifted and missing target",
			build: map[string]string{"drifted": "good", "missing": "back"},
			twice: true,
			mutate: func(t *testing.T, buildDir, targetDir *paths.Path) {
				writeFile(t, targetDir.Join("drifted"), "tampered")
				if err := targetDir.Join("missing").Remove(); err != nil {
					t.Fatalf("remove target: %v", err)
				}
			},
			wantChanged: true,
			wantFiles:   map[string]string{"drifted": "good", "missing": "back"},
		},
		{
			name:  "removes stale files",
			build: map[string]string{"keep": "k", "drop": "d"},
			twice: true,
			mutate: func(t *testing.T, buildDir, targetDir *paths.Path) {
				if err := buildDir.Join("drop").Remove(); err != nil {
					t.Fatalf("remove from build: %v", err)
				}
			},
			wantChanged:  true,
			wantFiles:    map[string]string{"keep": "k"},
			wantGone:     []string{"drop"},
			wantManifest: []string{"keep"},
		},
		{
			name:         "installs disable link as symlink",
			build:        map[string]string{"systemd-detect-virt.apparmor.d": "profile\n"},
			buildLinks:   map[string]string{"disable/systemd-detect-virt": "../systemd-detect-virt"},
			wantChanged:  true,
			wantFiles:    map[string]string{"systemd-detect-virt.apparmor.d": "profile\n"},
			wantLinks:    map[string]string{"disable/systemd-detect-virt": "../systemd-detect-virt"},
			wantManifest: []string{"disable/systemd-detect-virt", "systemd-detect-virt.apparmor.d"},
		},
		{
			name:        "noop when link unchanged",
			buildLinks:  map[string]string{"disable/x": "../x"},
			twice:       true,
			wantChanged: false,
			wantLinks:   map[string]string{"disable/x": "../x"},
		},
		{
			name:         "keeps an already disabled profile link",
			build:        map[string]string{"hostname.apparmor.d": "profile\n"},
			buildLinks:   map[string]string{"disable/hostname": "../hostname"},
			targetLinks:  map[string]string{"disable/hostname": "/etc/apparmor.d/hostname"},
			wantChanged:  true,
			wantFiles:    map[string]string{"hostname.apparmor.d": "profile\n"},
			wantLinks:    map[string]string{"disable/hostname": "/etc/apparmor.d/hostname"},
			wantManifest: []string{"hostname.apparmor.d"},
		},
		{
			name:       "repairs link left as a regular file",
			buildLinks: map[string]string{"disable/x": "../x"},
			twice:      true,
			mutate: func(t *testing.T, buildDir, targetDir *paths.Path) {
				if err := targetDir.Join("disable/x").Remove(); err != nil {
					t.Fatalf("remove target link: %v", err)
				}
				writeFile(t, targetDir.Join("disable/x"), "")
			},
			wantChanged: true,
			wantLinks:   map[string]string{"disable/x": "../x"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildDir := paths.New(t.TempDir())
			targetDir := paths.New(t.TempDir())
			stateDir := paths.New(t.TempDir())
			writeFiles(t, buildDir, tt.build)
			writeLinks(t, buildDir, tt.buildLinks)
			writeFiles(t, targetDir, tt.target)
			writeLinks(t, targetDir, tt.targetLinks)

			changed, err := installProfiles(buildDir, targetDir, stateDir)
			if err != nil {
				t.Fatalf("install() error = %v", err)
			}
			if tt.twice {
				if tt.mutate != nil {
					tt.mutate(t, buildDir, targetDir)
				}
				changed, err = installProfiles(buildDir, targetDir, stateDir)
				if err != nil {
					t.Fatalf("install() second run error = %v", err)
				}
			}
			if changed != tt.wantChanged {
				t.Errorf("install() changed = %v, want %v", changed, tt.wantChanged)
			}
			for rel, want := range tt.wantFiles {
				got, err := targetDir.Join(rel).ReadFileAsString()
				if err != nil {
					t.Errorf("read target %s: %v", rel, err)
					continue
				}
				if got != want {
					t.Errorf("target %s = %q, want %q", rel, got, want)
				}
			}
			for rel, want := range tt.wantLinks {
				got, err := os.Readlink(targetDir.Join(rel).String())
				if err != nil {
					t.Errorf("readlink target %s: %v", rel, err)
					continue
				}
				if got != want {
					t.Errorf("target %s link = %q, want %q", rel, got, want)
				}
			}
			for _, rel := range tt.wantGone {
				if targetDir.Join(rel).Exist() {
					t.Errorf("target %s still exists, want removed", rel)
				}
			}
			if tt.wantManifest != nil {
				manifest := readManifest(stateDir)
				keys := make([]string, 0, len(manifest))
				for k := range manifest {
					keys = append(keys, k)
				}
				slices.Sort(keys)
				if !slices.Equal(keys, tt.wantManifest) {
					t.Errorf("manifest keys = %v, want %v", keys, tt.wantManifest)
				}
			}
		})
	}
}

func TestInstall_Errors(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) (buildDir, targetDir, stateDir *paths.Path)
	}{
		{
			name: "missing build dir",
			setup: func(t *testing.T) (*paths.Path, *paths.Path, *paths.Path) {
				return tempPath(t, "missing"), paths.New(t.TempDir()), paths.New(t.TempDir())
			},
		},
		{
			name: "unreadable build file",
			setup: func(t *testing.T) (*paths.Path, *paths.Path, *paths.Path) {
				buildDir := paths.New(t.TempDir())
				writeFile(t, buildDir.Join("p"), "x")
				chmod(t, buildDir.Join("p"), 0o000)
				return buildDir, paths.New(t.TempDir()), paths.New(t.TempDir())
			},
		},
		{
			name: "undeletable drifted target file",
			setup: func(t *testing.T) (*paths.Path, *paths.Path, *paths.Path) {
				buildDir := paths.New(t.TempDir())
				targetDir := paths.New(t.TempDir())
				stateDir := paths.New(t.TempDir())
				writeFile(t, buildDir.Join("p"), "x")
				if _, err := installProfiles(buildDir, targetDir, stateDir); err != nil {
					t.Fatalf("first install: %v", err)
				}
				writeFile(t, targetDir.Join("p"), "drifted")
				chmod(t, targetDir, 0o555)
				return buildDir, targetDir, stateDir
			},
		},
		{
			name: "unreadable target file",
			setup: func(t *testing.T) (*paths.Path, *paths.Path, *paths.Path) {
				buildDir := paths.New(t.TempDir())
				targetDir := paths.New(t.TempDir())
				stateDir := paths.New(t.TempDir())
				writeFile(t, buildDir.Join("p"), "x")
				if _, err := installProfiles(buildDir, targetDir, stateDir); err != nil {
					t.Fatalf("first install: %v", err)
				}
				chmod(t, targetDir.Join("p"), 0o000)
				return buildDir, targetDir, stateDir
			},
		},
		{
			name: "target subdir blocked by a file",
			setup: func(t *testing.T) (*paths.Path, *paths.Path, *paths.Path) {
				buildDir := paths.New(t.TempDir())
				targetDir := paths.New(t.TempDir())
				writeFile(t, buildDir.Join("sub/p"), "x")
				writeFile(t, targetDir.Join("sub"), "a file, not a directory")
				return buildDir, targetDir, paths.New(t.TempDir())
			},
		},
		{
			name: "read-only target dir",
			setup: func(t *testing.T) (*paths.Path, *paths.Path, *paths.Path) {
				buildDir := paths.New(t.TempDir())
				targetDir := paths.New(t.TempDir())
				writeFile(t, buildDir.Join("p"), "x")
				chmod(t, targetDir, 0o555)
				return buildDir, targetDir, paths.New(t.TempDir())
			},
		},
		{
			name: "undeletable stale file",
			setup: func(t *testing.T) (*paths.Path, *paths.Path, *paths.Path) {
				buildDir := paths.New(t.TempDir())
				targetDir := paths.New(t.TempDir())
				stateDir := paths.New(t.TempDir())
				writeFiles(t, buildDir, map[string]string{"keep": "k", "drop": "d"})
				if _, err := installProfiles(buildDir, targetDir, stateDir); err != nil {
					t.Fatalf("first install: %v", err)
				}
				if err := buildDir.Join("drop").Remove(); err != nil {
					t.Fatalf("remove from build: %v", err)
				}
				chmod(t, targetDir, 0o555)
				return buildDir, targetDir, stateDir
			},
		},
		{
			name: "state dir is a file",
			setup: func(t *testing.T) (*paths.Path, *paths.Path, *paths.Path) {
				buildDir := paths.New(t.TempDir())
				stateDir := tempPath(t, "state")
				writeFile(t, buildDir.Join("p"), "x")
				writeFile(t, stateDir, "not a directory")
				return buildDir, paths.New(t.TempDir()), stateDir
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildDir, targetDir, stateDir := tt.setup(t)
			if _, err := installProfiles(buildDir, targetDir, stateDir); err == nil {
				t.Errorf("install() error = nil, wantErr true")
			}
		})
	}
}

func TestAaStatus(t *testing.T) {
	tests := []struct {
		name       string
		manifest   map[string]string
		target     map[string]string
		unreadable []string // target files made unreadable
		wantErr    bool
	}{
		{
			name: "no manifest",
		},
		{
			name: "detects drift and missing",
			manifest: map[string]string{
				"ok":      sha256Hex("ok"),
				"drifted": sha256Hex("old-content"),
				"missing": sha256Hex("whatever"),
			},
			target: map[string]string{"ok": "ok", "drifted": "new-content"},
		},
		{
			name:       "unreadable target",
			manifest:   map[string]string{"p": sha256Hex("x")},
			target:     map[string]string{"p": "x"},
			unreadable: []string{"p"},
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetDir := paths.New(t.TempDir())
			stateDir := paths.New(t.TempDir())
			writeFiles(t, targetDir, tt.target)
			if tt.manifest != nil {
				if err := writeManifest(stateDir, tt.manifest); err != nil {
					t.Fatalf("writeManifest: %v", err)
				}
			}
			for _, rel := range tt.unreadable {
				chmod(t, targetDir.Join(rel), 0o000)
			}
			err := aaStatus(stateDir, targetDir, &conf{mode: "complain", include: "default"})
			if (err != nil) != tt.wantErr {
				t.Errorf("aaStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAaList(t *testing.T) {
	tests := []struct {
		name     string
		manifest map[string]string
		wantErr  bool
	}{
		{name: "empty manifest"},
		{name: "sorted paths", manifest: map[string]string{"b.conf": "h", "a.conf": "h"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateDir := paths.New(t.TempDir())
			if tt.manifest != nil {
				if err := writeManifest(stateDir, tt.manifest); err != nil {
					t.Fatalf("writeManifest: %v", err)
				}
			}
			err := aaList(stateDir, paths.New("/target"))
			if (err != nil) != tt.wantErr {
				t.Errorf("aaList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAaUninstall(t *testing.T) {
	tests := []struct {
		name           string
		manifest       map[string]string
		target         map[string]string
		readOnlyTarget bool // make the target dir read-only
		readOnlyState  bool // make the state dir read-only
		wantChanged    bool
		wantErr        bool
	}{
		{
			name: "no manifest",
		},
		{
			name:        "removes files and manifest",
			manifest:    map[string]string{"p1": sha256Hex("x"), "sub/p2": sha256Hex("y")},
			target:      map[string]string{"p1": "x", "sub/p2": "y"},
			wantChanged: true,
		},
		{
			name:           "undeletable target file",
			manifest:       map[string]string{"p": sha256Hex("x")},
			target:         map[string]string{"p": "x"},
			readOnlyTarget: true,
			wantErr:        true,
		},
		{
			name:          "undeletable manifest",
			manifest:      map[string]string{"p": "h"}, // listed file already gone
			readOnlyState: true,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			targetDir := paths.New(t.TempDir())
			stateDir := paths.New(t.TempDir())
			writeFiles(t, targetDir, tt.target)
			if tt.manifest != nil {
				if err := writeManifest(stateDir, tt.manifest); err != nil {
					t.Fatalf("writeManifest: %v", err)
				}
			}
			if tt.readOnlyTarget {
				chmod(t, targetDir, 0o555)
			}
			if tt.readOnlyState {
				chmod(t, stateDir, 0o555)
			}

			changed, err := aaUninstall(stateDir, targetDir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("aaUninstall() error = %v, wantErr %v", err, tt.wantErr)
			}
			if changed != tt.wantChanged {
				t.Errorf("aaUninstall() changed = %v, want %v", changed, tt.wantChanged)
			}
			if tt.wantErr || !tt.wantChanged {
				return
			}
			for rel := range tt.target {
				if targetDir.Join(rel).Exist() {
					t.Errorf("target %s still exists after uninstall", rel)
				}
			}
			if stateDir.Join(manifestFile).Exist() {
				t.Error("manifest file still exists after uninstall")
			}
		})
	}
}
