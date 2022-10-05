#!/usr/bin/make -f

PKGNAME := apparmor.d

.PHONY: install lint archlinux debian ubuntu whonix clean

all:
	@echo "Nothing to do."

install:
	@echo "Nothing to do."

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
		${PKGNAME}_*.* .build
