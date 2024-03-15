# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

source "qemu" "ubuntu-server" {
  disk_image         = true
  iso_url            = "https://cloud-images.ubuntu.com/${var.release.ubuntu.codename}/current/${var.release.ubuntu.codename}-server-cloudimg-amd64.img"
  iso_checksum       = "file:https://cloud-images.ubuntu.com/${var.release.ubuntu.codename}/current/SHA256SUMS"
  iso_target_path    = "${var.iso_dir}/ubuntu-cloudimg-amd64.img"
  cpu_model          = "host"
  cpus               = 4
  memory             = 2048
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
  output_directory   = "${var.output}/"
  vm_name            = "${var.prefix}${source.name}.qcow2"
  boot_wait          = "10s"
  firmware           = var.firmware
  shutdown_command   = "echo ${var.password} | sudo -S /sbin/shutdown -hP now"
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

source "qemu" "ubuntu-server24" {
  disk_image         = true
  iso_url            = "https://cloud-images.ubuntu.com/${var.release.ubuntu24.codename}/current/${var.release.ubuntu24.codename}-server-cloudimg-amd64.img"
  iso_checksum       = "file:https://cloud-images.ubuntu.com/${var.release.ubuntu24.codename}/current/SHA256SUMS"
  iso_target_path    = "${var.iso_dir}/ubuntu-${var.release.ubuntu24.codename}-cloudimg-amd64.img"
  cpu_model          = "host"
  cpus               = 4
  memory             = 2048
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
  output_directory   = "${var.output}/"
  vm_name            = "${var.prefix}${source.name}.qcow2"
  boot_wait          = "10s"
  firmware           = var.firmware
  shutdown_command   = "echo ${var.password} | sudo -S /sbin/shutdown -hP now"
  cd_label           = "cidata"
  cd_content = {
    "meta-data" = ""
    "user-data" = templatefile("${path.cwd}/packer/init/ubuntu-server.user-data.yml",
      {
        username = "${var.username}"
        password = "${var.password}"
        ssh_key  = file("${var.ssh_publickey}")
        hostname = "${var.prefix}${source.name}"
      }
    )
  }
}

source "qemu" "ubuntu-desktop" {
  disk_image         = true
  iso_url            = "https://cloud-images.ubuntu.com/${var.release.ubuntu.codename}/current/${var.release.ubuntu.codename}-server-cloudimg-amd64.img"
  iso_checksum       = "file:https://cloud-images.ubuntu.com/${var.release.ubuntu.codename}/current/SHA256SUMS"
  iso_target_path    = "${var.iso_dir}/ubuntu-cloudimg-amd64.img"
  cpu_model          = "host"
  cpus               = 6
  memory             = 4096
  disk_size          = var.disk_size
  accelerator        = "kvm"
  headless           = true
  ssh_username       = var.username
  ssh_password       = var.password
  ssh_port           = 22
  ssh_wait_timeout   = "10000s"
  disk_compression   = true
  disk_detect_zeroes = "unmap"
  disk_discard       = "unmap"
  output_directory   = "${var.output}/"
  vm_name            = "${var.prefix}${source.name}.qcow2"
  boot_wait          = "10s"
  firmware           = var.firmware
  shutdown_command   = "echo ${var.password} | sudo -S /sbin/shutdown -hP now"
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

source "qemu" "ubuntu-desktop24" {
  disk_image         = true
  iso_url            = "https://cloud-images.ubuntu.com/${var.release.ubuntu24.codename}/current/${var.release.ubuntu24.codename}-server-cloudimg-amd64.img"
  iso_checksum       = "file:https://cloud-images.ubuntu.com/${var.release.ubuntu24.codename}/current/SHA256SUMS"
  iso_target_path    = "${var.iso_dir}/ubuntu-${var.release.ubuntu24.codename}-cloudimg-amd64.img"
  cpu_model          = "host"
  cpus               = 6
  memory             = 4096
  disk_size          = var.disk_size
  accelerator        = "kvm"
  headless           = false
  ssh_username       = var.username
  ssh_password       = var.password
  ssh_port           = 22
  ssh_wait_timeout   = "10000s"
  disk_compression   = true
  disk_detect_zeroes = "unmap"
  disk_discard       = "unmap"
  output_directory   = "${var.output}/"
  vm_name            = "${var.prefix}${source.name}.qcow2"
  boot_wait          = "10s"
  firmware           = var.firmware
  shutdown_command   = "echo ${var.password} | sudo -S /sbin/shutdown -hP now"
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
