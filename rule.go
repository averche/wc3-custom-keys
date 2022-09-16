package main

import (
	"fmt"
	"regexp"
)

type action uint8

const (
	keep         action = 0
	command      action = 1
	hotkey       action = 2
	hotkey2      action = 3
	replaceOne   action = 4
	replaceTwo   action = 5
	replaceThree action = 6
)

type match uint8

const (
	matchFalse   match = 0
	matchTrue    match = 1
	matchHotkey  match = 2
	matchCommand match = 3
)

// Rule encapusulates a regular expression and an associated action
type Rule struct {
	action action
	regex  *regexp.Regexp
}

// generateExpressions returns a set of regular rules with correspodning actions
func generateExpressions() []Rule {
	var rules []Rule

	rules = append(rules, Rule{ // command
		action: command,
		regex:  regexp.MustCompile(`^\[[\w]*\][ \t]*$`),
	})
	rules = append(rules, Rule{ // Hotkey
		action: hotkey,
		regex:  regexp.MustCompile(`^(?P<name>Hotkey=)(?P<hotkey>\w+)(,\w+){0,2}[ \t]*$`),
	})
	rules = append(rules, Rule{ // Researchhotkey
		action: hotkey2,
		regex:  regexp.MustCompile(`^(?P<name>Researchhotkey=)(?P<hotkey>\w+)(,\w+){0,2}[ \t]*$`),
	})
	rules = append(rules, Rule{ // Unhotkey
		action: hotkey2,
		regex:  regexp.MustCompile(`^(?P<name>Unhotkey=)(?P<hotkey>\w+)(,\w+){0,2}[ \t]*$`),
	})
	rules = append(rules, Rule{ // comment
		action: keep,
		regex:  regexp.MustCompile(`^\/\/.*$`),
	})
	rules = append(rules, Rule{ // empty
		action: keep,
		regex:  regexp.MustCompile(`^[ \t]*$`),
	})
	rules = append(rules, Rule{ // Awakentip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Awakentip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	rules = append(rules, Rule{ // Awakentip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Awakentip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Awakentip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Awakentip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Researchtip=t(i)p [Level %d]
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1> - \[\|cffffcc00Level %d\|r\])"?[ \t]*$`),
	})
	rules = append(rules, Rule{ // Researchtip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Researchtip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	rules = append(rules, Rule{ // Researchtip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Researchtip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Researchtip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Revivetip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Revivetip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	rules = append(rules, Rule{ // Revivetip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Revivetip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Revivetip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Revivetip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Untip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Untip="?\|cffc3dbff[\w \-\!\.]{2,}\|r"?$`),
	})
	rules = append(rules, Rule{ // Untip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Untip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	rules = append(rules, Rule{ // Untip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Untip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Untip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Untip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Tip=t(i)p1,t(i)p2,t(i)p3
		action: replaceThree,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1> - \[\|cffffcc00Level 1\|r\],)` +
			`(?P<p3>[\w \-\!\.]*)\|cffffcc00(?P<key2>\w)\|r(?P<p4>[\w \-\!\.]*)(?P<l2> - \[\|cffffcc00Level 2\|r\],)` +
			`(?P<p5>[\w \-\!\.]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.]*)(?P<l3> - \[\|cffffcc00Level 3\|r\])"?[ \t]*$`),
	})
	rules = append(rules, Rule{ // Tip=t(i)p1,t(i)p2,t(i)p3
		action: replaceThree,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1>,)` +
			`(?P<p3>[\w \-\!\.]*)\|cffffcc00(?P<key2>\w)\|r(?P<p4>[\w \-\!\.]*)(?P<l2>,)` +
			`(?P<p5>[\w \-\!\.]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.]*)(?P<l3>)"?$`),
	})
	rules = append(rules, Rule{ // Tip=t(i)p1,t(i)p2
		action: replaceTwo,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1> - \[\|cffffcc00Level 1\|r\],)` +
			`(?P<p5>[\w \-\!\.]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.]*)(?P<l3> - \[\|cffffcc00Level 2\|r\])"?[ \t]*$`),
	})
	rules = append(rules, Rule{ // Tip=t(i)p1,t(i)p2
		action: replaceTwo,
		regex: regexp.MustCompile(`^(?P<name>Tip=)"?` +
			`(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key1>\w)\|r(?P<p2>[\w \-\!\.]*)(?P<l1>,)` +
			`(?P<p5>[\w \-\!\.]*)\|cffffcc00(?P<key3>\w)\|r(?P<p6>[\w \-\!\.]*)(?P<l3>)"?$`),
	})
	rules = append(rules, Rule{ // Tip=tip (E)
		action: keep,
		regex:  regexp.MustCompile(`^Tip=[\w \-\!\.]* \(\|cffffcc00\w+\|r\)[\w \-\!\.]*$`),
	})
	rules = append(rules, Rule{ // Tip=t(i)p
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Tip=)"?(?P<p1>[\w \-\!\.]*)\|cffffcc00(?P<key>\w)\|r(?P<p2>[\w \-\!\.]*)"?$`),
	})
	rules = append(rules, Rule{ // Tip=tip
		action: replaceOne,
		regex:  regexp.MustCompile(`^(?P<name>Tip=)"?(?P<p1>[\w \-\!\.]*)"?$`),
	})

	return rules
}

func (e *Rule) extract(line string) string {
	return e.regex.ReplaceAllString(line, "$2")
}

// replace returns a string modified according to the regex and the action
func (e *Rule) replace(line string, key string) string {
	if key == "" {
		return line
	}

	switch e.action {

	case keep:
		return line

	case command:
		return line

	case hotkey, hotkey2:
		return fmt.Sprintf("%s%s", e.regex.ReplaceAllString(line, "$1"), key)

	case replaceOne:
		return fmt.Sprintf(e.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)"), key) + e.regex.ReplaceAllString(line, "$5")

	case replaceTwo:
		return fmt.Sprintf(e.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)$5$6$7$8 (|cffffcc00%s|r)"), key, key)

	case replaceThree:
		return fmt.Sprintf(e.regex.ReplaceAllString(line, "$1$2$3$4 (|cffffcc00%s|r)$5$6$7$8 (|cffffcc00%s|r)$9$10$11$12 (|cffffcc00%s|r)$13$14"), key, key, key)
	}

	return "<< ERROR >>"
}

// matches returns MatchTrue if the line matches the regex or a more specific match (MatchCommand/MatchHotkey)
func (e *Rule) matches(line string) match {
	if !e.regex.MatchString(line) {
		return matchFalse
	}

	switch e.action {
	case command:
		return matchCommand

	case hotkey, hotkey2:
		return matchHotkey
	}

	return matchTrue
}
