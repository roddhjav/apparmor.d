#!/usr/bin/make -f
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2022 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

DESTDIR ?= /
BUILD := .build
PKGNAME := apparmor.d
DISTRIBUTION := $(shell lsb_release --id --short)
VERSION := 0.$(shell git rev-list --count HEAD)-1
P = $(notdir $(wildcard ${BUILD}/apparmor.d/*))

.PHONY: all install auto local $(P) lint archlinux debian ubuntu whonix clean

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

auto:
	@[ ${DISTRIBUTION} = Arch ] || exit 0; \
		makepkg --syncdeps --install --cleanbuild --force
	@[ ${DISTRIBUTION} = Ubuntu ] || exit 0; \
		dch --newversion="${VERSION}" --urgency=medium --distribution=stable --controlmaint "Release ${VERSION}";  \
		dpkg-buildpackage -b -d --no-sign;  \
        sudo dpkg -i "../apparmor.d_${VERSION}_all.deb"; \
		make clean
	@[ ${DISTRIBUTION} = openSUSE ] || exit 0; \
		make local

local:
	@./configure --complain
	@make
	@sudo make install
	@sudo systemctl restart apparmor || sudo systemctl status apparmor

ABSTRACTIONS = $(shell find ${BUILD}/apparmor.d/abstractions/ -type f -printf "%P\n")
TUNABLES = $(shell find ${BUILD}/apparmor.d/tunables/ -type f -printf "%P\n")
$(P):
	@[ -f ${BUILD}/aa-log ] || exit 0; install -Dm755 ${BUILD}/aa-log ${DESTDIR}/usr/bin/aa-log
	@for file in ${ABSTRACTIONS}; do \
		install -Dm0644 "${BUILD}/apparmor.d/abstractions/$${file}" "${DESTDIR}/etc/apparmor.d/abstractions/$${file}"; \
	done;
	@for file in ${TUNABLES}; do \
		install -Dm0644 "${BUILD}/apparmor.d/tunables/$${file}" "${DESTDIR}/etc/apparmor.d/tunables/$${file}"; \
	done;
	@echo "Warning: profile dependencies fallback to unconfined."
	@for file in ${@}; do \
		grep 'rPx' "${BUILD}/apparmor.d/$${file}"; \
		sed -i -e "s/rPx/rPUx/g" "${BUILD}/apparmor.d/$${file}"; \
		install -Dvm0644 "${BUILD}/apparmor.d/$${file}" "${DESTDIR}/etc/apparmor.d/$${file}"; \
	done;
	@systemctl restart apparmor || systemctl status apparmor

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
		debian/.debhelper debian/debhelper* debian/*.debhelper debian/${PKGNAME} \
		${PKGNAME}-*.pkg.tar.zst.sig ${PKGNAME}-*.pkg.tar.zst \
		${PKGNAME}_*.* ${BUILD}
