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
  echo "$pkgver + 0.01" | bc
}

prepare() {
  git clone "$startdir" "$srcdir/$pkgname"
  cd "$srcdir/$pkgname"

  ./configure --distribution=archlinux
}

package() {
  local _build='.build/apparmor.d'
  cd "$srcdir/$pkgname"

  # Install all files from root/
  mapfile -t root < <(find root -type f -printf "%P\n")
  for file in "${root[@]}"; do
    install -Dm0644 "root/$file" "$pkgdir/$file"
  done

  # Install all files from $_build
  mapfile -t build < <(find "$_build/" -type f -printf "%P\n")
  for file in "${build[@]}"; do
    install -Dm0644 "$_build/$file" "$pkgdir/etc/apparmor.d/$file"
  done

  # Ensure some systemd services do not start before apparmor rules are loaded
  for path in systemd/*; do
    service=$(basename "$path")
    install -Dm0644 "$path" \
      "$pkgdir/usr/lib/systemd/system/$service.d/apparmor.conf"
  done

  # Set special access rights
  chmod 0755 "$pkgdir"/usr/bin/*
}
