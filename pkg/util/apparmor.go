// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func cmd(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s failed: %w", command, err)
	}
	return nil
}

func ReloadProfiles(files paths.PathList) error {
	args := []string{"--replace"}
	for _, file := range files {
		args = append(args, file.String())
	}
	return cmd("apparmor_parser", args...)
}

func ReloadAppArmor() error {
	_ = cmd("apparmor_parser", "--purge-cache")

	isActive := cmd("systemctl", "is-active", "--quiet", "apparmor.service") == nil
	var err error
	if isActive {
		err = cmd("systemctl", "reload", "apparmor.service")
	} else {
		err = cmd("systemctl", "start", "apparmor.service")
	}

	if err != nil {
		err2 := cmd("journalctl", "--no-pager", "--since=-5m", "--unit", "apparmor.service")
		if err2 != nil {
			return err2
		}
		return fmt.Errorf("failed to reload apparmor service: %w", err)
	}
	return nil
}
