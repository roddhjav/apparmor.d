# apparmor.d - Full set of apparmor profiles
# Copyright (c) 2023 SUSE LLC
# Copyright (c) 2023 Christian Boltz
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Warning: for development only, use https://build.opensuse.org/package/show/home:cboltz/apparmor.d for production use.

Name:           apparmor.d
Version:        0.0001
Release:        1%{?dist}
Summary:        Set of over 1500 AppArmor profiles
License:        GPL-2.0-only
URL:            https://github.com/roddhjav/apparmor.d
Source0:        %{name}-%{version}.tar.gz
Requires:       apparmor-profiles
BuildRequires:  distribution-release
BuildRequires:  just
BuildRequires:  golang-packaging
BuildRequires:  apparmor-profiles

%description
AppArmor.d is a set of over 1500 AppArmor profiles whose aim is to confine most Linux based applications and processes.

%prep
%autosetup

%build
just complain

%install
just destdir="%{buildroot}" install

%posttrans
apparmor_parser --purge-cache
%restart_on_update apparmor

%files
%license LICENSE
%doc README.md
%config /etc/apparmor.d/
/usr/bin/aa-log

%dir /usr/lib/systemd/system/*.service.d
/usr/lib/systemd/system/*.service.d/apparmor.conf
%dir /usr/lib/systemd/user/*.service.d
/usr/lib/systemd/user/*.service.d/apparmor.conf

/usr/share/bash-completion/completions/aa-log

%dir /usr/share/zsh
%dir /usr/share/zsh/site-functions
/usr/share/zsh/site-functions/_aa-log.zsh

%doc %{_mandir}/man8/aa-log.8.gz

%changelog
