#!/usr/bin/make -f
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage:
#   make archlinux flavor=gnome
#   vagrant up arch-gnome
#   vagrant ssh archl-gnome

# Build variables
flavor ?= 
disk ?= 10G

VERSION := 0.$(shell git rev-list --count HEAD)-1
PACKERS = $(notdir $(patsubst %.pkr.hcl,%,$(wildcard packer/*.pkr.hcl)))

.PHONY: ${PACKERS} lint

$(PACKERS):
	make --directory=../ package dist=${@}
	packer build -force -var version=${VERSION} \
		-var disk_size=${disk} -var flavor="${flavor}" \
		-only=main.qemu.${@}-${flavor} packer/

lint:
	@packer fmt packer/
	@packer validate packer/
