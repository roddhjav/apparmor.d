# apparmor.d - Full set of apparmor profiles
# Extended system directories definition
# Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
# SPDX-License-Identifier: GPL-2.0-only

# To allow extended personalisation without breaking everything.
# All apparmor profiles should always use the variables defined here.

# Single hexadecimal character
@{h}=[0-9a-fA-F]

# Single alphanumeric character
@{c}=[0-9a-zA-Z]

# Up to 10 digits (0-9999999999)
@{int}=[0-9]{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}{[0-9],}

# Any six characters
@{rand6}=@{c}@{c}@{c}@{c}@{c}@{c}

# Any eight characters
@{rand8}=@{c}@{c}@{c}@{c}@{c}@{c}@{c}@{c}

# Any ten characters
@{rand10}=@{c}@{c}@{c}@{c}@{c}@{c}@{c}@{c}@{c}@{c}

# MD5 hash
@{md5}=@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}

# Universally unique identifier
@{uuid}=@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}[-_]@{h}@{h}@{h}@{h}[-_]@{h}@{h}@{h}@{h}[-_]@{h}@{h}@{h}@{h}[-_]@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}@{h}

# Hexadecimal
@{hex}=@{h}*@{h}

# Shortcut for PCI device
@{pci_id}=@{h}@{h}@{h}@{h}:@{h}@{h}:@{h}@{h}.@{h}
@{pci_bus}=pci@{h}@{h}@{h}@{h}:@{h}@{h}
@{pci}=@{pci_bus}/@{pci_id}{,/@{pci_id}}{,/@{pci_bus}/@{pci_id}{,/@{pci_id}}}

# Date and time
@{date}=[0-2][0-9][0-9][0-9]-[01][0-9]-[0-3][0-9]
@{time}={[0-2],}[0-9]-[0-5][0-9]-[0-6][0-9]

# @{MOUNTDIRS} is a space-separated list of where user mount directories
# are stored, for programs that must enumerate all mount directories on a
# system.
@{MOUNTDIRS}=/media/* @{run}/media/* /mnt/

# @{MOUNTS} is a space-separated list of all user mounted directories.
@{MOUNTS}=@{MOUNTDIRS}/*/

# Common places for binaries and libraries across distributions
@{bin}=/{,usr/}{,s}bin
@{lib}=/{,usr/}lib{,exec,32,64}
