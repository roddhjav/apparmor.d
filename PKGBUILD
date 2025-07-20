# Maintainer: Alexandre Pujol <alexandre@pujol.io>
# shellcheck disable=SC2034,SC2154,SC2164

# Warning: for development only, use https://aur.archlinux.org/packages/apparmor.d-git for production use.

pkgname=apparmor.d
pkgver=0.001
pkgrel=1
pkgdesc="Full set of apparmor profiles"
arch=('x86_64' 'armv6h' 'armv7h' 'aarch64')
url="https://github.com/roddhjav/apparmor.d"
license=('GPL-2.0-only')
depends=('apparmor>=4.1.0' 'apparmor<5.0.0')
makedepends=('go' 'git' 'rsync' 'just')
conflicts=("$pkgname-git")

pkgver() {
  cd "$srcdir/$pkgname"
  echo "0.$(git rev-list --count HEAD)"
}

prepare() {
  rsync -a --delete "$startdir" "$srcdir"
}

build() {
  cd "$srcdir/$pkgname"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
  export DISTRIBUTION=arch
  just complain
}

package() {
  cd "$srcdir/$pkgname"
  just destdir="$pkgdir" install
}
