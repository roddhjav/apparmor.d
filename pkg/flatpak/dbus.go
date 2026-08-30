// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package flatpak

import (
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/directive"
	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

var (
	reProfileName = regexp.MustCompile(`^profile\s+(\S+)`)
	reDbusBind    = regexp.MustCompile(`^\s*dbus bind\s+.*name=(\S+)`)

	// Matches one or more trailing @{var} or {alt} at the very end of a name.
	// E.g., "Sources@{int}", "kwalletd{,5,6}", "secret{,s}" but NOT
	// "{S,s}ecret{,s}" where a literal follows the first alternation.
	reTrailingSuffix = regexp.MustCompile(`(@\{[^}]+\}|\{[^}]+\})+$`)

	// dbusLabelMap maps exact dbus names to profile labels.
	dbusLabelMap map[string]string

	// dbusLabelPrefixMap maps dbus name prefixes (from names with trailing
	// AppArmor variables/alternations) to profile labels.
	dbusLabelPrefixMap map[string]string

	// dbusLabelOnce guards the lazy scan of the profile tree filling the two
	// maps above. Scanning on first lookup and not on package load keeps the
	// cost out of the callers that never resolve a label.
	dbusLabelOnce sync.Once

	// dbusLabelDefault is a manually maintained fallback map for dbus names
	// that no profile declares ownership of via #aa:dbus own.
	dbusLabelDefault = map[string]string{
		"com.canonical.AppMenu.Registrar":          `"@{pp_dbusmenu}"`,
		"org.gnome.evolution.dataserver.Calendar8": "evolution-calendar-factory",
		"org.gnome.evolution.dataserver.Sources5":  "evolution-source-registry",
	}

	dbusInterfaceMap = map[string]string{
		"ca.desrt.dconf":                              "bus/session/ca.desrt.dconf.Writer",
		"org.a11y.Bus":                                "dbus-accessibility",
		"org.freedesktop.Avahi":                       "avahi-observe",
		"org.freedesktop.FileManager1":                "bus/session/org.freedesktop.FileManager1",
		"org.freedesktop.impl.portal.PermissionStore": "bus/session/org.freedesktop.impl.portal.PermissionStore",
		"org.freedesktop.NetworkManager":              "network-manager-observe",
		"org.freedesktop.Notifications":               "notifications",
		"org.freedesktop.ScreenSaver":                 "screensaver",
		"org.freedesktop.secrets":                     "secrets-service",
		"org.freedesktop.UPower":                      "upower-observe",
		"org.gnome.SessionManager":                    "session-manager",
		"org.gnome.SettingsDaemon.MediaKeys":          "mediakeys",
		"org.gnome.SettingsDaemon":                    "mediakeys",
		"org.gtk.vfs.*":                               "gvfs",
		"org.gtk.vfs":                                 "gvfs",
	}
)

// lookupDbusLabel looks up a dbus name in the label maps. It tries in order:
//  1. Exact match from scanned profiles
//  2. Longest prefix match (from names with trailing patterns like @{int})
//  3. Fallback default map (manually maintained)
//  4. "unconfined" with a warning
func lookupDbusLabel(name string) string {
	dbusLabelOnce.Do(func() {
		dbusLabelMap, dbusLabelPrefixMap = buildDbusLabelMap()
	})
	if label, ok := dbusLabelMap[name]; ok {
		return label
	}
	bestPrefix := ""
	for prefix := range dbusLabelPrefixMap {
		if strings.HasPrefix(name, prefix) && len(prefix) > len(bestPrefix) {
			bestPrefix = prefix
		}
	}
	if bestPrefix != "" {
		return dbusLabelPrefixMap[bestPrefix]
	}
	if label, ok := dbusLabelDefault[name]; ok {
		return label
	}
	log.Printf("WARNING: no profile owns dbus name %s, using unconfined", name)
	return "unconfined"
}

// buildDbusLabelMap scans all profile files for `dbus bind` rules and
// builds maps from dbus names/prefixes to profile labels. Names with trailing
// AppArmor patterns (@{int}, {,5,6}, etc.) are stripped to prefixes. When
// multiple profiles own the same name, the label uses alternation.
func buildDbusLabelMap() (map[string]string, map[string]string) {
	exactOwners := map[string][]string{}
	prefixOwners := map[string][]string{}

	files, err := aa.MagicRoot.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	if err != nil {
		return ownersToLabelMap(exactOwners), ownersToLabelMap(prefixOwners)
	}
	for _, file := range files {
		scanProfileForDbusOwn(file, exactOwners, prefixOwners)
	}

	return ownersToLabelMap(exactOwners), ownersToLabelMap(prefixOwners)
}

func ownersToLabelMap(owners map[string][]string) map[string]string {
	labelMap := make(map[string]string, len(owners))
	for name, profiles := range owners {
		profiles = util.RemoveDuplicate(profiles)
		slices.Sort(profiles)
		if len(profiles) == 1 {
			labelMap[name] = profiles[0]
		} else {
			labelMap[name] = `"{` + strings.Join(profiles, ",") + `}"`
		}
	}
	return labelMap
}

// scanProfileForDbusOwn parses a single profile file, extracting the profile
// name and any #aa:dbus own directives, and populates the owners maps. The
// regexes only run on the few lines that can possibly match.
func scanProfileForDbusOwn(file *paths.Path, exactOwners, prefixOwners map[string][]string) {
	lines, err := file.ReadFileAsLines()
	if err != nil {
		return
	}

	var profileName string
	for _, line := range lines {
		if strings.HasPrefix(line, "profile") {
			if m := reProfileName.FindStringSubmatch(line); m != nil {
				profileName = m[1]
			}
		}
		if profileName != "" && strings.Contains(line, "dbus bind") {
			if m := reDbusBind.FindStringSubmatch(line); m != nil {
				addDbusBind(m[1], profileName, exactOwners, prefixOwners)
			}
		}
	}
}

func addDbusBind(dbusName, profileName string, exactOwners, prefixOwners map[string][]string) {
	// Strip the trailing apparmor rule terminator captured by the regex.
	dbusName = strings.TrimRight(dbusName, ",")

	// No pattern characters: exact match
	if !strings.ContainsAny(dbusName, "@{*") {
		exactOwners[dbusName] = append(exactOwners[dbusName], profileName)
		return
	}

	// Skip names with @{var} variables (not expandable without variable defs)
	if strings.Contains(dbusName, "@{") {
		// But still try to extract a prefix from trailing @{var}
		prefix := reTrailingSuffix.ReplaceAllString(dbusName, "")
		if prefix != "" && prefix != dbusName && !strings.ContainsAny(prefix, "@{*") {
			prefixOwners[prefix] = append(prefixOwners[prefix], profileName)
		}
		return
	}

	// Expand {alt,ernation} patterns. Alts may contain glob wildcards (e.g.
	// the common `{,.*}` suffix), in which case the result is a prefix match.
	for _, name := range expandBraces(dbusName) {
		if strings.Contains(name, "*") {
			prefix := strings.TrimRight(name, "*")
			if prefix != "" && !strings.ContainsAny(prefix, "*@{") {
				prefixOwners[prefix] = append(prefixOwners[prefix], profileName)
			}
			continue
		}
		exactOwners[name] = append(exactOwners[name], profileName)
	}
}

// expandBraces expands simple brace alternations in a string.
// E.g. "org.a11y.{B,b}us" -> ["org.a11y.Bus", "org.a11y.bus"]
// Handles multiple brace groups and nested-free patterns.
func expandBraces(s string) []string {
	idx := strings.Index(s, "{")
	if idx == -1 {
		return []string{s}
	}
	end := strings.Index(s[idx:], "}")
	if end == -1 {
		return []string{s}
	}
	end += idx

	prefix := s[:idx]
	suffix := s[end+1:]
	alts := strings.Split(s[idx+1:end], ",")

	var results []string
	for _, alt := range alts {
		results = append(results, expandBraces(prefix+alt+suffix)...)
	}
	return results
}

func generateDbusRules(action, bus, name string) (aa.Rules, error) {
	args := action + " bus=" + bus + " name=" + name
	if action != "own" {
		args += " label=" + lookupDbusLabel(name)
	}
	opt := directive.NewOption(nil, []string{"", "dbus", args})

	drtv := directive.NewDbus()
	if _, err := drtv.SanityCheck(opt); err != nil {
		return nil, fmt.Errorf("error in dbus directive for interface %s: %v", name, err)
	}

	switch action {
	case "own":
		return drtv.Own(opt.ArgMap), nil
	case "talk":
		return drtv.Talk(opt.ArgMap), nil
	case "see":
		return drtv.See(opt.ArgMap), nil
	default:
		return nil, fmt.Errorf("unknown dbus action: %s", action)
	}
}
