# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2026 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

abi <abi/5.0>,

include <tunables/global>

profile flatpak.trigger flags=(attach_disconnected) {
  include <abstractions/base>
  include <abstractions/consoles>
  include <abstractions/nameservice-strict>

  @{bin}/bsdtar  ix,
  @{bin}/cp      ix,
  @{bin}/install ix,
  @{bin}/sed     ix,

  @{bin}/gtk{,4}-update-icon-cache  Px -> flatpak.trigger//&gtk-update-icon-cache,
  @{bin}/update-desktop-database    ix,
  @{bin}/update-mime-database       Px -> flatpak.trigger//&update-mime-database,

  /usr/share/flatpak/triggers/desktop-database.trigger ix,
  /usr/share/flatpak/triggers/gtk-icon-cache.trigger   ix,
  /usr/share/flatpak/triggers/mime-database.trigger    ix,

  @{system_share_dirs}/** r,
  @{system_share_dirs}/*ubuntu/applications/.mimeinfo.cache.* rw,
  @{system_share_dirs}/*ubuntu/applications/mimeinfo.cache w,
  @{system_share_dirs}/applications/.mimeinfo.cache.* w,
  @{system_share_dirs}/applications/mimeinfo.cache w,
  @{system_share_dirs}/icons/**/.icon-theme.cache rw,
  @{system_share_dirs}/icons/**/icon-theme.cache w,
  @{system_share_dirs}/icons/hicolor/index.theme w,
  @{system_share_dirs}/mime/{,**} w,

  @{user_share_dirs}/** r,
  @{user_share_dirs}/.mimeinfo.cache.* w,
  @{user_share_dirs}/**/.icon-theme.cache w,
  @{user_share_dirs}/**/icon-theme.cache w,
  @{user_share_dirs}/applications/.mimeinfo.cache.* w,
  @{user_share_dirs}/applications/mimeinfo.cache w,
  @{user_share_dirs}/mime/{,**} w,
  @{user_share_dirs}/mimeinfo.cache w,

  include if exists <local/flatpak.trigger>
}

# vim:syntax=apparmor
