// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"testing"

	"github.com/arduino/go-paths-helper"
)

const (
	Archlinux = `NAME="Arch Linux"
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

	Ubuntu = `PRETTY_NAME="Ubuntu 22.04.2 LTS"
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

	Debian = `PRETTY_NAME="Debian GNU/Linux 11 (bullseye)"
NAME="Debian GNU/Linux"
VERSION_ID="11"
VERSION="11 (bullseye)"
VERSION_CODENAME=bullseye
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"`

	OpenSUSETumbleweed = `ID="opensuse-tumbleweed"
ID_LIKE="opensuse suse"
VERSION_ID="20230404"
PRETTY_NAME="openSUSE Tumbleweed"
ANSI_COLOR="0;32"
CPE_NAME="cpe:/o:opensuse:tumbleweed:20230404"
BUG_REPORT_URL="https://bugs.opensuse.org"
HOME_URL="https://www.opensuse.org/"
DOCUMENTATION_URL="https://en.opensuse.org/Portal:Tumbleweed"
LOGO="distributor-logo-Tumbleweed"`

	ArcoLinux = `NAME=ArcoLinux
ID=arcolinux
ID_LIKE=arch
BUILD_ID=rolling
ANSI_COLOR="0;36"
HOME_URL="https://arcolinux.info/"
SUPPORT_URL="https://arcolinuxforum.com/"
BUG_REPORT_URL="https://github.com/arcolinux"
LOGO=arcolinux-hello`

	Fedora = `NAME="Fedora Linux"
VERSION="37 (Workstation Edition)"
ID=fedora
VERSION_ID=37
VERSION_CODENAME=""
PLATFORM_ID="platform:f37"
PRETTY_NAME="Fedora Linux 37 (Workstation Edition)"
ANSI_COLOR="0;38;2;60;110;180"
LOGO=fedora-logo-icon`
)

func Test_getSupportedDistribution(t *testing.T) {
	tests := []struct {
		name      string
		osRelease string
		want      string
	}{
		{
			name:      "Archlinux",
			osRelease: Archlinux,
			want:      "arch",
		},
		{
			name:      "Ubuntu",
			osRelease: Ubuntu,
			want:      "ubuntu",
		},
		{
			name:      "Debian",
			osRelease: Debian,
			want:      "debian",
		},
		{
			name:      "OpenSUSE Tumbleweed",
			osRelease: OpenSUSETumbleweed,
			want:      "opensuse",
		},
		// {
		// 	name:      "Fedora",
		// 	osRelease: Fedora,
		// 	want:      "fedora",
		// },
	}

	osReleaseFile = "/tmp/os-release"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := paths.New(osReleaseFile).WriteFile([]byte(tt.osRelease))
			if err != nil {
				return
			}
			got := getSupportedDistribution()
			if got != tt.want {
				t.Errorf("getSupportedDistribution() = %v, want %v", got, tt.want)
			}
		})
	}
}
