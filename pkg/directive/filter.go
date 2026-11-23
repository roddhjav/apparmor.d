// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2026 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package directive

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

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

func compare(refValue any, prefix string, arg string) (bool, error) {
	pattern := fmt.Sprintf(`^%s(==|[<>]=?|=)(.+)$`, prefix)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(arg)
	if len(matches) < 3 {
		return false, nil
	}
	operator := matches[1]
	targetStr := matches[2]

	var res bool
	switch refValue := refValue.(type) {
	case int:
		targetValue, err := strconv.Atoi(targetStr)
		if err != nil {
			return false, fmt.Errorf("invalid %s comparison %q: %w", prefix, arg, err)
		}
		res = cmp(refValue, operator, targetValue)

	case float64:
		targetValue, err := strconv.ParseFloat(targetStr, 64)
		if err != nil {
			return false, fmt.Errorf("invalid %s comparison %q: %w", prefix, arg, err)
		}
		res = cmp(refValue, operator, targetValue)

	default:
		panic("unsupported type")
	}

	return res, nil
}

func filterRuleForUs(c *tasks.TaskConfig, opt *Option) (bool, error) {
	for _, arg := range opt.ArgList {
		switch {
		case c.RBAC && arg == "RBAC":
			return true, nil
		case c.Test && arg == "test":
			return true, nil
		case arg == tasks.Distribution:
			return true, nil
		case arg == tasks.Family:
			return true, nil
		case strings.HasPrefix(arg, "abi"):
			if ok, err := compare(c.ABI, "abi", arg); ok || err != nil {
				return ok, err
			}
		case strings.HasPrefix(arg, "apparmor"):
			if ok, err := compare(c.Version, "apparmor", arg); ok || err != nil {
				return ok, err
			}
		}
	}
	return false, nil
}

func filter(c *tasks.TaskConfig, only bool, opt *Option, profile string) (string, error) {
	forUs, err := filterRuleForUs(c, opt)
	if err != nil {
		return "", err
	}
	if only == forUs {
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
	return filter(d.TaskConfig, true, opt, profile)
}

func (d FilterExclude) Apply(opt *Option, profile string) (string, error) {
	return filter(d.TaskConfig, false, opt, profile)
}
