FROM debian:11

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

# hadolint ignore=DL3008
RUN echo 'deb http://deb.debian.org/debian bullseye-backports main contrib non-free' >> /etc/apt/sources.list && \
    apt-get update -y && apt-get -qq -y --no-install-recommends upgrade && \
    apt-get -qq -y --no-install-recommends install \
        build-essential devscripts debhelper fakeroot config-package-dev \
		git lsb-release rsync && \
    apt-get -qq -y --no-install-recommends install \
        golang-1.19-go -t bullseye-backports && \
    apt-get -qy autoremove && \
    apt-get -qq --purge remove -y .\*-doc$ && \
    apt-get clean && \
    rm -rf /usr/share/doc /usr/share/man /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    ln -s /usr/lib/go-1.19/bin/go /usr/bin/go && \
    useradd -m -s /bin/bash -u 1000 build && \
    chown -R build:build /home/build

USER build
CMD ["/bin/bash"]
