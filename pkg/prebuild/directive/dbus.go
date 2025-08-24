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
				"common bus=<bus> name=<name> label=<profile>",
			},
		}},
	)
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
	case "common":
		r = d.common(opt.ArgMap)
	}

	aa.IndentationLevel = strings.Count(
		strings.SplitN(opt.Raw, Keyword, 1)[0], aa.Indentation,
	)
	generatedDbus := r.String()
	lenDbus := len(generatedDbus)
	generatedDbus = generatedDbus[:lenDbus-1]
	profile = strings.ReplaceAll(profile, opt.Raw, generatedDbus)
	return profile, nil
}

func (d Dbus) sanityCheck(opt *Option) (string, error) {
	if len(opt.ArgList) < 1 {
		return "", fmt.Errorf("unknown dbus action: %s in %s", opt.Name, opt.File)
	}
	action := opt.ArgList[0]
	if action != "own" && action != "talk" && action != "common" {
		return "", fmt.Errorf("unknown dbus action: %s in %s", opt.Name, opt.File)
	}

	if _, present := opt.ArgMap["name"]; !present {
		return "", fmt.Errorf("missing name for 'dbus: %s' in %s", action, opt.File)
	}
	if _, present := opt.ArgMap["bus"]; !present {
		return "", fmt.Errorf("missing bus for '%s' in %s", opt.ArgMap["name"], opt.File)
	}
	if _, present := opt.ArgMap["label"]; !present && action == "talk" {
		return "", fmt.Errorf("missing label for '%s' in %s", opt.ArgMap["name"], opt.File)
	}

	// Set default values
	if _, present := opt.ArgMap["path"]; !present {
		opt.ArgMap["path"] = "/" + strings.ReplaceAll(opt.ArgMap["name"], ".", "/") + "{,/**}"
	}
	opt.ArgMap["name"] += "{,.*}"
	return action, nil
}

func getInterfaces(rules map[string]string) []string {
	var interfaces []string
	if _, present := rules["interface"]; present {
		interfaces = []string{rules["interface"]}
	} else {
		interfaces = []string{rules["name"]}
	}

	if _, present := rules["interface+"]; present {
		interfaces = append(interfaces, rules["interface+"])
	}
	return interfaces
}

func (d Dbus) own(rules map[string]string) aa.Rules {
	interfaces := getInterfaces(rules)

	res := aa.Rules{
		&aa.Include{
			IsMagic: true, Path: "abstractions/bus/" + rules["bus"] + "/own",
		},
		&aa.Dbus{
			Access: []string{"bind"}, Bus: rules["bus"], Name: rules["name"],
		},
	}

	// Interfaces
	for _, iface := range interfaces {
		res = append(res,
			&aa.Dbus{
				Access: []string{"receive"}, Bus: rules["bus"], Path: rules["path"],
				Interface: iface,
				PeerName:  `"@{busname}"`,
			},
			&aa.Dbus{
				Access: []string{"send"}, Bus: rules["bus"], Path: rules["path"],
				Interface: iface,
				PeerName:  `"{@{busname},org.freedesktop.DBus}"`,
			},
		)
	}

	res = append(res,
		// DBus.Properties
		&aa.Dbus{
			Access: []string{"send", "receive"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.Properties",
			Member:    "{Get,GetAll,Set,PropertiesChanged}",
			PeerName:  `"{@{busname},org.freedesktop.DBus}"`,
		},

		// DBus.Introspectable
		&aa.Dbus{
			Access: []string{"receive"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.Introspectable",
			Member:    "Introspect",
			PeerName:  `"@{busname}"`,
		},

		// DBus.ObjectManager
		&aa.Dbus{
			Access: []string{"receive"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.ObjectManager",
			Member:    "GetManagedObjects",
			PeerName:  `"{@{busname},` + rules["name"] + `}"`,
		},
		&aa.Dbus{
			Access: []string{"send"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.ObjectManager",
			Member:    "{InterfacesAdded,InterfacesRemoved}",
			PeerName:  `"{@{busname},org.freedesktop.DBus}"`,
		},
	)
	return res
}

func (d Dbus) talk(rules map[string]string) aa.Rules {
	interfaces := getInterfaces(rules)
	res := aa.Rules{}

	// Interfaces
	for _, iface := range interfaces {
		res = append(res, &aa.Dbus{
			Access: []string{"send", "receive"}, Bus: rules["bus"], Path: rules["path"],
			Interface: iface,
			PeerName:  `"{@{busname},` + rules["name"] + `}"`, PeerLabel: rules["label"],
		})
	}

	res = append(res,
		// DBus.Properties
		&aa.Dbus{
			Access: []string{"send", "receive"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.Properties",
			Member:    "{Get,GetAll,Set,PropertiesChanged}",
			PeerName:  `"{@{busname},` + rules["name"] + `}"`, PeerLabel: rules["label"],
		},

		// DBus.Introspectable
		&aa.Dbus{
			Access: []string{"send"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.Introspectable",
			Member:    "Introspect",
			PeerName:  `"{@{busname},` + rules["name"] + `}"`, PeerLabel: rules["label"],
		},

		// DBus.ObjectManager
		&aa.Dbus{
			Access: []string{"send"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.ObjectManager",
			Member:    "GetManagedObjects",
			PeerName:  `"{@{busname},` + rules["name"] + `}"`, PeerLabel: rules["label"],
		},
		&aa.Dbus{
			Access: []string{"receive"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.ObjectManager",
			Member:    "{InterfacesAdded,InterfacesRemoved}",
			PeerName:  `"{@{busname},` + rules["name"] + `}"`, PeerLabel: rules["label"],
		},
	)
	return res
}

func (d Dbus) common(rules map[string]string) aa.Rules {
	res := aa.Rules{

		// DBus.Properties: read all properties from the interface
		&aa.Comment{
			Base: aa.Base{
				Comment:    " DBus.Properties: read all properties from the interface",
				IsLineRule: true,
			},
		},
		&aa.Dbus{
			Access: []string{"send"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.Properties",
			Member:    "{Get,GetAll}",
			PeerName:  `"{@{busname},` + rules["name"] + `}"`, PeerLabel: rules["label"],
		},
		nil,

		// DBus.Properties: receive property changed events
		&aa.Comment{
			Base: aa.Base{
				Comment:    " DBus.Properties: receive property changed events",
				IsLineRule: true,
			},
		},
		&aa.Dbus{
			Access: []string{"receive"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.Properties",
			Member:    "PropertiesChanged",
			PeerName:  `"{@{busname},` + rules["name"] + `}"`, PeerLabel: rules["label"],
		},
		nil,

		// DBus.Introspectable: allow clients to introspect the service
		&aa.Comment{
			Base: aa.Base{
				Comment:    " DBus.Introspectable: allow clients to introspect the service",
				IsLineRule: true,
			},
		},
		&aa.Dbus{
			Access: []string{"send"}, Bus: rules["bus"], Path: rules["path"],
			Interface: "org.freedesktop.DBus.Introspectable",
			Member:    "Introspect",
			PeerName:  `"{@{busname},` + rules["name"] + `}"`, PeerLabel: rules["label"],
		},
	}
	return res
}
