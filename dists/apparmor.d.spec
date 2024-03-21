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
BuildRequires:  golang-packaging

%description
AppArmor.d is a set of over 1500 AppArmor profiles whose aim is to confine most Linux based applications and processes.

%prep
%autosetup

%build
%make_build

%install
%make_install

%posttrans
rm -f /var/cache/apparmor/* 2>/dev/null
systemctl is-active -q apparmor && systemctl reload apparmor ||:

%files
%license LICENSE
%doc README.md
%exclude /etc/apparmor.d/libvirtd
%exclude /etc/apparmor.d/unix-chkpwd
%exclude /etc/apparmor.d/virt-aa-helper
%config /etc/apparmor.d/
/usr/bin/aa-log

%dir /usr/lib/systemd/system/dbus-broker.service.d
%dir /usr/lib/systemd/system/dbus.service.d
%dir /usr/lib/systemd/system/haveged.service.d
%dir /usr/lib/systemd/system/multipathd.service.d
%dir /usr/lib/systemd/system/pcscd.service.d
%dir /usr/lib/systemd/system/systemd-journald.service.d
%dir /usr/lib/systemd/system/systemd-networkd.service.d
%dir /usr/lib/systemd/system/systemd-timesyncd.service.d
%dir /usr/lib/systemd/user/at-spi-dbus-bus.service.d
%dir /usr/lib/systemd/user/dbus-broker.service.d
%dir /usr/lib/systemd/user/dbus.service.d
%dir /usr/lib/systemd/user/org.freedesktop.IBus.session.GNOME.service.d
%dir /usr/share/zsh
%dir /usr/share/zsh/site-functions

/usr/lib/systemd/system/dbus-broker.service.d/apparmor.conf
/usr/lib/systemd/system/dbus.service.d/apparmor.conf
/usr/lib/systemd/system/haveged.service.d/apparmor.conf
/usr/lib/systemd/system/multipathd.service.d/apparmor.conf
/usr/lib/systemd/system/pcscd.service.d/apparmor.conf
/usr/lib/systemd/system/systemd-journald.service.d/apparmor.conf
/usr/lib/systemd/system/systemd-networkd.service.d/apparmor.conf
/usr/lib/systemd/system/systemd-timesyncd.service.d/apparmor.conf
/usr/lib/systemd/user/at-spi-dbus-bus.service.d/apparmor.conf
/usr/lib/systemd/user/dbus-broker.service.d/apparmor.conf
/usr/lib/systemd/user/dbus.service.d/apparmor.conf
/usr/lib/systemd/user/org.freedesktop.IBus.session.GNOME.service.d/apparmor.conf
/usr/share/bash-completion/completions/aa-log
/usr/share/zsh/site-functions/_aa-log.zsh

%changelog
