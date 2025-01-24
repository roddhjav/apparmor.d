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

variable "ssh_publickey" {
  description = "Path to the ssh public key"
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}

variable "cpus" {
  description = "Default CPU of the VM"
  type        = string
  default     = "6"
}

variable "ram" {
  description = "Default RAM of the VM"
  type        = string
  default     = "4096"
}

variable "disk_size" {
  description = "Disk size of the VM to build"
  type        = string
  default     = "40G"
}

variable "iso_dir" {
  description = "Original ISO file directory"
  type        = string
  default     = "/var/lib/libvirt/images"
}

variable "base_dir" {
  description = "Final packer image output directory"
  type        = string
  default     = "/var/lib/libvirt/images"
}

variable "firmware" {
  description = "Path to the UEFI firmware"
  type        = string
  default     = "/usr/share/edk2/x64/OVMF_CODE.fd"
}

variable "output" {
  description = "Output build directory"
  type        = string
  default     = "/tmp/packer"
}

variable "prefix" {
  description = "Image name prefix"
  type        = string
  default     = "aa-"
}

variable "flavor" {
  description = "Distribution flavor to use (server, desktop, gnome, kde...)"
  type        = string
  default     = ""
}

variable "release" {
  description = "Distribution metadata to use"
  type = map(object({
    codename = string
    version  = string
  }))
  default = {
    "ubuntu22" : {
      codename = "jammy",
      version  = "22.04.2",
    },
    "ubuntu24" : {
      codename = "noble",
      version  = "24.04",
    },
    "debian" : {
      codename = "bookworm",
      version  = "12",
    }
    "opensuse" : {
      codename = "tumbleweed",
      version  = "",
    }
    "fedora" : {
      codename = "40",
      version  = "1.14",
    }
  }
}
