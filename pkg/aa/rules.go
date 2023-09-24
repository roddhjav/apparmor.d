// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

// Qualifier to apply extra settings to a rule
type Qualifier struct {
	Audit       bool
	AccessType  string
	Owner       bool
	NoNewPrivs  bool
	FileInherit bool
}

func NewQualifier(owner, noNewPrivs, fileInherit bool) Qualifier {
	return Qualifier{
		Audit:       false,
		AccessType:  "",
		Owner:       owner,
		NoNewPrivs:  noNewPrivs,
		FileInherit: fileInherit,
	}
}

func NewCapability(log map[string]string, noNewPrivs, fileInherit bool) Capability {
	return Capability{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Name:      log["capname"],
	}
}

func NewNetwork(log map[string]string, noNewPrivs, fileInherit bool) Network {
	return Network{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		AddressExpr: AddressExpr{
			Source:      log["laddr"],
			Destination: log["faddr"],
			Port:        log["lport"],
		},
		Domain:   log["family"],
		Type:     log["sock_type"],
		Protocol: log["protocol"],
	}
}

func NewFile(log map[string]string, noNewPrivs, fileInherit bool) File {
	owner := false
	if log["fsuid"] == log["ouid"] && log["OUID"] != "root" {
		owner = true
	}
	return File{
		Qualifier: NewQualifier(owner, noNewPrivs, fileInherit),
		Path:      log["name"],
		Access:    maskToAccess[log["requested_mask"]],
		Target:    log["target"],
	}
}

func NewSignal(log map[string]string, noNewPrivs, fileInherit bool) Signal {
	return Signal{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    maskToAccess[log["requested_mask"]],
		Set:       log["signal"],
		Peer:      log["peer"],
	}
}

func NewPtrace(log map[string]string, noNewPrivs, fileInherit bool) Ptrace {
	return Ptrace{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    maskToAccess[log["requested_mask"]],
		Peer:      log["peer"],
	}
}

func NewUnix(log map[string]string, noNewPrivs, fileInherit bool) Unix {
	return Unix{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    maskToAccess[log["requested_mask"]],
		Type:      log["sock_type"],
		Protocol:  log["protocol"],
		Address:   log["addr"],
		Label:     log["peer_label"],
		Attr:      log["attr"],
		Opt:       log["opt"],
		Peer:      log["peer"],
		PeerAddr:  log["peer_addr"],
	}
}

func NewMount(log map[string]string, noNewPrivs, fileInherit bool) Mount {
	return Mount{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		MountConditions: MountConditions{
			Fs:      "",
			Op:      "",
			FsType:  log["fstype"],
			Options: []string{},
		},
		Source:     log["srcname"],
		MountPoint: log["name"],
	}
}

func NewDbus(log map[string]string, noNewPrivs, fileInherit bool) Dbus {
	return Dbus{
		Qualifier: NewQualifier(false, noNewPrivs, fileInherit),
		Access:    log["mask"],
		Bus:       log["bus"],
		Name:      log["name"],
		Path:      log["path"],
		Interface: log["interface"],
		Member:    log["member"],
		Label:     log["peer_label"],
	}
// Preamble specific rules

type Abi struct {
	Path    string
	IsMagic bool
}

func (r Abi) Less(other Abi) bool {
	if r.Path == other.Path {
		return r.IsMagic == other.IsMagic
	}
	return r.Path < other.Path
}

func (r Abi) Equals(other Abi) bool {
	return r.Path == other.Path && r.IsMagic == other.IsMagic
}

type Alias struct {
	Path          string
	RewrittenPath string
}

func (r Alias) Less(other Alias) bool {
	if r.Path == other.Path {
		return r.RewrittenPath < other.RewrittenPath
	}
	return r.Path < other.Path
}

func (r Alias) Equals(other Alias) bool {
	return r.Path == other.Path && r.RewrittenPath == other.RewrittenPath
}
