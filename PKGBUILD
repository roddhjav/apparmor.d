# Maintainer: Alexandre Pujol <alexandre@pujol.io>
# shellcheck disable=SC2034,SC2154,SC2164

pkgname=apparmor.d
pkgver=21
pkgrel=1
pkgdesc="Full set of apparmor profiles"
arch=("any")
url="https://gitlab.com/archlex/hardening/$pkgname"
license=('GPL2')
depends=('apparmor')
makedepends=('bc')

pkgver() {
  echo "$pkgver + .1" | bc
}

prepare() {
  git clone "$startdir" "$srcdir/$pkgname"
  cd "$srcdir/$pkgname"

  ./configure --distribution=archlinux
}

package() {
  local _root='_build'
  cd "$srcdir/$pkgname"

  # Install all files from root/
  cp --recursive --preserve=mode,ownership,timestamps "$_root/root/"* "$pkgdir/"

  # Install all files from apparmor.d/
  install -d "$pkgdir"/etc/apparmor.d/
  cp --recursive --preserve=mode,ownership,timestamps \
    $_root/apparmor.d/* "$pkgdir"/etc/apparmor.d/

  # Ensure some systemd services do not start before apparmor rules are loaded
  for path in systemd/*; do
    service=$(basename "$path")
    install -Dm0644 "$path" \
      "$pkgdir/usr/lib/systemd/system/$service.d/apparmor.conf"
  done
}
