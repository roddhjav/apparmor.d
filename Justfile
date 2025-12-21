# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2025 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: `just`
# See https://apparmor.pujol.io/development/ for more information.

# Globally override any variables
set allow-duplicate-variables
import? '~/.aa.just'

# Build settings

destdir := "/"
build := ".build"
pkgdest := `pwd` / ".pkg"
pkgname := "apparmor.d"
gpgkey := "06A26D531D56C42D66805049C5469996F0DF68EC"
sign := "false"

# Prebuild options, only used for the `dev` install target
opt := "complain"

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

# Path to the UEFI firmware
firmware := if path_exists("/usr/share/edk2/x64/OVMF.4m.fd") == "true" {
	"/usr/share/edk2/x64/OVMF.4m.fd"
} else {
	"/usr/share/ovmf/OVMF.fd"
}

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
    opt          " + BLUE + "# Prebuild option, only used for the dev install target (default: " + opt + ")" + NORMAL + "

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
	systemctl restart apparmor.service || journalctl -xeu apparmor.service

# Prebuild, install, and load a dev profile
[group('install')]
dev +names:
	go run ./cmd/prebuild --{{opt}}
	for file in {{names}}; do \
		sudo install -Dm644 -v {{build}}/apparmor.d/$file /etc/apparmor.d/$file; \
	done
	sudo systemctl restart apparmor.service || sudo journalctl -xeu apparmor.service

# Build the package on Arch Linux
[group('packages')]
build-pkg: (_ensure_pkgdest)
	@PKGDEST={{pkgdest}} BUILDDIR=/tmp/makepkg makepkg --syncdeps --force --cleanbuild --noconfirm --noprogressbar

# Build the package on Debian
[group('packages')]
build-dpkg: (_ensure_pkgdest)
	dpkg-buildpackage -b -d {{ if sign == "true" { "--sign-key=" + gpgkey } else { "--no-sign" } }}
	lintian --color always --display-info --pedantic --tag-display-limit 0 || true
	mv ../{{pkgname}}*.deb {{pkgdest}}/

# Build the package on OpenSUSE
[group('packages')]
build-rpm: (_ensure_pkgdest)
	#!/usr/bin/env bash
	set -eu -o pipefail
	RPMBUILD_ROOT=$(mktemp -d /tmp/$PKGNAME.XXXXXX)
	ARCH=$(uname -m)
	VERSION="$(just version)"
	readonly RPMBUILD_ROOT ARCH VERSION

	mkdir -p "$RPMBUILD_ROOT"/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS/tmp}
	cp -p "dists/$PKGNAME.spec" "$RPMBUILD_ROOT/SPECS"
	tar -czf "$RPMBUILD_ROOT/SOURCES/$PKGNAME-$VERSION.tar.gz" --transform "s,^,$PKGNAME-$VERSION/," ./*

	cd "$RPMBUILD_ROOT"
	rpmbuild -bb --define "_topdir $RPMBUILD_ROOT" "SPECS/$PKGNAME.spec"

	mv "$RPMBUILD_ROOT/RPMS/$ARCH/"*.rpm "{{pkgdest}}/"
	rm -rf "$RPMBUILD_ROOT"

# Build & install apparmor.d on Arch based systems
[group('packages')]
pkg name="": (build-pkg)
	@sudo pacman -U --noconfirm \
		{{pkgdest}}/{{pkgname}}{{ if name != "" { "." + name } else { "" } }}-`just version`*.pkg.tar.zst 

# Build & install apparmor.d on Debian based systems
[group('packages')]
dpkg name="": (build-dpkg)
	@sudo dpkg -i  \
		{{pkgdest}}/{{pkgname}}{{ if name != "" { "." + name } else { "" } }}_`just version`*.deb

# Build & install apparmor.d on OpenSUSE based systems
[group('packages')]
rpm name="": (build-rpm)
	@sudo rpm -ivh --force \
		{{pkgdest}}/{{pkgname}}{{ if name != "" { "." + name } else { "" } }}-`just version`*.rpm

# Run the linters
[group('linter')]
lint:
	golangci-lint run
	packer fmt tests/packer/
	packer validate --syntax-only tests/packer/
	shellcheck --shell=bash \
		PKGBUILD dists/*.sh tests/check.sh \
		tests/packer/*.sh tests/packer/src/aa-update \
		tests/autopkgtest/autopkgtest.sh debian/common.postinst debian/common.postrm

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
		debian/.debhelper debian/debhelper* debian/*.debhelper debian/{{pkgname}}* \
		debian/*.substvars debian/*.debhelper debian/files \
		{{pkgdest}}/{{pkgname}}* {{pkgdest}}/ubuntu {{pkgdest}}/debian \
		{{pkgdest}}/archlinux {{pkgdest}}/opensuse \
		{{build}} coverage.out .logs/autopkgtest/ site .cache

# Build the package in a clean OCI container
[group('packages')]
package dist release="" flavor="":
	bash dists/docker.sh {{dist}} {{release}} {{flavor}}

# Build all packages in a clean OCI container
[group('packages')]
packages: (clean)
	#!/usr/bin/env bash
	set -eu -o pipefail
	declare -A matrix=(
		# ["archlinux"]="-"
		["debian"]="13"
		["ubuntu"]="24.04 25.10 26.04"
		# ["opensuse"]="-"
	)
	for dist in "${!matrix[@]}"; do
		IFS=' ' read -r -a releases <<< "${matrix[$dist]}"
		for release in "${releases[@]}"; do
			bash dists/docker.sh $dist $release
		done
	done

# Build the VM image
[group('vm')]
img dist release flavor: (package dist release flavor)
	#!/usr/bin/env bash
	set -eu -o pipefail
	RELEASE="{{release}}"
	[[ "$RELEASE" == "-" ]] && RELEASE=""
	mkdir -p {{base_dir}}
	packer build -force \
		-var dist={{dist}} \
		-var release="$RELEASE" \
		-var flavor={{flavor}} \
		-var prefix={{prefix}} \
		-var username={{username}} \
		-var password={{password}} \
		-var ssh_publickey={{ssh_publickey}} \
		-var disk_size={{disk_size}} \
		-var cpus={{vcpus}} \
		-var ram={{ram}} \
		-var firmware={{firmware}} \
		-var base_dir={{base_dir}} \
		-var output_dir={{output_dir}} \
		tests/packer/

# Create the machine
[group('vm')]
create osinfo flavor:
	@cp -f {{base_dir}}/{{prefix}}{{osinfo}}-{{flavor}}.qcow2 {{vm}}/{{prefix}}{{osinfo}}-{{flavor}}.qcow2
	@virt-install {{c}} \
		--import \
		--name {{prefix}}{{osinfo}}-{{flavor}} \
		--vcpus {{vcpus}} \
		--ram {{ram}} \
		--machine q35 \
		{{ if osinfo == "archlinux" { "" } else { "--boot uefi" } }} \
		--memorybacking source.type=memfd,access.mode=shared \
		--disk path={{vm}}/{{prefix}}{{osinfo}}-{{flavor}}.qcow2,format=qcow2,bus=virtio \
		--filesystem "`pwd`,0a31bc478ef8e2461a4b1cc10a24cc4",accessmode=passthrough,driver.type=virtiofs \
		--os-variant "{{ if osinfo == "opensuse" { "opensusetumbleweed" } else { osinfo } }}" \
		--graphics spice \
		--audio id=1,type=spice \
		--sound model=ich9 \
		--noautoconsole

# Start a machine
[group('vm')]
up osinfo flavor:
	@virsh {{c}} start {{prefix}}{{osinfo}}-{{flavor}}

# Stops the machine
[group('vm')]
halt osinfo flavor:
	@virsh {{c}} shutdown {{prefix}}{{osinfo}}-{{flavor}}

# Reboot the machine
[group('vm')]
reboot osinfo flavor:
	@virsh {{c}} reboot {{prefix}}{{osinfo}}-{{flavor}}

# Destroy the machine
[group('vm')]
destroy osinfo flavor:
	@virsh {{c}} destroy {{prefix}}{{osinfo}}-{{flavor}} || true
	@virsh {{c}} undefine {{prefix}}{{osinfo}}-{{flavor}} --nvram
	@rm -fv {{vm}}/{{prefix}}{{osinfo}}-{{flavor}}.qcow2

# List all snapshots for a machine
[group('vm')]
snapshots osinfo flavor:
	@virsh {{c}} snapshot-list {{prefix}}{{osinfo}}-{{flavor}}

# Snapshot a machine
[group('vm')]
snapshot osinfo flavor snapname:
    @virsh {{c}} snapshot-create-as {{prefix}}{{osinfo}}-{{flavor}} --name {{snapname}}

# Restore a machine to a specified snapshot
[group('vm')]
restore osinfo flavor snapname:
    @virsh {{c}} snapshot-revert {{prefix}}{{osinfo}}-{{flavor}} {{snapname}}

# Delete a specified snapshot from a machine
[group('vm')]
delete osinfo flavor snapname:
	@virsh {{c}} snapshot-delete {{prefix}}{{osinfo}}-{{flavor}} {{snapname}}

# Connect to the machine
[group('vm')]
ssh osinfo flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{osinfo}} {{flavor}}`

# Mount the shared directory on the machine
[group('vm')]
mount osinfo flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{osinfo}} {{flavor}}` \
		sh -c 'mount | grep 0a31bc478ef8e2461a4b1cc10a24cc4 || sudo mount 0a31bc478ef8e2461a4b1cc10a24cc4'

# Unmout the shared directory on the machine
[group('vm')]
umount osinfo flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{osinfo}} {{flavor}}` \
		sh -c 'true; sudo umount /home/{{username}}/Projects/apparmor.d || true'

# List the machines
[group('vm')]
list:
	@printf "{{BOLD}} %-4s %-22s %s{{NORMAL}}\n" "Id" "OsInfo-Flavor" "State"
	@virsh {{c}} list --all | grep {{prefix}} | sed 's/{{prefix}}//g'

# List the VM images
[group('vm')]
images:
	#!/usr/bin/env bash
	set -eu -o pipefail
	mkdir -p {{base_dir}}
	(
		printf "{{BOLD}}%s %s %s %s %s{{NORMAL}}\n" "OsInfo" "Flavor" "Size" "Date"
		find {{base_dir}} -iname '{{prefix}}*' -type f -printf "%f %k %Tb %Td %TH:%TM\n" | sort | awk '
			{
				split($1, item, "-")
				split(item[3], flavor, "\\.")
				if ($2>=1048576) {
					printf("%s %s %.1fGB %s %s %s\n", item[2], flavor[1], $2/1048576, $3, $4, $5)
				} else {
					printf("%s %s %.fMB %s %s %s\n", item[2], flavor[1], $2/1024, $3, $4, $5)
				}
			}
			'
	) | column -t

# List the VM images that can be created
[group('vm')]
available:
	#!/usr/bin/env bash
	set -eu -o pipefail
	(
		printf "{{BOLD}}%s %s %s{{NORMAL}}\n" "Distribution" "Release" "Flavor"
		find tests/cloud-init -iname '*.user-data.yml' -type f -printf "%f\n" | sort | awk '
			{
				split($1, item, "-")
				match(item[1], /^([a-z]+)([0-9.]*?)$/, osinfo)
				release = (osinfo[2] == "" ? "-" : osinfo[2])
				split(item[2], flavor, "\\.")
				printf("%s %s %s\n", osinfo[1], release, flavor[1])
			}
			'
	) | column -t

# Run the unit tests
[group('tests')]
tests:
	@go test ./cmd/... -v -cover -coverprofile=coverage.out
	@go test ./pkg/... -v -cover -coverprofile=coverage.out
	@go tool cover -func=coverage.out

# Run the autopkgtest tests
[group('tests')]
autopkgtest osinfo:
	@PREFIX='{{prefix}}' VM_DIR='{{vm}}' \
	USER='{{username}}' PASSWORD='{{password}}' SSH_OPT='{{sshopt}}' \
		bash tests/autopkgtest/autopkgtest.sh run {{osinfo}}

# Update the apparmor.d package on the test machine
[group('tests')]
autopkgtest-update dist release:
	just up {{dist}}{{release}} test || true
	just package {{dist}} {{release}} test
	scp {{sshopt}} {{pkgdest}}/{{dist}}/{{release}}/{{pkgname}}_*.deb \
		{{username}}@`just _get_ip {{dist}}{{release}} test`:/home/{{username}}/Projects/
	ssh {{sshopt}} {{username}}@`just _get_ip {{dist}}{{release}} test` \
		sudo dpkg -i /home/{{username}}/Projects/{{pkgname}}_*.deb
	just halt {{dist}}{{release}} test

_autopkgtest-log-merge:
	@mkdir -p .logs/autopkgtest
	@cat .logs/autopkgtest/aa-log-* > .logs/autopkgtest/merged.log

# Report all collected logs
[group('tests')]
autopkgtest-log: (_autopkgtest-log-merge)
	@aa-log --file .logs/autopkgtest/merged.log

# Report all generated rules
[group('tests')]
autopkgtest-rules: (_autopkgtest-log-merge)
	@aa-log --rules --file .logs/autopkgtest/merged.log

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
tests-init osinfo flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{osinfo}} {{flavor}}` \
		just --justfile /home/{{username}}/Projects/apparmor.d/Justfile init

# Synchronize the integration tests (machine)
[group('tests')]
tests-sync osinfo flavor:
	@ssh {{sshopt}} {{username}}@`just _get_ip {{osinfo}} {{flavor}}` \
		rsync -a --delete /home/{{username}}/Projects/apparmor.d/tests/ /home/{{username}}/Projects/tests/

# Re-synchronize the integration tests (machine)
[group('tests')]
tests-resync osinfo flavor: (mount osinfo flavor) \
	(tests-sync osinfo flavor) \
	(umount osinfo flavor)

# Run the integration tests (machine)
[group('tests')]
tests-run osinfo flavor name="": (tests-resync osinfo flavor)
	ssh {{sshopt}} {{username}}@`just _get_ip {{osinfo}} {{flavor}}` \
		bats --recursive --pretty --timing --print-output-on-failure \
			/home/{{username}}/Projects/tests/integration/{{name}}

# Get the current apparmor.d release version
[group('version')]
version:
	@bash -c 'source PKGBUILD && echo "$pkgver"'

# Create a new version number from the current release
[group('version')]
version-new:
	@bash -c 'source PKGBUILD && awk -v ver="$pkgver" "BEGIN {printf \"%.4f\n\", ver + 0.0001}"'

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
	git add PKGBUILD "dists/{{pkgname}}.spec" debian/changelog
	git commit -S -m "Release version $version"
	git tag -a "v$version" -m "{{pkgname}} v$version" --local-user={{gpgkey}}

# Create a release archive
[group('release')]
archive:
	#!/usr/bin/env bash
	set -eu -o pipefail
	version=`just version`
	mkdir -p {{pkgdest}}/release/ 
	git archive \
		--format=tar.gz \
		--prefix={{pkgname}}-$version/ \
		--output={{pkgdest}}/release/{{pkgname}}-$version.tar.gz \
		v$version
	gpg --armor --default-key {{gpgkey}} --detach-sig {{pkgdest}}/release/{{pkgname}}-$version.tar.gz
	gpg --verify {{pkgdest}}/release/{{pkgname}}-$version.tar.gz.asc

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


_ensure_pkgdest:
	@mkdir -p {{pkgdest}}

_get_ip osinfo flavor:
	@virsh --quiet --readonly {{c}} domifaddr {{prefix}}{{osinfo}}-{{flavor}} | \
		head -1 | \
		grep -E -o '([[:digit:]]{1,3}\.){3}[[:digit:]]{1,3}'
