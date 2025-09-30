# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: `just`
# See https://apparmor.pujol.io/development/ for more information.

# Build settings

destdir := "/"
build := ".build"
pkgdest := `pwd` / ".pkg"
pkgname := "apparmor.d"
gpgkey := "06A26D531D56C42D66805049C5469996F0DF68EC"

# The following variables are only  used for the development and test VM

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

usage := "
Build variables available:
    build        " + BLUE + "# Build directory (default: " + build + ")" + NORMAL + "
    destdir      " + BLUE + "# Installation destination (default: " + destdir + ")" + NORMAL + "
    pkgdest      " + BLUE + "# Package output directory (default: " + pkgdest + ")" + NORMAL + "

Development variables available:
    username     " + BLUE + "# VM username (default: " + username + ")" + NORMAL + "
    password     " + BLUE + "# VM password (default: " + password + ")" + NORMAL + "
    disk_size    " + BLUE + "# VM disk size (default: " + disk_size + ")" + NORMAL + "
    vcpus        " + BLUE + "# VM CPU (default: " + vcpus + ")" + NORMAL + "
    ram          " + BLUE + "# VM RAM (default: " + ram + ")" + NORMAL + "

See https://apparmor.pujol.io/development/ for more information."

# Show this help message
help:
	@just --list --unsorted
	@printf "%s\n" "{{usage}}"

# Build the go programs
[group('build')]
build:
	@go build -o {{build}}/ ./cmd/aa-log
	@go build -o {{build}}/ ./cmd/prebuild

# Prebuild the profiles in enforced mode
[group('build')]
enforce: build
	@./{{build}}/prebuild --buildir {{build}}

# Prebuild the profiles in enforce mode (test)
[group('build')]
enforce-test: build
	@./{{build}}/prebuild --buildir {{build}} --test

# Prebuild the profiles in complain mode
[group('build')]
complain: build
	./{{build}}/prebuild --buildir {{build}} --complain

# Prebuild the profiles in complain mode (test)
[group('build')]
complain-test: build
	@./{{build}}/prebuild --buildir {{build}} --complain --test

# Prebuild the profiles in FSP mode
[group('build')]
fsp: build
	@./{{build}}/prebuild --buildir {{build}} --full

# Prebuild the profiles in FSP mode (complain)
[group('build')]
fsp-complain: build
	@./{{build}}/prebuild --buildir {{build}} --complain --full

# Prebuild the profiles in FSP mode (debug)
[group('build')]
fsp-debug: build
	@./{{build}}/prebuild --buildir {{build}} --complain --full --debug

# Prebuild the profiles in server mode
[group('build')]
server: build
	@./{{build}}/prebuild --buildir {{build}} --server

# Prebuild the profiles in server mode (complain)
[group('build')]
server-complain: build
	@./{{build}}/prebuild --buildir {{build}} --server --complain

# Prebuild the profiles in server FSP mode
[group('build')]
server-fsp: build
	@./{{build}}/prebuild --buildir {{build}} --server --full

# Prebuild the profiles in server FSP mode (complain)
[group('build')]
server-fsp-complain: build
	@./{{build}}/prebuild --buildir {{build}} --server --full --complain

# Prebuild the profiles in server FSP mode (debug)
[group('build')]
server-fsp-debug: build
	@./{{build}}/prebuild --buildir {{build}} --server --full --complain --debug

# Install prebuild profiles
[group('install')]
install:
	#!/usr/bin/env bash
	set -eu -o pipefail
	install -Dm0755 {{build}}/aa-log {{destdir}}/usr/bin/aa-log
	mapfile -t share < <(find "{{build}}/share" -type f -not -name "*.md" -printf "%P\n")
	for file in "${share[@]}"; do
		install -Dm0644 "{{build}}/share/$file" "{{destdir}}/usr/share/$file"
	done
	mapfile -t aa < <(find "{{build}}/apparmor.d" -type f -printf "%P\n")
	for file in "${aa[@]}"; do
		install -Dm0644 "{{build}}/apparmor.d/$file" "{{destdir}}/etc/apparmor.d/$file"
	done
	mapfile -t links < <(find "{{build}}/apparmor.d" -type l -printf "%P\n")
	for file in "${links[@]}"; do
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

# Locally install prebuild profiles
[group('install')]
local +names:
	#!/usr/bin/env bash
	set -eu -o pipefail
	install -Dm0755 {{build}}/aa-log {{destdir}}/usr/bin/aa-log
	mapfile -t abs < <(find "{{build}}/apparmor.d/abstractions" -type f -printf "%P\n")
	for file in "${abs[@]}"; do
		install -Dm0644 "{{build}}/apparmor.d/abstractions/$file" "{{destdir}}/etc/apparmor.d/abstractions/$file"
	done;
	mapfile -t tunables < <(find "{{build}}/apparmor.d/tunables" -type f -printf "%P\n")
	for file in "${tunables[@]}"; do
		install -Dm0644 "{{build}}/apparmor.d/tunables/$file" "{{destdir}}/etc/apparmor.d/tunables/$file"
	done;
	echo "Warning: profile dependencies fallback to unconfined."
	for file in {{names}}; do
		grep -Ei 'rPx|rpx' "{{build}}/apparmor.d/$file" || true
		sed -i -e "s/rPx/rPUx/g" "{{build}}/apparmor.d/$file"
		install -Dvm0644 "{{build}}/apparmor.d/$file" "{{destdir}}/etc/apparmor.d/$file"
	done;
	systemctl restart apparmor || sudo journalctl -xeu apparmor.service

# Prebuild, install, and load a dev profile
[group('install')]
dev name:
	go run ./cmd/prebuild --complain --file `find apparmor.d -iname {{name}}`
	sudo install -Dm644 {{build}}/apparmor.d/{{name}} /etc/apparmor.d/{{name}}
	sudo apparmor_parser --write-cache --replace /etc/apparmor.d/{{name}}

# Build & install apparmor.d on Arch based systems
[group('packages')]
pkg:
	@makepkg --syncdeps --install --cleanbuild --force --noconfirm

# Build & install apparmor.d on Debian based systems
[group('packages')]
dpkg:
	@bash dists/build.sh dpkg
	@sudo dpkg -i {{pkgdest}}/{{pkgname}}_*.deb

# Build & install apparmor.d on OpenSUSE based systems
[group('packages')]
rpm:
	@bash dists/build.sh rpm
	@sudo rpm -ivh --force {{pkgdest}}/{{pkgname}}-*.rpm

# Run the unit tests
[group('tests')]
tests:
	@go test ./cmd/... -v -cover -coverprofile=coverage.out
	@go test ./pkg/... -v -cover -coverprofile=coverage.out
	@go tool cover -func=coverage.out

# Run the linters
[group('linter')]
lint:
	golangci-lint run
	packer fmt tests/packer/
	packer validate --syntax-only tests/packer/
	shellcheck --shell=bash \
		PKGBUILD dists/build.sh dists/docker.sh tests/check.sh \
		tests/packer/init.sh tests/packer/src/aa-update tests/packer/clean.sh \
		debian/{{pkgname}}.postinst debian/{{pkgname}}.postrm

# Run style checks on the profiles
[group('linter')]
check:
	@bash tests/check.sh

# Generate the man pages
[group('docs')]
man:
	@pandoc -t man -s -o share/man/man8/aa-log.8 share/man/man8/aa-log.md

# Build the documentation
[group('docs')]
docs:
	@ENABLED_GIT_REVISION_DATE=false MKDOCS_OFFLINE=true mkdocs build --strict

# Serve the documentation
[group('docs')]
serve:
	@ENABLED_GIT_REVISION_DATE=false MKDOCS_OFFLINE=false mkdocs serve

# Remove all build artifacts
clean:
	@rm -rf \
		debian/.debhelper debian/debhelper* debian/*.debhelper debian/{{pkgname}} \
		{{pkgdest}}/{{pkgname}}* {{build}} coverage.out

# Build the package in a clean OCI container
[group('packages')]
package dist:
	#!/usr/bin/env bash
	set -eu -o pipefail
	dist="{{dist}}"
	version=""
	if [[ $dist =~ ubuntu([0-9]+) ]]; then
		version="${BASH_REMATCH[1]}.04"
		dist="ubuntu"
	elif [[ $dist == debian* ]]; then
		version="trixie"
		dist="debian"
	fi
	bash dists/docker.sh $dist $version

# Build the VM image
[group('vm')]
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

# Create the machine
[group('vm')]
create dist flavor:
	@cp -f {{base_dir}}/{{prefix}}{{dist}}-{{flavor}}.qcow2 {{vm}}/{{prefix}}{{dist}}-{{flavor}}.qcow2
	@virt-install {{c}} \
		--import \
		--name {{prefix}}{{dist}}-{{flavor}} \
		--vcpus {{vcpus}} \
		--ram {{ram}} \
		--machine q35 \
		{{ if dist == "archlinux" { "" } else { "--boot uefi" } }} \
		--memorybacking source.type=memfd,access.mode=shared \
		--disk path={{vm}}/{{prefix}}{{dist}}-{{flavor}}.qcow2,format=qcow2,bus=virtio \
		--filesystem "`pwd`,0a31bc478ef8e2461a4b1cc10a24cc4",accessmode=passthrough,driver.type=virtiofs \
		--os-variant "`just _get_osinfo {{dist}}`" \
		--graphics spice \
		--audio id=1,type=spice \
		--sound model=ich9 \
		--noautoconsole

# Start a machine
[group('vm')]
up dist flavor:
	@virsh {{c}} start {{prefix}}{{dist}}-{{flavor}}

# Stops the machine
[group('vm')]
halt dist flavor:
	@virsh {{c}} shutdown {{prefix}}{{dist}}-{{flavor}}

# Reboot the machine
[group('vm')]
reboot dist flavor:
	@virsh {{c}} reboot {{prefix}}{{dist}}-{{flavor}}

# Destroy the machine
[group('vm')]
destroy dist flavor:
	@virsh {{c}} destroy {{prefix}}{{dist}}-{{flavor}} || true
	@virsh {{c}} undefine {{prefix}}{{dist}}-{{flavor}} --nvram
	@rm -fv {{vm}}/{{prefix}}{{dist}}-{{flavor}}.qcow2

# Connect to the machine
[group('vm')]
ssh dist flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{dist}} {{flavor}}`

# Mount the shared directory on the machine
[group('vm')]
mount dist flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{dist}} {{flavor}}` \
		sh -c 'mount | grep 0a31bc478ef8e2461a4b1cc10a24cc4 || sudo mount 0a31bc478ef8e2461a4b1cc10a24cc4'

# Unmout the shared directory on the machine
[group('vm')]
umount dist flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{dist}} {{flavor}}` \
		sh -c 'true; sudo umount /home/{{username}}/Projects/apparmor.d || true'

# List the machines
[group('vm')]
list:
	@printf "{{BOLD}} %-4s %-22s %s{{NORMAL}}\n" "Id" "Distribution-Flavor" "State"
	@virsh {{c}} list --all | grep {{prefix}} | sed 's/{{prefix}}//g'

# List the VM images
[group('vm')]
images:
	#!/usr/bin/env bash
	set -eu -o pipefail
	mkdir -p {{base_dir}}
	ls -lh {{base_dir}} | awk '
	BEGIN {
		printf("{{BOLD}}%-18s %-10s %-5s %s{{NORMAL}}\n", "Distribution", "Flavor", "Size", "Date")
	}
	{
		if ($9 ~ /^{{prefix}}.*\.qcow2$/) {
			split($9, arr, "-|\\.")
			printf("%-18s %-10s %-5s %s %s %s\n", arr[2], arr[3], $5, $6, $7, $8)
		}
	}
	'

# List the VM images that can be created
[group('vm')]
available:
	#!/usr/bin/env bash
	set -eu -o pipefail
	ls -lh tests/cloud-init | awk '
	BEGIN {
		printf("{{BOLD}}%-18s %s{{NORMAL}}\n", "Distribution", "Flavor")
	}
	{
		if ($9 ~ /^.*\.user-data.yml$/) {
			split($9, arr, "-|\\.")
			printf("%-18s %s\n", arr[1], arr[2])
		}
	}
	'

# Install dependencies for the integration tests
[group('tests')]
init:
	@bash tests/requirements.sh

# Run the integration tests
[group('tests')]
integration name="":
	bats --recursive --timing --print-output-on-failure tests/integration/{{name}}

# Install dependencies for the integration tests (machine)
[group('tests')]
tests-init dist flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{dist}} {{flavor}}` \
		just --justfile /home/{{username}}/Projects/apparmor.d/Justfile init

# Synchronize the integration tests (machine)
[group('tests')]
tests-sync dist flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{dist}} {{flavor}}` \
		rsync -a --delete /home/{{username}}/Projects/apparmor.d/tests/ /home/{{username}}/Projects/tests/

# Re-synchronize the integration tests (machine)
[group('tests')]
tests-resync dist flavor: (mount dist flavor) \
	(tests-sync dist flavor) \
	(umount dist flavor)

# Run the integration tests (machine)
[group('tests')]
tests-run dist flavor name="": (tests-resync dist flavor)
	ssh {{sshopt}} {{username}}@`just _get_ip {{dist}} {{flavor}}` \
		bats --recursive --pretty --timing --print-output-on-failure \
			/home/{{username}}/Projects/tests/integration/{{name}}

# Get the current apparmor.d release version
[group('version')]
version:
	@bash -c 'source PKGBUILD && echo "$pkgver"'

# Create a new version number from the current release
[group('version')]
version-new:
	@bash -c 'source PKGBUILD && echo $(echo "$pkgver" | awk "{print \$1 + 0.0001}")'

# Create a new release
[group('release')]
release: tests lint commit archive publish

# Write the new release version to package files & commit
[group('release')]
commit:
	#!/usr/bin/env bash
	set -eu -o pipefail
	version=`just version-new`
	cat > debian/changelog.tmp <<-EOF
		{{pkgname}} (${version}-1) stable; urgency=medium

		* Release {{pkgname}} v${version}

		-- $(git config user.name) <$(git config user.email)>  $(date -R)

	EOF
	cat debian/changelog >> debian/changelog.tmp
	mv debian/changelog.tmp debian/changelog
	sed -i "s/^pkgver=.*/pkgver=$version/" PKGBUILD
	sed -i "s/^Version:.*/Version:        $version/" "dists/{{pkgname}}.spec"
	echo git add PKGBUILD "dists/{{pkgname}}.spec" debian/changelog
	echo git commit -S -m "Release version $version"

# Create a release archive
[group('release')]
archive:
	#!/usr/bin/env bash
	set -eu -o pipefail
	version=`just version`
	git tag -a "v$version" -m "{{pkgname}} v$version" --local-user={{gpgkey}}
	git archive \
		--format=tar.gz \
		--prefix={{pkgname}}-$version/ \
		--output={{pkgdest}}/{{pkgname}}-$version.tar.gz \
		v$version
	gpg --armor --default-key {{gpgkey}} --detach-sig {{pkgdest}}/{{pkgname}}-$version.tar.gz
	gpg --verify {{pkgdest}}/{{pkgname}}-$version.tar.gz.asc

# Publish the new release on Github
[group('release')]
publish:
	#!/usr/bin/env bash
	set -eu -o pipefail
	owner="roddhjav"
	version=`just version`
	git push origin main --tags
	gh release create "v$version" --notes-from-tag --repo $owner/{{pkgname}}
	gh release upload "v$version" --repo $owner/{{pkgname}} \
		{{pkgdest}}/{{pkgname}}-$version.tar.gz \
		{{pkgdest}}/{{pkgname}}-$version.tar.gz.asc

_get_ip dist flavor:
	@virsh --quiet --readonly {{c}} domifaddr {{prefix}}{{dist}}-{{flavor}} | \
		head -1 | \
		grep -E -o '([[:digit:]]{1,3}\.){3}[[:digit:]]{1,3}'

_get_osinfo dist:
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
