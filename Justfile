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

base_dir := home_dir() / ".libvirt/base"
vm := home_dir() / ".vm"
output := base_dir / "packer"
prefix := "aa-"
c := "--connect=qemu:///system"
sshopt := "-o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no"

[doc('Show this help message')]
default:
	@echo -e "Integration environment helper for apparmor.d\n"
	@just --list --unsorted
	@echo -e "\nSee https://apparmor.pujol.io/development/vm/ for more information."

[doc('Build the apparmor.d package')]
package dist:
    #!/usr/bin/env bash
    set -eu -o pipefail
    dist="{{dist}}"
    [[ $dist =~ ubuntu* ]] && dist=ubuntu
    [[ $dist =~ debian* ]] && dist=debian
    make package dist=$dist

[doc('Build the image')]
img dist flavor: (package dist)
	@mkdir -p {{base_dir}}
	packer build -force \
		-var dist={{dist}} \
		-var flavor={{flavor}} \
		-var prefix={{prefix}} \
		-var base_dir={{base_dir}} \
		-var output={{output}} \
		tests/packer/

[doc('Create the machine')]
vm dist flavor:
	@cp -f {{base_dir}}/{{prefix}}{{dist}}-{{flavor}}.qcow2 {{vm}}/{{prefix}}{{dist}}-{{flavor}}.qcow2
	virt-install {{c}} \
		--import \
		--name {{prefix}}{{dist}}-{{flavor}} \
		--vcpus 6 \
		--ram 4096 \
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

[doc('Destroy the machine')]
destroy dist flavor:
	@virsh {{c}} destroy {{prefix}}{{dist}}-{{flavor}} || true
	@virsh {{c}} undefine {{prefix}}{{dist}}-{{flavor}} --nvram
	@rm -fv {{vm}}/{{prefix}}{{dist}}-{{flavor}}.qcow2

[doc('Connect to the machine')]
ssh dist flavor:
	@ssh {{sshopt}} user@`just get_ip {{dist}} {{flavor}}`

[doc('List the machines')]
list:
	@echo -e '\033[1m Id   Name                    State\033[0m'
	@virsh {{c}} list --all | grep {{prefix}}

[doc('List the machine images')]
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

[doc('Run the linters')]
lint:
	@packer fmt tests/packer/
	@packer validate --syntax-only tests/packer/

[doc('Remove the machine images')]
clean:
	@rm -fv {{base_dir}}/{{prefix}}*.qcow2

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
