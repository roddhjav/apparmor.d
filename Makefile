#!/usr/bin/make -f
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2022 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

DESTDIR ?= /
BUILD := .build
PKGNAME := apparmor.d

.PHONY: all install lint archlinux debian ubuntu whonix clean

all:
	@go build -o ${BUILD}/ ./cmd/aa-log

ROOT = $(shell find "${BUILD}/root" -type f -printf "%P\n")
PROFILES = $(shell find "${BUILD}/apparmor.d" -type f -printf "%P\n")
install:
	@install -Dm755 ${BUILD}/aa-log ${DESTDIR}/usr/bin/aa-log
	@for file in ${ROOT}; do \
		install -Dm0644 "${BUILD}/root/$${file}" "${DESTDIR}/$${file}"; \
	done;
	@for file in ${PROFILES}; do \
		install -Dm0644 "${BUILD}/apparmor.d/$${file}" "${DESTDIR}/etc/apparmor.d/$${file}"; \
	done;
	@for file in systemd/system/*; do \
		service="$$(basename "$$file")"; \
		install -Dm0644 "$${file}" "${DESTDIR}/usr/lib/systemd/system/$${service}.d/apparmor.conf"; \
	done;
	@for file in systemd/user/*; do \
		service="$$(basename "$$file")"; \
		install -Dm0644 "$${file}" "${DESTDIR}/usr/lib/systemd/user/$${service}.d/apparmor.conf"; \
	done

lint:
	@shellcheck --shell=bash \
		PKGBUILD configure pick dists/build/build.sh \
		debian/${PKGNAME}.postinst debian/${PKGNAME}.postrm

archlinux:
	@bash dists/build/build.sh archlinux

debian:
	@bash dists/build/build.sh debian

ubuntu:
	@bash dists/build/build.sh ubuntu

whonix:
	@bash dists/build/build.sh whonix

clean:
	@rm -rf \
		debian/.debhelper debian/debhelper* debian/*.debhelper \
		${PKGNAME}-*.pkg.tar.zst.sig ${PKGNAME}-*.pkg.tar.zst \
		${PKGNAME}_*.* ${BUILD}
