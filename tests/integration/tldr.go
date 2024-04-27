// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package integration

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

type Tldr struct {
	Url    string      // Tldr download url
	Dir    *paths.Path // Tldr cache directory
	Ignore []string    // List of ignored software
}

func NewTldr(dir *paths.Path) Tldr {
	return Tldr{
		Url: "https://github.com/tldr-pages/tldr/archive/refs/heads/main.tar.gz",
		Dir: dir,
	}
}

// Download and extract the tldr pages into the cache directory
func (t Tldr) Download() error {
	gzPath := t.Dir.Parent().Join("tldr.tar.gz")
	if !gzPath.Exist() {
		resp, err := http.Get(t.Url)
		if err != nil {
			return fmt.Errorf("downloading %s: %w", t.Url, err)
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

// Parse the tldr pages and return a list of scenarios
func (t Tldr) Parse() (*TestSuite, error) {
	testSuite := NewTestSuite()
	files, _ := t.Dir.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	for _, path := range files {
		content, err := path.ReadFile()
		if err != nil {
			return nil, err
		}
		raw := string(content)
		t := &Test{
			Name:      strings.TrimSuffix(path.Base(), ".md"),
			Root:      false,
			Arguments: map[string]string{},
			Commands:  []Command{},
		}
		if strings.Contains(raw, "sudo") {
			t.Root = true
		}
		rawTests := strings.Split(raw, "\n-")[1:]
		for _, test := range rawTests {
			res := strings.Split(test, "\n")
			dsc := strings.ReplaceAll(strings.Trim(res[0], " "), ":", "")
			cmd := strings.Trim(strings.Trim(res[2], "`"), " ")
			if t.Root {
				cmd = strings.ReplaceAll(cmd, "sudo ", "")
			}
			t.Commands = append(t.Commands, Command{
				Description: dsc,
				Cmd:         cmd,
			})
		}
		testSuite.Tests = append(testSuite.Tests, *t)
	}
	return testSuite, nil
}
