# Maintainer: Alexandre Pujol <alexandre@pujol.io>
# shellcheck disable=SC2034,SC2154,SC2164

pkgname=apparmor.d
pkgver=0.001
pkgrel=1
pkgdesc="Full set of apparmor profiles"
arch=("x86_64")
url="https://github.com/roddhjav/$pkgname"
license=('GPL2')
depends=('apparmor')
makedepends=('go' 'git')

pkgver() {
  cd "$srcdir/$pkgname"
  echo "0.$(git rev-list --count HEAD)"
}

prepare() {
  git clone "$startdir" "$srcdir/$pkgname"
  cd "$srcdir/$pkgname"

  ./configure
}

build() {
  cd "$srcdir/$pkgname/"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
  go build -o .build/ ./cmd/aa-log
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

  # Internal tool
  install -Dm755 .build/aa-log "$pkgdir"/usr/bin/aa-log
}
