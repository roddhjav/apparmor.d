#!/usr/bin/make -f
# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2022-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

DESTDIR ?= /
BUILD ?= .build
PKGDEST ?= ${PWD}/.pkg
PKGNAME := apparmor.d
P = $(filter-out dpkg,$(notdir $(wildcard ${BUILD}/apparmor.d/*)))

.PHONY: all build enforce full install local $(P) dev package pkg dpkg rpm tests lint check manual docs serve clean

all: build
	@./${BUILD}/prebuild --complain

build:
	@go build -o ${BUILD}/ ./cmd/aa-log
	@go build -o ${BUILD}/ ./cmd/prebuild

enforce: build
	@./${BUILD}/prebuild

full: build
	@./${BUILD}/prebuild --complain --full

SHARE = $(shell find "${BUILD}/share" -type f -not -name "*.md" -printf "%P\n")
PROFILES = $(shell find "${BUILD}/apparmor.d" -type f -printf "%P\n")
DISABLES = $(shell find "${BUILD}/apparmor.d" -type l -printf "%P\n")
install:
	@install -Dm0755 ${BUILD}/aa-log ${DESTDIR}/usr/bin/aa-log
	@for file in ${SHARE}; do \
		install -Dm0644 "${BUILD}/share/$${file}" "${DESTDIR}/usr/share/$${file}"; \
	done;
	@for file in ${PROFILES}; do \
		install -Dm0644 "${BUILD}/apparmor.d/$${file}" "${DESTDIR}/etc/apparmor.d/$${file}"; \
	done;
	@for file in ${DISABLES}; do \
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

local:
	@make
	@sudo make install
	@sudo systemctl restart apparmor || sudo systemctl status apparmor

ABSTRACTIONS = $(shell find ${BUILD}/apparmor.d/abstractions/ -type f -printf "%P\n")
TUNABLES = $(shell find ${BUILD}/apparmor.d/tunables/ -type f -printf "%P\n")
$(P):
	@install -Dm0755 ${BUILD}/aa-log ${DESTDIR}/usr/bin/aa-log
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

name ?= 
dev:
	@go run ./cmd/prebuild --complain --file $(shell find apparmor.d -iname ${name})
	@sudo install -Dm644 ${BUILD}/apparmor.d/${name} /etc/apparmor.d/${name}
	@sudo systemctl restart apparmor || systemctl status apparmor

dist ?= archlinux
package:
	@bash dists/docker.sh ${dist}

pkg:
	@makepkg --syncdeps --install --cleanbuild --force --noconfirm

dpkg:
	@bash dists/build.sh dpkg
	@sudo dpkg -i ${PKGDEST}/${PKGNAME}_*.deb

rpm:
	@bash dists/build.sh rpm
	@sudo rpm -ivh --force  ${PKGDEST}/${PKGNAME}-*.rpm

tests:
	@go test ./cmd/... -v -cover -coverprofile=coverage.out
	@go test ./pkg/... -v -cover -coverprofile=coverage.out
	@go tool cover -func=coverage.out

lint:
	@golangci-lint run
	@make --directory=tests lint
	@shellcheck --shell=bash \
		PKGBUILD dists/build.sh dists/docker.sh tests/check.sh \
		tests/packer/init/init.sh tests/packer/src/aa-update tests/packer/init/clean.sh \
		debian/${PKGNAME}.postinst debian/${PKGNAME}.postrm

check:
	@bash tests/check.sh

manual:
	@pandoc -t man -s -o root/usr/share/man/man8/aa-log.8 root/usr/share/man/man8/aa-log.md

docs:
	@ENABLED_GIT_REVISION_DATE=false MKDOCS_OFFLINE=true mkdocs build --strict

serve:
	@ENABLED_GIT_REVISION_DATE=false MKDOCS_OFFLINE=false mkdocs serve

clean:
	@rm -rf \
		debian/.debhelper debian/debhelper* debian/*.debhelper debian/${PKGNAME} \
		.pkg/${PKGNAME}* ${BUILD} coverage.out
