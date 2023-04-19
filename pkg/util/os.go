// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package util

import (
	"fmt"
	"os"
	"strings"

	"github.com/arduino/go-paths-helper"
)

var osReleaseFile = "/etc/os-release"

var firstPartyDists = []string{"arch", "debian", "ubuntu", "opensuse", "whonix"}

func GetSupportedDistribution() (string, error) {
	dist, present := os.LookupEnv("DISTRIBUTION")
	if present {
		return dist, nil
	}

	lines, err := paths.New(osReleaseFile).ReadFileAsLines()
	if err != nil {
		return "", err
	}

	id := ""
	id_like := ""
	for _, line := range lines {
		item := strings.Split(line, "=")
		if item[0] == "ID" {
			id = strings.Split(strings.Trim(item[1], "\""), " ")[0]
		} else if item[0] == "ID_LIKE" {
			id_like = strings.Split(strings.Trim(item[1], "\""), " ")[0]
		}
	}

	if InSlice(id, firstPartyDists) {
		return id, nil
	} else if InSlice(id_like, firstPartyDists) {
		return id_like, nil
	}
	return id, fmt.Errorf("%s is not a supported distribution", id)
}
