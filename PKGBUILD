# Maintainer: Alexandre Pujol <alexandre@pujol.io>
# shellcheck disable=SC2034,SC2154,SC2164

# Warning: for development only, use https://aur.archlinux.org/packages/apparmor.d-git for production use.

pkgbase=apparmor.d
pkgname=(
  apparmor.d
  apparmor.d.enforced
  # apparmor.d.fsp apparmor.d.fsp.enforced
  # apparmor.d.server apparmor.d.server.enforced
  # apparmor.d.server.fsp apparmor.d.server.fsp.enforced
)
pkgver=0.4900
pkgrel=1
pkgdesc="Full set of apparmor profiles"
arch=('x86_64' 'armv6h' 'armv7h' 'aarch64')
url="https://github.com/roddhjav/apparmor.d"
license=('GPL-2.0-only')
depends=('apparmor>=4.1.0' 'apparmor<5.0.0')
makedepends=('go' 'git' 'rsync' 'just')

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
  export GOPATH="${srcdir}"
  export GOFLAGS="-buildmode=pie -trimpath -ldflags=-linkmode=external -mod=readonly -modcacherw -tags=dev"
  export DISTRIBUTION=arch
  local -A modes=(
    # Mapping of modes to just build target.
    [default]=complain
    [enforced]=enforce
    # [fsp]=fsp-complain
    # [fsp.enforced]=fsp
    # [server]=server-complain
    # [server.enforced]=server
    # [server.fsp]=server-fsp-complain
    # [server.fsp.enforced]=server-fsp
  )
  for mode in "${!modes[@]}"; do
    just build=".build/$mode" "${modes[$mode]}"
  done
}

_conflicts() {
  local mode="$1"
  local pattern=".$mode"
  if [[ "$mode" == "default" ]]; then
    pattern=""
  else
    echo "$pkgbase"
  fi
  for pkg in "${pkgname[@]}"; do
    if [[ "$pkg" == "${pkgbase}${pattern}" ]]; then
      continue
    fi
    echo "$pkg"
  done
}

_install() {
  local mode="${1:?}"
  cd "$srcdir/$pkgbase"
  just build=".build/$mode" destdir="$pkgdir" install
}

package_apparmor.d() {
  mode=default
  pkgdesc="$pkgdesc (complain mode)"
  mapfile -t conflicts < <(_conflicts $mode)
  _install $mode
}

package_apparmor.d.enforced() {
  mode=enforced
  pkgdesc="$pkgdesc (enforced mode)"
  mapfile -t conflicts < <(_conflicts $mode)
  _install $mode
}

package_apparmor.d.fsp() {
  mode="fsp"
  pkgdesc="$pkgdesc (FSP mode)"
  mapfile -t conflicts < <(_conflicts $mode)
  _install $mode
}

package_apparmor.d.fsp.enforced() {
  mode="fsp.enforced"
  pkgdesc="$pkgdesc (FSP enforced mode)"
  mapfile -t conflicts < <(_conflicts $mode)
  _install $mode
}

package_apparmor.d.server() {
  mode="server"
  pkgdesc="$pkgdesc (server complain mode)"
  mapfile -t conflicts < <(_conflicts $mode)
  _install $mode
}

package_apparmor.d.server.enforced() {
  mode="server.enforced"
  pkgdesc="$pkgdesc (server enforced mode)"
  mapfile -t conflicts < <(_conflicts $mode)
  _install $mode
}

package_apparmor.d.server.fsp() {
  mode="server.fsp"
  pkgdesc="$pkgdesc (server FSP complain mode)"
  mapfile -t conflicts < <(_conflicts $mode)
  _install $mode
}

package_apparmor.d.server.fsp.enforced() {
  mode="server.fsp.enforced"
  pkgdesc="$pkgdesc (server FSP enforced mode)"
  mapfile -t conflicts < <(_conflicts $mode)
  _install $mode
}
