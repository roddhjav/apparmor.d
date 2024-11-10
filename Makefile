#!/usr/bin/make -f
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2022-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

DESTDIR ?= /
BUILD ?= .build
PKGDEST ?= ${PWD}/.pkg
PKGNAME := apparmor.d
PROFILE = $(filter-out dpkg,$(notdir $(wildcard ${BUILD}/apparmor.d/*)))
PROFILES = profiles-apparmor.d profiles-other $(patsubst dists/packages/%,profiles-%,$(basename $(wildcard dists/packages/*.conf)))

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

.PHONY: full
full: build
	@./${BUILD}/prebuild --complain --full

.PHONY: packages
packages: clean build
	@./${BUILD}/prebuild --complain --packages 

# Install apparmor.d
.PHONY: install
install: install-bin install-share install-systemd profiles-apparmor.d

# Install apparmor.d.base
.PHONY: install-base
install-base: install-bin install-share install-systemd profiles-base

.PHONY: install-bin
install-bin:
	@install -Dm0755 ${BUILD}/aa-log ${DESTDIR}/usr/bin/aa-log

.PHONY: install-share
install-share:
	@for file in $(shell find "${BUILD}/share" -type f -not -name "*.md" -printf "%P\n"); do \
		install -Dm0644 "${BUILD}/share/$${file}" "${DESTDIR}/usr/share/$${file}"; \
	done;

.PHONY: install-systemd
install-systemd:
	@for file in ${BUILD}/systemd/system/*; do \
		service="$$(basename "$${file}")"; \
		install -Dm0644 "$${file}" "${DESTDIR}/usr/lib/systemd/system/$${service}.d/apparmor.conf"; \
	done;
	@for file in ${BUILD}/systemd/user/*; do \
		service="$$(basename "$${file}")"; \
		install -Dm0644 "$${file}" "${DESTDIR}/usr/lib/systemd/user/$${service}.d/apparmor.conf"; \
	done

# Install all profiles for a given (sub)package
.PHONY: $(PROFILES)
$(PROFILES):
	@for file in $(shell find "${BUILD}/$(patsubst profiles-%,%,$@)" -type f -printf "%P\n"); do \
		install -Dm0644 "${BUILD}/$(patsubst profiles-%,%,$@)/$${file}" "${DESTDIR}/etc/apparmor.d/$${file}"; \
	done;
	@for file in $(shell find "${BUILD}/$(patsubst profiles-%,%,$@)" -type l -printf "%P\n"); do \
		mkdir -p "${DESTDIR}/etc/apparmor.d/disable"; \
		cp -d "${BUILD}/$(patsubst profiles-%,%,$@)/$${file}" "${DESTDIR}/etc/apparmor.d/$${file}"; \
	done;

# Partial install (not recommended)
.PHONY: $(PROFILE)
$(PROFILE): install-bin
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

.PHONY: package
dist ?= archlinux
package:
	@bash dists/docker.sh ${dist}

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

.PHONY: tests
tests:
	@go test ./cmd/... -v -cover -coverprofile=coverage.out
	@go test ./pkg/... -v -cover -coverprofile=coverage.out
	@go tool cover -func=coverage.out

.PHONY: lint
lint:
	@golangci-lint run
	@make --directory=tests lint
	@shellcheck --shell=bash \
		PKGBUILD dists/build.sh dists/docker.sh tests/check.sh \
		tests/packer/init/init.sh tests/packer/src/aa-update tests/packer/init/clean.sh \
		debian/${PKGNAME}.postinst debian/${PKGNAME}.postrm

.PHONY: check
check:
	@bash tests/check.sh

.PHONY: bats
bats:
	@bats --print-output-on-failure tests/bats/

.PHONY: manual
manual:
	@pandoc -t man -s -o root/usr/share/man/man8/aa-log.8 root/usr/share/man/man8/aa-log.md

.PHONY: docs
docs:
	@ENABLED_GIT_REVISION_DATE=false MKDOCS_OFFLINE=true mkdocs build --strict

.PHONY: serve
serve:
	@ENABLED_GIT_REVISION_DATE=false MKDOCS_OFFLINE=false mkdocs serve

.PHONY: clean
clean:
	@rm -rf \
		debian/.debhelper debian/debhelper* debian/*.debhelper debian/${PKGNAME} \
		.pkg/${PKGNAME}* ${BUILD} coverage.out
