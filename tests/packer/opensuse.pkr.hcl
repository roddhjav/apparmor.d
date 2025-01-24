# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# TODO: Fully automate the creation of the base image

source "qemu" "opensuse" {
  disk_image         = true
  iso_url            = "${var.base_dir}/base-tumbleweed-gnome.qcow2"
  iso_checksum       = "sha256:223ed62160ef4f1a4f21b69c574f552a07eee6ef66cf66eef2b49c5a7c4864f4"
  iso_target_path    = "${var.base_dir}/base-tumbleweed-gnome.qcow2"
  cpu_model          = "host"
  cpus               = 6
  memory             = 4096
  disk_size          = var.disk_size
  accelerator        = "kvm"
  headless           = false
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
  shutdown_command   = "echo ${var.password} | sudo shutdown -hP now"
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
