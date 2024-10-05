//
// This file is part of PathsHelper library.
//
// Copyright 2023 Arduino AG (http://www.arduino.cc/)
//
// PathsHelper library is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation; either version 2 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
//
// As a special exception, you may use this file as part of a free software
// library without restriction.  Specifically, if other files instantiate
// templates or use macros or inline functions from this file, or you compile
// this file and link it with other files to produce an executable, this
// file does not by itself cause the resulting executable to be covered by
// the GNU General Public License.  This exception does not however
// invalidate any other reasons why the executable file might be covered by
// the GNU General Public License.
//

//go:build !windows

package paths

import (
	"os/exec"
	"syscall"
)

func tellCommandNotToSpawnShell(_ *exec.Cmd) {
	// no op
}

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
