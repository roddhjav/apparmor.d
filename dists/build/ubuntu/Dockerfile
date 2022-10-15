FROM ubuntu:22.04

ENV DEBIAN_FRONTEND=noninteractive \
    TERM=xterm

# hadolint ignore=DL3008
RUN apt-get update -y && apt-get -qq -y --no-install-recommends upgrade && \
    apt-get -qq -y --no-install-recommends install \
        build-essential devscripts debhelper fakeroot \
		git lsb-release rsync golang-go && \
    apt-get -qy autoremove && \
    apt-get -qq --purge remove -y .\*-doc$ && \
    apt-get clean && \
    rm -rf /usr/share/doc /usr/share/man /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    useradd -m -s /bin/bash -u 1000 build && \
    chown -R build:build /home/build

USER build
CMD ["/bin/bash"]
