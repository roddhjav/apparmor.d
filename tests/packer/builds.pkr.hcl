# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

build {
  sources = [
    "source.qemu.archlinux-gnome",
    "source.qemu.archlinux-kde",
    "source.qemu.debian-server",
    "source.qemu.debian-gnome",
    "source.qemu.opensuse-kde",
    "source.qemu.ubuntu-desktop",
    "source.qemu.ubuntu-server",
  ]

  # Upload local files
  provisioner "file" {
    destination = "/tmp"
    sources     = ["${path.cwd}/packer/src"]
  }

  provisioner "file" {
    only        = ["qemu.archlinux-gnome", "qemu.archlinux-kde"]
    destination = "/tmp/src/"
    sources     = ["${path.cwd}/../apparmor.d-${var.version}-1-x86_64.pkg.tar.zst"]
  }

  provisioner "file" {
    only        = ["qemu.opensuse-*"]
    destination = "/tmp/src/"
    sources     = ["${path.cwd}/../apparmor.d-${var.version}-1.x86_64.rpm"]
  }

  provisioner "file" {
    only        = ["qemu.debian-server", "qemu.debian-gnome", "qemu.ubuntu-server", "qemu.ubuntu-desktop"]
    destination = "/tmp/src/"
    sources     = ["${path.cwd}/../apparmor.d_${var.version}-1_amd64.deb"]
  }

  # Wait for cloud-init to finish
  provisioner "shell" {
    except          = ["qemu.opensuse-*"]
    execute_command = "echo '${var.password}' | sudo -S sh -c '{{ .Vars }} {{ .Path }}'"
    inline = [
      "while [ ! -f /var/lib/cloud/instance/boot-finished ]; do echo 'Waiting for Cloud-Init...'; sleep 20; done",
      "cloud-init clean", # Remove logs and artifacts so cloud-init can re-run
    ]
  }

  # Install local files and config
  provisioner "shell" {
    script          = "${path.cwd}/packer/init/init.sh"
    execute_command = "echo '${var.password}' | sudo -S sh -c '{{ .Vars }} {{ .Path }}'"
  }

  # Minimize the image
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
