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
func (g *Group) Apply(rules []rule) error {
	for i, line := range g.Lines {
		for _, r := range rules {
			match := r.matches(line)
			if match == matchCommand || match == matchHotkey || match == matchTrue {
				l, err := r.replace(line, g.Hotkey)
				if err != nil {
					return fmt.Errorf("error @ %q: %w", line, err)
				}
				g.Lines[i] = l
				break
			}
		}
	}

	return nil
}

// Print outputs the group to the given writer
func (g *Group) Print(w io.Writer) {
	for _, line := range g.Lines {
		fmt.Fprintln(w, line)
	}
}
