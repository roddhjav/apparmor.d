#!/usr/bin/env bash
# Build the package for Archlinux/openSUSE/Debian/Ubuntu
# Copyright (C) 2022-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Usage: just [ dpkg | pkg | rpm ]

set -eu -o pipefail

readonly COMMAND="$1"
readonly OUTPUT="$PWD/.pkg"
readonly PKGNAME=apparmor.d
VERSION="0.$(git rev-list --count HEAD)"
readonly VERSION

main() {
    case "$COMMAND" in
    pkg)
        PKGDEST="$OUTPUT" BUILDDIR=/tmp/makepkg makepkg --syncdeps --force --cleanbuild --noconfirm --noprogressbar
        ;;

    dpkg)
        dch --newversion="$VERSION-1" --urgency=medium --distribution="$(lsb_release -sc)" --controlmaint "Release $VERSION-1"
        dpkg-buildpackage -b -d --no-sign
        lintian || true
        mv ../"${PKGNAME}_${VERSION}-1"_*.deb "$OUTPUT"
        ;;

    rpm)
        RPMBUILD_ROOT=$(mktemp -d /tmp/$PKGNAME.XXXXXX)
        ARCH=$(uname -m)
        readonly RPMBUILD_ROOT ARCH

        mkdir -p "$RPMBUILD_ROOT"/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS/tmp}
        cp -p "dists/$PKGNAME.spec" "$RPMBUILD_ROOT/SPECS"
        tar -czf "$RPMBUILD_ROOT/SOURCES/$PKGNAME-$VERSION.tar.gz" --transform "s,^,$PKGNAME-$VERSION/," ./*

        cd "$RPMBUILD_ROOT"
        sed -i "s/^Version:.*/Version:        $VERSION/" "SPECS/$PKGNAME.spec"
        rpmbuild -bb --define "_topdir $RPMBUILD_ROOT" "SPECS/$PKGNAME.spec"

        mv "$RPMBUILD_ROOT/RPMS/$ARCH/"*.rpm "$OUTPUT"
        rm -rf "$RPMBUILD_ROOT"
        ;;

    *) ;;
    esac
}

mkdir -p "$OUTPUT"
main "$@"
