// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package tasks

import (
	"reflect"
	"testing"

	"github.com/roddhjav/apparmor.d/pkg/paths"
)

const (
	osReleaseArchlinux = `NAME="Arch Linux"
PRETTY_NAME="Arch Linux"
ID=arch
BUILD_ID=rolling
ANSI_COLOR="38;2;23;147;209"
HOME_URL="https://archlinux.org/"
DOCUMENTATION_URL="https://wiki.archlinux.org/"
SUPPORT_URL="https://bbs.archlinux.org/"
BUG_REPORT_URL="https://bugs.archlinux.org/"
PRIVACY_POLICY_URL="https://terms.archlinux.org/docs/privacy-policy/"
LOGO=archlinux-logo`

	osReleaseUbuntu = `PRETTY_NAME="Ubuntu 22.04.2 LTS"
NAME="Ubuntu"
VERSION_ID="22.04"
VERSION="22.04.2 LTS (Jammy Jellyfish)"
VERSION_CODENAME=jammy
ID=ubuntu
ID_LIKE=debian
HOME_URL="https://www.ubuntu.com/"
SUPPORT_URL="https://help.ubuntu.com/"
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
UBUNTU_CODENAME=jammy`

	osReleaseDebian = `PRETTY_NAME="Debian GNU/Linux 11 (bullseye)"
NAME="Debian GNU/Linux"
VERSION_ID="11"
VERSION="11 (bullseye)"
VERSION_CODENAME=bullseye
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"`

	osReleaseOpenSUSETumbleweed = `ID="opensuse-tumbleweed"
ID_LIKE="opensuse suse"
VERSION_ID="20230404"
PRETTY_NAME="openSUSE Tumbleweed"
ANSI_COLOR="0;32"
CPE_NAME="cpe:/o:opensuse:tumbleweed:20230404"
BUG_REPORT_URL="https://bugs.opensuse.org"
HOME_URL="https://www.opensuse.org/"
DOCUMENTATION_URL="https://en.opensuse.org/Portal:Tumbleweed"
LOGO="distributor-logo-Tumbleweed"`

	osReleaseFedora = `NAME="Fedora Linux"
VERSION="37 (Workstation Edition)"
ID=fedora
VERSION_ID=37
VERSION_CODENAME=""
PLATFORM_ID="platform:f37"
PRETTY_NAME="Fedora Linux 37 (Workstation Edition)"
ANSI_COLOR="0;38;2;60;110;180"
LOGO=fedora-logo-icon`

	osReleaseNeon = `PRETTY_NAME="KDE neon 6.0"
NAME="KDE neon"
VERSION_ID="22.04"
VERSION="6.0"
VERSION_CODENAME=jammy
ID=neon
ID_LIKE="ubuntu debian"
HOME_URL="https://neon.kde.org/"
SUPPORT_URL="https://neon.kde.org/"
BUG_REPORT_URL="https://bugs.kde.org/"
PRIVACY_POLICY_URL="https://kde.org/privacypolicy/"
UBUNTU_CODENAME=jammy
LOGO=start-here-kde-neon`
)

func Test_getOSRelease(t *testing.T) {
	tests := []struct {
		name      string
		osRelease string
		want      map[string]string
	}{
		{
			name:      "Archlinux",
			osRelease: osReleaseArchlinux,
			want: map[string]string{
				"NAME":               "Arch Linux",
				"PRETTY_NAME":        "Arch Linux",
				"ID":                 "arch",
				"BUILD_ID":           "rolling",
				"ANSI_COLOR":         "38;2;23;147;209",
				"HOME_URL":           "https://archlinux.org/",
				"DOCUMENTATION_URL":  "https://wiki.archlinux.org/",
				"SUPPORT_URL":        "https://bbs.archlinux.org/",
				"BUG_REPORT_URL":     "https://bugs.archlinux.org/",
				"PRIVACY_POLICY_URL": "https://terms.archlinux.org/docs/privacy-policy/",
				"LOGO":               "archlinux-logo",
			},
		},
		{
			name:      "Ubuntu",
			osRelease: osReleaseUbuntu,
			want: map[string]string{
				"PRETTY_NAME":        "Ubuntu 22.04.2 LTS",
				"NAME":               "Ubuntu",
				"VERSION_ID":         "22.04",
				"VERSION":            "22.04.2 LTS (Jammy Jellyfish)",
				"VERSION_CODENAME":   "jammy",
				"ID":                 "ubuntu",
				"ID_LIKE":            "debian",
				"HOME_URL":           "https://www.ubuntu.com/",
				"SUPPORT_URL":        "https://help.ubuntu.com/",
				"BUG_REPORT_URL":     "https://bugs.launchpad.net/ubuntu/",
				"PRIVACY_POLICY_URL": "https://www.ubuntu.com/legal/terms-and-policies/privacy-policy",
				"UBUNTU_CODENAME":    "jammy",
			},
		},
	}
	osReleaseFile = "/tmp/os-release"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := paths.New(osReleaseFile).WriteFile([]byte(tt.osRelease))
			if err != nil {
				return
			}
			if got := getOSRelease(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOSRelease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDistribution(t *testing.T) {
	tests := []struct {
		name      string
		osRelease string
		want      string
	}{
		{
			name:      "Archlinux",
			osRelease: osReleaseArchlinux,
			want:      "arch",
		},
		{
			name:      "Ubuntu",
			osRelease: osReleaseUbuntu,
			want:      "ubuntu",
		},
		{
			name:      "Debian",
			osRelease: osReleaseDebian,
			want:      "debian",
		},
		{
			name:      "OpenSUSE Tumbleweed",
			osRelease: osReleaseOpenSUSETumbleweed,
			want:      "opensuse",
		},
		{
			name:      "Fedora",
			osRelease: osReleaseFedora,
			want:      "fedora",
		},
		{
			name:      "Neon",
			osRelease: osReleaseNeon,
			want:      "ubuntu",
		},
	}

	osReleaseFile = "/tmp/os-release"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := paths.New(osReleaseFile).WriteFile([]byte(tt.osRelease))
			if err != nil {
				return
			}
			Release = getOSRelease()
			got := getDistribution()
			if got != tt.want {
				t.Errorf("getDistribution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getFamily(t *testing.T) {
	tests := []struct {
		name string
		dist string
		want string
	}{
		{
			name: "Archlinux",
			dist: "arch",
			want: "pacman",
		},
		{
			name: "Ubuntu",
			dist: "ubuntu",
			want: "apt",
		},
		{
			name: "Debian",
			dist: "debian",
			want: "apt",
		},
		{
			name: "OpenSUSE Tumbleweed",
			dist: "opensuse",
			want: "zypper",
		},
		{
			name: "Neon",
			dist: "neon",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Distribution = tt.dist
			if got := getFamily(); got != tt.want {
				t.Errorf("getFamily() = %v, want %v", got, tt.want)
			}
		})
	}
}
