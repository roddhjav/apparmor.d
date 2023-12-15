// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/aa"
)

// Build the profiles with the following directive applied
var (
	Directives = []DirectiveFunc{
		DirectiveDbus,
	}
	DirectiveMsg = map[string]string{
		"DirectiveDbus":   "DBus directive applied",
	}
)

var (
	regDbus   = regexp.MustCompile(`(?m)# dbus: (.*)$`)
)

type DirectiveFunc func(*paths.Path, string) string

// Apply dbus directive
//
// Example of supported dbus directive:
// # dbus: own bus=session name=org.freedesktop.FileManager1
// # dbus: talk name=org.freedesktop.login1 label=systemd-logind
func DirectiveDbus(file *paths.Path, profile string) string {
	const lenHeader = 12 // len("profile {\n  ")

	var p *aa.AppArmorProfile
	for _, match := range regDbus.FindAllStringSubmatch(profile, -1) {
		origin := match[0]
		rules := map[string]string{}

		fields := strings.Fields(match[1])
		action, fields := fields[0], fields[1:]
		for _, t := range fields {
			tmp := strings.Split(t, "=")
			rules[tmp[0]] = tmp[1]
		}

		rules = sanitizeDbusRule(file, action, rules)
		switch action {
		case "own":
			p = dbusOwn(rules)
		case "talk":
			p = dbusTalk(rules)
		default:
			panic(fmt.Sprintf("Unknown dbus directive: %s in %s", action, file))
		}

		generatedDbus := p.String()
		lenDbus := len(generatedDbus)
		generatedDbus = generatedDbus[lenHeader : lenDbus-3]
		profile = strings.Replace(profile, origin, generatedDbus, -1)
	}
	return profile
}

func sanitizeDbusRule(file *paths.Path, action string, rules map[string]string) map[string]string {
	// Sanity check
	if _, present := rules["name"]; !present {
		panic(fmt.Sprintf("Missing name for 'dbus: own' in %s", file))
	}
	if _, present := rules["bus"]; !present {
		panic(fmt.Sprintf("Missing bus for '%s' in %s", rules["name"], file))
	}

	if _, present := rules["label"]; !present && action == "talk" {
		panic(fmt.Sprintf("Missing label for '%s' in %s", rules["name"], file))
	}

	// Set default values
	if _, present := rules["path"]; !present {
		rules["path"] = "/" + strings.Replace(rules["name"], ".", "/", -1) + "{,/**}"
	}
	if _, present := rules["interface"]; !present {
		rules["interface"] = "org.freedesktop.DBus.{Properties,ObjectManager}"
	}
	rules["name"] += "{,.*}"
	return rules
}

func dbusOwn(rules map[string]string) *aa.AppArmorProfile {
	interfaces := []string{rules["name"], rules["interface"]}
	p := &aa.AppArmorProfile{}
	p.Rules = append(p.Rules, &aa.Dbus{
		Access: "bind", Bus: rules["bus"], Name: rules["name"],
	})
	for _, iface := range interfaces {
		p.Rules = append(p.Rules, &aa.Dbus{
			Access:    "receive",
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			Name:      `":1.@{int}"`,
		})
	}
	for _, iface := range interfaces {
		p.Rules = append(p.Rules, &aa.Dbus{
			Access:    "send",
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			Name:      `"{:1.@{int},org.freedesktop.DBus}"`,
		})
	}
	return p
}

func dbusTalk(rules map[string]string) *aa.AppArmorProfile {
	interfaces := []string{rules["name"], rules["interface"]}
	p := &aa.AppArmorProfile{}
	for _, iface := range interfaces {
		p.Rules = append(p.Rules, &aa.Dbus{
			Access:    "send",
			Bus:       rules["bus"],
			Name:      `"{:1.@{int},` + rules["name"] + `}"`,
			Path:      rules["path"],
			Interface: iface,
			Label:     rules["label"],
		})
	}
	for _, iface := range interfaces {
		p.Rules = append(p.Rules, &aa.Dbus{
			Access:    "receive",
			Bus:       rules["bus"],
			Name:      `"{:1.@{int},` + rules["name"] + `}"`,
			Path:      rules["path"],
			Interface: iface,
			Label:     rules["label"],
		})
	}
	return p
}

