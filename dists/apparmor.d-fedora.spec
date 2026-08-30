# apparmor.d - Full set of apparmor profiles
# Copyright (c) 2023 SUSE LLC
# Copyright (c) 2023 Christian Boltz
# Copyright (C) 2023-2026 Alexandre Pujol <alexandre@pujol.io>
# Copyright (c) 2026 CatPieLeaf <catpieleaf@proton.me>
# SPDX-License-Identifier: GPL-2.0-only

Name:           apparmor.d
Version:        0.4912.0
Release:        1%{?dist}
Summary:        Full set of AppArmor policies
License:        GPL-2.0-only
URL:            https://github.com/roddhjav/apparmor.d
Source0:        %{name}-%{version}.tar.gz

Requires:       apparmor-profiles
Requires:       apparmor-parser
Requires:       apparmor-utils
BuildRequires:  just
BuildRequires:  golang
BuildRequires:  systemd-rpm-macros

%description
AppArmor.d is a set of over 1500 AppArmor profiles whose aim is to confine most Linux based applications and processes.

%package base
Summary:        Full set of AppArmor policies (base abstractions, tunables, and booleans)
BuildArch:      noarch

%description base
apparmor.d-base is a set of base abstractions, tunables, and booleans defined by apparmor.d.

%package tools
Summary:        Full set of AppArmor policies (userland toolings)

%description tools
apparmor.d-tools is a set of userland toolings to help manage AppArmor profiles defined in apparmor.d.

%prep
%autosetup

%build
just prebuild

%install
just destdir="%{buildroot}" install-profiles
just destdir="%{buildroot}" install-base
just destdir="%{buildroot}" install-tools

%posttrans
apparmor_parser --purge-cache || :
%restart_on_update apparmor

%transfiletriggerin tools -- /usr /etc /opt
/usr/bin/aa-install --install || :

%files
%license LICENSE
%doc README.md
/usr/share/apparmor.d/
%config /etc/apparmor.d/disable/hostname

%dir /usr/lib/systemd/system/*.service.d
/usr/lib/systemd/system/*.service.d/apparmor.conf
%dir /usr/lib/systemd/user/*.service.d
/usr/lib/systemd/user/*.service.d/apparmor.conf

%files base
%config /etc/apparmor.d/abstractions
%config /etc/apparmor.d/tunables

%files tools
/usr/bin/aa-*

%dir /usr/share/apparmor
/usr/share/apparmor/modes
/usr/share/apparmor/flags.d
/usr/share/apparmor/ignore.d
/usr/share/apparmor/overwrite.d/

%dir /usr/share/zsh
%dir /usr/share/zsh/site-functions
/usr/share/zsh/site-functions/_aa-*.zsh
/usr/share/bash-completion/completions/aa-*
%doc %{_mandir}/man1/aa-*.1.gz
%doc %{_mandir}/man8/aa-*.8.gz

%changelog
