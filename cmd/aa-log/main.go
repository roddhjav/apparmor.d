// aa-log - Review AppArmor generated messages
// Copyright (C) 2021-2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/roddhjav/apparmor.d/pkg/logs"
	"github.com/roddhjav/apparmor.d/pkg/util"
)

const usage = `aa-log [-h] [--systemd] [--file file] [profile]

    Review AppArmor generated messages in a colorful way. Supports logs from
    auditd, systemd, syslog as well as dbus session events.

    It can be given an optional profile name to filter the output with.

    Default logs are read from '/var/log/audit/audit.log'. Other files in 
    '/var/log/audit/' can easily be checked: 'aa-log -f 1' parses 'audit.log.1' 

Options:
    -h, --help         Show this help message and exit.
    -f, --file FILE    Set a logfile or a suffix to the default log file.
    -s, --systemd      Parse systemd logs from journalctl.

`

// Command line options
var (
	help    bool
	path    string
	systemd bool
)

func aaLog(logger string, path string, profile string) error {
	var err error
	var file io.Reader

	switch logger {
	case "auditd":
		file, err = logs.GetAuditLogs(path)
	case "systemd":
		file, err = logs.GetJournalctlLogs(path, !util.InSlice(path, logs.LogFiles))
	default:
		err = fmt.Errorf("Logger %s not supported.", logger)
	}
	if err != nil {
		return err
	}
	aaLogs := logs.NewApparmorLogs(file, profile)
	fmt.Print(aaLogs.String())
	return nil
}

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.StringVar(&path, "f", "", "Set a logfile or a suffix to the default log file.")
	flag.StringVar(&path, "file", "", "Set a logfile or a suffix to the default log file.")
	flag.BoolVar(&systemd, "s", false, "Parse systemd logs from journalctl.")
	flag.BoolVar(&systemd, "systemd", false, "Parse systemd logs from journalctl.")
}

func main() {
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	profile := ""
	if len(flag.Args()) >= 1 {
		profile = flag.Args()[0]
	}

	logger := "auditd"
	if systemd {
		logger = "systemd"
	}

	logfile := logs.GetLogFile(path)
	err := aaLog(logger, logfile, profile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
