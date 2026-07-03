---
title: Base
tags:
  - tunables
  - default
---

## system

### @{d}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L12 "View source"){ .abs-source }

Any digit

```
@{d}=[0-9]
```
### @{l}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L15 "View source"){ .abs-source }

Any letter

```
@{l}=[a-zA-Z]
```
### @{c}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L18 "View source"){ .abs-source }

Single alphanumeric character

```
@{c}=[0-9a-zA-Z]
```
### @{w}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L21 "View source"){ .abs-source }

Word character: matches any letter, digit or underscore.

```
@{w}=[a-zA-Z0-9_]
```
### @{h}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L24 "View source"){ .abs-source }

Single hexadecimal character

```
@{h}=[0-9a-fA-F]
```
### @{int}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L27 "View source"){ .abs-source }

Integer up to 10 digits (0-9999999999)

```
@{int}=@{d}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}{@{d},}
```
### @{hex}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L30 "View source"){ .abs-source }

hexadecimal, alphanumeric and word up to 64 characters

```
@{hex}=@{h}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}{@{h},}
@{rand}=@{c}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}{@{c},}
@{word}=@{w}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}{@{w},}
```
### @{u8}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L35 "View source"){ .abs-source }

Unsigned integer over 8 bits (0...255)

```
@{u8}=[0-9]{[0-9],} 1[0-9][0-9] 2[0-4][0-9] 25[0-5]
```
### @{u16}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L38 "View source"){ .abs-source }

Unsigned integer over 16 bits (0...65,535 5 digits)

```
@{u16}={@{d},[1-9]@{d},[1-9][@{d}@{d},[1-9]@{d}@{d}@{d},[1-6]@{d}@{d}@{d}@{d}}
```
### @{u32}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L41 "View source"){ .abs-source }

Unsigned integer over 32 bits (0...4,294,967,295 10 digits)

```
@{u32}={@{d},[1-9]@{d},[1-9]@{d}@{d},[1-9]@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-4]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}}
```
### @{u64}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L44 "View source"){ .abs-source }

Unsigned integer over 64 bits (0...18,446,744,073,709,551,615 20 digits).

```
@{u64}={@{d},[1-9]@{d},[1-9]@{d}@{d},[1-9]@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},[1-9]@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d},1@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}@{d}}
```
### @{int2}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L47 "View source"){ .abs-source }

Any x digits characters

```
@{int2}=@{d}@{d}
@{int4}=@{int2}@{int2}
@{int6}=@{int4}@{int2}
@{int8}=@{int4}@{int4}
@{int9}=@{int8}@{d}
@{int10}=@{int8}@{int2}
@{int12}=@{int8}@{int4}
@{int15}=@{int8}@{int4}@{int2}@{d}
@{int16}=@{int8}@{int8}
@{int32}=@{int16}@{int16}
@{int64}=@{int32}@{int32}
```
### @{hex2}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L60 "View source"){ .abs-source }

Any x hexadecimal characters

```
@{hex2}=@{h}@{h}
@{hex4}=@{hex2}@{hex2}
@{hex6}=@{hex4}@{hex2}
@{hex8}=@{hex4}@{hex4}
@{hex9}=@{hex8}@{h}
@{hex10}=@{hex8}@{hex2}
@{hex12}=@{hex8}@{hex4}
@{hex15}=@{hex8}@{hex4}@{hex2}@{h}
@{hex16}=@{hex8}@{hex8}
@{hex32}=@{hex16}@{hex16}
@{hex38}=@{hex32}@{hex6}
@{hex64}=@{hex32}@{hex32}
```
### @{rand2}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L74 "View source"){ .abs-source }

Any x alphanumeric characters

```
@{rand2}=@{c}@{c}
@{rand4}=@{rand2}@{rand2}
@{rand6}=@{rand4}@{rand2}
@{rand8}=@{rand4}@{rand4}
@{rand9}=@{rand8}@{c}
@{rand10}=@{rand8}@{rand2}
@{rand12}=@{rand8}@{rand4}
@{rand15}=@{rand8}@{rand4}@{rand2}@{c}
@{rand16}=@{rand8}@{rand8}
@{rand32}=@{rand16}@{rand16}
@{rand64}=@{rand32}@{rand32}
```
### @{word2}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L87 "View source"){ .abs-source }

Any x word characters

```
@{word2}=@{w}@{w}
@{word4}=@{word2}@{word2}
@{word6}=@{word4}@{word2}
@{word8}=@{word4}@{word4}
@{word9}=@{word8}@{w}
@{word10}=@{word8}@{word2}
@{word12}=@{word8}@{word4}
@{word15}=@{word8}@{word4}@{word2}@{w}
@{word16}=@{word8}@{word8}
@{word32}=@{word16}@{word16}
@{word64}=@{word32}@{word32}
```
### @{pci_bus}

[:material-file-eye:](https://gitlab.com/apparmor/apparmor/-/blob/9a402e9e6a04b45dd788c7dfb0106dae45443b01/profiles/apparmor.d/tunables/system#L100 "View source"){ .abs-source }

Shortcut for PCI bus (e.g., /sys/devices/@{pci_bus}/**)

```
@{pci_bus}=pci@{hex4}:@{hex2}
```
