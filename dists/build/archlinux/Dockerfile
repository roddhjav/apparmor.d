FROM archlinux:base-devel

RUN pacman -Syu --noconfirm --noprogressbar --quiet \
        devtools git pacman-contrib \
        go git rsync lsb-release && \
    paccache -r -k 0 && \
    pacman -Rscn --noconfirm --noprogressbar pacman-contrib && \
    useradd -m -s /bin/bash -u 1000 build && \
    echo "build ALL=NOPASSWD: ALL" >> /etc/sudoers && \
    chown -R build:build /home/build

USER build
CMD ["/bin/bash"]
