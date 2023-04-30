# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

build {
  name = "main"
  sources = [
    "source.qemu.archlinux-gnome",
    "source.qemu.archlinux-kde",
    "source.qemu.debian-server",
    "source.qemu.opensuse-kde",
    "source.qemu.ubuntu-server",
  ]

  provisioner "file" {
    destination = "/tmp"
    sources     = ["${path.cwd}/packer/src"]
  }

  provisioner "file" {
    only        = ["qemu.archlinux-gnome", "qemu.archlinux-kde"]
    destination = "/tmp/src/"
    sources     = ["${path.cwd}/../apparmor.d-${var.version}-x86_64.pkg.tar.zst"]
  }

  provisioner "file" {
    only        = ["qemu.debian-server", "qemu.ubuntu-server", "qemu.ubuntu-desktop"]
    destination = "/tmp/src/"
    sources     = ["${path.cwd}/../apparmor.d_${var.version}_all.deb"]
  }

  provisioner "shell" {
    execute_command = "echo '${var.password}' | sudo -S sh -c '{{ .Vars }} {{ .Path }}'"
    inline = [
      "while [ ! -f /var/lib/cloud/instance/boot-finished ]; do echo 'Waiting for Cloud-Init...'; sleep 20; done",
      "cloud-init clean", # Remove logs and artifacts so cloud-init can re-run
    ]
  }

  provisioner "shell" {
    script          = "${path.cwd}/packer/init/init.sh"
    execute_command = "echo '${var.password}' | sudo -S sh -c '{{ .Vars }} {{ .Path }}'"
  }

  provisioner "shell" {
    script          = "${path.cwd}/packer/init/clean.sh"
    execute_command = "echo '${var.password}' | sudo -S sh -c '{{ .Vars }} {{ .Path }}'"
  }

  post-processor "vagrant" {
    output = "${var.base_dir}/packer_${var.prefix}${source.name}.box"
  }

  post-processor "shell-local" {
    inline = [
      "vagrant box add --force --name ${var.prefix}${source.name} ${var.base_dir}/packer_${var.prefix}${source.name}.box"
    ]
  }

}