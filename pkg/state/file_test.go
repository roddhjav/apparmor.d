// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package state

import (
	"slices"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

const stateFixture = `# System states
#
# Comment block.

# AppArmor ABI
@{ABI} = 5

# The display server in use.
@{DS} = wayland

#${FSP} = false

include if exists <tunables/multiarch.d/state.d>
`

func writeState(t *testing.T, content string) *paths.Path {
	t.Helper()
	p := paths.New(t.TempDir()).Join("state")
	if err := p.WriteFile([]byte(content)); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
	return p
}

func TestKey(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "bare name", in: "DS", want: "@{DS}"},
		{name: "wrapped name", in: "@{DS}", want: "@{DS}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Key(tt.in); got != tt.want {
				t.Errorf("Key() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		missing bool
		wantErr bool
	}{
		{name: "existing file"},
		{name: "missing file", missing: true, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := writeState(t, stateFixture)
			if tt.missing {
				p = paths.New(t.TempDir()).Join("absent")
			}
			_, err := Load(p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFile_Get(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		want    string
		wantOk  bool
	}{
		{name: "existing bare", varName: "ABI", want: "5", wantOk: true},
		{name: "existing wrapped", varName: "@{DS}", want: "wayland", wantOk: true},
		{name: "commented-out is not defined", varName: "FSP"},
		{name: "unknown", varName: "NOPE"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := Load(writeState(t, stateFixture))
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			got, ok := f.Get(tt.varName)
			if got != tt.want || ok != tt.wantOk {
				t.Errorf("Get() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestFile_Names(t *testing.T) {
	f, err := Load(writeState(t, stateFixture))
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	want := []string{"@{ABI}", "@{DS}"}
	if got := f.Names(); !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

func TestFile_SetSave(t *testing.T) {
	tests := []struct {
		name      string
		varName   string
		value     string
		wantFound bool
		want      string // full content after Save, in the aa canonical format
	}{
		{
			name: "replace value in place", varName: "DS", value: "x11", wantFound: true,
			want: `# System states

# Comment block.
# AppArmor ABI
@{ABI} = 5
# The display server in use.
@{DS} = x11
#${FSP} = false

include if exists <tunables/multiarch.d/state.d>
`,
		},
		{
			name: "unknown name untouched", varName: "NOPE", value: "1",
			want: `# System states

# Comment block.
# AppArmor ABI
@{ABI} = 5
# The display server in use.
@{DS} = wayland
#${FSP} = false

include if exists <tunables/multiarch.d/state.d>
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := writeState(t, stateFixture)
			f, err := Load(p)
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if found := f.Set(tt.varName, tt.value); found != tt.wantFound {
				t.Fatalf("Set() = %v, want %v", found, tt.wantFound)
			}
			if err := f.Save(); err != nil {
				t.Fatalf("Save() error = %v", err)
			}
			want := tt.want
			if want == "" {
				want = stateFixture
			}
			got, err := p.ReadFileAsString()
			if err != nil {
				t.Fatalf("read state: %v", err)
			}
			if got != want {
				t.Errorf("Save() wrote %q, want %q", got, want)
			}
		})
	}
}

func TestFile_Add(t *testing.T) {
	tests := []struct {
		name    string
		initial string // "" means NewFile
		varName string
		value   string
		want    string
	}{
		{
			name: "append to new file", varName: "with_heroic", value: "true",
			want: "@{with_heroic} = true\n",
		},
		{
			name: "append keeps trailing newline", initial: "@{A} = 1\n",
			varName: "B", value: "2", want: "@{A} = 1\n@{B} = 2\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := paths.New(t.TempDir()).Join("state.d/aa-config")
			var f *File
			if tt.initial == "" {
				f = NewFile(p)
			} else {
				writeFileHelper(t, p, tt.initial)
				var err error
				if f, err = Load(p); err != nil {
					t.Fatalf("Load() error = %v", err)
				}
			}
			f.Add(tt.varName, tt.value)
			if err := f.Save(); err != nil {
				t.Fatalf("Save() error = %v", err)
			}
			got, err := p.ReadFileAsString()
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			if got != tt.want {
				t.Errorf("Add()+Save() wrote %q, want %q", got, tt.want)
			}
		})
	}
}

func writeFileHelper(t *testing.T, p *paths.Path, content string) {
	t.Helper()
	if err := p.Parent().MkdirAll(); err != nil {
		t.Fatalf("mkdir %s: %v", p.Parent(), err)
	}
	if err := p.WriteFile([]byte(content)); err != nil {
		t.Fatalf("write %s: %v", p, err)
	}
}
