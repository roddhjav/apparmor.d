// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

type Tldr struct {
	URL    string      // Tldr download url
	Dir    *paths.Path // Tldr cache directory
	Ignore []string    // List of ignored software
}

func NewTldr(dir *paths.Path) Tldr {
	return Tldr{
		URL: "https://github.com/tldr-pages/tldr/archive/refs/heads/main.tar.gz",
		Dir: dir,
	}
}

// Download and extract the tldr pages into the cache directory
func (t Tldr) Download() error {
	gzPath := t.Dir.Parent().Join("tldr.tar.gz")
	if !gzPath.Exist() {
		resp, err := http.Get(t.URL)
		if err != nil {
			return fmt.Errorf("downloading %s: %w", t.URL, err)
		}
		defer resp.Body.Close()

		out, err := gzPath.Create()
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			return err
		}
	}

	pages := []string{"tldr-main/pages/linux", "tldr-main/pages/common"}
	return extratTo(gzPath, t.Dir, pages)
}

// Parse the tldr pages and return a list of tests
func (t Tldr) Parse() (Tests, error) {
	tests := make(Tests, 0)
	files, _ := t.Dir.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	for _, path := range files {
		content, err := path.ReadFile()
		if err != nil {
			return nil, err
		}
		raw := string(content)
		t := Test{
			Name:     strings.TrimSuffix(path.Base(), ".md"),
			Commands: []Command{},
		}
		rawTests := strings.Split(raw, "\n-")[1:]
		for _, test := range rawTests {
			res := strings.Split(test, "\n")
			dsc := strings.ReplaceAll(strings.Trim(res[0], " "), ":", "")
			cmd := strings.Trim(strings.Trim(res[2], "`"), " ")
			t.Commands = append(t.Commands, Command{
				Description: dsc,
				Cmd:         cmd,
			})
		}
		tests = append(tests, t)
	}
	return tests, nil
}

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
func extratTo(src *paths.Path, dst *paths.Path, subfolders []string) error {
	gzIn, err := src.Open()
	if err != nil {
		return fmt.Errorf("opening %s: %w", src, err)
	}
	defer gzIn.Close()

	in, err := gzip.NewReader(gzIn)
	if err != nil {
		return fmt.Errorf("decoding %s: %w", src, err)
	}
	defer in.Close()

	if err := dst.MkdirAll(); err != nil {
		return fmt.Errorf("creating %s: %w", src, err)
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
				return fmt.Errorf("creating %s: %w", file.Name(), err)
			}
			if _, err := io.Copy(file, tarIn); err != nil {
				return fmt.Errorf("extracting %s: %w", file.Name(), err)
			}
			file.Close()
		}
	}
	return nil
}
