#!/usr/bin/make -f
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2022-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

DESTDIR ?= /
BUILD ?= .build
PKGDEST ?= ${PWD}/.pkg
PKGNAME := apparmor.d
PROFILES = $(filter-out dpkg,$(notdir $(wildcard ${BUILD}/apparmor.d/*)))

.PHONY: all
all: build
	@./${BUILD}/prebuild --complain

.PHONY: build
build:
	@go build -o ${BUILD}/ ./cmd/aa-log
	@go build -o ${BUILD}/ ./cmd/prebuild

.PHONY: enforce
enforce: build
	@./${BUILD}/prebuild

.PHONY: fsp
fsp: build
	@./${BUILD}/prebuild --full

.PHONY: fsp-complain
fsp-complain: build
	@./${BUILD}/prebuild --complain --full

.PHONY: install
install:
	@install -Dm0755 ${BUILD}/aa-log ${DESTDIR}/usr/bin/aa-log
	@for file in $(shell find "${BUILD}/share" -type f -not -name "*.md" -printf "%P\n"); do \
		install -Dm0644 "${BUILD}/share/$${file}" "${DESTDIR}/usr/share/$${file}"; \
	done;
	@for file in $(shell find "${BUILD}/apparmor.d" -type f -printf "%P\n"); do \
		install -Dm0644 "${BUILD}/apparmor.d/$${file}" "${DESTDIR}/etc/apparmor.d/$${file}"; \
	done;
	@for file in $(shell find "${BUILD}/apparmor.d" -type l -printf "%P\n"); do \
		mkdir -p "${DESTDIR}/etc/apparmor.d/disable"; \
		cp -d "${BUILD}/apparmor.d/$${file}" "${DESTDIR}/etc/apparmor.d/$${file}"; \
	done;
	@for file in ${BUILD}/systemd/system/*; do \
		service="$$(basename "$$file")"; \
		install -Dm0644 "$${file}" "${DESTDIR}/usr/lib/systemd/system/$${service}.d/apparmor.conf"; \
	done;
	@for file in ${BUILD}/systemd/user/*; do \
		service="$$(basename "$$file")"; \
		install -Dm0644 "$${file}" "${DESTDIR}/usr/lib/systemd/user/$${service}.d/apparmor.conf"; \
	done


.PHONY: $(PROFILES)
$(PROFILES):
	@install -Dm0755 ${BUILD}/aa-log ${DESTDIR}/usr/bin/aa-log
	@for file in $(shell find ${BUILD}/apparmor.d/abstractions/ -type f -printf "%P\n"); do \
		install -Dm0644 "${BUILD}/apparmor.d/abstractions/$${file}" "${DESTDIR}/etc/apparmor.d/abstractions/$${file}"; \
	done;
	@for file in $(shell find ${BUILD}/apparmor.d/tunables/ -type f -printf "%P\n"); do \
		install -Dm0644 "${BUILD}/apparmor.d/tunables/$${file}" "${DESTDIR}/etc/apparmor.d/tunables/$${file}"; \
	done;
	@echo "Warning: profile dependencies fallback to unconfined."
	@for file in ${@}; do \
		grep 'rPx' "${BUILD}/apparmor.d/$${file}"; \
		sed -i -e "s/rPx/rPUx/g" "${BUILD}/apparmor.d/$${file}"; \
		install -Dvm0644 "${BUILD}/apparmor.d/$${file}" "${DESTDIR}/etc/apparmor.d/$${file}"; \
	done;
	@systemctl restart apparmor || sudo journalctl -xeu apparmor.service

.PHONY: dev
name ?= 
dev:
	@go run ./cmd/prebuild --complain --file $(shell find apparmor.d -iname ${name})
	@sudo install -Dm644 ${BUILD}/apparmor.d/${name} /etc/apparmor.d/${name}
	@sudo systemctl restart apparmor || sudo journalctl -xeu apparmor.service

.PHONY: pkg
pkg:
	@makepkg --syncdeps --install --cleanbuild --force --noconfirm

.PHONY: dpkg
dpkg:
	@bash dists/build.sh dpkg
	@sudo dpkg -i ${PKGDEST}/${PKGNAME}_*.deb

.PHONY: rpm
rpm:
	@bash dists/build.sh rpm
	@sudo rpm -ivh --force  ${PKGDEST}/${PKGNAME}-*.rpm

.PHONY: check
check:
	@bash tests/check.sh

.PHONY: integration
integration:
	@bats --recursive --timing --print-output-on-failure tests/integration/
