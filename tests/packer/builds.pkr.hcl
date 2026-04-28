# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

locals {
  name   = "${var.prefix}${var.dist}${var.release}-${var.flavor}"
  osinfo = "${var.dist}${var.release}"
  group  = contains(["debian", "ubuntu"], var.dist) ? "sudo" : "wheel"
}

source "qemu" "default" {
  disk_image         = true
  iso_url            = var.DM[local.osinfo].img_url
  iso_checksum       = "file:${var.DM[local.osinfo].img_checksum}"
  iso_target_path    = pathexpand("${var.iso_dir}/${basename("${var.DM[local.osinfo].img_url}")}")
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
  output_directory   = pathexpand(var.output_dir)
  vm_name            = "${local.name}.qcow2"
  boot_wait          = "10s"
  firmware           = pathexpand(var.firmware)
  shutdown_command   = "echo ${var.password} | sudo -S /sbin/shutdown -hP now"
  cd_label           = "cidata"
  cd_content = {
    "meta-data" = ""
    "user-data" = format("%s\n%s\n%s",
      templatefile("${path.cwd}/tests/cloud-init/common.yml",
        {
          username = var.username
          password = var.password
          ssh_key  = file(var.ssh_publickey)
          hostname = regex_replace(local.name, "\\.", "")
          group    = local.group
        }
      ),
      file("${path.cwd}/tests/cloud-init/${regex_replace(local.osinfo, "[0-9.]*$", "")}.yml"),
      file("${path.cwd}/tests/cloud-init/${local.osinfo}-${var.flavor}.user-data.yml")
    )
  }
}

build {
  sources = [
    "source.qemu.default",
  ]

  # Upload artifacts
  provisioner "file" {
    destination = "/tmp/"
    sources = [
      "${path.cwd}/tests/packer/src/",
      "${path.cwd}/tests/packer/init.sh",
      "${path.cwd}/tests/packer/clean.sh",
      "${path.cwd}/.pkg/",
    ]
  }

  # Full system provisioning
  provisioner "shell" {
    execute_command = "echo '${var.password}' | sudo -S sh -c '{{ .Vars }} {{ .Path }}'"
    inline = [
      # Wait for cloud-init to finish
      "while [ ! -f /var/lib/cloud/instance/boot-finished ]; do echo 'Waiting for Cloud-Init...'; sleep 20; done",

      # Ensure cloud-init is successful
      "cloud-init status || cloud-init collect-logs --tarfile /root/cloud-init.tar.gz",

      # Remove logs and artifacts so cloud-init can re-run
      "cloud-init clean || true",

      # Install local files and config
      "bash /tmp/init.sh",

      # Minimize the image
      "bash /tmp/clean.sh",
    ]
  }

  post-processor "shell-local" {
    inline = [
      "mv ${var.output_dir}/${local.name}.qcow2 ${var.base_dir}/${local.name}.qcow2",
    ]
  }

}
