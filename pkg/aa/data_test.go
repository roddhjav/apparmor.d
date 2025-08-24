// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

var (
	// Comment
	comment1 = &Comment{Base: Base{Comment: "comment", IsLineRule: true}}
	comment2 = &Comment{Base: Base{Comment: "another comment", IsLineRule: true}}

	// Abi
	abi1 = &Abi{IsMagic: true, Path: "abi/4.0"}
	abi2 = &Abi{IsMagic: true, Path: "abi/3.0"}

	// Alias
	alias1 = &Alias{Path: "/mnt/usr", RewrittenPath: "/usr"}
	alias2 = &Alias{Path: "/mnt/var", RewrittenPath: "/var"}

	// Include
	include1      = &Include{IsMagic: true, Path: "abstraction/base"}
	include2      = &Include{IsMagic: false, Path: "abstraction/base"}
	includeLocal1 = &Include{IfExists: true, IsMagic: true, Path: "local/foo"}

	// Variable
	variable1 = &Variable{Name: "bin", Values: []string{"/{,usr/}{,s}bin"}, Define: true}
	variable2 = &Variable{Name: "exec_path", Values: []string{"@{bin}/foo", "@{lib}/foo"}, Define: true}

	// All
	all1 = &All{}
	all2 = &All{Base: Base{Comment: "comment"}}

	// Rlimit
	rlimit1 = &Rlimit{Key: "nproc", Op: "<=", Value: "200"}
	rlimit2 = &Rlimit{Key: "cpu", Op: "<=", Value: "2"}
	rlimit3 = &Rlimit{Key: "nproc", Op: "<", Value: "2"}

	// Userns
	userns1 = &Userns{Create: true}
	userns2 = &Userns{}

	// Capability
	capability1Log = map[string]string{
		"apparmor":   "ALLOWED",
		"class":      "cap",
		"operation":  "capable",
		"capname":    "net_admin",
		"capability": "12",
		"profile":    "pkexec",
		"comm":       "pkexec",
	}
	capability1 = &Capability{Names: []string{"net_admin"}}
	capability2 = &Capability{Names: []string{"sys_ptrace"}}

	// Network
	network1Log = map[string]string{
		"apparmor":       "ALLOWED",
		"class":          "net",
		"operation":      "create",
		"family":         "netlink",
		"profile":        "sddm-greeter",
		"sock_type":      "raw",
		"protocol":       "15",
		"requested_mask": "create",
		"denied_mask":    "create",
		"comm":           "sddm-greeter",
	}
	network3Log = map[string]string{
		"apparmor":  "ALLOWED",
		"class":     "net",
		"operation": "sendmsg",
		"info":      "failed af match",
		"error":     "-13",
		"profile":   "unattended-upgrade",
		"comm":      "unattended-upgr",
		"laddr":     "127.0.0.1",
		"lport":     "57007",
		"faddr":     "127.0.0.53",
		"saddr":     "127.0.0.1",
		"src":       "57007",
		"fport":     "53",
		"sock_type": "dgram",
		"protocol":  "17",
		"requested": "send",
		"denied":    "send",
	}
	network1 = &Network{Domain: "netlink", Type: "raw", Protocol: "15"}
	network2 = &Network{Domain: "inet", Type: "dgram"}
	network3 = &Network{
		Base:         Base{Comment: " failed af match"},
		LocalAddress: LocalAddress{IP: "127.0.0.1", Port: "57007"},
		PeerAddress:  PeerAddress{IP: "127.0.0.53", Port: "53", Src: "127.0.0.1"},
		Type:         "dgram",
		Protocol:     "17",
	}

	// Mount
	mount1Log = map[string]string{
		"apparmor":  "ALLOWED",
		"class":     "mount",
		"operation": "mount",
		"info":      "failed perms check",
		"error":     "-13",
		"profile":   "dockerd",
		"name":      "/var/lib/docker/overlay2/opaque-bug-check1209538631/merged/",
		"comm":      "dockerd",
		"fstype":    "overlay",
		"srcname":   "overlay",
	}
	mount2Log = map[string]string{
		"apparmor":  "ALLOWED",
		"class":     "mount",
		"operation": "mount",
		"info":      "failed perms check",
		"error":     "-13",
		"profile":   "loupe",
		"name":      "/newroot/dev/tty",
		"comm":      "bwrap",
		"srcname":   "/oldroot/dev/tty",
		"flags":     "rw, rbind",
	}
	mount1 = &Mount{
		Base:            Base{Comment: " failed perms check"},
		MountConditions: MountConditions{FsType: "overlay"},
		Source:          "overlay",
		MountPoint:      "/var/lib/docker/overlay2/opaque-bug-check1209538631/merged/",
	}
	mount2 = &Mount{
		Base:            Base{Comment: " failed perms check"},
		MountConditions: MountConditions{Options: []string{"rw", "rbind"}},
		Source:          "/oldroot/dev/tty",
		MountPoint:      "/newroot/dev/tty",
	}

	// Remount
	remount1 = &Remount{MountPoint: "/"}
	remount2 = &Remount{MountPoint: "/{,**}/"}

	// Umount
	umount1Log = map[string]string{
		"apparmor":  "ALLOWED",
		"class":     "mount",
		"operation": "umount",
		"profile":   "systemd",
		"name":      "/",
		"comm":      "(ostnamed)",
	}
	umount1 = &Umount{MountPoint: "/"}
	umount2 = &Umount{MountPoint: "/oldroot/"}

	// PivotRoot
	pivotroot1Log = map[string]string{
		"apparmor":  "ALLOWED",
		"class":     "mount",
		"profile":   "systemd",
		"operation": "pivotroot",
		"comm":      "(ostnamed)",
		"name":      "@{run}/systemd/mount-rootfs/",
		"srcname":   "@{run}/systemd/mount-rootfs/",
	}
	pivotroot1 = &PivotRoot{
		OldRoot: "@{run}/systemd/mount-rootfs/",
		NewRoot: "@{run}/systemd/mount-rootfs/",
	}
	pivotroot2 = &PivotRoot{
		OldRoot:       "@{run}/systemd/mount-rootfs/",
		NewRoot:       "/newroot",
		TargetProfile: "brwap",
	}
	pivotroot3 = &PivotRoot{
		NewRoot: "/newroot",
	}

	// Change Profile
	changeprofile1Log = map[string]string{
		"apparmor":  "ALLOWED",
		"class":     "file",
		"profile":   "systemd",
		"operation": "change_onexec",
		"comm":      "(systemd)",
		"name":      "systemd-user",
		"target":    "systemd-user",
	}
	changeprofile1 = &ChangeProfile{ProfileName: "systemd-user"}
	changeprofile2 = &ChangeProfile{ProfileName: "brwap"}
	changeprofile3 = &ChangeProfile{ExecMode: "safe", Exec: "/bin/bash", ProfileName: "brwap//default"}

	// Mqueue
	mqueue1 = &Mqueue{Access: []string{"r"}, Type: "posix", Name: "/"}
	mqueue2 = &Mqueue{Access: []string{"r"}, Type: "sysv", Name: "/"}

	// IO Uring
	iouring1 = &IOUring{Access: []string{"sqpoll"}, Label: "foo"}
	iouring2 = &IOUring{Access: []string{"override_creds"}}

	// Signal
	signal1Log = map[string]string{
		"apparmor":       "ALLOWED",
		"class":          "signal",
		"profile":        "firefox",
		"operation":      "signal",
		"comm":           "49504320492F4F20506172656E74",
		"requested_mask": "receive",
		"denied_mask":    "receive",
		"signal":         "kill",
		"peer":           "firefox//&firejail-default",
	}
	signal1 = &Signal{
		Access: []string{"receive"},
		Set:    []string{"kill"},
		Peer:   "firefox//&firejail-default",
	}
	signal2 = &Signal{
		Access: []string{"receive"},
		Set:    []string{"up"},
		Peer:   "firefox//&firejail-default",
	}

	// Ptrace
	ptrace1Log = map[string]string{
		"apparmor":       "ALLOWED",
		"class":          "ptrace",
		"profile":        "xdg-document-portal",
		"operation":      "ptrace",
		"comm":           "pool-/usr/lib/x",
		"requested_mask": "read",
		"denied_mask":    "read",
		"peer":           "nautilus",
	}
	ptrace2Log = map[string]string{
		"apparmor":       "DENIED",
		"class":          "ptrace",
		"operation":      "ptrace",
		"comm":           "systemd-journal",
		"requested_mask": "readby",
		"denied_mask":    "readby",
		"peer":           "systemd-journald",
	}
	ptrace1 = &Ptrace{Access: []string{"read"}, Peer: "nautilus"}
	ptrace2 = &Ptrace{Access: []string{"readby"}, Peer: "systemd-journald"}

	// Unix
	unix1Log = map[string]string{
		"apparmor":       "ALLOWED",
		"class":          "unix",
		"family":         "unix",
		"operation":      "file_perm",
		"profile":        "gsettings",
		"comm":           "dbus-daemon",
		"requested_mask": "send receive",
		"addr":           "none",
		"peer_addr":      "@/tmp/dbus-AaKMpxzC4k",
		"peer":           "dbus-daemon",
		"denied_mask":    "send receive",
		"sock_type":      "stream",
		"protocol":       "0",
	}
	unix1 = &Unix{
		Access:    []string{"send", "receive"},
		Type:      "stream",
		Protocol:  "0",
		Address:   "none",
		PeerAddr:  "@/tmp/dbus-AaKMpxzC4k",
		PeerLabel: "dbus-daemon",
	}
	unix2 = &Unix{
		Base:   Base{FileInherit: true},
		Access: []string{"receive"},
		Type:   "stream",
	}

	// Dbus
	dbus1Log = map[string]string{
		"apparmor":   "ALLOWED",
		"operation":  "dbus_method_call",
		"bus":        "session",
		"path":       "/org/gtk/vfs/metadata",
		"interface":  "org.gtk.vfs.Metadata",
		"member":     "Remove",
		"name":       ":1.15",
		"mask":       "receive",
		"label":      "gvfsd-metadata",
		"peer_pid":   "3888",
		"peer_label": "tracker-extract",
	}
	dbus2Log = map[string]string{
		"apparmor":  "ALLOWED",
		"operation": "dbus_bind",
		"bus":       "session",
		"name":      "org.gnome.evolution.dataserver.Sources5",
		"mask":      "bind",
		"pid":       "3442",
		"label":     "evolution-source-registry",
	}
	dbus1 = &Dbus{
		Access:    []string{"receive"},
		Bus:       "session",
		Path:      "/org/gtk/vfs/metadata",
		Interface: "org.gtk.vfs.Metadata",
		Member:    "Remove",
		PeerName:  ":1.15",
		PeerLabel: "tracker-extract",
	}
	dbus2 = &Dbus{
		Access: []string{"bind"},
		Bus:    "session",
		Name:   "org.gnome.evolution.dataserver.Sources5",
	}
	dbus3 = &Dbus{
		Access: []string{"bind"},
		Bus:    "session",
		Name:   "org.gnome.evolution.dataserver",
	}

	// File
	file1Log = map[string]string{
		"apparmor":       "ALLOWED",
		"operation":      "open",
		"class":          "file",
		"profile":        "cupsd",
		"name":           "/usr/share/poppler/cMap/Identity-H",
		"comm":           "gs",
		"requested_mask": "r",
		"denied_mask":    "r",
		"fsuid":          "209",
		"FSUID":          "cups",
		"ouid":           "0",
		"OUID":           "root",
	}
	file2Log = map[string]string{
		"apparmor":       "ALLOWED",
		"operation":      "open",
		"class":          "file",
		"profile":        "gsd-print-notifications",
		"name":           "@{PROC}/4163/cgroup",
		"comm":           "gsd-print-notif",
		"requested_mask": "r",
		"denied_mask":    "r",
		"fsuid":          "1000",
		"FSUID":          "user",
		"ouid":           "1000",
		"OUID":           "user",
		"error":          "-1",
	}
	file1 = &File{Path: "/usr/share/poppler/cMap/Identity-H", Access: []string{"r"}}
	file2 = &File{
		Base:   Base{NoNewPrivs: true},
		Owner:  true,
		Path:   "@{PROC}/4163/cgroup",
		Access: []string{"r"},
	}

	// Link
	link1Log = map[string]string{
		"apparmor":       "ALLOWED",
		"operation":      "link",
		"class":          "file",
		"profile":        "mkinitcpio",
		"name":           "/tmp/mkinitcpio.QDWtza/early@{lib}/firmware/i915/dg1_dmc_ver2_02.bin.zst",
		"comm":           "cp",
		"requested_mask": "l",
		"denied_mask":    "l",
		"fsuid":          "0",
		"ouid":           "0",
		"target":         "/tmp/mkinitcpio.QDWtza/root@{lib}/firmware/i915/dg1_dmc_ver2_02.bin.zst",
		"FSUID":          "root",
		"OUID":           "root",
	}
	link3Log = map[string]string{
		"apparmor":       "ALLOWED",
		"operation":      "link",
		"class":          "file",
		"profile":        "dolphin",
		"name":           "@{user_config_dirs}/kiorc",
		"comm":           "dolphin",
		"requested_mask": "l",
		"denied_mask":    "l",
		"fsuid":          "1000",
		"ouid":           "1000",
		"target":         "@{user_config_dirs}/#3954",
	}
	link1 = &Link{
		Path:   "/tmp/mkinitcpio.QDWtza/early@{lib}/firmware/i915/dg1_dmc_ver2_02.bin.zst",
		Target: "/tmp/mkinitcpio.QDWtza/root@{lib}/firmware/i915/dg1_dmc_ver2_02.bin.zst",
	}
	link2 = &Link{
		Owner:  true,
		Path:   "@{user_config_dirs}/powerdevilrc{,.@{rand6}}",
		Target: "@{user_config_dirs}/#@{int}",
	}
	link3 = &Link{
		Owner:  true,
		Path:   "@{user_config_dirs}/kiorc",
		Target: "@{user_config_dirs}/#3954",
	}

	// Profile
	profile1 = &Profile{
		Header: Header{
			Name:        "sudo",
			Attachments: []string{},
			Attributes:  map[string]string{},
			Flags:       []string{},
		},
	}
	profile2 = &Profile{
		Header: Header{
			Name:        "systemctl",
			Attachments: []string{},
			Attributes:  map[string]string{},
			Flags:       []string{},
		},
	}

	// Hat
	hat1 = &Hat{Name: "user"}
	hat2 = &Hat{Name: "root"}
)
