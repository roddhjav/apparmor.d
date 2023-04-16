# apparmor.d - Full set of apparmor profiles
# Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
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

variable "ssh_privatekey" {
  description = "Path to the ssh private key"
  type        = string
  default     = "~/.ssh/id_ed25519"
}

variable "disk_size" {
  description = "Disk size of the App VM to build"
  type        = string
  default     = "10G"
}

variable "iso_dir" {
  description = "Original ISO file directory"
  type        = string
  default     = "/var/lib/libvirt/images"
}

variable "output" {
  description = "Output build directory"
  type        = string
  default     = "/tmp/packer"
}

variable "prefix" {
  description = "Image name prefix"
  type        = string
  default     = "aa"
}

variable "version" {
  description = "apparmor.d version"
  type        = string
  default     = "0.001"
}

variable "release" {
  description = "Distribution release to use"
  type        = map(string)
  default = {
    "ubuntu" : "jammy",    # 22.04 LTS
    "debian" : "bullseye", # 11
    "opensuse" : "9",
  }
}

variable "flavor" {
  description = "Distribution flavor to use (-desktop, -gnome, -kde...)"
  type        = string
  default     = ""
}
