// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/roddhjav/apparmor.d/pkg/prebuild"
	"github.com/roddhjav/apparmor.d/pkg/tasks"
)

type FilterOnly struct {
	tasks.BaseTask
}

type FilterExclude struct {
	tasks.BaseTask
}

// NewFilterOnly creates a new FilterOnly directive.
func NewFilterOnly() *FilterOnly {
	return &FilterOnly{
		BaseTask: tasks.BaseTask{
			Keyword: "only",
			Msg:     "Only directive applied",
			Help:    []string{"filters..."},
		},
	}
}

// NewFilterExclude creates a new FilterExclude directive.
func NewFilterExclude() *FilterExclude {
	return &FilterExclude{
		BaseTask: tasks.BaseTask{
			Keyword: "exclude",
			Msg:     "Exclude directive applied",
			Help:    []string{"filters..."},
		},
	}
}

func cmp[T float64 | int](refValue T, operator string, value T) bool {
	var res bool
	switch operator {
	case "<":
		res = refValue < value
	case "<=":
		res = refValue <= value
	case ">":
		res = refValue > value
	case ">=":
		res = refValue >= value
	case "==", "=":
		res = refValue == value
	}
	return res
}

func compare(refValue any, prefix string, arg string) bool {
	pattern := fmt.Sprintf(`^%s(==|[<>]=?|=)(.+)$`, prefix)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(arg)
	if len(matches) < 3 {
		return false
	}
	operator := matches[1]
	targetStr := matches[2]

	var res bool
	switch refValue := refValue.(type) {
	case int:
		targetValue, err := strconv.Atoi(targetStr)
		if err != nil {
			panic(err)
		}
		res = cmp(refValue, operator, targetValue)

	case float64:
		targetValue, err := strconv.ParseFloat(targetStr, 64)
		if err != nil {
			panic(err)
		}
		res = cmp(refValue, operator, targetValue)

	default:
		panic("unsupported type")
	}

	return res
}

func filterRuleForUs(opt *Option) bool {
	for _, arg := range opt.ArgList {
		var res bool
		if prebuild.RBAC && arg == "RBAC" {
			res = true
		}
		if prebuild.Test && arg == "test" {
			res = true
		}
		if arg == tasks.Distribution {
			res = true
		}
		if arg == tasks.Family {
			res = true
		}
		if strings.HasPrefix(arg, "abi") {
			res = compare(prebuild.ABI, "abi", arg)
		}
		if strings.HasPrefix(arg, "apparmor") {
			res = compare(prebuild.Version, "apparmor", arg)
		}

		if res {
			return true
		}
	}
	return false
}

func filter(only bool, opt *Option, profile string) (string, error) {
	if only && filterRuleForUs(opt) {
		return opt.Clean(profile), nil
	}
	if !only && !filterRuleForUs(opt) {
		return opt.Clean(profile), nil
	}

	if opt.IsInline() {
		profile = strings.ReplaceAll(profile, opt.Raw, "")
	} else {
		regRemoveParagraph := regexp.MustCompile(`(?s)` + opt.Raw + `\n.*?\n\n`)
		profile = regRemoveParagraph.ReplaceAllString(profile, "")
	}
	return profile, nil
}

func (d FilterOnly) Apply(opt *Option, profile string) (string, error) {
	return filter(true, opt, profile)
}

func (d FilterExclude) Apply(opt *Option, profile string) (string, error) {
	return filter(false, opt, profile)
}
