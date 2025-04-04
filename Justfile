# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Integration environment for apparmor.d
#
# Usage:
#   just
#   just img ubuntu24 server
#   just vm ubuntu24 server
#   just up ubuntu24 server
#   just ssh ubuntu24 server
#   just halt ubuntu24 server
#   just destroy ubuntu24 server
#   just list
#   just images
#   just available
#   just clean

# Admin username
username := "user"

# Default admin password
password := "user"

# Disk size of the VM to build
disk_size := "40G"

# Virtual machine CPU
vcpus := "6"

# Virtual machine RAM
ram := "4096"

# Path to the ssh key
ssh_keyname := "id_ed25519"
ssh_privatekey := home_dir() / ".ssh/" + ssh_keyname
ssh_publickey := ssh_privatekey + ".pub" 

# Where the VM are stored
vm := home_dir() / ".vm"

# Where the VM images are stored
base_dir := home_dir() / ".libvirt/base"

# Where the packer temporary output is stored
output_dir := base_dir / "packer"

# SSH options
sshopt := "-i " + ssh_privatekey + " -o IdentitiesOnly=yes -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

# Libvirt connection address
c := "--connect=qemu:///system"

# VM prefix
prefix := "aa-"

# Build setings
destdir := "/"
build := ".build"
pkgdest := `pwd` / ".pkg/dist"
pkgname := "apparmor.d"

[doc('Show this help message')]
default:
	@echo -e "Integration environment helper for apparmor.d\n"
	@just --list --unsorted
	@echo -e "\nSee https://apparmor.pujol.io/development/vm/ for more information."

[doc('Build the go programs')]
build:
	@go build -o {{build}}/ ./cmd/aa-log
	@go build -o {{build}}/ ./cmd/prebuild

[doc('Prebuild the profiles in enforced mode')]
enforce: build
	@./{{build}}/prebuild

[doc('Prebuild the profiles in complain mode')]
complain: build
	@./{{build}}/prebuild --complain

[doc('Prebuild the profiles in FSP mode')]
fsp: build
	@./{{build}}/prebuild --complain --full

[doc('Install the profiles')]
install:
	#!/usr/bin/env bash
	set -eu -o pipefail
	install -Dm0755 {{build}}/aa-log {{destdir}}/usr/bin/aa-log
	install -Dm0644 systemd/aa-fix.service {{destdir}}/usr/lib/systemd/system/aa-fix.service
	for file in $(find "{{build}}/share" -type f -not -name "*.md" -printf "%P\n"); do
		install -Dm0644 "{{build}}/share/$file" "{{destdir}}/usr/share/$file"
	done
	for file in $(find "{{build}}/apparmor.d" -type f -printf "%P\n"); do
		install -Dm0644 "{{build}}/apparmor.d/$file" "{{destdir}}/etc/apparmor.d/$file"
	done
	for file in $(find "{{build}}/apparmor.d" -type l -printf "%P\n"); do
		mkdir -p "{{destdir}}/etc/apparmor.d/disable"
		cp -d "{{build}}/apparmor.d/$file" "{{destdir}}/etc/apparmor.d/$file"
	done
	for file in "{{build}}/systemd/system/"*; do
		service="$(basename "$file")"
		install -Dm0644 "$file" "{{destdir}}/usr/lib/systemd/system/$service.d/apparmor.conf"
	done
	for file in "{{build}}/systemd/user/"*; do
		service="$(basename "$file")"
		install -Dm0644 "$file" "{{destdir}}/usr/lib/systemd/user/$service.d/apparmor.conf"
	done

[doc('Build & install apparmor.d on Arch based systems')]
pkg:
	@makepkg --syncdeps --install --cleanbuild --force --noconfirm

[doc('Build & install apparmor.d on Debian based systems')]
dpkg:
	@bash dists/build.sh dpkg
	@sudo dpkg -i {{pkgdest}}/{{pkgname}}_*.deb

[doc('Build & install apparmor.d on OpenSUSE based systems')]
rpm:
	@bash dists/build.sh rpm
	@sudo rpm -ivh --force  {{pkgdest}}/{{pkgname}}-*.rpm

[doc('Run the unit tests')]
tests:
	@go test ./cmd/... -v -cover -coverprofile=coverage.out
	@go test ./pkg/... -v -cover -coverprofile=coverage.out
	@go tool cover -func=coverage.out

[doc('Run the linters')]
lint:
	golangci-lint run
	packer fmt tests/packer/
	packer validate --syntax-only tests/packer/
	shellcheck --shell=bash \
		PKGBUILD dists/build.sh dists/docker.sh tests/check.sh \
		tests/packer/init.sh tests/packer/src/aa-update tests/packer/clean.sh \
		debian/{{pkgname}}.postinst debian/{{pkgname}}.postrm

[doc('Run style checks on the profiles')]
check:
	@bash tests/check.sh

[doc('Generate the man pages')]
man:
	@pandoc -t man -s -o share/man/man8/aa-log.8 share/man/man8/aa-log.md

[doc('Build the documentation')]
docs:
	@ENABLED_GIT_REVISION_DATE=false MKDOCS_OFFLINE=true mkdocs build --strict

[doc('Serve the documentation')]
serve:
	@ENABLED_GIT_REVISION_DATE=false MKDOCS_OFFLINE=false mkdocs serve

[doc('Remove all build artifacts')]
clean:
	@rm -rf \
		debian/.debhelper debian/debhelper* debian/*.debhelper debian/{{pkgname}} \
		.pkg/{{pkgname}}* {{build}} coverage.out

[doc('Build the apparmor.d package')]
package dist:
	#!/usr/bin/env bash
	set -eu -o pipefail
	dist="{{dist}}"
	version=""
	if [[ $dist =~ ubuntu([0-9]+) ]]; then
		version="${BASH_REMATCH[1]}.04"
		dist="ubuntu"
	elif [[ $dist =~ debian([0-9]+) ]]; then
		version="${BASH_REMATCH[1]}"
		dist="debian"
	fi
	bash dists/docker.sh $dist $version

[doc('Build the image')]
img dist flavor: (package dist)
	@mkdir -p {{base_dir}}
	packer build -force \
		-var dist={{dist}} \
		-var flavor={{flavor}} \
		-var prefix={{prefix}} \
		-var username={{username}} \
		-var password={{password}} \
		-var ssh_publickey={{ssh_publickey}} \
		-var disk_size={{disk_size}} \
		-var cpus={{vcpus}} \
		-var ram={{ram}} \
		-var base_dir={{base_dir}} \
		-var output_dir={{output_dir}} \
		tests/packer/

[doc('Create the machine')]
create dist flavor:
	@cp -f {{base_dir}}/{{prefix}}{{dist}}-{{flavor}}.qcow2 {{vm}}/{{prefix}}{{dist}}-{{flavor}}.qcow2
	@virt-install {{c}} \
		--import \
		--name {{prefix}}{{dist}}-{{flavor}} \
		--vcpus {{vcpus}} \
		--ram {{ram}} \
		--machine q35 \
		--boot uefi \
		--memorybacking source.type=memfd,access.mode=shared \
		--disk path={{vm}}/{{prefix}}{{dist}}-{{flavor}}.qcow2,format=qcow2,bus=virtio \
		--filesystem "`pwd`,0a31bc478ef8e2461a4b1cc10a24cc4",accessmode=passthrough,driver.type=virtiofs \
		--os-variant "`just get_osinfo {{dist}}`" \
		--graphics spice \
		--audio id=1,type=spice \
		--sound model=ich9 \
		--noautoconsole

[doc('Start a machine')]
up dist flavor:
	@virsh {{c}} start {{prefix}}{{dist}}-{{flavor}}

[doc('Stops the machine')]
halt dist flavor:
	@virsh {{c}} shutdown {{prefix}}{{dist}}-{{flavor}}

[doc('Reboot the machine')]
reboot dist flavor:
	@virsh {{c}} reboot {{prefix}}{{dist}}-{{flavor}}

[doc('Destroy the machine')]
destroy dist flavor:
	@virsh {{c}} destroy {{prefix}}{{dist}}-{{flavor}} || true
	@virsh {{c}} undefine {{prefix}}{{dist}}-{{flavor}} --nvram
	@rm -fv {{vm}}/{{prefix}}{{dist}}-{{flavor}}.qcow2

[doc('Connect to the machine')]
ssh dist flavor:
	@ssh {{sshopt}} {{username}}@`just get_ip {{dist}} {{flavor}}`

[doc('List the machines')]
list:
	@echo -e '\033[1m Id   Distribution Flavor  State\033[0m'
	@virsh {{c}} list --all | grep {{prefix}} | sed 's/{{prefix}}//g'

[doc('List the images')]
images:
	#!/usr/bin/env bash
	set -eu -o pipefail
	ls -lh {{base_dir}} | awk '
	BEGIN {
		printf("\033[1m%-18s %-10s %-5s %s\033[0m\n", "Distribution", "Flavor", "Size", "Date")
	}
	{
		if ($9 ~ /^{{prefix}}.*\.qcow2$/) {
			split($9, arr, "-|\\.")
			printf("%-18s %-10s %-5s %s %s %s\n", arr[2], arr[3], $5, $6, $7, $8)
		}
	}
	'

[doc('List the machine that can be created')]
available:
	#!/usr/bin/env bash
	set -eu -o pipefail
	ls -lh tests/cloud-init | awk '
	BEGIN {
		printf("\033[1m%-18s %s\033[0m\n", "Distribution", "Flavor")
	}
	{
		if ($9 ~ /^.*\.user-data.yml$/) {
			split($9, arr, "-|\\.")
			printf("%-18s %s\n", arr[1], arr[2])
		}
	}
	'

[doc('Run the integration tests on the machine')]
integration dist flavor:
	@ssh {{sshopt}} user@`just get_ip {{dist}} {{flavor}}` \
		cp -rf /home/user/Projects/apparmor.d/tests/integration/ /home/user/Projects
	@ssh {{sshopt}} user@`just get_ip {{dist}} {{flavor}}` \
		sudo umount /home/user/Projects/apparmor.d
	@ssh {{sshopt}} user@`just get_ip {{dist}} {{flavor}}` \
		@bats --recursive --timing --print-output-on-failure Projects/integration/



get_ip dist flavor:
	@virsh --quiet --readonly {{c}} domifaddr {{prefix}}{{dist}}-{{flavor}} | \
		head -1 | \
		grep -E -o '([[:digit:]]{1,3}\.){3}[[:digit:]]{1,3}'

get_osinfo dist:
    #!/usr/bin/env python3
    osinfo = {
        "archlinux": "archlinux",
        "debian12": "debian12",
        "debian13": "debian13",
        "ubuntu22": "ubuntu22.04",
        "ubuntu24": "ubuntu24.04",
        "ubuntu25": "ubuntu25.04", 
        "opensuse": "opensusetumbleweed",
    }
    print(osinfo.get("{{dist}}", "{{dist}}"))
