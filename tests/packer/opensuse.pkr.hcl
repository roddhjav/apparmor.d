# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# TODO: Fully automate the creation of the base image
# To save some dev time, 'base_opensuse_kde' is manually created from the opensuse iso with:
# - KDE
# - username/password defined in the variables
# - cloud-init installed and enabled

source "qemu" "opensuse-kde" {
  disk_image         = true
  iso_url            = "${var.iso_dir}/base_opensuse_kde.qcow2"
  iso_checksum       = "sha256:62a174725bdf26981d15969e53461b89359f7763450cbfd3e258d4035731279b"
  iso_target_path    = "${var.iso_dir}/base_opensuse_kde.qcow2"
  cpus               = 6
  memory             = 4096
  disk_size          = "${var.disk_size}"
  accelerator        = "kvm"
  headless           = false
  ssh_username       = "${var.username}"
  ssh_password       = "${var.password}"
  ssh_port           = 22
  ssh_wait_timeout   = "1000s"
  disk_compression   = true
  disk_detect_zeroes = "unmap"
  disk_discard       = "unmap"
  output_directory   = "${var.iso_dir}/packer/"
  vm_name            = "${var.prefix}${source.name}.qcow2"
  boot_wait          = "10s"
  firmware           = "${var.firmware}"
  shutdown_command   = "echo ${var.password} | sudo shutdown -hP now"
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
