// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package flatpak

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"gopkg.in/ini.v1"
)

type FlatpakMetadata struct {
	Application `ini:"Application"`
	Context     `ini:"Context"`
	DbusSession []Dbus      `ini:"-"`
	DbusSystem  []Dbus      `ini:"-"`
	Extensions  []Extension `ini:"-"`
	rootdir     *paths.Path
}

type Application struct {
	Name         string `ini:"name"`
	Runtime      string `ini:"runtime"`
	SDK          string `ini:"sdk"`
	Base         string `ini:"base,omitempty"`
	Command      string `ini:"command"`
	Abstractions []string
}

type Context struct {
	Shared      []string `ini:"shared,omitempty,allowshadow"`
	Sockets     []string `ini:"sockets,omitempty,allowshadow"`
	Devices     []string `ini:"devices,omitempty,allowshadow"`
	Filesystems []string `ini:"filesystems,omitempty,allowshadow"`
	Persistent  []string `ini:"persistent,omitempty,allowshadow"`
	Features    []string `ini:"features,omitempty,allowshadow"`
}

type Extension struct {
	Name           string `ini:"-"`
	Directory      string `ini:"directory"`
	AddLdPath      string `ini:"add-ld-path"`
	Version        string `ini:"version"`
	Autodelete     bool   `ini:"autodelete"`
	LocalSubset    bool   `ini:"local-subset"`
	NoAutodownload bool   `ini:"no-autodownload"`
}

type Dbus struct {
	Interface string `ini:"-"`
	Action    string `ini:"-"`
}

func parseSemicolonList(section *ini.Section, key string) []string {
	k := section.Key(key)
	if k == nil || k.String() == "" {
		return nil
	}
	return strings.Split(strings.TrimRight(k.String(), ";"), ";")
}

// LoadFlatpakMetadata reads and parses the metadata file for a Flatpak application,
// returning a structured FlatpakMetadata object.
func LoadFlatpakMetadata(path *paths.Path) (*FlatpakMetadata, error) {
	ini.PrettyFormat = false
	ini.PrettyEqual = false
	ini.PrettySection = false
	cfg, err := ini.LoadSources(
		ini.LoadOptions{
			Insensitive:              false,
			AllowBooleanKeys:         true,
			AllowShadows:             true,
			IgnoreInlineComment:      true,
			KeyValueDelimiters:       "=",
			KeyValueDelimiterOnWrite: "=",
		},
		path.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	meta := &FlatpakMetadata{}
	if err := cfg.MapTo(meta); err != nil {
		return nil, fmt.Errorf("failed to map INI to struct: %w", err)
	}

	// Manually parse the Context section as semicolon-separated lists
	contextSection := cfg.Section("Context")
	if contextSection != nil {
		meta.Shared = parseSemicolonList(contextSection, "shared")
		meta.Sockets = parseSemicolonList(contextSection, "sockets")
		meta.Devices = parseSemicolonList(contextSection, "devices")
		meta.Filesystems = parseSemicolonList(contextSection, "filesystems")
		meta.Persistent = parseSemicolonList(contextSection, "persistent")
		meta.Features = parseSemicolonList(contextSection, "features")
	}

	// Manually map extensions as they can have varying section names
	for _, section := range cfg.Sections() {
		// Identify an extension section
		switch {
		case strings.HasPrefix(section.Name(), "Extension "):
			var ext Extension
			if err := section.MapTo(&ext); err != nil {
				log.Printf("Could not map section %s to Extension: %v", section.Name(), err)
				continue
			}
			ext.Name = strings.TrimPrefix(section.Name(), "Extension ")
			meta.Extensions = append(meta.Extensions, ext)

		case section.Name() == "Session Bus Policy":
			for _, key := range section.Keys() {
				meta.DbusSession = append(meta.DbusSession, Dbus{
					Interface: key.Name(),
					Action:    key.String(),
				})
			}

		case section.Name() == "System Bus Policy":
			for _, key := range section.Keys() {
				meta.DbusSystem = append(meta.DbusSystem, Dbus{
					Interface: key.Name(),
					Action:    key.String(),
				})
			}
		}
	}

	meta.rootdir = paths.New("/var/lib/flatpak/app", meta.Name, "current/active/files")
	return meta, nil
}

// resolveNegations processes a list of items and removes any item that has a
// corresponding negation (prefixed with '!'). Both the negation and the
// negated item are removed from the result.
func resolveNegations(items []string) []string {
	negated := make(map[string]bool)
	for _, item := range items {
		if strings.HasPrefix(item, "!") {
			negated[strings.TrimPrefix(item, "!")] = true
		}
	}
	if len(negated) == 0 {
		return items
	}
	res := make([]string, 0, len(items))
	for _, item := range items {
		if strings.HasPrefix(item, "!") {
			continue
		}
		if negated[item] {
			continue
		}
		res = append(res, item)
	}
	return res
}

// LoadOverride checks if an override file exists for the application and, if so,
// loads it and merges its contents with the existing metadata, giving precedence
// to the override values.
func (f *FlatpakMetadata) LoadOverride(dir *paths.Path) error {
	path := dir.Join(f.Name)
	if !path.Exist() {
		return nil
	}
	override, err := LoadFlatpakMetadata(path)
	if err != nil {
		return err
	}

	// Merge Context fields from override, then resolve negations
	f.Shared = resolveNegations(append(f.Shared, override.Shared...))
	f.Sockets = resolveNegations(append(f.Sockets, override.Sockets...))
	f.Devices = resolveNegations(append(f.Devices, override.Devices...))
	f.Filesystems = resolveNegations(append(f.Filesystems, override.Filesystems...))
	f.Persistent = resolveNegations(append(f.Persistent, override.Persistent...))
	f.Features = resolveNegations(append(f.Features, override.Features...))

	// Merge DBUS policies
	f.DbusSession = append(f.DbusSession, override.DbusSession...)
	f.DbusSystem = append(f.DbusSystem, override.DbusSystem...)

	// Merge Extensions
	f.Extensions = append(f.Extensions, override.Extensions...)
	return nil
}

// EntryPoint retrieves the real entry point of the application
func (f *FlatpakMetadata) EntryPoint() string {
	command := f.rootdir.Join("bin").Join(f.Command)
	if ok, _ := command.IsSymlink(); ok {
		link, _ := os.Readlink(command.String())
		if filepath.IsAbs(link) {
			return link
		}
		return "/app/bin/" + link
	}
	return "/app/bin/" + f.Command
}

// AppBin retrieves additional binary path available within /app/
func (f *FlatpakMetadata) AppBin() []string {
	res := []string{}
	ignore := []string{
		"bin", "etc", "lib", "lib32", "lib64", "libexec", "share", "include", "share",
	}
	dirs, _ := f.rootdir.ReadDir(paths.FilterDirectories())
	for _, dir := range dirs {
		if !slices.Contains(ignore, dir.Base()) {
			res = append(res, dir.Base())
		}
	}
	return res
}

// Requirements retrieves additional requirements (baseapp, sockets, devices...)
// based on information it can collect in installation files
func (f *FlatpakMetadata) Requirements() {
	if f.Base != "" {
		return
	}

	patterns := map[string]string{
		"java":                    "abstractions/java",
		"chrome_":                 "app/org.chromium.Chromium",
		"com.valvesoftware.Steam": "app/com.valvesoftware.Steam",
	}
	_ = filepath.WalkDir(f.rootdir.String(), func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		name := d.Name()
		for pattern, base := range patterns {
			if strings.HasPrefix(name, pattern) {
				if strings.HasPrefix(base, "app/") {
					f.Base = base
				} else {
					f.Abstractions = append(f.Abstractions, base)
				}
				return filepath.SkipAll
			}
		}
		return nil
	})
}

// Parts retrieves the different parts of the application ID (TLD, Vendor, Product, Name)
func (f *FlatpakMetadata) Parts() (string, string, string, string) {
	var name string
	parts := strings.Split(f.Name, ".")
	if len(parts) < 4 {
		name = strings.ToLower(parts[2])
	} else {
		name = strings.ToLower(parts[3])
	}
	return parts[0], parts[1], parts[2], name
}
