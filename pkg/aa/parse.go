// apparmor.d - Full set of apparmor profiles
// Copyright (C) 2021-2024 Alexandre Pujol <alexandre@pujol.io>
// SPDX-License-Identifier: GPL-2.0-only

package aa

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

const (
	tokALLOW        = "allow"
	tokAUDIT        = "audit"
	tokDENY         = "deny"
	tokPROMPT       = "prompt"
	tokPRIORITY     = "priority"
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
	newRuleMap = map[string]func(q Qualifier, r rule) (Rule, error){
		ABI.Tok():           newAbi,
		ALIAS.Tok():         newAlias,
		ALL.Tok():           newAll,
		"set":               newRlimit,
		USERNS.Tok():        newUserns,
		CAPABILITY.Tok():    newCapability,
		NETWORK.Tok():       newNetwork,
		MOUNT.Tok():         newMount,
		UMOUNT.Tok():        newUmount,
		REMOUNT.Tok():       newRemount,
		MQUEUE.Tok():        newMqueue,
		IOURING.Tok():       newIOUring,
		PIVOTROOT.Tok():     newPivotRoot,
		CHANGEPROFILE.Tok(): newChangeProfile,
		SIGNAL.Tok():        newSignal,
		PTRACE.Tok():        newPtrace,
		UNIX.Tok():          newUnix,
		DBUS.Tok():          newDbus,
		FILE.Tok():          newFile,
		LINK.Tok():          newLink,
	}

	tok = map[Kind]string{
		COMMENT:  "#",
		VARIABLE: "@{",
		HAT:      "^",
	}
	openBlocks  = []rune{tokOPENPAREN, tokOPENBRACE, tokOPENBRACKET}
	closeBlocks = []rune{tokCLOSEPAREN, tokCLOSEBRACE, tokCLOSEBRACKET}

	inHeader              = false
	regParagraph          = regexp.MustCompile(`(?s).*?\n\n|$`)
	regVariableDefinition = regexp.MustCompile(`@{(.*)}\s*[+=]+\s*(.*)`)
)

const (
	CONTENT   Kind = "content"
	RAW       Kind = "raw"
	QUALIFIER Kind = "qualifier"
)

// Intermediate token for the representation of apparmor blocks such as
// profile, subprofile, hat, qualifier and condition.
type block struct {
	kind Kind
	raw  string
	next *block
}

// Split a raw input block string into tokens by '{', '}', but ignore
// variables and second level blocks.
func tokenizeBlock(input string) ([]*block, error) {
	if len(input) > 0 && input[0] == tokOPENBRACE {
		return nil, fmt.Errorf("wrong block format: %s", input)
	}

	blocks := []*block{}
	blockCounter := 0
	blockStack := []rune{}
	blockRecored := false
	blockStart := 0
	blockEnd := 0
	blockContentStart := 0
	blockContentStartBkp := 0
	blockContentEnd := 0
	for idx, r := range input {
		switch r {
		case tokOPENBRACE:
			blockStack = append(blockStack, r)

			// Block rules starts with ' {', ignore nested blocks and variables
			if len(blockStack) == 1 {
				ignore := false

				// Ignore the block if it is inside a variable definition
				if input[idx-1] == '@' {
					ignore = true
				} else {
					i := idx
					for i > 0 && input[i] != '\n' {
						i--
					}
					j := idx
					for j < len(input) && input[j] != '\n' {
						j++
					}
					line := input[i:j]
					match := regVariableDefinition.FindStringSubmatch(line)
					if len(match) > 0 {
						ignore = true
					} else if input[idx-1] != ' ' {
						ignore = true
					}
				}

				if !ignore {
					blockStart = idx
					blockRecored = true
				}
			}

		case tokCLOSEBRACE:
			if len(blockStack) <= 0 {
				return nil, fmt.Errorf("unbalanced block, missing '{' for '} at: }%s",
					input[blockContentStart:idx])
			}

			if len(blockStack) == 1 && blockRecored {
				blockRecored = false
				blockEnd = idx
				blockContentStartBkp = blockContentStart
				blockContentEnd = blockStart
				blockContentRaw := input[blockContentStart:blockContentEnd]
				blockContentStart = blockEnd + 1

				// Collect the block header
				i := len(blockContentRaw) - 1
				for i > 0 && blockContentRaw[i] != '\n' {
					i--
				}
				blockHeader := strings.Trim(blockContentRaw[i:], "\n ")
				blockHeader = strings.Trim(blockHeader, "\n\t ")

				// Ignore commented block, restore previous id values
				if len(blockHeader) > 0 && blockHeader[0] == '#' {
					blockContentStart = blockContentStartBkp
					blockStack = blockStack[:len(blockStack)-1]
					continue
				}

				// Collect out of block content (preamble, profile content outside of sub blocks)
				blockContentRaw = blockContentRaw[:i]
				if blockContentRaw != "" {
					blocks = append(blocks, &block{
						kind: CONTENT,
						raw:  blockContentRaw,
					})
				}

				// Collect the block content
				blockRaw := input[blockStart+1 : blockEnd]
				blockRaw = strings.Trim(blockRaw, "\n\t ")
				var kind Kind
				switch {
				case strings.HasPrefix(blockHeader, PROFILE.Tok()), isAARE(blockHeader):
					kind = PROFILE
				case strings.HasPrefix(blockHeader, HAT.Tok()),
					strings.HasPrefix(blockHeader, HAT.String()):
					kind = HAT
				case blockHeader == IF.Tok():
					kind = IF
				case blockHeader == ELSE.Tok():
					kind = ELSE
				default:
					fmt.Printf("blockRaw: %v\n", blockRaw)
					return nil, fmt.Errorf("unrecognized block type: %s", blockHeader)
				}
				blocks = append(blocks, &block{
					kind: kind,
					raw:  blockHeader,
					next: &block{
						kind: RAW,
						raw:  blockRaw,
					},
				})
			}
			blockStack = blockStack[:len(blockStack)-1]
		}
	}

	if blockCounter != 0 {
		return nil, fmt.Errorf("unbalanced block, missing '{' or '}': %s",
			input[blockContentEnd:len(input)-1])
	}
	if len(blocks) == 0 {
		// No block found, it can be a tunable/abstraction file.
		blocks = append(blocks, &block{
			kind: CONTENT,
			raw:  input,
		})
	}
	return blocks, nil
}

func parseBlock(b *block) (Rules, error) {
	var res Rules
	var rrr Rules
	var err error

	switch b.kind {
	case CONTENT:
		// Line rules
		var raw string
		raw, res, err = parseLineRules(false, b.raw)
		if err != nil {
			return nil, err
		}

		// Comma rules
		rules, err := parseCommaRules(raw)
		if err != nil {
			return nil, err
		}
		rrr, err = newRules(rules)
		if err != nil {
			return nil, err
		}

		res = append(res, rrr...)
		for _, r := range res {
			if r.Constraint() == PreambleRule {
				return nil, fmt.Errorf("Rule not allowed in block: %s", r)
			}
		}

	case RAW:
		var blocks []*block
		blocks, err = tokenizeBlock(b.raw)
		if err != nil {
			return nil, err
		}
		for _, block := range blocks {
			rrr, err = parseBlock(block)
			if err != nil {
				return nil, err
			}
			res = append(res, rrr...)
		}
		return res, nil

	case PROFILE:
		header, err := newHeader(parseRule(b.raw))
		if err != nil {
			return nil, err
		}
		rules, err := parseBlock(b.next)
		if err != nil {
			return nil, err
		}
		profile := &Profile{
			Header: header,
			Rules:  rules,
		}
		res = append(res, profile)

	case HAT:
		hat, err := newHat(parseRule(b.raw))
		if err != nil {
			return nil, err
		}
		rules, err := parseBlock(b.next)
		if err != nil {
			return nil, err
		}
		hat.Rules = rules
		res = append(res, hat)

	case IF, ELSE:
		// Not implemented yet

	}
	return res, nil
}

// Parse the line rule from a raw string.
func parseLineRules(isPreamble bool, input string) (string, Rules, error) {
	var res Rules
	var r Rule
	var err error

	for _, line := range strings.Split(input, "\n") {
		tmp := strings.TrimLeft(line, "\t ")
		switch {
		case strings.HasPrefix(tmp, COMMENT.Tok()):
			r, err = newComment(rule{kv{comment: tmp[1:]}})
			if err != nil {
				return "", nil, err
			}
			res = append(res, r)
			input = strings.Replace(input, line, "", 1)

		case strings.HasPrefix(tmp, INCLUDE.Tok()):
			r, err = newInclude(parseRule(line)[1:])
			if err != nil {
				return "", nil, err
			}
			res = append(res, r)
			input = strings.Replace(input, line, "", 1)

		case strings.HasPrefix(tmp, VARIABLE.Tok()) && isPreamble:
			r, err = newVariable(parseRule(line))
			if err != nil {
				return "", nil, err
			}
			res = append(res, r)
			input = strings.Replace(input, line, "", 1)
		}
	}
	return input, res, nil
}

// Parse the comma rules from a raw string. It splits rules string into tokens
// separated by "," but ignore comma inside comments, quotes, brakets, and parentheses.
// Warning: the input string should only contain comma rules.
// Return a pre-parsed rule struct representation of a profile/block.
func parseCommaRules(input string) ([]rule, error) {
	rules := []rule{}
	blockStart := 0
	blockCounter := 0
	comment := false
	aare := false
	canHaveInlineComment := false
	size := len(input)
	for idx, r := range input {
		switch r {
		case tokOPENBRACE, tokOPENBRACKET, tokOPENPAREN:
			if !comment {
				blockCounter++
			}

		case tokCLOSEBRACE, tokCLOSEBRACKET, tokCLOSEPAREN:
			if !comment {
				blockCounter--
			}

		case '#':
			if !comment && canHaveInlineComment {
				comment = true
				blockStart = idx + 1
			}

		case '\n':
			if comment {
				comment = !comment
				if canHaveInlineComment {
					commentRaw := input[blockStart:idx]

					// Inline comments belong to the previous rule (in the same line)
					lastRule := rules[len(rules)-1]
					lastRule[len(lastRule)-1].comment = commentRaw

					// Ignore the collected comment for the next rule
					blockStart = idx
				}
			}
			canHaveInlineComment = false

		case tokCOLON:
			if blockCounter == 0 && !comment {
				if idx+1 < size && !strings.ContainsRune(" \n", rune(input[idx+1])) {
					// Colon in AARE, it is valid, not a separator
					aare = true
				}
				if !aare {
					ruleRaw := input[blockStart:idx]
					ruleRaw = strings.Trim(ruleRaw, "\n ")
					rules = append(rules, parseRule(ruleRaw))
					blockStart = idx + 1
					canHaveInlineComment = true
				}
				aare = false

			}
		}
	}
	return rules, nil
}

func parseParagraph(input string) (Rules, error) {
	// Line rules
	var raw string
	raw, res, err := parseLineRules(false, input)
	if err != nil {
		return nil, err
	}

	// Comma rules
	rules, err := parseCommaRules(raw)
	if err != nil {
		return nil, err
	}
	rrr, err := newRules(rules)
	if err != nil {
		return nil, err
	}

	res = append(res, rrr...)
	// for _, r := range res {
	// 	if r.Constraint() == PreambleRule {
	// 		return nil, fmt.Errorf("Rule not allowed in block: %s", r)
	// 	}
	// }
	return res, nil
}

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
func tokenizeRule(str string) []string {
	var currentToken strings.Builder
	isVariable, wasTokPLUS, quoted := false, false, false

	blockStack := []rune{}
	tokens := make([]string, 0, len(str)/2)
	if inHeader && len(str) > 2 && str[0:2] == VARIABLE.Tok() {
		isVariable = true
	}
	for _, r := range str {
		switch {
		case (r == ' ' || r == '\t' || r == '\n') && len(blockStack) == 0 && !quoted:
			// Split on space/tab/newline if not in a block or quoted
			if currentToken.Len() != 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}

		case (r == '+' || r == '=') && len(blockStack) == 0 && !quoted && isVariable:
			// Handle variable assignment
			if currentToken.Len() != 0 {
				tokens = append(tokens, currentToken.String())
				currentToken.Reset()
			}
			if wasTokPLUS {
				tokens[len(tokens)-1] = tokPLUS + tokEQUAL
			} else {
				tokens = append(tokens, string(r))
			}
			wasTokPLUS = (r == '+')

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

// Parse a string into a rule struct.
// The input string must be a tokenised raw line rule (using parseCommaRules or
// parseBlock).
// Example:
//
//	`unix (send receive) type=stream addr="@/tmp/.ICE[0-9]-unix/19 5" peer=(label=gnome-shell, addr=none)`
//
// Returns:
//
//	  rule{
//		   {key: "unix"}, {key: "send"}, {key: "receive"},
//		   {key: "type", values: rule{{key: "stream"}}},
//		   {key: "addr", values: rule{
//			 {key: `"@/tmp/.ICE[0-9]*-unix/19 5"`},
//		   }},
//		   {key: "peer", values: rule{
//			 {key: "label", values: rule{{key: `"@{p_systemd}"`}}},
//			 {key: "addr", values: rule{{key: "none"}}},
//		   }},
//	   },
func parseRule(str string) rule {
	res := make(rule, 0, len(str)/2)
	tokens := tokenizeRule(str)

	inAare := len(tokens) > 0 && (isAARE(tokens[0]) || tokens[0] == tokOWNER)
	for idx, token := range tokens {
		switch {
		case token == tokEQUAL, token == tokPLUS+tokEQUAL, token == tokLESS+tokEQUAL: // Variable & Rlimit
			res = append(res, kv{key: token})

		case strings.Contains(token, "=") && !inAare: // Map
			items := strings.SplitN(token, "=", 2)
			key := items[0]
			if len(items) > 1 {
				values := strings.Trim(items[1], ",")
				if strings.Contains(values, "=") || !strings.ContainsAny(values, ", ") {
					values = strings.Trim(values, "()\n")
				}
				res = append(res, kv{key: key, values: parseRule(values)})
			} else {
				res = append(res, kv{key: key})
			}

		case strings.Contains(token, "(") && !inAare: // List
			token = strings.Trim(token, "()\n")
			var sep string
			switch {
			case strings.Contains(token, ","):
				sep = ","
			case strings.Contains(token, " "):
				sep = " "
			}

			var values rule
			if sep == "" {
				values = append(values, kv{key: token})
			} else {
				for _, v := range strings.Split(token, sep) {
					values = append(values, kv{
						key: strings.Trim(v, " "),
					})
				}
			}
			res = append(res, values...)

		case strings.HasPrefix(token, COMMENT.Tok()): // Comment
			if idx > 0 && idx < len(tokens)-1 {
				res[len(res)-1].comment = " " + strings.Join(tokens[idx+1:], " ")
				return res
			}

		default: // Single value
			token = strings.Trim(token, "\n")
			res = append(res, kv{key: token})
		}
	}
	return res
}

// Intermediate token for the representation of a rule. All comma and line
// rules are parsed into this structure. Then, they are converted into the actual
// rule struct using basic constructor functions.
type rule []kv

type kv struct {
	key     string
	values  rule
	comment string
}

// Get return the value of a key from a rule.
//
// Example:
//
//	`include <tunables/global>`
//
// Gives:
//
//	Get(0): "include"
//	Get(1): "<tunables/global>"
func (r rule) Get(idx int) string {
	return r[idx].key
}

// GetString return string representation of a rule.
//
// Example:
//
//	`profile foo @{exec_path} flags=(complain attach_disconnected)`
//
// Gives:
//
//	"profile foo @{exec_path}"
func (r rule) GetString() string {
	return strings.Join(r.GetSlice(), " ")
}

// GetSlice return a slice of all non map value of a rule.
//
// Example:
//
//	`profile foo @{exec_path} flags=(complain attach_disconnected)`
//
// Gives:
//
//	[]string{"profile", "foo", "@{exec_path}"}
func (r rule) GetSlice() []string {
	res := []string{}
	for _, kv := range r {
		if kv.values == nil {
			res = append(res, kv.key)
		}
	}
	return res
}

// GetAsMap return a map of slice of all map value of a rule.
//
// Example:
//
//	`profile foo @{exec_path} flags=(complain attach_disconnected)`
//
// Gives:
//
//	map[string]string{"flags": {"complain", "attach_disconnected"}}
func (r rule) GetAsMap() map[string][]string {
	res := map[string][]string{}
	for _, kv := range r {
		if kv.values != nil {
			res[kv.key] = kv.values.GetSlice()
		}
	}
	return res
}

// GetValues return the values from a key.
//
// Example:
//
//	`dbus receive peer=(name=:1.3, label=power-profiles-daemon)`
//
// Gives:
//
//	  GetValues("peer"):
//		rule{
//			{key: "name", values: rule{{Key: ":1.3"}}},
//			{key: "label", values: rule{{Key: "power-profiles-daemon"}}},
//		}},
func (r rule) GetValues(key string) rule {
	for _, kv := range r {
		if kv.key == key {
			return kv.values
		}
	}
	return nil
}

// GetValuesAsSlice return the values from a key as a slice.
//
// Example:
//
//	`mount options=(rw silent rprivate) -> /oldroot/`
//
// Gives:
//
//	GetValuesAsSlice("options"):
//	 []string{"rw", "silent", "rprivate"}
func (r rule) GetValuesAsSlice(key string) []string {
	return r.GetValues(key).GetSlice()
}

// GetValuesAsString return the values from a key as a string.
//
// Example:
//
//	`signal (receive) set=(term) peer=at-spi-bus-launcher`
//
// Gives:
//
//	GetValuesAsString("peer"): "at-spi-bus-launcher"
func (r rule) GetValuesAsString(key string) string {
	return r.GetValues(key).GetString()
}

// String return a generic representation of a rule.
func (r rule) String() string {
	var res strings.Builder
	for _, kv := range r {
		if kv.values == nil {
			if res.Len() > 0 {
				res.WriteString(" ")
			}
			res.WriteString(kv.key)
		} else {
			res.WriteString(" " + kv.key)
			v := strings.TrimLeft(kv.values.String(), " ")
			if strings.Contains(v, " ") || strings.Contains(v, "=") {
				res.WriteString("=(" + v + ")")
			} else {
				res.WriteString("=" + v)
			}
		}
		if kv.comment != "" {
			res.WriteString(COMMENT.Tok() + " " + kv.comment)
		}
	}
	return res.String()
}

func isAARE(str string) bool {
	if len(str) < 1 {
		return false
	}
	switch str[0] {
	case '@', '/', '"':
		return true
	default:
		return false
	}
}

// Convert a slice of internal rules to a slice of ApparmorRule.
func newRules(rules []rule) (Rules, error) {
	var err error
	var r Rule
	res := make(Rules, 0, len(rules))

	for _, rule := range rules {
		if len(rule) == 0 {
			return nil, fmt.Errorf("empty rule")
		}

		owner := false
		q := Qualifier{}
	qualifier:
		switch rule.Get(0) {
		// File & Link prefix
		case tokOWNER:
			owner = true
			rule = rule[1:]
			goto qualifier
		// Qualifier
		case tokPRIORITY:
			q.Priority = rule.GetValues(tokPRIORITY).GetString()
			rule = rule[1:]
			goto qualifier
		case tokALLOW, tokDENY, tokPROMPT:
			q.AccessType = rule.Get(0)
			rule = rule[1:]
			goto qualifier
		case tokAUDIT:
			q.Audit = true
			rule = rule[1:]
			goto qualifier

		default:
			// Line rules
			if newRule, ok := newRuleMap[rule.Get(0)]; ok {
				r, err = newRule(q, rule[1:])
				if err != nil {
					return nil, err
				}
				if owner && r.Kind() == LINK {
					r.(*Link).Owner = owner
				}
				res = append(res, r)
			} else {
				raw := rule.Get(0)
				if raw != "" {
					// File
					if isAARE(raw) || owner {
						r, err = newFile(q, rule)
						if err != nil {
							return nil, err
						}
						r.(*File).Owner = owner
						res = append(res, r)
					} else {
						fmt.Printf("Unknown rule: %s", rule)
						// return nil, fmt.Errorf("Unknown rule: %s", rule)
					}
				} else {
					return nil, fmt.Errorf("unrecognized rule: %s", rule)
				}
			}
		}
	}
	return res, nil
}

func (f *AppArmorProfileFile) parsePreamble(preamble string) error {
	var err error
	inHeader = true

	// Line rules
	preamble, lineRules, err := parseLineRules(true, preamble)
	if err != nil {
		return err
	}
	f.Preamble = append(f.Preamble, lineRules...)

	// Comma rules
	r, err := parseCommaRules(preamble)
	if err != nil {
		return err
	}
	commaRules, err := newRules(r)
	if err != nil {
		return err
	}
	f.Preamble = append(f.Preamble, commaRules...)

	for _, r := range f.Preamble {
		if r.Constraint() == BlockRule {
			f.Preamble = nil
			return fmt.Errorf("Rule not allowed in preamble: %s", r)
		}
	}
	inHeader = false
	return err
}

// Parse an apparmor profile file.
//
// Warning: It is purposely an uncomplete parser for apparmor profile,
// it is only aimed for internal tooling purpose. For "simplicity", it is not
// using antlr / participle. It is only used for experimental feature in the
// apparmor.d project. Technically, it is more a scanner than a parser.
//
// Very basic:
//   - Only supports parsing of preamble and profile headers.
//   - Stop at the first profile header.
//   - Does not support multiline coma rules.
//   - Does not support multiple profiles by file.
//
// Current use case:
//   - Parse include and tunables
//   - Parse variable in profile preamble and in tunable files
//   - Parse (sub) profiles header to edit flags
func (f *AppArmorProfileFile) Parse(input string) (int, error) {
	var raw strings.Builder
	rawHeader := ""
	nb := 0

done:
	for i, line := range strings.Split(input, "\n") {
		tmp := strings.TrimLeft(line, "\t ")
		switch {
		case tmp == "":
			continue
		case strings.HasPrefix(tmp, PROFILE.Tok()):
			rawHeader = strings.TrimRight(tmp, "{")
			nb = i
			break done
		case strings.HasPrefix(tmp, HAT.String()), strings.HasPrefix(tmp, HAT.Tok()):
			nb = i
			break done
		default:
			raw.WriteString(tmp + "\n")
		}
	}

	if err := f.parsePreamble(raw.String()); err != nil {
		return nb, err
	}
	if rawHeader != "" {
		header, err := newHeader(parseRule(rawHeader))
		if err != nil {
			return nb, err
		}
		profile := &Profile{Header: header}
		f.Profiles = append(f.Profiles, profile)
	}
	return nb, nil
}

// ParseRules parses apparmor profile rules by paragraphs
func ParseRules(input string) (ParaRules, []string, error) {
	paragraphRules := ParaRules{}
	paragraphs := []string{}

	for _, match := range regParagraph.FindAllStringSubmatch(input, -1) {
		if len(match[0]) == 0 {
			continue
		}

		// Ignore blocks header
		tmp := strings.TrimLeft(match[0], "\t ")
		tmp = strings.TrimRight(tmp, "\n")
		var paragraph string
		switch {
		case strings.HasPrefix(tmp, PROFILE.Tok()):
			_, paragraph, _ = strings.Cut(match[0], "\n")
		case strings.HasPrefix(tmp, HAT.String()), strings.HasPrefix(tmp, HAT.Tok()):
			_, paragraph, _ = strings.Cut(match[0], "\n")
		case strings.HasSuffix(tmp, "}"):
			paragraph = strings.Replace(match[0], "}\n", "\n", 1)
		default:
			paragraph = match[0]
		}

		paragraphs = append(paragraphs, paragraph)
		rules, err := parseParagraph(paragraph)
		if err != nil {
			return nil, nil, err
		}
		paragraphRules = append(paragraphRules, rules)
	}

	return paragraphRules, paragraphs, nil
}
