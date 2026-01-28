# Maintainer: Alexandre Pujol <alexandre@pujol.io>
# shellcheck disable=SC2034,SC2154,SC2164

# Warning: for development only, use https://aur.archlinux.org/packages/apparmor.d-git for production use.

pkgbase=apparmor.d
pkgname=(
  apparmor.d
  # apparmor.d-base
  # apparmor.d-tools
)
pkgver=0.4902
pkgrel=1
pkgdesc="Full set of apparmor profiles"
arch=('x86_64' 'armv6h' 'armv7h' 'aarch64')
url="https://github.com/roddhjav/apparmor.d"
license=('GPL-2.0-only')
depends=('apparmor>=4.1.3')
makedepends=('go' 'rsync' 'just')

prepare() {
  rsync -a --delete "$startdir" "$srcdir"
}

build() {
  cd "$srcdir/$pkgbase"
  export CGO_CPPFLAGS="${CPPFLAGS}"
  export CGO_CFLAGS="${CFLAGS}"
  export CGO_CXXFLAGS="${CXXFLAGS}"
  export CGO_LDFLAGS="${LDFLAGS}"
  export GOPATH="${srcdir}"
  export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw -tags=dev"
  export DISTRIBUTION=arch
  just complain
}

package_apparmor.d() {
  cd "$srcdir/$pkgbase"
  just destdir="$pkgdir" install
}
