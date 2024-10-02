// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prepare

import (
	"fmt"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

type Configure struct {
	prebuild.Base
	OneFile bool
}

func init() {
	RegisterTask(&Configure{
		Base: prebuild.Base{
			Keyword: "configure",
			Msg:     "Set distribution specificities",
		},
		OneFile: false,
	})
}

func (p Configure) Apply() ([]string, error) {
	res := []string{}

	if prebuild.ABI == 4 {
		if err := OverwriteUpstreamProfile(p.OneFile); err != nil {
			return res, err
		}
	}
	switch prebuild.Distribution {
	case "arch", "opensuse":

	case "ubuntu":
		if err := prebuild.DebianHide.Init(); err != nil {
			return res, err
		}

		if prebuild.ABI == 3 {
			if err := util.CopyTo(prebuild.DistDir.Join("ubuntu"), prebuild.RootApparmord); err != nil {
				return res, err
			}
		}

	case "debian", "whonix":
		if err := prebuild.DebianHide.Init(); err != nil {
			return res, err
		}

		// Copy Debian specific abstractions
		if err := util.CopyTo(prebuild.DistDir.Join("ubuntu"), prebuild.RootApparmord); err != nil {
			return res, err
		}

	default:
		return []string{}, fmt.Errorf("%s is not a supported distribution", prebuild.Distribution)

	}
	return res, nil
}

// Overwrite upstream profile: disable upstream & rename ours
func OverwriteUpstreamProfile(oneFile bool) error {
	const ext = ".apparmor.d"
	disableDir := prebuild.RootApparmord.Join("disable")
	if err := disableDir.Mkdir(); err != nil {
		return err
	}

	path := prebuild.DistDir.Join("overwrite")
	if !path.Exist() {
		return fmt.Errorf("%s not found", path)
	}
	for _, name := range util.MustReadFileAsLines(path) {
		origin := prebuild.RootApparmord.Join(name)
		dest := prebuild.RootApparmord.Join(name + ext)
		if !dest.Exist() && oneFile {
			continue
		}
		if err := origin.Rename(dest); err != nil {

			return err
		}
		originRel, err := origin.RelFrom(dest)
		if err != nil {
			return err
		}
		if err := os.Symlink(originRel.String(), disableDir.Join(name).String()); err != nil {
			return err
		}
	}
	return nil
}
