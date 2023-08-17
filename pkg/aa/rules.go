// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

// Preamble section of a profile
type Preamble struct {
	Abi              []Abi
	PreambleIncludes []Include
	Aliases          []Alias
	Variables        []Variable
}

// Profile section of a profile
type Profile struct {
	Name        string
	Attachments []string
	Attributes  []string
	Flags       []string
	Rules
}

type Rules struct {
	Includes      []Include
	Rlimit        []Rlimit
	Userns        Userns
	Capability    []Capability
	Network       []Network
	Mount         []Mount
	Umount        []Umount
	Remount       []Remount
	PivotRoot     []PivotRoot
	ChangeProfile []ChangeProfile
	Unix          []Unix
	Ptrace        []Ptrace
	Signal        []Signal
	Dbus          []Dbus
	File          []File
}


// Qualifier to apply extra settings to a rule
type Qualifier struct {
	Audit       bool
	AccessType  string
	Owner       bool
	NoNewPrivs  bool
	FileInherit bool
}

// Preamble rules

type Abi struct {
	AbsPath   string
	MagicPath string
}

type Alias struct {
	Path          string
	RewrittenPath string
}

type Include struct {
	IfExists  bool
	AbsPath   string
	MagicPath string
}

type Variable struct {
	Name   string
	Values []string
}

// Profile rules

type Rlimit struct {
	Key   string
	Op    string
	Value string
}

type Userns struct {
	Qualifier
	Create bool
}

type Capability struct {
	Qualifier
	Name string
}

type AddressExpr struct {
	Source      string
	Destination string
	Port        string
}

type Network struct {
	Qualifier
	Domain   string
	Type     string
	Protocol string
	AddressExpr
}

type MountConditions struct {
	Fs      string
	Op      string
	FsType  string
	Options []string
}

type Mount struct {
	Qualifier
	MountConditions
	Source     string
	MountPoint string
}

type Umount struct {
	Qualifier
	MountConditions
	MountPoint string
}

type Remount struct {
	Qualifier
	MountConditions
	MountPoint string
}

type PivotRoot struct {
	Qualifier
	OldRoot       string
	NewRoot       string
	TargetProfile string
}

type ChangeProfile struct {
	ExecMode    string
	Exec        string
	ProfileName string
}

type IOUring struct {
	Qualifier
	Access string
	Label  string
}

type Signal struct {
	Qualifier
	Access string
	Set    string
	Peer   string
}

type Ptrace struct {
	Qualifier
	Access string
	Peer   string
}

type Unix struct {
	Qualifier
	Access   string
	Type     string
	Protocol string
	Address  string
	Label    string
	Attr     string
	Opt      string
	Peer     string
	PeerAddr string
}

type Mqueue struct {
	Qualifier
	Access string
	Type   string
	Label  string
}

type Dbus struct {
	Qualifier
	Access    string
	Bus       string
	Name      string
	Path      string
	Interface string
	Member    string
	Label     string
}

type File struct {
	Qualifier
	Path   string
	Access string
	Target string
}

