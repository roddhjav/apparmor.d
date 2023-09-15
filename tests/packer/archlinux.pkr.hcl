# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

source "qemu" "archlinux-gnome" {
  disk_image         = true
  iso_url            = "https://geo.mirror.pkgbuild.com/images/latest/Arch-Linux-x86_64-cloudimg.qcow2"
  iso_checksum       = "file:https://geo.mirror.pkgbuild.com/images/latest/Arch-Linux-x86_64-cloudimg.qcow2.SHA256"
  iso_target_path    = "${var.iso_dir}/archlinux-cloudimg-amd64.img"
  cpus               = 6
  memory             = 4096
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
  vm_name            = "${var.prefix}${source.name}.qcow2"
  boot_wait          = "10s"
  shutdown_command   = "echo ${var.password} | sudo -S shutdown -hP now"
  cd_label           = "cidata"
  cd_content = {
    "meta-data" = ""
    "user-data" = templatefile("${path.cwd}/packer/init/${source.name}.user-data.yml",
      {
        username = "${var.username}"
        password = "${var.password}"
        ssh_key  = file("${var.ssh_publickey}")
        hostname = "${var.prefix}${source.name}"
      }
    )
  }
}

source "qemu" "archlinux-kde" {
  disk_image         = true
  iso_url            = "https://geo.mirror.pkgbuild.com/images/latest/Arch-Linux-x86_64-cloudimg.qcow2"
  iso_checksum       = "file:https://geo.mirror.pkgbuild.com/images/latest/Arch-Linux-x86_64-cloudimg.qcow2.SHA256"
  iso_target_path    = "${var.iso_dir}/archlinux-cloudimg-amd64.img"
  cpus               = 6
  memory             = 4096
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
  vm_name            = "${var.prefix}${source.name}.qcow2"
  boot_wait          = "10s"
  shutdown_command   = "echo ${var.password} | sudo -S shutdown -hP now"
  cd_label           = "cidata"
  cd_content = {
    "meta-data" = ""
    "user-data" = templatefile("${path.cwd}/packer/init/${source.name}.user-data.yml",
      {
        username = "${var.username}"
        password = "${var.password}"
        ssh_key  = file("${var.ssh_publickey}")
        hostname = "${var.prefix}${source.name}"
      }
    )
  }
}
