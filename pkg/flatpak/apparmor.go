// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Package flatpak generates AppArmor profiles for Flatpak applications by
// reading their metadata files and translating sandbox permissions
// (filesystems, devices, sockets, D-Bus access) into AppArmor rules.
//
// Experimental and currently not enabled: the API and generated rules are
// subject to change.
package flatpak

import (
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
)

type profileID int

const (
	pApp profileID = iota
	pDbusProxy
)

var (
	allowedFeatures = []string{"devel", "multiarch", "bluetooth", "canbus", "per-app-dev-shm"}
	allowedShared   = []string{"network", "ipc"}
	allowedSockets  = []string{
		"x11", "wayland", "fallback-x11", "pulseaudio", "session-bus",
		"system-bus", "ssh-auth", "pcsc", "cups", "gpg-agent", "inherit-wayland-socket",
	}
	allowedDevices = []string{"dri", "input", "usb", "kvm", "all", "shm"}

	fsMap = map[string]string{
		"xdg-desktop":      "@{user_desktop_dirs}",
		"xdg-documents":    "@{user_documents_dirs}",
		"xdg-download":     "@{user_download_dirs}",
		"xdg-music":        "@{user_music_dirs}",
		"xdg-pictures":     "@{user_pictures_dirs}",
		"xdg-public-share": "@{user_publicshare_dirs}",
		"xdg-videos":       "@{user_videos_dirs}",
		"xdg-templates":    "@{user_templates_dirs}",
		"xdg-config":       "@{user_config_dirs}",
		"xdg-cache":        "@{user_cache_dirs}",
		"xdg-data":         "@{user_share_dirs}",
	}
)

type FlatpakAppArmorProfile struct {
	aa.AppArmorProfileFile
	Metadata *FlatpakMetadata
	FileName string
}

// ProfileName returns the name of the profile generated for an application.
func ProfileName(appID string) string {
	return "flatpak." + appID
}

func NewFlatpakAppArmorProfile(meta *FlatpakMetadata, mode string) *FlatpakAppArmorProfile {
	appID := meta.Name
	tld, vendor, product, name := meta.Parts()

	// Define profile names
	profileName := ProfileName(appID)
	profileDbusProxyName := "flatpak.dbus." + appID

	att := "/att/" + profileName
	appFlags := []string{"attach_disconnected.path=" + att, "mediate_deleted"}
	dbusFlags := []string{"attach_disconnected.path=" + att}
	if mode == "complain" {
		appFlags = append(appFlags, "complain")
		dbusFlags = append(dbusFlags, "complain")
	}

	return &FlatpakAppArmorProfile{
		FileName: profileName,
		Metadata: meta,
		AppArmorProfileFile: aa.AppArmorProfileFile{
			Kind: aa.ProfileKind,
			Preamble: aa.Rules{
				&aa.Comment{Base: aa.Base{Comment: " apparmor.d - Full set of apparmor profiles", IsLineRule: true}},
				&aa.Comment{Base: aa.Base{Comment: " Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>", IsLineRule: true}},
				&aa.Comment{Base: aa.Base{Comment: " SPDX-License-Identifier: GPL-2.0-only", IsLineRule: true}},
				&aa.Comment{Base: aa.Base{Comment: " vim:syntax=apparmor", IsLineRule: true}},
				&aa.Comment{Base: aa.Base{Comment: "", IsLineRule: true}},
				&aa.Comment{Base: aa.Base{Comment: " This profile has been generated with aa-flatpak.", IsLineRule: true}},
				nil,
				&aa.Abi{IsMagic: true, Path: "abi/5.0"},
				&aa.Include{IsMagic: true, Path: "tunables/global"},
				&aa.Variable{Name: "tld", Values: []string{tld}, Define: true},
				&aa.Variable{Name: "vendor", Values: []string{vendor}, Define: true},
				&aa.Variable{Name: "product", Values: []string{product}, Define: true},
				&aa.Variable{Name: "name", Values: []string{name}, Define: true},
				&aa.Variable{Name: "appid", Values: []string{appID}, Define: true},
				&aa.Variable{Name: "command", Values: []string{meta.Command}, Define: true},
				&aa.Variable{Name: "profile_app", Values: []string{profileName}, Define: true},
				&aa.Variable{Name: "profile_dbus", Values: []string{profileDbusProxyName}, Define: true},
				&aa.Variable{Name: "att", Values: []string{att}, Define: true},
			},
			Profiles: []*aa.Profile{
				{
					Header: aa.Header{
						Name:  profileName,
						Flags: appFlags,
					},
					Rules: aa.Rules{
						&aa.Include{IsMagic: true, Path: "abstractions/attached/base"},
						&aa.Include{IsMagic: true, Path: "abstractions/flatpak/base"},
					},
				},
				{
					Header: aa.Header{
						Name:  profileDbusProxyName,
						Flags: dbusFlags,
					},
					Rules: aa.Rules{
						&aa.Include{IsMagic: true, Path: "abstractions/attached/base"},
						&aa.Include{IsMagic: true, Path: "abstractions/flatpak/dbus-proxy"},
					},
				},
			},
		},
	}
}

func (p *FlatpakAppArmorProfile) Generate() error {
	// Profile App
	app := p.Profiles[pApp]
	app.Rules = slices.Concat(app.Rules,
		p.addRuntime(),
		p.addBaseApp(),
		p.addAbstractions(),
		p.addShared(),
		p.addSockets(),
		p.addDevices(),
		p.addFeatures(),
		p.addFileSystem(),
		p.finalize(app.Name),
	)
	app.Merge(nil)
	app.Sort()
	if err := app.Validate(); err != nil {
		return fmt.Errorf("validation error in profile %s: %w", app.Name, err)
	}

	// Profile dbus proxy
	proxy := p.Profiles[pDbusProxy]
	dbus, err := p.addDbus()
	if err != nil {
		return err
	}
	proxy.Rules = slices.Concat(proxy.Rules, dbus, p.finalize(proxy.Name))
	proxy.Merge(nil)
	// proxy.Sort()
	if err := proxy.Validate(); err != nil {
		return fmt.Errorf("validation error in profile %s: %w", proxy.Name, err)
	}
	return nil
}

func (p *FlatpakAppArmorProfile) addRuntime() aa.Rules {
	runtime, _, _ := strings.Cut(p.Metadata.Runtime, "/")
	platform := strings.TrimSuffix(runtime, ".Platform")
	return aa.Rules{
		&aa.Include{IsMagic: true, Path: "abstractions/flatpak/platform/" + platform},
	}
}

func (p *FlatpakAppArmorProfile) addBaseApp() aa.Rules {
	if len(p.Metadata.Base) == 0 {
		return aa.Rules{}
	}
	baseapp := strings.Split(p.Metadata.Base, "/")[1]
	baseapp = strings.TrimSuffix(baseapp, ".BaseApp")
	return aa.Rules{
		&aa.Include{IsMagic: true, Path: "abstractions/flatpak/baseapp/" + baseapp},
	}
}

func (p *FlatpakAppArmorProfile) addAbstractions() aa.Rules {
	rules := aa.Rules{}
	for _, abs := range p.Metadata.Abstractions {
		rules = append(rules, &aa.Include{IsMagic: true, Path: abs})
	}
	return rules
}

func addIncludes(items, allowed []string, prefix string) aa.Rules {
	rules := aa.Rules{}
	for _, item := range items {
		if !slices.Contains(allowed, item) {
			log.Printf("Warning: unknown %s item: %s", prefix, item)
			continue
		}
		rules = append(rules, &aa.Include{
			IsMagic: true, Path: "abstractions/flatpak/" + prefix + "/" + item,
		})
	}
	return rules
}

func (p *FlatpakAppArmorProfile) addShared() aa.Rules {
	return addIncludes(p.Metadata.Shared, allowedShared, "shared")
}

func (p *FlatpakAppArmorProfile) addSockets() aa.Rules {
	return addIncludes(p.Metadata.Sockets, allowedSockets, "sockets")
}

func (p *FlatpakAppArmorProfile) addDevices() aa.Rules {
	return addIncludes(p.Metadata.Devices, allowedDevices, "devices")
}

func (p *FlatpakAppArmorProfile) addFeatures() aa.Rules {
	return addIncludes(p.Metadata.Features, allowedFeatures, "features")
}

func (p *FlatpakAppArmorProfile) addDbus() (aa.Rules, error) {
	rules, err := generateDbusRules("own", "session", p.Metadata.Name)
	if err != nil {
		return nil, err
	}

	for _, policy := range []struct {
		bus  string
		list []Dbus
	}{
		{"session", p.Metadata.DbusSession},
		{"system", p.Metadata.DbusSystem},
	} {
		for _, dbus := range policy.list {
			if dbus.Action != "own" && dbus.Action != "talk" {
				continue
			}
			if abs, present := dbusInterfaceMap[dbus.Interface]; present && dbus.Action == "talk" {
				rules = append(rules, &aa.Include{
					Base:    aa.Base{Comment: " " + dbus.Interface},
					IsMagic: true, Path: "abstractions/" + abs,
				})
				continue
			}
			r, err := generateDbusRules(dbus.Action, policy.bus, dbus.Interface)
			if err != nil {
				return nil, err
			}
			rules = append(rules, r...)
		}
	}
	return rules, nil
}

// parseFilesystemSpec parses a Flatpak filesystem specification like "path:ro" or "path:create"
func parseFilesystemSpec(fs string) (string, aa.AccessMask) {
	basePath, mode, _ := strings.Cut(fs, ":")
	if mode == "ro" {
		return basePath, aa.MustAccess(aa.FILE, "r")
	}
	return basePath, aa.MustAccess(aa.FILE, "r", "w", "l", "k") // Default access is read-write
}

// extractPathComponents parses paths that might have XDG prefix with subdirectory
// For example: xdg-documents/reports or xdg-data/applications
func extractPathComponents(fs string) (string, string) {
	if strings.HasPrefix(fs, "/") {
		return fs, ""
	}
	if base, sub, found := strings.Cut(fs, "/"); found {
		return base, sub
	}
	return fs, ""
}

func allowDirectory(owner bool, dirname string, access aa.AccessMask) aa.Rules {
	dirname = strings.TrimRight(dirname, "/")
	if access == aa.MustAccess(aa.FILE, "r") {
		return aa.Rules{
			&aa.File{Owner: owner, Path: dirname, Access: access}, // File
			&aa.File{Owner: owner, Path: dirname + "/{,**}", Access: access},
		}
	} else {
		r := &aa.File{Owner: owner, Path: dirname + "/**", Access: access}
		return aa.Rules{
			&aa.File{Owner: owner, Path: dirname, Access: aa.MustAccess(aa.FILE, "r", "w")}, // File
			&aa.File{Owner: owner, Path: dirname + "/", Access: aa.MustAccess(aa.FILE, "r")},
			r,
		}
	}
}

// TODO: Revisit this, it is much more complex than this.
// See https://docs.flatpak.org/en/latest/sandbox-permissions.html
func (p *FlatpakAppArmorProfile) addFileSystem() aa.Rules {
	rules := aa.Rules{}

	for _, fsSpec := range p.Metadata.Filesystems {
		// Parse access permissions (e.g., ":ro", ":create")
		fs, access := parseFilesystemSpec(fsSpec)

		// Handle paths that might have XDG prefix with subdirectory
		basePath, subPath := extractPathComponents(fs)

		// Handle standard XDG directories (xdg-documents, xdg-pictures, etc.)
		if variable, ok := fsMap[basePath]; ok {
			if subPath != "" {
				variable += "/" + subPath
			}
			rules = append(rules, allowDirectory(true, variable, access)...)
			continue
		}

		// Handle special cases
		switch basePath {
		case "host-etc":
			// Host /etc directory
			rules = append(rules, allowDirectory(true, "/etc", access)...)

		case "host", "host-os":
			// host grants the requested access, host-os is read only
			hostAccess := access
			if basePath == "host-os" {
				hostAccess = aa.MustAccess(aa.FILE, "r")
			}

			// Full host access to important directories
			if basePath == "host" {
				for _, dirname := range []string{"home", "media", "opt", "run/media", "srv"} {
					rules = append(rules, &aa.File{
						Qualifier: aa.Qualifier{Audit: true},
						Path:      "/run/host/" + dirname + "/{,**}",
						Access:    access,
					})
				}
			}

			// Access to system directories and libraries
			for _, dirname := range []string{"/usr", "/bin", "/sbin", "/lib32", "/lib64", "/lib"} {
				rules = append(rules, allowDirectory(true, dirname, hostAccess)...)
			}

			// Special files needed for dynamic linking and alternatives
			for _, filename := range []string{"/etc/ld.so.cache", "/etc/alternatives"} {
				rules = append(rules, &aa.File{
					Owner: true, Path: filename, Access: hostAccess,
				})
			}

		case "home":
			// Home directory access (excluding ~/.var/app)
			rules = append(rules, allowDirectory(true, "@{HOME}/", access)...)
			// TODO: other app dir are excluded
			// &aa.File{Owner: true, Path: "@{HOME}/.var/app/{,**}", Access: []string{"k"}},

		case "xdg-run":
			// XDG runtime directory access with subdirectory support
			path := "@{run}/user/@{uid}/"
			if subPath != "" {
				path += subPath
			}
			rules = append(rules, allowDirectory(true, path, access)...)

		// External drive access locations
		case "/media", "/run/media", "/mnt":
			rules = append(rules, allowDirectory(true, basePath, access)...)

		default:
			// Handle absolute paths
			if strings.HasPrefix(basePath, "/") {
				// Handle absolute paths that aren't special cases
				path := basePath
				if subPath != "" {
					path += "/" + subPath
				}
				rules = append(rules, allowDirectory(false, path, access)...)

			} else if strings.HasPrefix(basePath, "~") {
				path := "@{HOME}/"
				if subPath != "" {
					path += subPath
				}
				rules = append(rules, allowDirectory(true, path, access)...)
			}
		}
	}

	for _, dirname := range p.Metadata.Persistent {
		path := "@{HOME}/" + dirname
		if dirname == "." {
			path = "@{HOME}/"
			rules = append(rules, &aa.Include{
				IsMagic: true, Path: "abstractions/deny-sensitive-home",
			})
		}
		rules = append(rules, allowDirectory(true, path, aa.MustAccess(aa.FILE, "r", "w", "k"))...)
	}
	return rules
}

func (p *FlatpakAppArmorProfile) finalize(profile string) aa.Rules {
	return aa.Rules{
		nil,
		&aa.Include{IfExists: true, IsMagic: true, Path: "local/" + profile},
	}
}
