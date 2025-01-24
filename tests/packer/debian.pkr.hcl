# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

source "qemu" "debian" {
  disk_image         = true
  iso_url            = "https://cdimage.debian.org/images/cloud/${var.release.debian.codename}/latest/debian-${var.release.debian.version}-genericcloud-amd64.qcow2"
  iso_checksum       = "file:https://cdimage.debian.org/images/cloud/${var.release.debian.codename}/latest/SHA512SUMS"
  iso_target_path    = "${var.iso_dir}/debian-${var.release.debian.codename}-cloudimg-amd64.img"
  cpu_model          = "host"
  cpus               = var.cpus
  memory             = var.ram
  disk_size          = var.disk_size
  accelerator        = "kvm"
  headless           = true
  ssh_username       = var.username
  ssh_password       = var.password
  ssh_port           = 22
  ssh_wait_timeout   = "1000s"
  disk_compression   = true
  disk_detect_zeroes = "unmap"
  disk_discard       = "unmap"
  output_directory   = var.output
  vm_name            = "${var.prefix}${source.name}-${var.flavor}.qcow2"
  boot_wait          = "10s"
  firmware           = var.firmware
  shutdown_command   = "echo ${var.password} | sudo -S /sbin/shutdown -hP now"
  cd_label           = "cidata"
  cd_content = {
    "meta-data" = ""
    "user-data" = templatefile("${path.cwd}/cloud-init/${source.name}-${var.flavor}.user-data.yml",
      {
        username = "${var.username}"
        password = "${var.password}"
        ssh_key  = file("${var.ssh_publickey}")
        hostname = "${var.prefix}${source.name}"
      }
    )
  }
}
