// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"slices"
	"strings"
)

const (
	tokARROW        = "->"
	tokEQUAL        = "="
	tokLESS         = "<"
	tokPLUS         = "+"
	tokCLOSEBRACE   = '}'
	tokCLOSEBRACKET = ']'
	tokCLOSEPAREN   = ')'
	tokCOLON        = ','
	tokOPENBRACE    = '{'
	tokOPENBRACKET  = '['
	tokOPENPAREN    = '('
)

var (
	newRuleMap = map[string]func([]string) (Rule, error){
		COMMENT.Tok(): newComment,
		ABI.Tok():     newAbi,
		ALIAS.Tok():   newAlias,
		INCLUDE.Tok(): newInclude,
	}

	tok = map[Kind]string{
		COMMENT:  "#",
		VARIABLE: "@{",
	}
	openBlocks  = []rune{tokOPENPAREN, tokOPENBRACE, tokOPENBRACKET}
	closeBlocks = []rune{tokCLOSEPAREN, tokCLOSEBRACE, tokCLOSEBRACKET}
)

// Split a raw input rule string into tokens by space or =, but ignore spaces
// within quotes, brakets, or parentheses.
//
// Example:
//
//	`owner @{user_config_dirs}/powerdevilrc{,.@{rand6}} rwl -> @{user_config_dirs}/#@{int}`
//
// Returns:
//
//	[]string{"owner", "@{user_config_dirs}/powerdevilrc{,.@{rand6}}", "rwl", "->", "@{user_config_dirs}/#@{int}"}
func tokenize(str string) []string {
	var currentToken strings.Builder
	var isVariable bool
	var quoted bool

	blockStack := []rune{}
	tokens := make([]string, 0, len(str)/2)
	if len(str) > 2 && str[0:2] == VARIABLE.Tok() {
		isVariable = true
	}
	for _, r := range str {
		switch {
		case (r == ' ' || r == '\t') && len(blockStack) == 0 && !quoted:
			// Split on space/tab if not in a block or quoted
			if currentToken.Len() != 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}

		case (r == '=' || r == '+') && len(blockStack) == 0 && !quoted && isVariable:
			// Handle variable assignment
			if currentToken.Len() != 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			tokens = append(tokens, string(r))

		case r == '"' && len(blockStack) == 0:
			quoted = !quoted
			currentToken.WriteRune(r)

		case slices.Contains(openBlocks, r):
			blockStack = append(blockStack, r)
			currentToken.WriteRune(r)

		case slices.Contains(closeBlocks, r):
			if len(blockStack) > 0 {
				blockStack = blockStack[:len(blockStack)-1]
			} else {
				panic(fmt.Sprintf("Unbalanced block, missing '{' or '}' on: %s\n", str))
			}
			currentToken.WriteRune(r)

		default:
			currentToken.WriteRune(r)
		}
	}
	if currentToken.Len() != 0 {
		tokens = append(tokens, currentToken.String())
	}
	return tokens
}

func tokenToSlice(token string) []string {
	res := []string{}
	token = strings.Trim(token, "()\n")
	if strings.ContainsAny(token, ", ") {
		var sep string
		switch {
		case strings.Contains(token, ","):
			sep = ","
		case strings.Contains(token, " "):
			sep = " "
		}
		for _, v := range strings.Split(token, sep) {
			res = append(res, strings.Trim(v, " "))
		}
	} else {
		res = append(res, token)
	}
	return res
}

func tokensStripComment(tokens []string) []string {
	res := []string{}
	for _, v := range tokens {
		if v == COMMENT.Tok() {
			break
		}
		res = append(res, v)
	}
	return res
}

// Convert a slice of internal rules to a slice of ApparmorRule.
func newRules(rules [][]string) (Rules, error) {
	var err error
	var r Rule
	res := make(Rules, 0, len(rules))

	for _, rule := range rules {
		if len(rule) == 0 {
			return nil, fmt.Errorf("Empty rule")
		}

		if newRule, ok := newRuleMap[rule[0]]; ok {
			r, err = newRule(rule)
			if err != nil {
				return nil, err
			}
			res = append(res, r)
		} else if strings.HasPrefix(rule[0], VARIABLE.Tok()) {
			r, err = newVariable(rule)
			if err != nil {
				return nil, err
			}
			res = append(res, r)
		} else {
			return nil, fmt.Errorf("Unrecognized rule: %s", rule)
		}
	}
	return res, nil
}

func (f *AppArmorProfileFile) parsePreamble(input []string) error {
	var err error
	var r Rule
	var rules Rules

	tokenizedRules := [][]string{}
	for _, line := range input {
		if strings.HasPrefix(line, COMMENT.Tok()) {
			r, err = newComment(strings.Split(line, " "))
			if err != nil {
				return err
			}
			rules = append(rules, r)
		} else {
			tokens := tokenize(line)
			tokenizedRules = append(tokenizedRules, tokens)
		}
	}

	rr, err := newRules(tokenizedRules)
	if err != nil {
		return err
	}
	f.Preamble = append(f.Preamble, rules...)
	f.Preamble = append(f.Preamble, rr...)
	return nil
}

// Parse an apparmor profile file.
//
// Only supports parsing of apparmor file preamble and profile headers.
//
// Warning: It is purposelly an uncomplete basic parser for apparmor profile,
// it is only aimed for internal tooling purpose. For "simplicity", it is not
// using antlr / participle. It is only used for experimental feature in the
// apparmor.d project.
//
// Stop at the first profile header. Does not support multiline coma rules.
//
// Current use case:
//
//   - Parse include and tunables
//   - Parse variable in profile preamble and in tunable files
//   - Parse (sub) profiles header to edit flags
func (f *AppArmorProfileFile) Parse(input string) error {
	rawHeader := ""
	rawPreamble := []string{}

done:
	for _, line := range strings.Split(input, "\n") {
		tmp := strings.TrimLeft(line, "\t ")
		tmp = strings.TrimRight(tmp, ",")
		switch {
		case tmp == "":
			continue
		case strings.HasPrefix(tmp, PROFILE.Tok()):
			rawHeader = tmp
			break done
		default:
			rawPreamble = append(rawPreamble, tmp)
		}
	}

	if err := f.parsePreamble(rawPreamble); err != nil {
		return err
	}
	if rawHeader != "" {
		header, err := newHeader(tokenize(rawHeader))
		if err != nil {
			return err
		}
		profile := &Profile{Header: header}
		f.Profiles = append(f.Profiles, profile)
	}
	return nil
}
