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
@{pci}=@{pci_bus}/**/

# Date and time
@{date}=[0-2][0-9][0-9][0-9]-[01][0-9]-[0-3][0-9]
@{time}={[0-2],}[0-9]-[0-5][0-9]-[0-6][0-9]

# @{MOUNTDIRS} is a space-separated list of where user mount directories
# are stored, for programs that must enumerate all mount directories on a
# system.
@{MOUNTDIRS}=/media/ @{run}/media/*/ /mnt/

# @{MOUNTS} is a space-separated list of all user mounted directories.
@{MOUNTS}=@{MOUNTDIRS}/*/

# Common places for binaries and libraries across distributions
@{bin}=/{,usr/}{,s}bin
@{lib}=/{,usr/}lib{,exec,32,64}

# accommodate different shells
@{unix_shell}={ba,da,fi,k,z}sh

# Image extensions
# bmp, jpg, jpeg, png, gif
@{ext_pictures}={[bB][mM][pP],[jJ][pP]{,[eE]}[gG],[pP][nN][gG],[gG][iI][fF],[wW][eE][bB][pP]}

# Audio/video extensions:
# a52, aac, ac3, mka, flac, mp1, mp2, mp3, mpc, oga, oma, wav, wv, wm, wma, 3g2, 3gp, 3gp2, 3gpp,
# asf, avi, divx, m1v, m2v, m4v, mkv, mov, mp4, mpa, mpe, mpg, mpeg, mpeg1, mpeg2, mpeg4, ogg, ogm,
# ogx, ogv, rm, rmvb, webm, wmv, wtv, mp2t, flv, m4a
@{ext_audiovideo}  = [aA]{52,[aA][cC],[cC]3}
@{ext_audiovideo} += [mM][kK][aA]
@{ext_audiovideo} += [fF][lL][aA][cC]
@{ext_audiovideo} += [mM][pP][123cC]
@{ext_audiovideo} += [oO][gGmM][aA]
@{ext_audiovideo} += [wW]{,[aA]}[vV]
@{ext_audiovideo} += [wW][mM]{,[aA]}
@{ext_audiovideo} +=  3[gG]{[2pP],[pP][2pP]}
@{ext_audiovideo} += [aA][sS][fF]
@{ext_audiovideo} += [aA][vV][iI]
@{ext_audiovideo} += [dD][iI][vV][xX]
@{ext_audiovideo} += [mM][124][vV]
@{ext_audiovideo} += [mM][kKoO][vV]
@{ext_audiovideo} += [mM][pP][4aAeEgG]
@{ext_audiovideo} += [mM][pP][eE][gG]{,[124]}
@{ext_audiovideo} += [oO][gG][gGmMxXvV]
@{ext_audiovideo} += [rR][mM]{,[vV][bB]}
@{ext_audiovideo} += [wW][eE][bB][mM]
@{ext_audiovideo} += [wW][mMtT][vV]
@{ext_audiovideo} += [mM][pP]2[tT]
@{ext_audiovideo} += [fF][lL][vV]
@{ext_audiovideo} += [mM]4[aA]
