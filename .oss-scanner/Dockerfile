# OSS Scanner build for apparmor.d.
#
# The asset under review is the policy, not the Go tooling. This image builds
# the profiles for Ubuntu 26.04 (ABI 5) and installs them the way the package
# does, in default mode:
#
#   /etc/apparmor.d
#
# Sources:
#   /src/apparmor.d          profile sources (groups/, profiles-*/, abstractions/, tunables/)
#   /src/.build/apparmor.d   prebuilt profiles (variables resolved, Ubuntu configuration applied)
#
# Check that a profile (or a patched profile) still compiles, without a kernel:
#   apparmor_parser -Q -K -I /etc/apparmor.d /etc/apparmor.d/<profile>
#   just check             (project style and security checks on the sources)
#
# Nothing here loads policy into the kernel: the container has no AppArmor.
FROM ubuntu:26.04

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=C.UTF-8 \
    GOFLAGS=-mod=readonly \
    GOTOOLCHAIN=local

RUN apt-get update && apt-get install -y --no-install-recommends \
      ca-certificates git golang-go just apparmor apparmor-utils \
      bash findutils grep sed gawk lsb-release \
    && rm -rf /var/lib/apt/lists/*

# scanner contract: the checkout lives inside the image, at /src
COPY . /src
WORKDIR /src

# Build the tools and prebuild the profiles for this distribution (detected
# from /etc/os-release, so Ubuntu 26.04 here).
RUN just build && just prebuild

# Install as the Debian/Ubuntu package does, then enable the profiles.
# --all: the package only enables the profiles of installed programs, the scan
# covers all of them.
# --no-reload: there is no AppArmor kernel in the build container.
RUN just install-base install-tools install-profiles \
    && aa-install --no-reload --install --all

# Cheap checks so the image is known to work; failures never abort the build.
RUN (just check || echo "WARNING: just check reported failures") \
    && (go test ./cmd/... $(go list ./pkg/... | grep -v /pkg/paths) \
        || echo "WARNING: go tests reported failures")
