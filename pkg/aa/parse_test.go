// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"github.com/roddhjav/apparmor.d/pkg/paths"
)

func Test_tokenizeRule(t *testing.T) {
	inHeader = true
	for _, tt := range testParseRules {
		t.Run(tt.name, func(t *testing.T) {
			if got := tokenizeRule(tt.raw); !reflect.DeepEqual(got, tt.tokens) {
				t.Errorf("tokenize() = %v, want %v", got, tt.tokens)
			}
		})
	}
}

func Test_parseRule(t *testing.T) {
	inHeader = true
	for _, tt := range testParseRules {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRule(tt.raw); !reflect.DeepEqual(got, tt.rule) {
				t.Errorf("parseRule() = %v, want %v", got, tt.rule)
			}
		})
	}
}

func Test_rule_Getter(t *testing.T) {
	for _, tt := range testParseRules {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wGetAsMap == nil {
				tt.wGetAsMap = map[string][]string{}
			}
			if tt.wGetSlice == nil {
				tt.wGetSlice = []string{}
			}

			if tt.getIdx > 0 {
				if got := tt.rule.Get(tt.getIdx); got != tt.wGet {
					t.Errorf("rule.Get() = %v, want %v", got, tt.wGet)
				}
			}

			if got := tt.rule.GetString(); got != tt.wGetString {
				t.Errorf("rule.GetString() = %v, want %v", got, tt.wGetString)
			}
			if got := tt.rule.GetSlice(); !reflect.DeepEqual(got, tt.wGetSlice) {
				t.Errorf("rule.GetSlice() = %v, want %v", got, tt.wGetSlice)
			}
			if got := tt.rule.GetAsMap(); !reflect.DeepEqual(got, tt.wGetAsMap) {
				t.Errorf("rule.GetAsMap() = %v, want %v", got, tt.wGetAsMap)
			}

			if tt.getKey != "" {
				if got := tt.rule.GetValues(tt.getKey); !reflect.DeepEqual(got, tt.wGetValues) {
					t.Errorf("rule.GetValues() = %v, want %v", got, tt.wGetValues)
				}
				if got := tt.rule.GetValuesAsSlice(tt.getKey); !reflect.DeepEqual(got, tt.wGetValuesAsSlice) {
					t.Errorf("rule.GetValuesAsSlice() = %v, want %v", got, tt.wGetValuesAsSlice)
				}
				if got := tt.rule.GetValuesAsString(tt.getKey); got != tt.wGetValuesAsString {
					t.Errorf("rule.GetValuesAsString() = %v, want %v", got, tt.wGetValuesAsString)
				}
			}

			if got := tt.rule.String(); got != tt.wString {
				t.Errorf("rule.String() = |%v|, want |%v|", got, tt.wString)
			}
		})
	}
}

func Test_parseLineRules(t *testing.T) {
	for _, tt := range testsLineRules {
		t.Run(tt.name, func(t *testing.T) {
			got, rules, err := parseLineRules(tt.isPreamble, tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseLineRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseLineRules() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(rules, tt.rules) {
				t.Errorf("parseLineRules() rules = %v, want %v", rules, tt.rules)
			}
		})
	}
}

func Test_parseCommaRules(t *testing.T) {
	for _, tt := range testsCommaRules {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCommaRules(tt.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCommaRules() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.rule) {
				t.Errorf("parseCommaRules() = %v, want %v", got, tt.rule)
			}
		})
	}
}

func Test_tokenizeBlock(t *testing.T) {
	for _, tt := range testBlocks {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tokenizeBlock(tt.raw)
			if (err != nil) != tt.wTokenizeErr {
				t.Errorf("tokenizeBlock() error = %v, wantErr %v", err, tt.wTokenizeErr)
				return
			}
			if !reflect.DeepEqual(got, tt.blocks) {
				t.Errorf("tokenizeBlock() = %v, want %v", got, tt.blocks)
			}
		})
	}
}

func Test_parseBlock(t *testing.T) {
	for _, tt := range testBlocks {
		t.Run(tt.name, func(t *testing.T) {
			for idx, b := range tt.blocks {
				var err error
				var got Rules
				want := tt.rules[idx]
				if idx == 0 && b.kind == CONTENT && (strings.HasPrefix(b.raw, "# Simple test") || strings.Contains(b.raw, "$")) {
					f := &AppArmorProfileFile{}
					err = f.parsePreamble(b.raw)
					got = f.Preamble
				} else {
					got, err = parseBlock(b)
				}
				if (err != nil) != tt.wParseBlockErr {
					t.Errorf("parseBlock() error = %v, wantErr %v", err, tt.wParseBlockErr)
					return
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("parseBlock() = %v, want %v", got, want)
				}
			}
		})
	}
}

func Test_newRules(t *testing.T) {
	for _, tt := range testParseRules {
		if tt.wRule == nil {
			continue
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := newRules([]rule{tt.rule})
			if (err != nil) != tt.wError {
				t.Errorf("newRules() error = %v, wantErr %v", err, tt.wError)
				return
			}
			if len(got) == 0 {
				return
			}
			if !reflect.DeepEqual(got[0], tt.wRule) {
				t.Errorf("newRules() = %v, want %v", got[0], tt.wRule)
			}
		})
	}
}

func Test_AppArmorProfileFile_Parse(t *testing.T) {
	for _, tt := range testBlocks {
		t.Run(tt.name, func(t *testing.T) {
			got := &AppArmorProfileFile{}
			nb, err := got.Parse(tt.raw)
			if (err != nil) != tt.wParseErr {
				t.Errorf("AppArmorProfileFile.Parse() error = %v, wantErr %v", err, tt.wParseErr)
			}
			if !reflect.DeepEqual(got, tt.apparmor) {
				t.Errorf("AppArmorProfileFile.Parse() = |%v|, want |%v|", got, tt.apparmor)
			}
			raw := strings.Join(strings.Split(tt.raw, "\n")[nb:], "\n")
			gotRules, _, err := ParseRules(raw)
			if (err != nil) != tt.wParseRulesErr {
				t.Errorf("ParseRules() error = %v, wantErr %v", err, tt.wParseRulesErr)
				return
			}
			if !reflect.DeepEqual(gotRules, tt.wRules) {
				t.Errorf("ParseRules() got = %v, want %v", gotRules, tt.wRules)
			}
		})
	}
}

func Test_AppArmorProfileFile_Scan(t *testing.T) {
	for _, tt := range testBlocks {
		t.Run(tt.name, func(t *testing.T) {
			if tt.apparmorAll == nil {
				tt.apparmorAll = tt.apparmor
			}
			got := &AppArmorProfileFile{}
			if err := got.Scan(tt.raw); (err != nil) != tt.wParseErr {
				t.Errorf("AppArmorProfileFile.Scan() error = %v, wantErr %v", err, tt.wParseErr)
			}
			if !reflect.DeepEqual(got, tt.apparmorAll) {
				t.Errorf("AppArmorProfileFile.Scan() = %v, want %v", got, tt.apparmorAll)
			}
		})
	}
	for _, tt := range testParser {
		t.Run(tt.name, func(t *testing.T) {
			got := &AppArmorProfileFile{}
			if err := got.Scan(tt.raw); (err != nil) != tt.wantErr {
				t.Errorf("AppArmorProfileFile.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppArmorProfileFile.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testReport struct {
	Success bool   `csv:"success"`
	Name    string `csv:"name"`
	Desc    string `csv:"desc"`
	Error   string `csv:"error"`
}

// Test the parser on our own profiles
func Test_Parser_ApparmorD(t *testing.T) {
	datadir := paths.New("../../apparmor.d/")
	files, err := datadir.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories(), paths.FilterOutNames("README.md"))
	if err != nil {
		panic(err)
	}

	reports := []*testReport{}
	templateFailure := "\033[0;31m[FAILED]\033[0m(\033[0;37m%s\033[0m): %v"
	templateSuccess := "\033[0;32m[SUCCESS]\033[0m(\033[0;37m%s\033[0m)"
	for _, parse := range []bool{true, false} {
		base := "Parse/"
		if !parse {
			base = "Scan/"
		}

		for _, file := range files {
			if !file.Exist() {
				panic(file.String() + " %s not found")
			}
			name, err := file.RelFrom(datadir)
			if err != nil {
				panic(err)
			}

			t.Run(base+name.String(), func(t *testing.T) {
				var err error
				r := &testReport{Name: name.String()}
				raw := file.MustReadFileAsString()

				p := &AppArmorProfileFile{}
				p.Kind = KindFromPath(file)
				if parse {
					_, err = p.Parse(raw)
				} else {
					err = p.Scan(raw)
				}
				if err != nil {
					r.Error = err.Error()
					reports = append(reports, r)
					t.Errorf(templateFailure, name, err)
					return
				}

				if err = p.Validate(); err != nil {
					r.Error = err.Error()
					reports = append(reports, r)
					t.Errorf(templateFailure, name, err)
					return
				}

				r.Success = true
				reports = append(reports, r)
				t.Logf(templateSuccess, name)
			})
		}

		success := 0
		for _, r := range reports {
			if r.Success {
				success++
			}
		}
		// [13/12/25]: 4750 tests, success: 4750, fail 0, success rate: 100%
		t.Logf("[TOTAL]: %d tests, success: %d, fail %d, success rate: %v%%\n",
			len(reports), success, len(reports)-success, (success*100.0)/len(reports))
	}
}

func Test_Parser_Upstream(t *testing.T) {
	// os.Setenv("WITH_UPSTREAM", "true")
	if os.Getenv("WITH_UPSTREAM") == "" {
		t.Skip("Skipping test in CI environment")
	}
	apparmorDir := os.Getenv("APPARMOR_DIR")
	if apparmorDir == "" {
		apparmorDir = "../../../apparmor"
	}

	regTestDescription := regexp.MustCompile(`(?m)^#=(DESCRIPTION|Description) (.+)$`)
	regTestResult := regexp.MustCompile(`(?m)^#=EXRESULT (.+)$`)
	regTestDisabled := regexp.MustCompile(`(?m)^#=DISABLED$`)
	getSettings := func(profile string) (string, bool, bool) {
		desc := regTestDescription.FindStringSubmatch(profile)
		result := regTestResult.FindStringSubmatch(profile)
		disabled := regTestDisabled.FindStringSubmatch(profile)
		if len(disabled) == 1 && disabled[0] == "#=DISABLED" {
			return "", false, true
		}
		if len(desc) == 3 && len(result) == 2 {
			return desc[2], result[1] == "FAIL", false
		}
		return "", false, false
	}

	datadir := paths.New(apparmorDir).Join("parser/tst/simple_tests")
	files, err := datadir.ReadDirRecursiveFiltered(nil, paths.FilterOutDirectories(),
		paths.FilterOutNames("readme"))
	if err != nil {
		panic(err)
	}

	reports := []*testReport{}
	templateFailure := "\033[0;31m[FAILED]\033[0m(\033[0;37m%s\033[0m): %v"
	templateSuccess := "\033[0;32m[SUCCESS]\033[0m(\033[0;37m%s\033[0m)"
	for _, file := range files {
		if !file.Exist() {
			panic(file.String() + " %s not found")
		}
		name, err := file.RelFrom(datadir)
		if err != nil {
			panic(err)
		}

		raw := file.MustReadFileAsString()
		desc, wantErr, isDisabled := getSettings(raw)
		if isDisabled {
			t.Logf("Skipping disabled test: %s", name)
			continue
		}
		t.Run(name.String(), func(t *testing.T) {
			r := &testReport{Name: name.String()}
			r.Desc = desc

			p := &AppArmorProfileFile{}
			err := p.Scan(raw)
			if err != nil {
				if wantErr {
					r.Success = true
					reports = append(reports, r)
					t.Logf(templateSuccess, name)
				} else {
					r.Error = err.Error()
					reports = append(reports, r)
					t.Errorf(templateFailure, name, err)
				}
				return
			}

			err = p.Validate()
			if err != nil {
				if wantErr {
					r.Success = true
					reports = append(reports, r)
					t.Logf(templateSuccess, name)
				} else {
					r.Error = err.Error()
					reports = append(reports, r)
					t.Errorf(templateFailure, name, err)
				}
				return
			}

			// No errors occurred
			if wantErr {
				r.Error = "expected error but got none"
				reports = append(reports, r)
				t.Errorf(templateFailure, name, "expected error but got none")
			} else {
				r.Success = true
				reports = append(reports, r)
				t.Logf(templateSuccess, name)
			}
		})
	}

	success := 0
	for _, r := range reports {
		if r.Success {
			success++
		}
	}

	// [01/06/24]: 1986 tests, success: 1242, fail 744, success rate: 62%
	// [13/12/25]: 2148 tests, success: 1722, fail 426, success rate: 80%
	t.Logf("[TOTAL]: %d tests, success: %d, fail %d, success rate: %v%%\n",
		len(reports), success, len(reports)-success, (success*100.0)/len(reports))
}

var (
	// Test cases for tokenizeRule, parseRule,rule getters, and newRules
	testParseRules = []struct {
		name               string
		raw                string
		tokens             []string
		rule               rule
		token              []string
		getIdx             int    // index of the rule to get
		getKey             string // key to get
		wGet               string
		wGetString         string
		wGetSlice          []string
		wGetAsMap          map[string][]string
		wGetValues         rule
		wGetValuesAsSlice  []string
		wGetValuesAsString string
		wString            string
		wRule              Rule
		wError             bool
	}{
		{
			name:    "empty",
			raw:     "",
			tokens:  []string{},
			rule:    rule{},
			wString: "",
			wError:  true,
		},
		{
			name:   "abi",
			raw:    `abi <abi/4.0>`,
			tokens: []string{"abi", "<abi/4.0>"},
			rule: rule{
				{key: "abi"}, {key: "<abi/4.0>"},
			},
			getIdx:     1,
			wGet:       "<abi/4.0>",
			wGetString: "abi <abi/4.0>",
			wGetSlice:  []string{"abi", "<abi/4.0>"},
			wString:    "abi <abi/4.0>",
			wRule:      &Abi{IsMagic: true, Path: "abi/4.0"},
		},
		{
			name:   "alias",
			raw:    `alias /mnt/usr -> /usr`,
			tokens: []string{"alias", "/mnt/usr", "->", "/usr"},
			rule: rule{
				{key: "alias"}, {key: "/mnt/usr"}, {key: "->"}, {key: "/usr"},
			},
			getIdx:     2,
			wGet:       "->",
			wGetString: "alias /mnt/usr -> /usr",
			wGetSlice:  []string{"alias", "/mnt/usr", "->", "/usr"},
			wString:    "alias /mnt/usr -> /usr",
			wRule:      &Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
		},
		{
			name:   "variable-1",
			raw:    `@{name} = torbrowser "tor browser"`,
			tokens: []string{"@{name}", "=", "torbrowser", `"tor browser"`},
			rule: rule{
				{key: "@{name}"}, {key: "="}, {key: "torbrowser"}, {key: `"tor browser"`},
			},
			getIdx:     2,
			wGet:       "torbrowser",
			wGetString: `@{name} = torbrowser "tor browser"`,
			wGetSlice:  []string{"@{name}", "=", "torbrowser", `"tor browser"`},
			wString:    `@{name} = torbrowser "tor browser"`,
		},
		{
			name:   "variable-2",
			raw:    `@{exec_path} += @{bin}/@{name}`,
			tokens: []string{"@{exec_path}", "+=", "@{bin}/@{name}"},
			rule: rule{
				{key: "@{exec_path}"}, {key: "+="}, {key: "@{bin}/@{name}"},
			},
			getIdx:     1,
			wGet:       "+=",
			wGetString: `@{exec_path} += @{bin}/@{name}`,
			wGetSlice:  []string{"@{exec_path}", "+=", "@{bin}/@{name}"},
			wString:    `@{exec_path} += @{bin}/@{name}`,
		},
		{
			name:   "variable-3",
			raw:    `@{empty}="dummy"`,
			tokens: []string{"@{empty}", "=", `"dummy"`},
			rule: rule{
				{key: "@{empty}"}, {key: "="}, {key: `"dummy"`},
			},
			getIdx:     1,
			wGet:       "=",
			wGetString: `@{empty} = "dummy"`,
			wGetSlice:  []string{"@{empty}", "=", `"dummy"`},
			wString:    `@{empty} = "dummy"`,
		},
		{
			name:   "variable-4",
			raw:    `@{XDG_PROJECTS_DIR}+="Git"`,
			tokens: []string{"@{XDG_PROJECTS_DIR}", "+=", `"Git"`},
			rule: rule{
				{key: "@{XDG_PROJECTS_DIR}"}, {key: "+="}, kv{key: "\"Git\""},
			},
			getIdx:     1,
			wGet:       "+=",
			wGetString: `@{XDG_PROJECTS_DIR} += "Git"`,
			wGetSlice:  []string{"@{XDG_PROJECTS_DIR}", "+=", `"Git"`},
			wString:    `@{XDG_PROJECTS_DIR} += "Git"`,
		},
		{
			name:   "header",
			raw:    `profile foo @{exec_path} xattrs=(security.tagged=allowed) flags=(complain attach_disconnected)`,
			tokens: []string{"profile", "foo", "@{exec_path}", "xattrs=(security.tagged=allowed)", "flags=(complain attach_disconnected)"},
			rule: rule{
				{key: "profile"},
				{key: "foo"},
				{key: "@{exec_path}"},
				{key: "xattrs", values: rule{
					{key: "security.tagged", values: rule{{key: "allowed"}}},
				}},
				{key: "flags", values: rule{
					{key: "complain"},
					{key: "attach_disconnected"},
				}},
			},
			getIdx:     2,
			getKey:     "flags",
			wGet:       "@{exec_path}",
			wGetString: "profile foo @{exec_path}",
			wGetSlice:  []string{"profile", "foo", "@{exec_path}"},
			wGetAsMap: map[string][]string{
				"flags":  {"complain", "attach_disconnected"},
				"xattrs": {},
			},
			wGetValues:         rule{{key: "complain"}, {key: "attach_disconnected"}},
			wGetValuesAsSlice:  []string{"complain", "attach_disconnected"},
			wGetValuesAsString: "complain attach_disconnected",
			wString:            "profile foo @{exec_path} xattrs=(security.tagged=allowed) flags=(complain attach_disconnected)",
		},
		{
			name:   "include",
			raw:    `include <tunables/global>`,
			tokens: []string{"include", "<tunables/global>"},
			rule: rule{
				{key: "include"}, {key: "<tunables/global>"},
			},
			getIdx:     1,
			wGet:       "<tunables/global>",
			wGetString: `include <tunables/global>`,
			wGetSlice:  []string{"include", "<tunables/global>"},
			wString:    `include <tunables/global>`,
			wRule:      &Include{IfExists: false, IsMagic: true, Path: "tunables/global"},
			wError:     true, // newRules only convert comma rules
		},
		{
			name:   "include-if-exists",
			raw:    `include if exists "/etc/apparmor.d/dummy"`,
			tokens: []string{"include", "if", "exists", `"/etc/apparmor.d/dummy"`},
			rule: rule{
				{key: "include"}, {key: "if"}, {key: "exists"}, {key: `"/etc/apparmor.d/dummy"`},
			},
			getIdx:     1,
			wGet:       "if",
			wGetString: `include if exists "/etc/apparmor.d/dummy"`,
			wGetSlice:  []string{"include", "if", "exists", `"/etc/apparmor.d/dummy"`},
			wString:    `include if exists "/etc/apparmor.d/dummy"`,
			wRule:      &Include{IfExists: true, IsMagic: false, Path: "/etc/apparmor.d/dummy"},
			wError:     true, // newRules only convert comma rules
		},
		{
			name:   "rlimit",
			raw:    `set rlimit nproc <= 200`,
			tokens: []string{"set", "rlimit", "nproc", "<=", "200"},
			rule: rule{
				{key: "set"}, {key: "rlimit"}, {key: "nproc"}, {key: "<="}, {key: "200"},
			},
			getIdx:     4,
			wGet:       "200",
			wGetString: `set rlimit nproc <= 200`,
			wGetSlice:  []string{"set", "rlimit", "nproc", "<=", "200"},
			wString:    `set rlimit nproc <= 200`,
			wRule:      &Rlimit{Key: "nproc", Op: "<=", Value: "200"},
		},
		{
			name:       "userns",
			raw:        `userns`,
			tokens:     []string{"userns"},
			rule:       rule{{key: "userns"}},
			wGetString: `userns`,
			wGetSlice:  []string{"userns"},
			wString:    `userns`,
			wRule:      &Userns{Create: true},
		},
		{
			name:   "capability",
			raw:    `capability dac_read_search`,
			tokens: []string{"capability", "dac_read_search"},
			rule: rule{
				{key: "capability"}, {key: "dac_read_search"},
			},
			getIdx:     1,
			wGet:       "dac_read_search",
			wGetString: `capability dac_read_search`,
			wGetSlice:  []string{"capability", "dac_read_search"},
			wString:    `capability dac_read_search`,
			wRule:      &Capability{Names: []string{"dac_read_search"}},
		},
		{
			name:   "network",
			raw:    `network netlink raw`,
			tokens: []string{"network", "netlink", "raw"},
			rule: rule{
				{key: "network"}, {key: "netlink"}, {key: "raw"},
			},
			getIdx:     1,
			wGet:       "netlink",
			wGetString: `network netlink raw`,
			wGetSlice:  []string{"network", "netlink", "raw"},
			wString:    `network netlink raw`,
			wRule:      &Network{Domain: "netlink", Type: "raw"},
		},
		{
			name:   "mount",
			raw:    `mount /{,**}`,
			tokens: []string{"mount", "/{,**}"},
			rule: rule{
				{key: "mount"}, {key: "/{,**}"},
			},
			getIdx:     1,
			wGet:       "/{,**}",
			wGetString: `mount /{,**}`,
			wGetSlice:  []string{"mount", "/{,**}"},
			wString:    `mount /{,**}`,
			wRule: &Mount{
				MountConditions: MountConditions{Options: []string{}},
				Source:          "/{,**}",
			},
		},
		{
			name:   "mount-2",
			raw:    `mount               options=(rw rbind)                 /tmp/newroot/ -> /tmp/newroot/`,
			tokens: []string{"mount", "options=(rw rbind)", "/tmp/newroot/", "->", "/tmp/newroot/"},
			rule: rule{
				{key: "mount"},
				{key: "options", values: rule{{key: "rw"}, {key: "rbind"}}},
				{key: "/tmp/newroot/"}, {key: "->"}, {key: "/tmp/newroot/"}},
			getIdx:             1,
			getKey:             "options",
			wGet:               "options",
			wGetString:         "mount /tmp/newroot/ -> /tmp/newroot/",
			wGetSlice:          []string{"mount", "/tmp/newroot/", "->", "/tmp/newroot/"},
			wGetAsMap:          map[string][]string{"options": {"rw", "rbind"}},
			wGetValues:         rule{{key: "rw"}, {key: "rbind"}},
			wGetValuesAsSlice:  []string{"rw", "rbind"},
			wGetValuesAsString: "rw rbind",
			wString:            "mount options=(rw rbind) /tmp/newroot/ -> /tmp/newroot/",
			wRule: &Mount{
				MountConditions: MountConditions{Options: []string{"rw", "rbind"}},
				Source:          "/tmp/newroot/",
				MountPoint:      "/tmp/newroot/",
			},
		},
		{
			name:   "mount-3",
			raw:    `mount               options=(rw silent rprivate)                     -> /oldroot/`,
			tokens: []string{"mount", "options=(rw silent rprivate)", "->", "/oldroot/"},
			rule: rule{
				{key: "mount"},
				{key: "options", values: rule{{key: "rw"}, {key: "silent"}, {key: "rprivate"}}},
				{key: "->"}, {key: "/oldroot/"},
			},
			getIdx:             3,
			getKey:             "options",
			wGet:               "/oldroot/",
			wGetString:         "mount -> /oldroot/",
			wGetSlice:          []string{"mount", "->", "/oldroot/"},
			wGetAsMap:          map[string][]string{"options": {"rw", "silent", "rprivate"}},
			wGetValues:         rule{{key: "rw"}, {key: "silent"}, {key: "rprivate"}},
			wGetValuesAsSlice:  []string{"rw", "silent", "rprivate"},
			wGetValuesAsString: "rw silent rprivate",
			wString:            "mount options=(rw silent rprivate) -> /oldroot/",
			wRule: &Mount{
				MountConditions: MountConditions{Options: []string{"rw", "rprivate", "silent"}},
				MountPoint:      "/oldroot/",
			},
		},
		{
			name:   "mount-4",
			raw:    `mount fstype=devpts options=(rw nosuid noexec)                devpts -> /newroot/dev/pts/`,
			tokens: []string{"mount", "fstype=devpts", "options=(rw nosuid noexec)", "devpts", "->", "/newroot/dev/pts/"},
			rule: rule{
				{key: "mount"},
				{key: "fstype", values: rule{{key: "devpts"}}},
				{key: "options", values: rule{{key: "rw"}, {key: "nosuid"}, {key: "noexec"}}},
				{key: "devpts"}, {key: "->"}, {key: "/newroot/dev/pts/"},
			},
			getIdx:     1,
			getKey:     "fstype",
			wGet:       "fstype",
			wGetString: "mount devpts -> /newroot/dev/pts/",
			wGetSlice:  []string{"mount", "devpts", "->", "/newroot/dev/pts/"},
			wGetAsMap: map[string][]string{
				"fstype":  {"devpts"},
				"options": {"rw", "nosuid", "noexec"},
			},
			wGetValues:         rule{{key: "devpts"}},
			wGetValuesAsSlice:  []string{"devpts"},
			wGetValuesAsString: "devpts",
			wString:            "mount fstype=devpts options=(rw nosuid noexec) devpts -> /newroot/dev/pts/",
			wRule: &Mount{
				MountConditions: MountConditions{
					FsType:  "devpts",
					Options: []string{"rw", "noexec", "nosuid"},
				},
				Source:     "devpts",
				MountPoint: "/newroot/dev/pts/",
			},
		},
		{
			name:   "signal",
			raw:    `signal (receive) set=(cont, term,winch) peer=at-spi-bus-launcher`,
			tokens: []string{"signal", "(receive)", "set=(cont, term,winch)", "peer=at-spi-bus-launcher"},
			rule: rule{
				{key: "signal"}, {key: "receive"},
				{key: "set", values: rule{{key: "cont"}, {key: "term"}, {key: "winch"}}},
				{key: "peer", values: rule{{key: "at-spi-bus-launcher"}}},
			},
			getIdx:     1,
			getKey:     "peer",
			wGet:       "receive",
			wGetString: "signal receive",
			wGetSlice:  []string{"signal", "receive"},
			wGetAsMap: map[string][]string{
				"peer": {"at-spi-bus-launcher"},
				"set":  {"cont", "term", "winch"},
			},
			wGetValues:         rule{{key: "at-spi-bus-launcher"}},
			wGetValuesAsSlice:  []string{"at-spi-bus-launcher"},
			wGetValuesAsString: "at-spi-bus-launcher",
			wString:            "signal receive set=(cont term winch) peer=at-spi-bus-launcher",
			wRule: &Signal{
				Access: []string{"receive"},
				Set:    []string{"cont", "term", "winch"},
				Peer:   "at-spi-bus-launcher",
			},
		},
		{
			name:   "ptrace",
			raw:    "ptrace peer=/bin/true peer=/sbin/init peer=MY_PROFILE",
			tokens: []string{"ptrace", "peer=/bin/true", "peer=/sbin/init", "peer=MY_PROFILE"},
			rule: rule{
				{key: "ptrace"},
				{key: "peer", values: rule{{key: "/bin/true"}}},
				{key: "peer", values: rule{{key: "/sbin/init"}}},
				{key: "peer", values: rule{{key: "MY_PROFILE"}}},
			},
			getIdx:     3,
			getKey:     "peer",
			wGet:       "peer",
			wGetString: "ptrace",
			wGetSlice:  []string{"ptrace"},
			wGetAsMap: map[string][]string{
				"peer": {"/bin/true", "/sbin/init", "MY_PROFILE"},
			},
			wGetValues:         rule{{key: "/bin/true"}, {key: "/sbin/init"}, {key: "MY_PROFILE"}},
			wGetValuesAsSlice:  []string{"/bin/true", "/sbin/init", "MY_PROFILE"},
			wGetValuesAsString: "/bin/true /sbin/init MY_PROFILE",
			wString:            "ptrace peer=/bin/true peer=/sbin/init peer=MY_PROFILE",
			wRule: &Ptrace{
				Peer: "/bin/true",
			},
			wError: true,
		},
		{
			name:   "unix-1",
			raw:    `unix (send receive) type=stream addr="@/tmp/.ICE[0-9]*-unix/19 5" peer=(label="@{p_systemd}", addr=none)`,
			tokens: []string{"unix", "(send receive)", "type=stream", "addr=\"@/tmp/.ICE[0-9]*-unix/19 5\"", "peer=(label=\"@{p_systemd}\", addr=none)"},
			rule: rule{
				{key: "unix"}, {key: "send"}, {key: "receive"},
				{key: "type", values: rule{{key: "stream"}}},
				{key: "addr", values: rule{
					{key: `"@/tmp/.ICE[0-9]*-unix/19 5"`},
				}},
				{key: "peer", values: rule{
					{key: "label", values: rule{{key: `"@{p_systemd}"`}}},
					{key: "addr", values: rule{{key: "none"}}},
				}},
			},
			getIdx:     3,
			getKey:     "peer",
			wGet:       "type",
			wGetString: "unix send receive",
			wGetSlice:  []string{"unix", "send", "receive"},
			wGetAsMap: map[string][]string{
				"addr": {`"@/tmp/.ICE[0-9]*-unix/19 5"`},
				"peer": {},
				"type": {"stream"},
			},
			wGetValues: rule{
				{key: "label", values: rule{{key: `"@{p_systemd}"`}}},
				{key: "addr", values: rule{{key: "none"}}},
			},
			wGetValuesAsSlice:  []string{},
			wGetValuesAsString: "",
			wString:            `unix send receive type=stream addr=("@/tmp/.ICE[0-9]*-unix/19 5") peer=(label="@{p_systemd}" addr=none)`,
			wRule: &Unix{
				Access:    []string{"send", "receive"},
				Type:      "stream",
				Address:   `"@/tmp/.ICE[0-9]*-unix/19 5"`,
				PeerLabel: `"@{p_systemd}"`,
				PeerAddr:  "none",
			},
		},
		{
			name: "unix-2",
			raw: `  unix (connect, receive, send)
			type=stream
			peer=(addr="@/tmp/ibus/dbus-????????")`,
			tokens: []string{"unix", "(connect, receive, send)", "type=stream", `peer=(addr="@/tmp/ibus/dbus-????????")`},
			rule: rule{
				{key: "unix"}, {key: "connect"}, {key: "receive"}, {key: "send"},
				{key: "type", values: rule{{key: "stream"}}},
				{key: "peer", values: rule{
					{key: "addr", values: rule{{key: `"@/tmp/ibus/dbus-????????"`}}},
				}},
			},
			getIdx:     4,
			getKey:     "type",
			wGet:       "type",
			wGetString: "unix connect receive send",
			wGetSlice:  []string{"unix", "connect", "receive", "send"},
			wGetAsMap: map[string][]string{
				"type": {"stream"},
				"peer": {},
			},
			wGetValues:         rule{{key: "stream"}},
			wGetValuesAsSlice:  []string{"stream"},
			wGetValuesAsString: "stream",
			wString:            `unix connect receive send type=stream peer=(addr="@/tmp/ibus/dbus-????????")`,
			wRule: &Unix{
				Access:   []string{"connect", "receive", "send"},
				Type:     "stream",
				PeerAddr: `"@/tmp/ibus/dbus-????????"`,
			},
		},
		{
			name: "dbus",
			raw:  `dbus receive bus=system path=/org/freedesktop/DBus interface=org.freedesktop.DBus member=AddMatch peer=(name=:1.3, label=power-profiles-daemon)`,
			tokens: []string{
				"dbus", "receive", "bus=system",
				"path=/org/freedesktop/DBus", "interface=org.freedesktop.DBus",
				"member=AddMatch", "peer=(name=:1.3, label=power-profiles-daemon)",
			},
			rule: rule{
				{key: "dbus"}, {key: "receive"},
				{key: "bus", values: rule{{key: "system"}}},
				{key: "path", values: rule{{key: "/org/freedesktop/DBus"}}},
				{key: "interface", values: rule{{key: "org.freedesktop.DBus"}}},
				{key: "member", values: rule{{key: "AddMatch"}}},
				{key: "peer", values: rule{
					{key: "name", values: rule{{key: ":1.3"}}},
					{key: "label", values: rule{{key: "power-profiles-daemon"}}},
				}},
			},
			getIdx:     2,
			getKey:     "path",
			wGet:       "bus",
			wGetString: "dbus receive",
			wGetSlice:  []string{"dbus", "receive"},
			wGetAsMap: map[string][]string{
				"bus":       {"system"},
				"interface": {"org.freedesktop.DBus"},
				"member":    {"AddMatch"},
				"path":      {"/org/freedesktop/DBus"},
				"peer":      {},
			},
			wGetValues:         rule{{key: "/org/freedesktop/DBus"}},
			wGetValuesAsSlice:  []string{"/org/freedesktop/DBus"},
			wGetValuesAsString: "/org/freedesktop/DBus",
			wString:            `dbus receive bus=system path=/org/freedesktop/DBus interface=org.freedesktop.DBus member=AddMatch peer=(name=:1.3 label=power-profiles-daemon)`,
			wRule: &Dbus{
				Access:    []string{"receive"},
				Bus:       "system",
				Path:      "/org/freedesktop/DBus",
				Interface: "org.freedesktop.DBus",
				Member:    "AddMatch",
				PeerName:  ":1.3",
				PeerLabel: "power-profiles-daemon",
			},
		},
		{
			name:   "file-1",
			raw:    `owner @{user_config_dirs}/powerdevilrc{,.@{rand6}} rwl -> @{user_config_dirs}/#@{int}`,
			tokens: []string{"owner", "@{user_config_dirs}/powerdevilrc{,.@{rand6}}", "rwl", "->", "@{user_config_dirs}/#@{int}"},
			rule: rule{
				{key: "owner"},
				{key: "@{user_config_dirs}/powerdevilrc{,.@{rand6}}"},
				{key: "rwl"},
				{key: "->"},
				{key: "@{user_config_dirs}/#@{int}"},
			},
			getIdx:     3,
			wGet:       "->",
			wGetString: "owner @{user_config_dirs}/powerdevilrc{,.@{rand6}} rwl -> @{user_config_dirs}/#@{int}",
			wGetSlice:  []string{"owner", "@{user_config_dirs}/powerdevilrc{,.@{rand6}}", "rwl", "->", "@{user_config_dirs}/#@{int}"},
			wString:    `owner @{user_config_dirs}/powerdevilrc{,.@{rand6}} rwl -> @{user_config_dirs}/#@{int}`,
			wRule: &File{
				Owner:  true,
				Path:   "@{user_config_dirs}/powerdevilrc{,.@{rand6}}",
				Access: []string{"r", "w", "l"},
				Target: "@{user_config_dirs}/#@{int}",
			},
		},
		{
			name:   "file-2",
			raw:    `@{sys}/devices/@{pci}/class r`,
			tokens: []string{"@{sys}/devices/@{pci}/class", "r"},
			rule: rule{
				{key: "@{sys}/devices/@{pci}/class"}, {key: "r"},
			},
			getIdx:     1,
			wGet:       "r",
			wGetString: `@{sys}/devices/@{pci}/class r`,
			wGetSlice:  []string{"@{sys}/devices/@{pci}/class", "r"},
			wString:    `@{sys}/devices/@{pci}/class r`,
			wRule:      &File{Path: "@{sys}/devices/@{pci}/class", Access: []string{"r"}},
		},
		{
			name:   "file-3",
			raw:    `owner @{PROC}/@{pid}/task/@{tid}/comm rw`,
			tokens: []string{"owner", "@{PROC}/@{pid}/task/@{tid}/comm", "rw"},
			rule: rule{
				{key: "owner"}, {key: "@{PROC}/@{pid}/task/@{tid}/comm"}, {key: "rw"},
			},
			getIdx:     1,
			wGet:       "@{PROC}/@{pid}/task/@{tid}/comm",
			wGetString: `owner @{PROC}/@{pid}/task/@{tid}/comm rw`,
			wGetSlice:  []string{"owner", "@{PROC}/@{pid}/task/@{tid}/comm", "rw"},
			wString:    `owner @{PROC}/@{pid}/task/@{tid}/comm rw`,
			wRule: &File{
				Owner:  true,
				Path:   "@{PROC}/@{pid}/task/@{tid}/comm",
				Access: []string{"r", "w"},
			},
		},
		{
			name:   "file-4",
			raw:    `owner /{var/,}tmp/#@{int} rw`,
			tokens: []string{"owner", "/{var/,}tmp/#@{int}", "rw"},
			rule: rule{
				{key: "owner"}, {key: "/{var/,}tmp/#@{int}"}, {key: "rw"},
			},
			getIdx:     2,
			wGet:       "rw",
			wGetString: `owner /{var/,}tmp/#@{int} rw`,
			wGetSlice:  []string{"owner", "/{var/,}tmp/#@{int}", "rw"},
			wString:    `owner /{var/,}tmp/#@{int} rw`,
			wRule: &File{
				Owner:  true,
				Path:   "/{var/,}tmp/#@{int}",
				Access: []string{"r", "w"},
			},
		},
		{
			name:   "file-5",
			raw:    `owner @{run}/user/@{uid}/gvfs/smb-share:server=*,share=**/ r`,
			tokens: []string{"owner", "@{run}/user/@{uid}/gvfs/smb-share:server=*,share=**/", "r"},
			rule: rule{
				kv{key: "owner"},
				kv{key: "@{run}/user/@{uid}/gvfs/smb-share:server=*,share=**/"},
				kv{key: "r"},
			},
			getIdx:     0,
			wGet:       "owner",
			wGetString: `owner @{run}/user/@{uid}/gvfs/smb-share:server=*,share=**/ r`,
			wGetSlice:  []string{"owner", "@{run}/user/@{uid}/gvfs/smb-share:server=*,share=**/", "r"},
			wString:    `owner @{run}/user/@{uid}/gvfs/smb-share:server=*,share=**/ r`,
			wRule: &File{
				Owner:  true,
				Path:   "@{run}/user/@{uid}/gvfs/smb-share:server=*,share=**/",
				Access: []string{"r"},
			},
		},
	}

	// Test cases for parseLineRules
	testsLineRules = []struct {
		name       string
		isPreamble bool
		raw        string
		want       string
		rules      Rules
		wantErr    bool
	}{
		{
			name:       "strin.aa",
			isPreamble: true,
			raw: `alias /mnt/usr -> /usr,
				include <tunables/global>
				include if exists "/etc/apparmor.d/global/dummy space"
				@{name} = torbrowser "tor browser"
				alias /mnt/{,usr.sbin.}mount.cifs
						->
						/sbin/mount.cifs,`,
			want: `alias /mnt/usr -> /usr,



				alias /mnt/{,usr.sbin.}mount.cifs
						->
						/sbin/mount.cifs,`,
			rules: Rules{
				&Include{IsMagic: true, Path: "tunables/global"},
				&Include{IfExists: true, Path: "/etc/apparmor.d/global/dummy space"},
				&Variable{Name: "name", Define: true, Values: []string{"torbrowser", "\"tor browser\""}},
			},
			wantErr: false,
		},
		{
			name:       "comment",
			isPreamble: true,
			raw: `
			# IsLineRule comment
			include <tunables/global> # comment included
			@{lib_dirs} = @{lib}/@{name} /opt/@{name} # comment in variable`,
			want: "\n\n\n",
			rules: Rules{
				&Comment{Base: Base{IsLineRule: true, Comment: " IsLineRule comment"}},
				&Include{
					Base:    Base{Comment: " comment included"},
					IsMagic: true, Path: "tunables/global",
				},
				&Variable{
					Base: Base{Comment: " comment in variable"},
					Name: "lib_dirs", Define: true,
					Values: []string{"@{lib}/@{name}", "/opt/@{name}"},
				},
			},
			wantErr: false,
		},
	}

	// Test cases for parseCommaRules
	testsCommaRules = []struct {
		name    string
		raw     string
		rule    []rule
		wantErr bool
	}{
		{
			name: "string.aa",
			raw: `alias /mnt/{,usr.sbin.}mount.cifs
			->
			/sbin/mount.cifs,
			network inet stream, # inline comment
	  		dbus bind bus=session name=org.mpris.MediaPlayer2,
			owner @{user_config_dirs}/powerdevilrc{,.@{rand6}} rwl -> @{user_config_dirs}/#@{int},
			@{sys}/class/ r,
			@{run}/udev/data/+pci:*  r,`,
			rule: []rule{
				{
					kv{key: "alias"},
					kv{key: "/mnt/{,usr.sbin.}mount.cifs"},
					kv{key: "->"},
					kv{key: "/sbin/mount.cifs"},
				},
				{
					kv{key: "network"}, kv{key: "inet"}, kv{key: "stream", comment: " inline comment"},
				},
				{
					kv{key: "dbus"}, kv{key: "bind"},
					kv{key: "bus", values: rule{kv{key: "session"}}},
					kv{key: "name", values: rule{kv{key: "org.mpris.MediaPlayer2"}}},
				},
				{
					kv{key: "owner"}, kv{key: "@{user_config_dirs}/powerdevilrc{,.@{rand6}}"}, kv{key: "rwl"},
					kv{key: "->"}, kv{key: "@{user_config_dirs}/#@{int}"},
				},
				{
					kv{key: "@{sys}/class/"}, kv{key: "r"},
				},
				{
					kv{key: "@{run}/udev/data/+pci:*"}, kv{key: "r"},
				},
			},
			wantErr: false,
		},
		{
			name: ",-in-aare",
			raw:  `@{sys}/cgroup/cpu,cpuacct/user.slice/cpu.cfs_quota_us r,`,
			rule: []rule{
				{
					kv{key: "@{sys}/cgroup/cpu,cpuacct/user.slice/cpu.cfs_quota_us"}, kv{key: "r"},
				},
			},
			wantErr: false,
		},
		{
			name: "#-in-aare",
			raw:  `owner @{user_config_dirs}/#** rw,`,
			rule: []rule{
				{
					kv{key: "owner"}, kv{key: "@{user_config_dirs}/#**"}, kv{key: "rw"},
				},
			},
			wantErr: false,
		},
	}

	// Test cases for tokenizeBlock, parseBlock, and Parse
	testBlocks = []struct {
		name           string
		raw            string
		blocks         []*block
		wTokenizeErr   bool
		rules          []Rules
		wParseBlockErr bool
		apparmor       *AppArmorProfileFile
		apparmorAll    *AppArmorProfileFile
		wParseErr      bool
		wRules         ParaRules
		wParseRulesErr bool
	}{
		{
			name:           "empty",
			raw:            "",
			blocks:         []*block{{kind: CONTENT, raw: ""}},
			wTokenizeErr:   false,
			rules:          []Rules{nil},
			wParseBlockErr: false,
			apparmor:       &AppArmorProfileFile{},
			wParseErr:      false,
			wRules:         ParaRules{},
			wParseRulesErr: false,
		},
		{
			name: "comment",
			raw: `
			# IsLineRule comment
			include <tunables/global> # comment included
			@{lib_dirs} = @{lib}/@{name} /opt/@{name} # comment in variable`,
			blocks: []*block{
				{
					kind: "content",
					raw:  "\n\t\t\t# IsLineRule comment\n\t\t\tinclude <tunables/global> # comment included\n\t\t\t@{lib_dirs} = @{lib}/@{name} /opt/@{name} # comment in variable",
				},
			},
			wTokenizeErr: false,
			rules: []Rules{
				{
					&Comment{Base: Base{IsLineRule: true, Comment: " IsLineRule comment"}},
					&Include{
						Base: Base{Comment: " comment included"},
						Path: "tunables/global", IsMagic: true,
					},
				},
			},
			wParseBlockErr: false,
			apparmor: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{Base: Base{IsLineRule: true, Comment: " IsLineRule comment"}},
					&Include{
						Base: Base{Comment: " comment included"},
						Path: "tunables/global", IsMagic: true,
					},
					&Variable{
						Base: Base{Comment: " comment in variable"},
						Name: "lib_dirs", Define: true,
						Values: []string{"@{lib}/@{name}", "/opt/@{name}"},
					},
				},
			},
			wParseErr:      false,
			wRules:         ParaRules{},
			wParseRulesErr: false,
		},
		{
			name: "cornercases",
			raw: `# Simple test
			include <tunables/global>

			# { commented block }
			@{name} = {D,d}ummy
			@{exec_path} = @{bin}/@{name}
			@{exec_path} += @{lib}/@{name}
			alias /mnt/{,usr.sbin.}mount.cifs -> /sbin/mount.cifs,
			@{coreutils} += gawk {,e,f}grep head
			profile @{exec_path} {
			}
			`,
			blocks: []*block{
				{
					kind: "content",
					raw:  "# Simple test\n\t\t\tinclude <tunables/global>\n\n\t\t\t# { commented block }\n\t\t\t@{name} = {D,d}ummy\n\t\t\t@{exec_path} = @{bin}/@{name}\n\t\t\t@{exec_path} += @{lib}/@{name}\n\t\t\talias /mnt/{,usr.sbin.}mount.cifs -> /sbin/mount.cifs,\n\t\t\t@{coreutils} += gawk {,e,f}grep head",
				},
				{
					kind: "profile",
					raw:  "profile @{exec_path}",
					next: &block{
						kind: "raw",
						raw:  "",
					},
				},
			},
			wTokenizeErr: false,
			rules: []Rules{
				{
					&Comment{Base: Base{IsLineRule: true, Comment: " Simple test"}},
					&Include{IsMagic: true, Path: "tunables/global"},
					&Comment{Base: Base{IsLineRule: true, Comment: " { commented block }"}},
					&Variable{Name: "name", Values: []string{"{D,d}ummy"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"@{bin}/@{name}"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"@{lib}/@{name}"}, Define: false},
					&Variable{Name: "coreutils", Values: []string{"gawk", "{,e,f}grep", "head"}, Define: false},
					&Alias{Path: "/mnt/{,usr.sbin.}mount.cifs", RewrittenPath: "/sbin/mount.cifs"},
				},
				{
					&Profile{
						Header: Header{
							Name:        "@{exec_path}",
							Attachments: []string{},
							Attributes:  map[string]string{},
							Flags:       []string{},
						},
					},
				},
			},
			wParseBlockErr: false,
			apparmor: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{Base: Base{IsLineRule: true, Comment: " Simple test"}},
					&Include{IsMagic: true, Path: "tunables/global"},
					&Comment{Base: Base{IsLineRule: true, Comment: " { commented block }"}},
					&Variable{Name: "name", Values: []string{"{D,d}ummy"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"@{bin}/@{name}"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"@{lib}/@{name}"}},
					&Variable{Name: "coreutils", Values: []string{"gawk", "{,e,f}grep", "head"}},
					&Alias{Path: "/mnt/{,usr.sbin.}mount.cifs", RewrittenPath: "/sbin/mount.cifs"},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "@{exec_path}",
							Attachments: []string{},
							Attributes:  map[string]string{},
							Flags:       []string{},
						},
					},
				},
			},
			wParseErr:      false,
			wRules:         ParaRules{},
			wParseRulesErr: false,
		},
		{
			name: "condition-1",
			raw: `
			$FOO=true
			$BAR = False

			/bin/true {
			  /bin/false rix,
			  if ${FOO} {
			    /bin/true rix,
			  }
			  /bin/true rix,
			  if ${BAR} {
			    /etc/shadow rw,
			  }
			  /bin/sh rix,
			}`,
			blocks: []*block{
				{
					kind: CONTENT,
					raw:  "\n\t\t\t$FOO=true\n\t\t\t$BAR = False\n",
				},
				{
					kind: PROFILE,
					raw:  "/bin/true",
					next: &block{
						kind: RAW,
						raw:  "/bin/false rix,\n\t\t\t  if ${FOO} {\n\t\t\t    /bin/true rix,\n\t\t\t  }\n\t\t\t  /bin/true rix,\n\t\t\t  if ${BAR} {\n\t\t\t    /etc/shadow rw,\n\t\t\t  }\n\t\t\t  /bin/sh rix,",
					},
				},
			},
			wTokenizeErr: false,
			rules: []Rules{
				{
					&Boolean{Name: "FOO", Value: true},
					&Boolean{Name: "BAR", Value: false},
				},
				{
					&Profile{
						Header: Header{
							Name:        "/bin/true",
							Attachments: []string{},
							Attributes:  map[string]string{},
							Flags:       []string{},
						},
						Rules: Rules{
							&File{Path: "/bin/false", Access: []string{"r", "ix"}},
							&Condition{
								Expression: "${FOO}",
								IfRules: Rules{
									&File{Path: "/bin/true", Access: []string{"r", "ix"}},
								},
							},
							&File{Path: "/bin/true", Access: []string{"r", "ix"}},
							&Condition{
								Expression: "${BAR}",
								IfRules: Rules{
									&File{Path: "/etc/shadow", Access: []string{"r", "w"}},
								},
							},
							&File{Path: "/bin/sh", Access: []string{"r", "ix"}},
						},
					},
				},
			},
			wParseBlockErr: false,
			apparmor: &AppArmorProfileFile{
				Preamble: Rules{
					&Boolean{Name: "FOO", Value: true},
					&Boolean{Name: "BAR", Value: false},
				},
			},
			apparmorAll: &AppArmorProfileFile{
				Preamble: Rules{
					&Boolean{Name: "FOO", Value: true},
					&Boolean{Name: "BAR", Value: false},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "/bin/true",
							Attachments: []string{},
							Attributes:  map[string]string{},
							Flags:       []string{},
						},
						Rules: Rules{
							&File{Path: "/bin/false", Access: []string{"r", "ix"}},
							&Condition{
								Expression: "${FOO}",
								IfRules: Rules{
									&File{Path: "/bin/true", Access: []string{"r", "ix"}},
								},
							},
							&File{Path: "/bin/true", Access: []string{"r", "ix"}},
							&Condition{
								Expression: "${BAR}",
								IfRules: Rules{
									&File{Path: "/etc/shadow", Access: []string{"r", "w"}},
								},
							},
							&File{Path: "/bin/sh", Access: []string{"r", "ix"}},
						},
					},
				},
			},
			wParseErr:      false,
			wRules:         ParaRules{nil},
			wParseRulesErr: false,
		},
		{
			name: "string.aa",
			raw:  testData.Join("string.aa").MustReadFileAsString(),
			blocks: []*block{
				{
					kind: CONTENT,
					raw:  "# Simple test profile for the AppArmorProfileFile.String() method\n\nabi <abi/4.0>,\n\nalias /mnt/usr -> /usr,\n\ninclude <tunables/global>\n\n@{exec_path} = @{bin}/foo @{lib}/foo",
				},
				{
					kind: PROFILE,
					raw:  "profile foo @{exec_path} xattrs=(security.tagged=allowed) flags=(complain attach_disconnected)",
					next: &block{
						kind: RAW,
						raw:  pStringAAContentStr,
					},
				},
			},
			wTokenizeErr: false,
			rules: []Rules{
				{
					&Comment{Base: Base{IsLineRule: true, Comment: " Simple test profile for the AppArmorProfileFile.String() method"}},
					&Include{Path: "tunables/global", IsMagic: true},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{"@{bin}/foo", "@{lib}/foo"},
					},
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
				},
				{
					&Profile{
						Header: Header{
							Name:        "foo",
							Attachments: []string{"@{exec_path}"},
							Attributes:  map[string]string{"security.tagged": "allowed"},
							Flags:       []string{"complain", "attach_disconnected"},
						},
						Rules: rulesStringAA,
					},
				},
			},
			wParseBlockErr: false,
			apparmor: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{Base: Base{Comment: " Simple test profile for the AppArmorProfileFile.String() method", IsLineRule: true}},
					&Include{IsMagic: true, Path: "tunables/global"},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{"@{bin}/foo", "@{lib}/foo"},
					},
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "foo",
							Attachments: []string{"@{exec_path}"},
							Attributes:  map[string]string{"security.tagged": "allowed"},
							Flags:       []string{"complain", "attach_disconnected"},
						},
					},
				},
			},
			apparmorAll: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{Base: Base{Comment: " Simple test profile for the AppArmorProfileFile.String() method", IsLineRule: true}},
					&Include{IsMagic: true, Path: "tunables/global"},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{"@{bin}/foo", "@{lib}/foo"},
					},
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "foo",
							Attachments: []string{"@{exec_path}"},
							Attributes:  map[string]string{"security.tagged": "allowed"},
							Flags:       []string{"complain", "attach_disconnected"},
						},
						Rules: rulesStringAA,
					},
				},
			},
			wParseErr: false,
			wRules: ParaRules{
				{
					&Include{IsMagic: true, Path: "abstractions/base"},
					&Include{IsMagic: true, Path: "abstractions/nameservice-strict"},
				},
				{
					&Rlimit{Key: "nproc", Op: "<=", Value: "200"},
				},
				{
					&Capability{Names: []string{"dac_read_search"}},
					&Capability{Names: []string{"dac_override"}},
				},
				{
					&Network{Domain: "inet", Type: "stream"},
					&Network{Domain: "inet6", Type: "stream"},
				},
				{
					&Mount{
						Base: Base{IsLineRule: false, Comment: " failed perms check"},
						MountConditions: MountConditions{
							FsType:  "fuse.portal",
							Options: []string{"rw", "rbind"},
						},
						Source:     "@{run}/user/@{uid}/",
						MountPoint: "/",
					},
				},
				{
					&Umount{
						MountConditions: MountConditions{Options: []string{}},
						MountPoint:      "@{run}/user/@{uid}/",
					},
				},
				{
					&Signal{
						Access: []string{"receive"},
						Set:    []string{"term"},
						Peer:   "at-spi-bus-launcher",
					},
				},
				{
					&Ptrace{Access: []string{"read"}, Peer: "nautilus"},
				},
				{
					&Unix{
						Access:    []string{"send", "receive"},
						Type:      "stream",
						Address:   "@/tmp/.ICE-unix/1995",
						PeerLabel: "gnome-shell",
					},
				},
				{
					&Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
					&Dbus{
						Access:    []string{"receive"},
						Bus:       "system",
						Path:      "/org/freedesktop/DBus",
						Interface: "org.freedesktop.DBus",
						Member:    "AddMatch",
						PeerName:  ":1.3",
						PeerLabel: "power-profiles-daemon",
					},
				},
				{
					&File{Path: "/opt/intel/oneapi/compiler/*/linux/lib/*.so./*", Access: []string{"m", "r"}},
					&File{Path: "@{PROC}/@{pid}/task/@{tid}/comm", Access: []string{"r", "w"}},
					&File{Path: "@{sys}/devices/@{pci}/class", Access: []string{"r"}},
				},
				{
					&Include{IfExists: true, IsMagic: true, Path: "local/foo"},
				},
			},
			wParseRulesErr: false,
		},
		{
			name: "string.aa/content",
			raw:  pStringAAContentStr,
			blocks: []*block{{
				kind: CONTENT,
				raw:  pStringAAContentStr,
			}},
			wTokenizeErr:   false,
			rules:          []Rules{rulesStringAA},
			wParseBlockErr: false,
			apparmor:       &AppArmorProfileFile{},
			wParseErr:      true,
			wRules: ParaRules{
				{
					&Include{IsMagic: true, Path: "abstractions/base"},
					&Include{IsMagic: true, Path: "abstractions/nameservice-strict"},
				},
				{
					&Rlimit{Key: "nproc", Op: "<=", Value: "200"},
				},
				{
					&Capability{Names: []string{"dac_read_search"}},
					&Capability{Names: []string{"dac_override"}},
				},
				{
					&Network{Domain: "inet", Type: "stream"},
					&Network{Domain: "inet6", Type: "stream"},
				},
				{
					&Mount{
						Base: Base{IsLineRule: false, Comment: " failed perms check"},
						MountConditions: MountConditions{
							FsType:  "fuse.portal",
							Options: []string{"rw", "rbind"},
						},
						Source:     "@{run}/user/@{uid}/",
						MountPoint: "/",
					},
				},
				{
					&Umount{
						MountConditions: MountConditions{Options: []string{}},
						MountPoint:      "@{run}/user/@{uid}/",
					},
				},
				{
					&Signal{
						Access: []string{"receive"},
						Set:    []string{"term"},
						Peer:   "at-spi-bus-launcher",
					},
				},
				{
					&Ptrace{Access: []string{"read"}, Peer: "nautilus"},
				},
				{
					&Unix{
						Access:    []string{"send", "receive"},
						Type:      "stream",
						Address:   "@/tmp/.ICE-unix/1995",
						PeerLabel: "gnome-shell",
					},
				},
				{
					&Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
					&Dbus{
						Access:    []string{"receive"},
						Bus:       "system",
						Path:      "/org/freedesktop/DBus",
						Interface: "org.freedesktop.DBus",
						Member:    "AddMatch",
						PeerName:  ":1.3",
						PeerLabel: "power-profiles-daemon",
					},
				},
				{
					&File{Path: "/opt/intel/oneapi/compiler/*/linux/lib/*.so./*", Access: []string{"m", "r"}},
					&File{Path: "@{PROC}/@{pid}/task/@{tid}/comm", Access: []string{"r", "w"}},
					&File{Path: "@{sys}/devices/@{pci}/class", Access: []string{"r"}},
				},
			},
			wParseRulesErr: false,
		},
		{
			name: "profile0.aa",
			raw:  testData.Join("profile0.aa").MustReadFileAsString(),
			blocks: []*block{
				{
					kind: PROFILE,
					raw:  "profile A flags=(attach_disconnected)",
					next: &block{
						kind: RAW,
						raw:  "/path/to/A mr,\n\n  profile B {\n    /path/to/B mr,\n  }\n  profile C {\n    /path/to/C mr,\n  }\n  profile D {\n    /path/to/D mr,\n  }",
					},
				},
			},
			wTokenizeErr: false,
			rules: []Rules{
				{&Profile{
					Header: Header{
						Name: "A", Attachments: []string{}, Attributes: map[string]string{},
						Flags: []string{"attach_disconnected"},
					},
					Rules: Rules{
						&File{Path: "/path/to/A", Access: []string{"m", "r"}},
						profileB,
						profileC,
						profileD,
					},
				}},
			},
			wParseBlockErr: false,
			apparmor: &AppArmorProfileFile{
				Profiles: []*Profile{
					{
						Header: Header{
							Name: "A", Attachments: []string{}, Attributes: map[string]string{},
							Flags: []string{"attach_disconnected"},
						},
					},
				},
			},
			apparmorAll: &AppArmorProfileFile{
				Profiles: []*Profile{
					{
						Header: Header{
							Name: "A", Attachments: []string{}, Attributes: map[string]string{},
							Flags: []string{"attach_disconnected"},
						},
						Rules: Rules{
							&File{Path: "/path/to/A", Access: []string{"m", "r"}},
							profileB,
							profileC,
							profileD,
						},
					},
				},
			},
			wParseErr: false,
			wRules: ParaRules{
				{
					&File{Path: "/path/to/A", Access: []string{"m", "r"}},
				},
			},
			wParseRulesErr: false,
		},
		{
			name: "profile1.aa",
			raw:  testData.Join("profile1.aa").MustReadFileAsString(),
			blocks: []*block{
				{
					kind: PROFILE,
					raw:  "profile A flags=(attach_disconnected)",
					next: &block{
						kind: RAW,
						raw:  "/path/to/A mr,\n\n  profile B {\n    /path/to/B mr,\n\n    profile C {\n      /path/to/C mr,\n\n      profile D {\n        /path/to/D mr,\n      }\n    }\n  }",
					},
				},
			},
			wTokenizeErr: false,
			rules: []Rules{
				{&Profile{
					Header: Header{
						Name: "A", Attachments: []string{}, Attributes: map[string]string{},
						Flags: []string{"attach_disconnected"},
					},
					Rules: Rules{
						&File{Path: "/path/to/A", Access: []string{"m", "r"}},
						&Profile{
							Header: Header{
								Name: "B", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
							},
							Rules: Rules{
								&File{Path: "/path/to/B", Access: []string{"m", "r"}},
								&Profile{
									Header: Header{
										Name: "C", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
									},
									Rules: Rules{
										&File{Path: "/path/to/C", Access: []string{"m", "r"}},
										&Profile{
											Header: Header{
												Name: "D", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
											},
											Rules: Rules{
												&File{Path: "/path/to/D", Access: []string{"m", "r"}},
											},
										},
									},
								},
							},
						},
					},
				}},
			},
			wParseBlockErr: false,
			apparmor: &AppArmorProfileFile{
				Profiles: []*Profile{
					{
						Header: Header{
							Name: "A", Attachments: []string{}, Attributes: map[string]string{},
							Flags: []string{"attach_disconnected"},
						},
					},
				},
			},
			apparmorAll: &AppArmorProfileFile{
				Profiles: []*Profile{
					{
						Header: Header{
							Name: "A", Attachments: []string{}, Attributes: map[string]string{},
							Flags: []string{"attach_disconnected"},
						},
						Rules: Rules{
							&File{Path: "/path/to/A", Access: []string{"m", "r"}},
							&Profile{
								Header: Header{
									Name: "B", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
								},
								Rules: Rules{
									&File{Path: "/path/to/B", Access: []string{"m", "r"}},
									&Profile{
										Header: Header{
											Name: "C", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
										},
										Rules: Rules{
											&File{Path: "/path/to/C", Access: []string{"m", "r"}},
											&Profile{
												Header: Header{
													Name: "D", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
												},
												Rules: Rules{
													&File{Path: "/path/to/D", Access: []string{"m", "r"}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wParseErr: false,
			wRules: ParaRules{
				{
					&File{Path: "/path/to/A", Access: []string{"m", "r"}},
				},
				{
					&File{Path: "/path/to/B", Access: []string{"m", "r"}},
				},
				{
					&File{Path: "/path/to/C", Access: []string{"m", "r"}},
				},
			},
			wParseRulesErr: false,
		},
		{
			name: "full.aa",
			raw:  testData.Join("full.aa").MustReadFileAsString(),
			blocks: []*block{
				{
					kind: CONTENT,
					raw:  "# Simple test profile with all rules used\n\nabi <abi/4.0>,\n\nalias /mnt/usr -> /usr,\n\ninclude <tunables/global> # optional: a comment\ninclude if exists \"/etc/apparmor.d/global/dummy space\"\n\n@{name}=torbrowser \"tor browser\" \n@{lib_dirs} = @{lib}/@{name} /opt/@{name} # another comment\n@{config_dirs} = @{HOME}/.mozilla/\n@{cache_dirs}=@{user_cache_dirs}/mozilla/\n\nalias /mnt/{,usr.sbin.}mount.cifs -> /sbin/mount.cifs,\n\n@{exec_path} = @{bin}/@{name} @{lib_dirs}/@{name}",
				},
				{
					kind: PROFILE,
					raw:  "profile foo @{exec_path} xattrs=(security.tagged=allowed) flags=(complain attach_disconnected)",
					next: &block{
						kind: RAW,
						raw:  "include <abstractions/base>\n  include <abstractions/nameservice-strict>\n  include \"/etc/apparmor.d/abstractions/dummy space\"\n\n  all,\n\n  set rlimit nproc <= 200,\n\n  userns,\n\n  capability dac_read_search,\n  capability dac_override,\n\n  network inet stream,\n  network netlink raw,\n\n  mount /{,**},\n  mount               options=(rw rbind)           /tmp/newroot/ -> /tmp/newroot/,\n  mount               options=(rw silent rprivate)               -> /oldroot/,\n  mount fstype=devpts options=(rw nosuid noexec)          devpts -> /newroot/dev/pts/,\n\n  remount /newroot/{,**},\n\n  umount @{run}/user/@{uid}/,\n\n  pivot_root oldroot=/tmp/oldroot/ /tmp/,\n\n  change_profile -> libvirt-@{uuid},\n\n  mqueue r type=posix /,\n\n  io_uring sqpoll label=foo,\n\n  signal (receive) set=(cont,term,winch) peer=at-spi-bus-launcher,\n\n  ptrace (read) peer=nautilus,\n\n  unix (send receive) type=stream addr=\"@/tmp/.ICE[0-9]-unix/19 5\" peer=(label=gnome-shell, addr=none),\n\n  dbus bind bus=session name=org.gnome.*,\n  dbus receive bus=system path=/org/freedesktop/DBus\n       interface=org.freedesktop.DBus\n       member=AddMatch\n       peer=(name=:1.3, label=power-profiles-daemon),\n\n  # A comment! before a paragraph of rules\n  \"/opt/Mullvad VPN/resources/*.so*\" mr,\n  \"/opt/Mullvad VPN/resources/*\" r,\n  \"/opt/Mullvad VPN/resources/openvpn\" rix,\n  /usr/share/gnome-shell/extensions/ding@rastersoft.com/{,*/}ding.js rPx,\n  /opt/intel/oneapi/compiler/*/linux/lib/*.so./* rm,\n\n  owner @{user_config_dirs}/powerdevilrc{,.@{rand6}} rwl -> @{user_config_dirs}/#@{int},\n  link @{user_config_dirs}/kiorc -> @{user_config_dirs}/#@{int},\n\n  @{run}/udev/data/+pci:* r,\n\n  @{sys}/devices/@{pci}/class r,\n\n  owner @{PROC}/@{pid}/task/@{tid}/comm rw,\n\n  ^action {\n    include <abstractions/base>\n    include if exists <local/foo_action>\n  }\n\n  profile systemctl {\n    include <abstractions/base>\n    include <abstractions/systemctl>\n\n    capability net_admin,\n  \n    include if exists <local/foo_systemctl>\n  }\n\n  profile sudo {\n    include <abstractions/base>\n    include <abstractions/app/sudo>\n\n    @{sh_path} rix,\n\n    include if exists <local/foo_sudo>\n  }\n\n  include if exists <local/foo>",
					},
				},
				{
					kind: CONTENT,
					raw:  "\n",
				},
				{
					kind: PROFILE,
					raw:  "profile foo2",
					next: &block{
						kind: RAW,
						raw:  "include <abstractions/base>\n\n  include if exists <local/foo2>",
					},
				},
			},
			wTokenizeErr: false,
			rules: []Rules{
				// block 0: preamble content (parsed via parsePreamble)
				{
					&Comment{Base: Base{IsLineRule: true, Comment: " Simple test profile with all rules used"}},
					&Include{Base: Base{Comment: " a comment", Optional: true}, IsMagic: true, Path: "tunables/global"},
					&Include{IfExists: true, Path: "/etc/apparmor.d/global/dummy space"},
					&Variable{Name: "name", Values: []string{"torbrowser", "\"tor browser\""}, Define: true},
					&Variable{Base: Base{Comment: " another comment"}, Name: "lib_dirs", Values: []string{"@{lib}/@{name}", "/opt/@{name}"}, Define: true},
					&Variable{Name: "config_dirs", Values: []string{"@{HOME}/.mozilla/"}, Define: true},
					&Variable{Name: "cache_dirs", Values: []string{"@{user_cache_dirs}/mozilla/"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"@{bin}/@{name}", "@{lib_dirs}/@{name}"}, Define: true},
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
					&Alias{Path: "/mnt/{,usr.sbin.}mount.cifs", RewrittenPath: "/sbin/mount.cifs"},
				},
				// block 1: profile foo (parseBlock returns Rules{&Profile{...}})
				{
					&Profile{
						Header: Header{
							Name:        "foo",
							Attachments: []string{"@{exec_path}"},
							Attributes:  map[string]string{"security.tagged": "allowed"},
							Flags:       []string{"complain", "attach_disconnected"},
						},
						Rules: rulesFullAA,
					},
				},
				// block 2: empty content (newline)
				nil,
				// block 3: profile foo2
				{
					&Profile{
						Header: Header{
							Name:        "foo2",
							Attachments: []string{},
							Attributes:  map[string]string{},
							Flags:       []string{},
						},
						Rules: Rules{
							&Include{IsMagic: true, Path: "abstractions/base"},
							&Include{IfExists: true, IsMagic: true, Path: "local/foo2"},
						},
					},
				},
			},
			apparmor: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{Base: Base{IsLineRule: true, Comment: " Simple test profile with all rules used"}},
					&Include{
						Base:    Base{Comment: " a comment", Optional: true},
						IsMagic: true, Path: "tunables/global",
					},
					&Include{IfExists: true, Path: "/etc/apparmor.d/global/dummy space"},
					&Variable{Name: "name", Values: []string{"torbrowser", "\"tor browser\""}, Define: true},
					&Variable{
						Base: Base{Comment: " another comment"}, Define: true,
						Name: "lib_dirs", Values: []string{"@{lib}/@{name}", "/opt/@{name}"},
					},
					&Variable{Name: "config_dirs", Values: []string{"@{HOME}/.mozilla/"}, Define: true},
					&Variable{Name: "cache_dirs", Values: []string{"@{user_cache_dirs}/mozilla/"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"@{bin}/@{name}", "@{lib_dirs}/@{name}"}, Define: true},
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
					&Alias{Path: "/mnt/{,usr.sbin.}mount.cifs", RewrittenPath: "/sbin/mount.cifs"},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "foo",
							Attachments: []string{"@{exec_path}"},
							Attributes:  map[string]string{"security.tagged": "allowed"},
							Flags:       []string{"complain", "attach_disconnected"},
						},
					},
				},
			},
			apparmorAll: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{Base: Base{IsLineRule: true, Comment: " Simple test profile with all rules used"}},
					&Include{
						Base:    Base{Comment: " a comment", Optional: true},
						IsMagic: true, Path: "tunables/global",
					},
					&Include{IfExists: true, Path: "/etc/apparmor.d/global/dummy space"},
					&Variable{Name: "name", Values: []string{"torbrowser", "\"tor browser\""}, Define: true},
					&Variable{
						Base: Base{Comment: " another comment"}, Define: true,
						Name: "lib_dirs", Values: []string{"@{lib}/@{name}", "/opt/@{name}"},
					},
					&Variable{Name: "config_dirs", Values: []string{"@{HOME}/.mozilla/"}, Define: true},
					&Variable{Name: "cache_dirs", Values: []string{"@{user_cache_dirs}/mozilla/"}, Define: true},
					&Variable{Name: "exec_path", Values: []string{"@{bin}/@{name}", "@{lib_dirs}/@{name}"}, Define: true},
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
					&Alias{Path: "/mnt/{,usr.sbin.}mount.cifs", RewrittenPath: "/sbin/mount.cifs"},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "foo",
							Attachments: []string{"@{exec_path}"},
							Attributes:  map[string]string{"security.tagged": "allowed"},
							Flags:       []string{"complain", "attach_disconnected"},
						},
						Rules: rulesFullAA,
					},
					{
						Header: Header{
							Name:        "foo2",
							Attachments: []string{},
							Attributes:  map[string]string{},
							Flags:       []string{},
						},
						Rules: Rules{
							&Include{IsMagic: true, Path: "abstractions/base"},
							&Include{IfExists: true, IsMagic: true, Path: "local/foo2"},
						},
					},
				},
			},
			wParseErr: false,
			wRules: ParaRules{
				{
					&Include{IsMagic: true, Path: "abstractions/base"},
					&Include{IsMagic: true, Path: "abstractions/nameservice-strict"},
					&Include{Path: "/etc/apparmor.d/abstractions/dummy space"},
				},
				{
					&All{},
				},
				{
					&Rlimit{Key: "nproc", Op: "<=", Value: "200"},
				},
				{
					&Userns{Create: true},
				},
				{
					&Capability{Names: []string{"dac_read_search"}},
					&Capability{Names: []string{"dac_override"}},
				},
				{
					&Network{Domain: "inet", Type: "stream"},
					&Network{Domain: "netlink", Type: "raw"},
				},
				{
					&Mount{
						MountConditions: MountConditions{Options: []string{}},
						Source:          "/{,**}",
					},
					&Mount{
						MountConditions: MountConditions{Options: []string{"rw", "rbind"}},
						Source:          "/tmp/newroot/",
						MountPoint:      "/tmp/newroot/",
					},
					&Mount{
						MountConditions: MountConditions{Options: []string{"rw", "rprivate", "silent"}},
						MountPoint:      "/oldroot/",
					},
					&Mount{
						MountConditions: MountConditions{
							FsType:  "devpts",
							Options: []string{"rw", "noexec", "nosuid"},
						},
						Source:     "devpts",
						MountPoint: "/newroot/dev/pts/",
					},
				},
				{
					&Remount{
						MountConditions: MountConditions{Options: []string{}},
						MountPoint:      "/newroot/{,**}",
					},
				},
				{
					&Umount{
						MountConditions: MountConditions{Options: []string{}},
						MountPoint:      "@{run}/user/@{uid}/",
					},
				},
				{
					&PivotRoot{OldRoot: "/tmp/oldroot/", NewRoot: "/tmp/"},
				},
				{
					&ChangeProfile{ProfileName: "libvirt-@{uuid}"},
				},
				{
					&Mqueue{Access: []string{"r"}, Type: "posix", Name: "/"},
				},
				{
					&IOUring{Access: []string{"sqpoll"}, Label: "foo"},
				},
				{
					&Signal{
						Access: []string{"receive"},
						Set:    []string{"cont", "term", "winch"},
						Peer:   "at-spi-bus-launcher",
					},
				},
				{
					&Ptrace{Access: []string{"read"}, Peer: "nautilus"},
				},
				{
					&Unix{
						Access:    []string{"send", "receive"},
						Type:      "stream",
						Address:   "\"@/tmp/.ICE[0-9]-unix/19 5\"",
						PeerLabel: "gnome-shell",
						PeerAddr:  "none",
					},
				},
				{
					&Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
					&Dbus{
						Access:    []string{"receive"},
						Bus:       "system",
						Path:      "/org/freedesktop/DBus",
						Interface: "org.freedesktop.DBus",
						Member:    "AddMatch",
						PeerName:  ":1.3",
						PeerLabel: "power-profiles-daemon",
					},
				},
				{
					&Comment{Base: Base{IsLineRule: true, Comment: " A comment! before a paragraph of rules"}},
					&File{
						Path:   "\"/opt/Mullvad VPN/resources/*.so*\"",
						Access: []string{"m", "r"},
					},
					&File{Path: "\"/opt/Mullvad VPN/resources/*\"", Access: []string{"r"}},
					&File{
						Path:   "\"/opt/Mullvad VPN/resources/openvpn\"",
						Access: []string{"r", "ix"},
					},
					&File{
						Path:   "/usr/share/gnome-shell/extensions/ding@rastersoft.com/{,*/}ding.js",
						Access: []string{"r", "Px"},
					},
					&File{
						Path:   "/opt/intel/oneapi/compiler/*/linux/lib/*.so./*",
						Access: []string{"m", "r"},
					},
				},
				{
					&File{
						Owner: true, Path: "@{user_config_dirs}/powerdevilrc{,.@{rand6}}",
						Access: []string{"r", "w", "l"}, Target: "@{user_config_dirs}/#@{int}",
					},
					&Link{Path: "@{user_config_dirs}/kiorc", Target: "@{user_config_dirs}/#@{int}"},
				},
				{
					&File{Path: "@{run}/udev/data/+pci:*", Access: []string{"r"}},
				},
				{
					&File{Path: "@{sys}/devices/@{pci}/class", Access: []string{"r"}},
				},
				{
					&File{Owner: true, Path: "@{PROC}/@{pid}/task/@{tid}/comm", Access: []string{"r", "w"}},
				},
				{
					&Include{IsMagic: true, Path: "abstractions/base"},
					&Include{IfExists: true, IsMagic: true, Path: "local/foo_action"},
				},
				{
					&Include{IsMagic: true, Path: "abstractions/base"},
					&Include{IsMagic: true, Path: "abstractions/systemctl"},
				},
				{
					&Include{IfExists: true, IsMagic: true, Path: "local/foo_systemctl"},
					&Capability{Names: []string{"net_admin"}},
				},
				{
					&Include{IsMagic: true, Path: "abstractions/base"},
					&Include{IsMagic: true, Path: "abstractions/app/sudo"},
				},
				{
					&File{Path: "@{sh_path}", Access: []string{"r", "ix"}},
				},
				{
					&Include{IfExists: true, IsMagic: true, Path: "local/foo_sudo"},
				},
				{
					&Include{IfExists: true, IsMagic: true, Path: "local/foo"},
				},
				{
					&Include{IsMagic: true, Path: "abstractions/base"},
				},
			},
			wParseRulesErr: false,
		},
	}

	// Additional test cases for Scan())
	testParser = []struct {
		name    string
		raw     string
		want    *AppArmorProfileFile
		wantErr bool
	}{
		{
			name: "parse.aa",
			raw:  testData.Join("parse.aa").MustReadFileAsString(),
			want: &AppArmorProfileFile{
				Preamble: Rules{
					&Comment{Base: Base{IsLineRule: true, Comment: " apparmor.d - Full set of apparmor profiles"}},
					&Comment{Base: Base{IsLineRule: true, Comment: " Copyright (C) 2024 Alexandre Pujol <alexandre@pujol.io>"}},
					&Comment{Base: Base{IsLineRule: true, Comment: " SPDX-License-Identifier: GPL-2.0-only"}},
					&Comment{Base: Base{IsLineRule: true, Comment: " TODO: Test rule with only ','"}},
					&Include{
						Base:    Base{Optional: true, Comment: " A nice message"},
						IsMagic: true, Path: "tunables/global",
					},
					&Include{IfExists: true, Path: "/etc/apparmor.d/global/dummy space"},
					&Variable{
						Name: "name", Define: true,
						Values: []string{"torbrowser", "\"tor browser\""},
					},
					&Variable{
						Base: Base{Comment: " This is a comment"},
						Name: "lib_dirs", Define: true,
						Values: []string{"@{lib}/@{name}", "/opt/@{name}"},
					},
					&Variable{
						Name: "exec_path", Define: true,
						Values: []string{"@{bin}/@{name}", "@{lib_dirs}/@{name}"},
					},
					&Abi{IsMagic: true, Path: "abi/4.0"},
					&Alias{Path: "/mnt/usr", RewrittenPath: "/usr"},
					&Alias{Path: "/mnt/{,usr.sbin.}mount.cifs", RewrittenPath: "/sbin/mount.cifs"},
				},
				Profiles: []*Profile{
					{
						Header: Header{
							Name:        "foo",
							Attachments: []string{"@{exec_path}"},
							Attributes:  map[string]string{"security.tagged": "allowed"},
							Flags:       []string{"complain", "attach_disconnected"},
						},
						Rules: Rules{
							&Include{IsMagic: true, Path: "abstractions/base"},
							&Comment{Base: Base{IsLineRule: true, Comment: " Oh my god, it's a comment! before a paragraph of rules"}},
							&Include{IfExists: true, IsMagic: true, Path: "local/foo"},
							&Remount{MountConditions: MountConditions{Options: []string{}}, MountPoint: "/newroot/{,**}"},
							&Unix{
								Access:    []string{"send", "receive"},
								Type:      "stream",
								Address:   "\"@/tmp/.ICE[0-9]-unix/19 5\"",
								PeerLabel: "gnome-shell",
								PeerAddr:  "none",
							},
							&Dbus{
								Base:   Base{Comment: " wfdwde"},
								Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*",
							},
							&Dbus{
								Access:    []string{"receive"},
								Bus:       "system",
								Path:      "/org/freedesktop/DBus",
								Interface: "org.freedesktop.DBus",
								Member:    "AddMatch",
								PeerName:  ":1.3",
								PeerLabel: "power-profiles-daemon",
							},
							&File{
								Base: Base{Comment: " To be able to read the /proc/ files"},
								Path: "\"/opt/Mullvad VPN/resources/*.so*\"", Access: []string{"m", "r"},
							},
							&File{Path: "/usr/share/gnome-shell/extensions/ding@rastersoft.com/{,*/}ding.js", Access: []string{"r", "Px"}},
							&File{Path: "@{bin}/zsh", Access: []string{"r", "ix"}},
							&File{Owner: true, Path: "/{var/,}tmp/#@{int}", Access: []string{"r", "w"}},
							&File{Owner: true, Path: "@{user_config_dirs}/powerdevilrc{,.@{rand6}}", Access: []string{"r", "w", "l"}, Target: "@{user_config_dirs}/#@{int}"},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	// Indirect test resources
	pStringAAContentStr = `include <abstractions/base>
  include <abstractions/nameservice-strict>

  set rlimit nproc <= 200,

  capability dac_read_search,
  capability dac_override,

  network inet stream,
  network inet6 stream,

  mount fstype=fuse.portal options=(rw rbind) @{run}/user/@{uid}/ -> /, # failed perms check

  umount @{run}/user/@{uid}/,

  signal receive set=term peer=at-spi-bus-launcher,

  ptrace read peer=nautilus,

  unix (send receive) type=stream addr=@/tmp/.ICE-unix/1995 peer=(label=gnome-shell),

  dbus bind bus=session name=org.gnome.*,
  dbus receive bus=system path=/org/freedesktop/DBus
       interface=org.freedesktop.DBus
       member=AddMatch
       peer=(name=:1.3, label=power-profiles-daemon),

  /opt/intel/oneapi/compiler/*/linux/lib/*.so./* rm,
  @{PROC}/@{pid}/task/@{tid}/comm rw,
  @{sys}/devices/@{pci}/class r,

  include if exists <local/foo>`
	rulesStringAA = Rules{
		&Include{IsMagic: true, Path: "abstractions/base"},
		&Include{IsMagic: true, Path: "abstractions/nameservice-strict"},
		&Include{IfExists: true, IsMagic: true, Path: "local/foo"},
		&Rlimit{Key: "nproc", Op: "<=", Value: "200"},
		&Capability{Names: []string{"dac_read_search"}},
		&Capability{Names: []string{"dac_override"}},
		&Network{Domain: "inet", Type: "stream"},
		&Network{Domain: "inet6", Type: "stream"},
		&Mount{
			Base: Base{Comment: " failed perms check"},
			MountConditions: MountConditions{
				FsType:  "fuse.portal",
				Options: []string{"rw", "rbind"},
			},
			Source:     "@{run}/user/@{uid}/",
			MountPoint: "/",
		},
		&Umount{
			MountConditions: MountConditions{Options: []string{}},
			MountPoint:      "@{run}/user/@{uid}/",
		},
		&Signal{
			Access: []string{"receive"},
			Set:    []string{"term"},
			Peer:   "at-spi-bus-launcher",
		},
		&Ptrace{Access: []string{"read"}, Peer: "nautilus"},
		&Unix{
			Access:    []string{"send", "receive"},
			Type:      "stream",
			Address:   "@/tmp/.ICE-unix/1995",
			PeerLabel: "gnome-shell",
		},
		&Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
		&Dbus{
			Access:    []string{"receive"},
			Bus:       "system",
			Name:      "",
			Path:      "/org/freedesktop/DBus",
			Interface: "org.freedesktop.DBus",
			Member:    "AddMatch",
			PeerName:  ":1.3",
			PeerLabel: "power-profiles-daemon",
		},
		&File{Path: "/opt/intel/oneapi/compiler/*/linux/lib/*.so./*", Access: []string{"m", "r"}},
		&File{Path: "@{PROC}/@{pid}/task/@{tid}/comm", Access: []string{"r", "w"}},
		&File{Path: "@{sys}/devices/@{pci}/class", Access: []string{"r"}},
	}
	rulesFullAA = Rules{
		&Include{IsMagic: true, Path: "abstractions/base"},
		&Include{IsMagic: true, Path: "abstractions/nameservice-strict"},
		&Include{Path: "/etc/apparmor.d/abstractions/dummy space"},
		&Comment{Base: Base{IsLineRule: true, Comment: " A comment! before a paragraph of rules"}},
		&All{},
		&Rlimit{Key: "nproc", Op: "<=", Value: "200"},
		&Userns{Create: true},
		&Capability{Names: []string{"dac_read_search"}},
		&Capability{Names: []string{"dac_override"}},
		&Network{Domain: "inet", Type: "stream"},
		&Network{Domain: "netlink", Type: "raw"},
		&Mount{MountConditions: MountConditions{Options: []string{}}, Source: "/{,**}"},
		&Mount{MountConditions: MountConditions{Options: []string{"rw", "rbind"}}, Source: "/tmp/newroot/", MountPoint: "/tmp/newroot/"},
		&Mount{MountConditions: MountConditions{Options: []string{"rw", "rprivate", "silent"}}, MountPoint: "/oldroot/"},
		&Mount{MountConditions: MountConditions{FsType: "devpts", Options: []string{"rw", "noexec", "nosuid"}}, Source: "devpts", MountPoint: "/newroot/dev/pts/"},
		&Remount{MountConditions: MountConditions{Options: []string{}}, MountPoint: "/newroot/{,**}"},
		&Umount{MountConditions: MountConditions{Options: []string{}}, MountPoint: "@{run}/user/@{uid}/"},
		&PivotRoot{OldRoot: "/tmp/oldroot/", NewRoot: "/tmp/"},
		&ChangeProfile{ProfileName: "libvirt-@{uuid}"},
		&Mqueue{Access: []string{"r"}, Type: "posix", Name: "/"},
		&IOUring{Access: []string{"sqpoll"}, Label: "foo"},
		&Signal{Access: []string{"receive"}, Set: []string{"cont", "term", "winch"}, Peer: "at-spi-bus-launcher"},
		&Ptrace{Access: []string{"read"}, Peer: "nautilus"},
		&Unix{Access: []string{"send", "receive"}, Type: "stream", Address: "\"@/tmp/.ICE[0-9]-unix/19 5\"", PeerLabel: "gnome-shell", PeerAddr: "none"},
		&Dbus{Access: []string{"bind"}, Bus: "session", Name: "org.gnome.*"},
		&Dbus{Access: []string{"receive"}, Bus: "system", Path: "/org/freedesktop/DBus", Interface: "org.freedesktop.DBus", Member: "AddMatch", PeerName: ":1.3", PeerLabel: "power-profiles-daemon"},
		&File{Path: "\"/opt/Mullvad VPN/resources/*.so*\"", Access: []string{"m", "r"}},
		&File{Path: "\"/opt/Mullvad VPN/resources/*\"", Access: []string{"r"}},
		&File{Path: "\"/opt/Mullvad VPN/resources/openvpn\"", Access: []string{"r", "ix"}},
		&File{Path: "/usr/share/gnome-shell/extensions/ding@rastersoft.com/{,*/}ding.js", Access: []string{"r", "Px"}},
		&File{Path: "/opt/intel/oneapi/compiler/*/linux/lib/*.so./*", Access: []string{"m", "r"}},
		&File{Owner: true, Path: "@{user_config_dirs}/powerdevilrc{,.@{rand6}}", Access: []string{"r", "w", "l"}, Target: "@{user_config_dirs}/#@{int}"},
		&Link{Path: "@{user_config_dirs}/kiorc", Target: "@{user_config_dirs}/#@{int}"},
		&File{Path: "@{run}/udev/data/+pci:*", Access: []string{"r"}},
		&File{Path: "@{sys}/devices/@{pci}/class", Access: []string{"r"}},
		&File{Owner: true, Path: "@{PROC}/@{pid}/task/@{tid}/comm", Access: []string{"r", "w"}},
		&Hat{
			Name:  "action",
			Flags: []string{},
			Rules: Rules{
				&Include{IsMagic: true, Path: "abstractions/base"},
				&Include{IfExists: true, IsMagic: true, Path: "local/foo_action"},
			},
		},
		&Profile{
			Header: Header{Name: "systemctl", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{}},
			Rules: Rules{
				&Include{IsMagic: true, Path: "abstractions/base"},
				&Include{IsMagic: true, Path: "abstractions/systemctl"},
				&Include{IfExists: true, IsMagic: true, Path: "local/foo_systemctl"},
				&Capability{Names: []string{"net_admin"}},
			},
		},
		&Profile{
			Header: Header{Name: "sudo", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{}},
			Rules: Rules{
				&Include{IsMagic: true, Path: "abstractions/base"},
				&Include{IsMagic: true, Path: "abstractions/app/sudo"},
				&Include{IfExists: true, IsMagic: true, Path: "local/foo_sudo"},
				&File{Path: "@{sh_path}", Access: []string{"r", "ix"}},
			},
		},
		&Include{IfExists: true, IsMagic: true, Path: "local/foo"},
	}
	profileB = &Profile{
		Header: Header{
			Name: "B", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
		},
		Rules: Rules{
			&File{Path: "/path/to/B", Access: []string{"m", "r"}},
		},
	}
	profileC = &Profile{
		Header: Header{
			Name: "C", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
		},
		Rules: Rules{
			&File{Path: "/path/to/C", Access: []string{"m", "r"}},
		},
	}
	profileD = &Profile{
		Header: Header{
			Name: "D", Attachments: []string{}, Attributes: map[string]string{}, Flags: []string{},
		},
		Rules: Rules{
			&File{Path: "/path/to/D", Access: []string{"m", "r"}},
		},
	}
)
