// This file is part of PathsHelper library.
// Copyright (C) 2018-2025 Arduino AG (http://www.arduino.cc/)
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package paths

import (
	"os"
)

// NullPath return the path to the /dev/null equivalent for the current OS
func NullPath() *Path {
	return New("/dev/null")
}

// TempDir returns the default path to use for temporary files
func TempDir() *Path {
	return New(os.TempDir()).Canonical()
}

// MkTempDir creates a new temporary directory in the directory
// dir with a name beginning with prefix and returns the path of
// the new directory. If dir is the empty string, TempDir uses the
// default directory for temporary files
func MkTempDir(dir, prefix string) (*Path, error) {
	path, err := os.MkdirTemp(dir, prefix)
	if err != nil {
		return nil, err
	}
	return New(path).Canonical(), nil
}

// MkTempFile creates a new temporary file in the directory dir with a name beginning with prefix,
// opens the file for reading and writing, and returns the resulting *os.File. If dir is nil,
// MkTempFile uses the default directory for temporary files (see paths.TempDir). Multiple programs
// calling TempFile simultaneously will not choose the same file. The caller can use f.Name() to
// find the pathname of the file. It is the caller's responsibility to remove the file when no longer needed.
func MkTempFile(dir *Path, prefix string) (*os.File, error) {
	tmpDir := ""
	if dir != nil {
		tmpDir = dir.String()
	}
	return os.CreateTemp(tmpDir, prefix)
}

// Getwd returns a rooted path name corresponding to the current
// directory.
func Getwd() (*Path, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return New(wd), nil
}
