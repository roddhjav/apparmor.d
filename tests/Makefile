#!/usr/bin/make -f
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage:
#   make archlinux flavor=gnome
#   vagrant up arch-gnome
#   vagrant ssh archl-gnome

# Build variables
flavor ?= 
disk ?= 10G

VERSION := 0.$(shell git rev-list --count HEAD)
BASE = archlinux debian ubuntu opensuse fedora

.PHONY: ${BASE} lint

$(BASE):
	@make --directory=../ package dist=${@}
	@packer build -force -var version=${VERSION} \
		-var disk_size=${disk} -var flavor="${flavor}" \
		-only=qemu.${@} packer/

lint:
	@packer fmt --check packer/
	@packer validate --syntax-only packer/
