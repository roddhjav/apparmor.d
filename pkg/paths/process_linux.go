// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

//go:build !windows

package paths

import (
	"os/exec"
	"syscall"
)

func tellCommandToStartOnNewProcessGroup(oscmd *exec.Cmd) {
	// https://groups.google.com/g/golang-nuts/c/XoQ3RhFBJl8

	// Start the process in a new process group.
	// This is needed to kill the process and its children
	// if we need to kill the process.
	if oscmd.SysProcAttr == nil {
		oscmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	oscmd.SysProcAttr.Setpgid = true
}

func kill(oscmd *exec.Cmd) error {
	// https://groups.google.com/g/golang-nuts/c/XoQ3RhFBJl8

	// Kill the process group
	pgid, err := syscall.Getpgid(oscmd.Process.Pid)
	if err != nil {
		return err
	}
	return syscall.Kill(-pgid, syscall.SIGKILL)
}
