package main

import (
	"fmt"
	"io"
)

// Group denotes a single hotkey group within the CustomKeys.txt file
type Group struct {
	Hotkey string
	Lines  []string
}

// Apply will apply the rules and replace all lines with the group their regex equivalents
func (g *Group) Apply(rules []rule) {
	for i, line := range g.Lines {
		for _, r := range rules {
			match := r.matches(line)
			if match == matchCommand || match == matchHotkey || match == matchTrue {
				g.Lines[i] = r.replace(line, g.Hotkey)
				break
			}
		}
	}
}

// Print outputs the group to the given writer
func (g *Group) Print(w io.Writer) {
	for _, line := range g.Lines {
		fmt.Fprintln(w, line)
	}
}
