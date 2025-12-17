# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023-2024 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# Variables definitions

variable "username" {
  description = "Admin username"
  type        = string
  default     = "user"
}

variable "password" {
  description = "Default admin password"
  type        = string
  default     = "user"
}

variable "cpus" {
  description = "Default CPU of the VM"
  type        = string
  default     = "6"
}

variable "ram" {
  description = "Default RAM of the VM"
  type        = string
  default     = "2048"
}

variable "disk_size" {
  description = "Disk size of the VM to build"
  type        = string
  default     = "40G"
}

variable "ssh_publickey" {
  description = "Path to the ssh public key"
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}

variable "iso_dir" {
  description = "Original ISO file directory"
  type        = string
  default     = "~/.libvirt/iso"
}

variable "base_dir" {
  description = "Final packer image output directory"
  type        = string
  default     = "~/.libvirt/base"
}

variable "output_dir" {
  description = "Output build directory"
  type        = string
  default     = "~/.libvirt/base/packer"
}

variable "firmware" {
  description = "Path to the UEFI firmware"
  type        = string
  default     = ""
}

variable "prefix" {
  description = "Image name prefix"
  type        = string
  default     = "aa-"
}

variable "dist" {
  description = "Distribution to target"
  type        = string
  default     = ""
}

variable "release" {
  description = "Release to target"
  type        = string
  default     = ""
}

variable "flavor" {
  description = "Distribution flavor to use (server, desktop, gnome, kde...)"
  type        = string
  default     = ""
}

variable "DM" {
  description = "Distribution Metadata to use"
  type = map(object({
    img_url      = string
    img_checksum = string
  }))
  default = {
    "archlinux" : {
      img_url      = "https://geo.mirror.pkgbuild.com/images/latest/Arch-Linux-x86_64-cloudimg.qcow2"
      img_checksum = "https://geo.mirror.pkgbuild.com/images/latest/Arch-Linux-x86_64-cloudimg.qcow2.SHA256"
    },
    "debian13" : {
      img_url      = "https://cdimage.debian.org/images/cloud/trixie/latest/debian-13-genericcloud-amd64.qcow2"
      img_checksum = "https://cdimage.debian.org/images/cloud/trixie/latest/SHA512SUMS"
    }
    "ubuntu24.04" : {
      img_url      = "https://cloud-images.ubuntu.com/noble/current/noble-server-cloudimg-amd64.img"
      img_checksum = "https://cloud-images.ubuntu.com/noble/current/SHA256SUMS"
    },
    "ubuntu25.10" : {
      img_url      = "https://cloud-images.ubuntu.com/questing/current/questing-server-cloudimg-amd64.img"
      img_checksum = "https://cloud-images.ubuntu.com/questing/current/SHA256SUMS"
    },
    "ubuntu26.04" : {
      img_url      = "https://cloud-images.ubuntu.com/resolute/current/resolute-server-cloudimg-amd64.img"
      img_checksum = "https://cloud-images.ubuntu.com/resolute/current/SHA256SUMS"
    },
    "opensuse" : {
      img_url      = "https://download.opensuse.org/tumbleweed/appliances/openSUSE-Tumbleweed-Minimal-VM.x86_64-Cloud.qcow2"
      img_checksum = "https://download.opensuse.org/tumbleweed/appliances/openSUSE-Tumbleweed-Minimal-VM.x86_64-Cloud.qcow2.sha256"
    }
  }
}
