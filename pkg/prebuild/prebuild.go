// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2023 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package prebuild

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/arduino/go-paths-helper"
	"github.com/roddhjav/apparmor.d/pkg/logging"
)

var (
	Distribution  string
	DistDir       *paths.Path
	Root          *paths.Path
	RootApparmord *paths.Path
	FlagDir       *paths.Path
)

func init() {
	DistDir = paths.New("dists")
	Root = paths.New(".build")
	FlagDir = DistDir.Join("flags")
	RootApparmord = Root.Join("apparmor.d")
	Distribution = getSupportedDistribution()
}

func getFctName(i any) string {
	tmp := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	res := strings.Split(tmp, ".")
	return res[len(res)-1]
}

func printPrepareMessage(name string, msg []string) {
	logging.Success("%v", PrepareMsg[name])
	logging.Indent = "   "
	for _, line := range msg {
		logging.Bullet("%s", line)
	}
	logging.Indent = ""
}

func printBuildMessage() {
	for _, fct := range Builds {
		name := getFctName(fct)
		logging.Success("%v", BuildMsg[name])
	}
	for _, fct := range Directives {
		name := getFctName(fct)
		logging.Success("%v", DirectiveMsg[name])
	}
}

func Prepare() error {
	for _, fct := range Prepares {
		msg, err := fct()
		if err != nil {
			return err
		}
		printPrepareMessage(getFctName(fct), msg)
	}
	return nil
}

func Build() error {
	files, _ := RootApparmord.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories())
	for _, file := range files {
		if !file.Exist() {
			continue
		}
		content, _ := file.ReadFile()
		profile := string(content)
		for _, fct := range Builds {
			profile = fct(profile)
		}
		for _, fct := range Directives {
			profile = fct(file, profile)
		}
		if err := file.WriteFile([]byte(profile)); err != nil {
			panic(err)
		}
	}
	printBuildMessage()
	return nil
}
