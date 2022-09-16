package main

import (
	"fmt"
	"io"
)

// group denotes a single hotkey group within the CustomKeys.txt file
type group struct {
	hotkey string
	lines  []string
}

// apply will modify each line in the group by applying the matching rule
func (g *group) apply(rules []rule) error {
	for i, line := range g.lines {
		for _, r := range rules {
			match := r.matches(line)
			if match == matchCommand || match == matchHotkey || match == matchTrue {
				l, err := r.replace(line, g.hotkey)
				if err != nil {
					return fmt.Errorf("error @ %q: %w", line, err)
				}
				g.lines[i] = l
				break
			}
		}
	}

	return nil
}

// print outputs the group to the given writer
func (g *group) print(w io.Writer) {
	for _, line := range g.lines {
		fmt.Fprintln(w, line)
	}
}
