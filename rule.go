package main

import (
	"fmt"
	"regexp"
)

type action uint8

const (
	actionKeep action = iota
	actionCommand
	actionHotkey
	actionHotkey2
	actionReplaceOne
	actionReplaceTwo
	actionReplaceThree
)

type match uint8

const (
	matchFalse match = iota
	matchTrue
	matchHotkey
	matchCommand
)

// rule encapusulates a regular expression and an associated action
type rule struct {
	action action
	regex  *regexp.Regexp
}

// rules returns a set of rules to apply to all lines
func rules() []rule {
	var rules []rule

	rules = append(rules, rule{ // command
		action: actionCommand,
		regex:  regexp.MustCompile(`^\[[\w]*\][ \t]*$`),
	})
	rules = append(rules, rule{ // Hotkey
		action: actionHotkey,
		regex:  regexp.MustCompile(`^(?P<name>Hotkey=)(?P<hotkey>\w+)(,\w+){0,2}[ \t]*$`),
	})
	rules = append(rules, rule{ // Researchhotkey
		action: actionHotkey2,
		regex:  regexp.MustCompile(`^(?P<name>Researchhotkey=)(?P<hotkey>\w+)(,\w+){0,2}[ \t]*$`),
	})
	rules = append(rules, rule{ // Unhotkey
		action: actionHotkey2,
		regex:  regexp.MustCompile(`^(?P<name>Unhotkey=)(?P<hotkey>\w+)(,\w+){0,2}[ \t]*$`),
	})
	rules = append(rules, rule{ // comment
		action: actionKeep,
		regex:  regexp.MustCompile(`^\/\/.*$`),
	})
	rules = append(rules, rule{ // empty
		action: actionKeep,
		regex:  regexp.MustCompile(`^[ \t]*$`),
	})
	rules = append(rules, rule{ // Awakentip=t(i)p
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Awakentip=)"?(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Awakentip=tip
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Awakentip=)"?(?P<p1>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Researchtip=t(i)p [Level %d]
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)(?P<l1> - \[\|cffffcc00Level %d\|r\])"?[ \t]*$`),
	})
	rules = append(rules, rule{ // Researchtip=t(i)p
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Researchtip=tip
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Revivetip=t(i)p
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Revivetip=)"?(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Revivetip=tip
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Revivetip=)"?(?P<p1>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Untip=t(i)p
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Untip=)"?(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Untip=tip
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Untip=)"?(?P<p1>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Tip=t(i)p1,t(i)p2,t(i)p3
		action: actionReplaceThree,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)(?P<l1> - \[\|cffffcc00Level 1\|r\],)` +
			`(?P<p3>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key2>\w)\|r(?P<p4>[\w \-\!\.\(\)]*)(?P<l2> - \[\|cffffcc00Level 2\|r\],)` +
			`(?P<p5>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.\(\)]*)(?P<l3> - \[\|cffffcc00Level 3\|r\])"?[ \t]*$`),
	})
	rules = append(rules, rule{ // Tip=t(i)p1,t(i)p2,t(i)p3
		action: actionReplaceThree,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)(?P<l1>,)` +
			`(?P<p3>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key2>\w)\|r(?P<p4>[\w \-\!\.\(\)]*)(?P<l2>,)` +
			`(?P<p5>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.\(\)]*)(?P<l3>)"?$`),
	})
	rules = append(rules, rule{ // Tip=t(i)p1,t(i)p2
		action: actionReplaceTwo,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)(?P<l1> - \[\|cffffcc00Level 1\|r\],)` +
			`(?P<p5>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.\(\)]*)(?P<l3> - \[\|cffffcc00Level 2\|r\])"?[ \t]*$`),
	})
	rules = append(rules, rule{ // Tip=t(i)p1,t(i)p2
		action: actionReplaceTwo,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)(?P<l1>,)` +
			`(?P<p5>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.\(\)]*)(?P<l3>)"?$`),
	})
	rules = append(rules, rule{ // Tip=t(i)p
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Tip=)"?(?P<p1>[\w \-\!\.\(\)]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.\(\)]*)"?$`),
	})
	rules = append(rules, rule{ // Tip=tip
		action: actionReplaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Tip=)"?(?P<p1>[\w \-\!\.\(\)]*)"?$`),
	})

	return rules
}

func (r *rule) extract(line string) string {
	return r.regex.ReplaceAllString(line, "$2")
}

// replace returns a string modified according to the rule
func (r *rule) replace(line string, key string) (string, error) {
	if key == "" {
		return line, nil
	}

	switch r.action {
	case actionKeep:
		return line, nil

	case actionCommand:
		return line, nil

	case actionHotkey, actionHotkey2:
		return fmt.Sprintf("%s%s", r.regex.ReplaceAllString(line, "$1"), key), nil

	case actionReplaceOne:
		return fmt.Sprintf(r.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)"), key) + r.regex.ReplaceAllString(line, "$5"), nil

	case actionReplaceTwo:
		return fmt.Sprintf(r.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)$5$6$7$8 (|cffffcc00%s|r)"), key, key), nil

	case actionReplaceThree:
		return fmt.Sprintf(r.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)$5$6$7$8 (|cffffcc00%s|r)$9$10$11$12 (|cffffcc00%s|r)$13$14"), key, key, key), nil
	}

	return "", fmt.Errorf("unknown action %d", r.action)
}

// matches returns MatchTrue if the line matches the regex or a more specific match (MatchCommand/MatchHotkey)
func (r *rule) matches(line string) match {
	if !r.regex.MatchString(line) {
		return matchFalse
	}

	switch r.action {
	case actionCommand:
		return matchCommand

	case actionHotkey, actionHotkey2:
		return matchHotkey
	}

	return matchTrue
}
