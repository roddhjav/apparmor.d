# Maintainer: Alexandre Pujol <alexandre@pujol.io>
# shellcheck disable=SC2034,SC2154,SC2164

# Warning: for development only, use https://aur.archlinux.org/packages/apparmor.d-git for production use.

pkgbase=apparmor.d
pkgname=(
  apparmor.d apparmor.d.base
  apparmor.d.other
)
pkgver=0.0001
pkgrel=1
pkgdesc="Full set of apparmor profiles (base)"
arch=("any")
url="https://github.com/roddhjav/apparmor.d"
license=('GPL2')
depends=('apparmor')
makedepends=('go' 'git' 'rsync')
conflicts=("$pkgbase-git" "$pkgbase")

pkgver() {
  cd "$srcdir/$pkgbase"
  echo "0.$(git rev-list --count HEAD)"
}

prepare() {
  rsync -a --delete "$startdir" "$srcdir"
}

build() {
  cd "$srcdir/$pkgbase"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw"
  make BUILD=.buid.all DISTRIBUTION=arch
  make packages DISTRIBUTION=arch
}

package_apparmor.d() {
  pkgdesc="Full set of apparmor profiles"
  arch=("$CARCH")
  conflicts=("${pkgname[@]:1}")
  cd "$srcdir/$pkgbase"
  make install BUILD=.buid.all DESTDIR="$pkgdir"
}

package_apparmor.d.base() {
  arch=("$CARCH")
  cd "$srcdir/$pkgbase"
  make install-base DESTDIR="$pkgdir"
}

package_apparmor.d.other() {
  pkgdesc="Full set of apparmor profiles (other)"
  depends=(apparmor.d.base)
  cd "$srcdir/$pkgbase"
  make profiles-other DESTDIR="$pkgdir"
}
