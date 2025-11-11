// aa-log - Review AppArmor generated messages
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/roddhjav/apparmor.d/pkg/logs"
)

const usage = `aa-log [-h] [--systemd] [--file file] [--load] [--rules | --raw] [--since] [--namespace] [profile]

    Review AppArmor generated messages in a colorful way. It supports logs from
    auditd, systemd, syslog as well as dbus session events.

    It can be given an optional profile name to filter the output with.

    Default logs are read from '/var/log/audit/audit.log'. Other files in
    '/var/log/audit/' can easily be checked: 'aa-log -f 1' parses 'audit.log.1'

    Logs written with 'aa-log' can be read again with 'aa-log -l'.

Options:
    -h, --help         Show this help message and exit.
    -f, --file FILE    Set a logfile or a suffix to the default log file.
    -s, --systemd      Parse systemd logs from journalctl.
    -n, --namespace NS Filter the logs to the specified namespace.
    -r, --rules        Convert the log into AppArmor rules.
    -R, --raw          Print the raw log without any formatting.
    -S, --since DATE   Show entries not older than the specified date.
    -l, --load         Load logs from the default aa-log output.

`

// Command line options
var (
	help      bool
	rules     bool
	path      string
	systemd   bool
	namespace string
	raw       bool
	since     string
	load      bool
)

func aaLog(logger string, path string, profile string, namespace string) error {
	var err error
	var file io.Reader

	start := time.Now()
	switch logger {
	case "auditd":
		file, err = logs.GetAuditLogs(path)
	case "systemd":
		file, err = logs.GetJournalctlLogs(path, since, !slices.Contains(logs.LogFiles, path))
	default:
		err = fmt.Errorf("logger %s not supported", logger)
	}
	if err != nil {
		return err
	}
	endRead := time.Now()

	if raw {
		fmt.Print(strings.Join(logs.GetApparmorLogs(file, profile, namespace), "\n") + "\n")
		return nil
	}

	var aaLogs logs.AppArmorLogs
	if load {
		aaLogs = logs.Load(file, profile, namespace)
	} else {
		aaLogs = logs.New(file, profile, namespace)
	}
	endParse := time.Now()
	if rules {
		profiles := aaLogs.ParseToProfiles()
		for _, p := range profiles {
			p.Merge(nil)
			p.Sort()
			p.Format()
			fmt.Print(p.String() + "\n\n")
		}
	} else {
		fmt.Print(aaLogs.String())
	}
	if withTime {
		printTiming(start, endRead, endParse, time.Now())
	}
	return nil
}

func init() {
	flag.BoolVar(&help, "h", false, "Show this help message and exit.")
	flag.BoolVar(&help, "help", false, "Show this help message and exit.")
	flag.StringVar(&path, "f", "", "Set a logfile or a suffix to the default log file.")
	flag.StringVar(&path, "file", "", "Set a logfile or a suffix to the default log file.")
	flag.BoolVar(&systemd, "s", false, "Parse systemd logs from journalctl.")
	flag.BoolVar(&systemd, "systemd", false, "Parse systemd logs from journalctl.")
	flag.BoolVar(&rules, "r", false, "Convert the log into AppArmor rules.")
	flag.BoolVar(&rules, "rules", false, "Convert the log into AppArmor rules.")
	flag.BoolVar(&raw, "R", false, "Print the raw log without any formatting.")
	flag.BoolVar(&raw, "raw", false, "Print the raw log without any formatting.")
	flag.StringVar(&since, "S", "", "Display logs since the START time.")
	flag.StringVar(&since, "since", "", "Display logs since the START time.")
	flag.BoolVar(&load, "l", false, "Load logs from the default aa-log output.")
	flag.BoolVar(&load, "load", false, "Load logs from the default aa-log output.")
	flag.StringVar(&namespace, "n", "", "Filter the logs to the specified namespace")
	flag.StringVar(&namespace, "namespace", "", "Filter the logs to the specified namespace")
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

	path = logs.SelectLogFile(path)
	err := aaLog(logger, path, profile, namespace)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
