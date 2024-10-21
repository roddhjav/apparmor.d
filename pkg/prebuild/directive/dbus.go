// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

// Dbus directive
//
// Example of supported directive:
// #aa:dbus own bus=session name=org.freedesktop.FileManager1
// #aa:dbus talk name=org.freedesktop.login1 label=systemd-logind
//
// See https://apparmor.pujol.io/development/dbus/#dbus-directive
//

package directive

import (
	"fmt"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/aa"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
)

var defaultInterfaces = []string{
	"org.freedesktop.DBus.Properties",
	"org.freedesktop.DBus.ObjectManager",
}

type Dbus struct {
	prebuild.Base
}

func init() {
	RegisterDirective(&Dbus{
		Base: prebuild.Base{
			Keyword: "dbus",
			Msg:     "Dbus directive applied",
			Help: []string{
				"own bus=<bus> name=<name> [interface=AARE] [path=AARE]",
				"talk bus=<bus> name=<name> label=<profile> [interface=AARE] [path=AARE]",
			},
		}},
	)
}

func setInterfaces(rules map[string]string) []string {
	interfaces := []string{rules["name"]}
	if _, present := rules["interface"]; present {
		interfaces = append(interfaces, rules["interface"])
	}
	interfaces = append(interfaces, defaultInterfaces...)
	return interfaces
}

func (d Dbus) Apply(opt *Option, profile string) (string, error) {
	var r aa.Rules

	action, err := d.sanityCheck(opt)
	if err != nil {
		return "", err
	}
	switch action {
	case "own":
		r = d.own(opt.ArgMap)
	case "talk":
		r = d.talk(opt.ArgMap)
	}

	aa.IndentationLevel = strings.Count(
		strings.SplitN(opt.Raw, Keyword, 1)[0], aa.Indentation,
	)
	generatedDbus := r.String()
	lenDbus := len(generatedDbus)
	generatedDbus = generatedDbus[:lenDbus-1]
	profile = strings.Replace(profile, opt.Raw, generatedDbus, -1)
	return profile, nil
}

func (d Dbus) sanityCheck(opt *Option) (string, error) {
	if len(opt.ArgList) < 1 {
		return "", fmt.Errorf("Unknown dbus action: %s in %s", opt.Name, opt.File)
	}
	action := opt.ArgList[0]
	if action != "own" && action != "talk" {
		return "", fmt.Errorf("Unknown dbus action: %s in %s", opt.Name, opt.File)
	}

	if _, present := opt.ArgMap["name"]; !present {
		return "", fmt.Errorf("Missing name for 'dbus: %s' in %s", action, opt.File)
	}
	if _, present := opt.ArgMap["bus"]; !present {
		return "", fmt.Errorf("Missing bus for '%s' in %s", opt.ArgMap["name"], opt.File)
	}
	if _, present := opt.ArgMap["label"]; !present && action == "talk" {
		return "", fmt.Errorf("Missing label for '%s' in %s", opt.ArgMap["name"], opt.File)
	}

	// Set default values
	if _, present := opt.ArgMap["path"]; !present {
		opt.ArgMap["path"] = "/" + strings.Replace(opt.ArgMap["name"], ".", "/", -1) + "{,/**}"
	}
	opt.ArgMap["name"] += "{,.*}"
	return action, nil
}

func (d Dbus) own(rules map[string]string) aa.Rules {
	interfaces := setInterfaces(rules)
	res := aa.Rules{}
	res = append(res, &aa.Dbus{
		Access: []string{"bind"}, Bus: rules["bus"], Name: rules["name"],
	})
	for _, iface := range interfaces {
		res = append(res, &aa.Dbus{
			Access:    []string{"receive"},
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			PeerName:  `":1.@{int}"`,
		})
	}
	for _, iface := range interfaces {
		res = append(res, &aa.Dbus{
			Access:    []string{"send"},
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			PeerName:  `"{:1.@{int},org.freedesktop.DBus}"`,
		})
	}
	res = append(res, &aa.Dbus{
		Access:    []string{"receive"},
		Bus:       rules["bus"],
		Path:      rules["path"],
		Interface: "org.freedesktop.DBus.Introspectable",
		Member:    "Introspect",
		PeerName:  `":1.@{int}"`,
	})
	return res
}

func (d Dbus) talk(rules map[string]string) aa.Rules {
	interfaces := setInterfaces(rules)
	res := aa.Rules{}
	for _, iface := range interfaces {
		res = append(res, &aa.Dbus{
			Access:    []string{"send"},
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			PeerName:  `"{:1.@{int},` + rules["name"] + `}"`,
			PeerLabel: rules["label"],
		})
	}
	for _, iface := range interfaces {
		res = append(res, &aa.Dbus{
			Access:    []string{"receive"},
			Bus:       rules["bus"],
			Path:      rules["path"],
			Interface: iface,
			PeerName:  `"{:1.@{int},` + rules["name"] + `}"`,
			PeerLabel: rules["label"],
		})
	}
	return res
}
