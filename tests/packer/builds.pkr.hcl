# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

build {
  sources = [
    "source.qemu.archlinux",
    "source.qemu.debian",
    "source.qemu.fedora",
    "source.qemu.opensuse",
    "source.qemu.ubuntu22",
    "source.qemu.ubuntu24",
  ]

  # Upload artifacts
  provisioner "file" {
    destination = "/tmp/"
    sources = [
      "${path.cwd}/packer/src/",
      "${path.cwd}/packer/init.sh",
      "${path.cwd}/packer/clean.sh",
      "${path.cwd}/../.pkg/",
    ]
  }

  # Full system provisioning
  provisioner "shell" {
    execute_command = "echo '${var.password}' | sudo -S sh -c '{{ .Vars }} {{ .Path }}'"
    inline = [
      # Wait for cloud-init to finish
      "while [ ! -f /var/lib/cloud/instance/boot-finished ]; do echo 'Waiting for Cloud-Init...'; sleep 20; done",

      # Ensure cloud-init is successful
      "cloud-init status",

      # Remove logs and artifacts so cloud-init can re-run
      "cloud-init clean",

      # Install local files and config
      "bash /tmp/init.sh",

      # Minimize the image
      "bash /tmp/clean.sh",
    ]
  }

  post-processor "vagrant" {
    output = "${var.base_dir}/packer_${var.prefix}${source.name}-${var.flavor}.box"
  }

  post-processor "shell-local" {
    inline = [
      "vagrant box add --force --name ${var.prefix}${source.name}-${var.flavor} ${var.base_dir}/packer_${var.prefix}${source.name}-${var.flavor}.box"
    ]
  }

}
