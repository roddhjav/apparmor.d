// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"path/filepath"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/pkg/errors"
)

// Either or not to extract the file
func toExtrat(name string, subfolders []string) bool {
	for _, subfolder := range subfolders {
		if strings.HasPrefix(name, subfolder) {
			return true
		}
	}
	return false
}

// Extract part of an archive to a destination directory
func ExtratTo(src *paths.Path, dst *paths.Path, subfolders []string) error {
	gzIn, err := src.Open()
	if err != nil {
		return errors.Wrapf(err, "opening %s", src)
	}
	defer gzIn.Close()

	in, err := gzip.NewReader(gzIn)
	if err != nil {
		return errors.Wrapf(err, "decoding %s", src)
	}
	defer in.Close()

	if err := dst.MkdirAll(); err != nil {
		return errors.Wrapf(err, "creating %s", src)
	}

	tarIn := tar.NewReader(in)
	for {
		header, err := tarIn.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg {
			if !toExtrat(header.Name, subfolders) {
				continue
			}
			path := dst.Join(filepath.Base(header.Name))
			file, err := path.Create()
			if err != nil {
				return errors.Wrapf(err, "creating %s", file.Name())
			}
			if _, err := io.Copy(file, tarIn); err != nil {
				return errors.Wrapf(err, "extracting %s", file.Name())
			}
			file.Close()
		}
	}
	return nil
}
