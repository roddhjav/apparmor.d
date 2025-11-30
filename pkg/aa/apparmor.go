// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

// MagicRoot is the default Apparmor magic directory: /etc/apparmor.d/.
var MagicRoot = paths.New("/etc/apparmor.d")

// FileKind represents an AppArmor file kind.
type FileKind uint8

const (
	ProfileKind FileKind = iota
	AbstractionKind
	TunableKind
)

func KindFromPath(file *paths.Path) FileKind {
	dirname := file.Parent().String()
	switch {
	case strings.Contains(dirname, "abstractions"):
		return AbstractionKind
	case strings.Contains(dirname, "tunables"):
		return TunableKind
	case strings.Contains(dirname, "local"):
		return AbstractionKind
	case strings.Contains(dirname, "mappings"):
		return AbstractionKind
	default:
		return ProfileKind
	}
}

// AppArmorProfileFiles represents a full set of apparmor profiles
type AppArmorProfileFiles map[string]*AppArmorProfileFile

// AppArmorProfileFile represents a full apparmor profile file.
// Warning: close to the BNF grammar of apparmor profile but not exactly the same (yet):
//   - Some rules are not supported yet (subprofile, hat...)
//   - The structure is simplified as it only aims at writing profile, not parsing it.
type AppArmorProfileFile struct {
	Preamble Rules
	Profiles []*Profile
}

func NewAppArmorProfile() *AppArmorProfileFile {
	return &AppArmorProfileFile{}
}

// DefaultTunables return a minimal working profile to build the profile
// It should not be used when loading file from /etc/apparmor.d
func DefaultTunables() *AppArmorProfileFile {
	return &AppArmorProfileFile{
		Preamble: Rules{
			&Variable{Name: "arch", Values: []string{"x86_64", "amd64", "i386"}, Define: true},
			&Variable{Name: "bin", Values: []string{"/{,usr/}bin"}, Define: true},
			&Variable{Name: "c", Values: []string{"[0-9a-zA-Z]"}, Define: true},
			&Variable{Name: "dpkg_script_ext", Values: []string{"config", "templates", "preinst", "postinst", "prerm", "postrm"}, Define: true},
			&Variable{Name: "etc_ro", Values: []string{"/{,usr/}etc/"}, Define: true},
			&Variable{Name: "HOME", Values: []string{"/home/*"}, Define: true},
			&Variable{Name: "int", Values: []string{"[0-9]{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}"}, Define: true},
			&Variable{Name: "int2", Values: []string{"[0-9][0-9]"}, Define: true},
			&Variable{Name: "lib", Values: []string{"/{,usr/}lib{,exec,32,64}"}, Define: true},
			&Variable{Name: "MOUNTS", Values: []string{"/media/*/", "/run/media/*/*/", "/mnt/*/"}, Define: true},
			&Variable{Name: "multiarch", Values: []string{"*-linux-gnu*"}, Define: true},
			&Variable{Name: "rand", Values: []string{"@{c}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}"}, Define: true}, // Up to 10 characters
			&Variable{Name: "run", Values: []string{"/run/", "/var/run/"}, Define: true},
			&Variable{Name: "uid", Values: []string{"{[0-9],[1-9][0-9],[1-9][0-9][0-9],[1-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9],[1-4][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9][0-9]}"}, Define: true},
			&Variable{Name: "user_cache_dirs", Values: []string{"/home/*/.cache"}, Define: true},
			&Variable{Name: "user_config_dirs", Values: []string{"/home/*/.config"}, Define: true},
			&Variable{Name: "user_share_dirs", Values: []string{"/home/*/.local/share"}, Define: true},
			&Variable{Name: "user", Values: []string{"[a-zA-Z_]{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}"}, Define: true},
			&Variable{Name: "version", Values: []string{"@{int}{.@{int},}{.@{int},}{-@{rand},}"}, Define: true},
			&Variable{Name: "w", Values: []string{"[a-zA-Z0-9_]"}, Define: true},
			&Variable{Name: "word", Values: []string{"@{w}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}"}, Define: true},
		},
	}
}

// String returns the formatted representation of a profile file as a string
func (f *AppArmorProfileFile) String() string {
	return renderTemplate("apparmor", f)
}

// Validate the profile file
func (f *AppArmorProfileFile) Validate() error {
	if err := f.Preamble.Validate(); err != nil {
		return err
	}
	for _, p := range f.Profiles {
		if err := p.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// GetDefaultProfile ensure a profile is always present in the profile file and
// return it, as a default profile.
func (f *AppArmorProfileFile) GetDefaultProfile() *Profile {
	if len(f.Profiles) == 0 {
		f.Profiles = append(f.Profiles, &Profile{})
	}
	return f.Profiles[0]
}

// Sort the rules in the profile
// Follow: https://apparmor.pujol.io/development/guidelines/#guidelines
func (f *AppArmorProfileFile) Sort() {
	for _, p := range f.Profiles {
		p.Sort()
	}
}

// MergeRules merge similar rules together.
// Steps:
//   - Remove identical rules
//   - Merge rule access. Eg: for same path, 'r' and 'w' becomes 'rw'
//
// Note: logs.regCleanLogs helps a lot to do a first cleaning
func (f *AppArmorProfileFile) MergeRules() {
	for _, p := range f.Profiles {
		p.Merge(nil)
	}
}

// Format the profile for better readability before printing it.
// Follow: https://apparmor.pujol.io/development/guidelines/#the-file-block
func (f *AppArmorProfileFile) Format() {
	for _, p := range f.Profiles {
		p.Format()
	}
}

// Merge merges two profiles together.
func (f *AppArmorProfileFile) Merge(other *AppArmorProfileFile) error {
	f.Preamble = append(f.Preamble, other.Preamble...)
	f.Profiles = append(f.Profiles, other.Profiles...)
	return nil
}

// Clean the profile file from comments
func (f *AppArmorProfileFile) Clean() {
	delete := []int{}
	for i, r := range f.Preamble {
		switch r.(type) {
		case *Comment:
			delete = append(delete, i)
		}
	}
	for i := len(delete) - 1; i >= 0; i-- {
		f.Preamble = slices.Delete(f.Preamble, delete[i], delete[i]+1)
	}
}
