// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package configure

import (
	"fmt"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

var (
	serverIgnorePatterns = []string{
		"include <abstractions/app/chromium>",
		"include <abstractions/app/firefox>",
		"include <abstractions/app/open>",
		"include <abstractions/common/desktop>",
		"include <abstractions/common/electron>",
		"include <abstractions/common/gnome>",
		"include <abstractions/cosmic>",
		"include <abstractions/desktop>",
		"include <abstractions/desktop>",
		"include <abstractions/freedesktop.org>",
		"include <abstractions/gnome-strict>",
		"include <abstractions/kde-strict>",
		"include <abstractions/lxqt>",
		"include <abstractions/xfce>",
	}
	serverIgnoreGroups = []string{
		"akonadi",
		"avahi",
		"bluetooth",
		"browsers",
		"cosmic",
		"cups",
		"display-manager",
		"flatpak",
		"freedesktop",
		"gnome",
		"gvfs",
		"hyprland",
		"kde",
		"lxqt",
		"steam",
		"xfce",
		"zed",
	}
)

type Server struct {
	tasks.Base
}

func init() {
	RegisterTask(&Server{
		Base: tasks.Base{
			Keyword: "server",
			Msg:     "Configure AppArmor for server",
		},
	})
}

func (p Server) Apply() ([]string, error) {
	res := []string{}

	// Ignore desktop related groups
	groupNb := 0
	for _, group := range serverIgnoreGroups {
		path := prebuild.RootApparmord.Join("groups", group)
		if path.IsDir() {
			if err := path.RemoveAll(); err != nil {
				return res, err
			}
			groupNb++
		} else {
			res = append(res, fmt.Sprintf("Group %s not found, ignoring", path))
		}
	}

	// Ignore profiles using a desktop related abstraction
	fileNb := 0
	files, _ := prebuild.RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	for _, file := range files {
		if !file.Exist() {
			continue
		}
		profile, err := file.ReadFileAsString()
		if err != nil {
			return res, err
		}
		for _, pattern := range serverIgnorePatterns {
			if strings.Contains(profile, pattern) {
				if err := file.RemoveAll(); err != nil {
					return res, err
				}
				fileNb++
				break
			}
		}
	}

	res = append(res, fmt.Sprintf("%d groups ignored", groupNb))
	res = append(res, fmt.Sprintf("%d profiles ignored", fileNb))
	return res, nil
}
