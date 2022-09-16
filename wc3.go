package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("CustomKeys.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

func generate(rules []Rule, r io.Reader, w io.Writer) error {
	var (
		current       Group
		currentHotkey string
	)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		matched := false

	innerloop:
		for _, e := range rules {
			switch e.matches(line) {

			case matchCommand:
				current.Hotkey = currentHotkey
				current.Apply(rules)
				current.Print(w)
				current = Group{Lines: []string{line}}
				currentHotkey = ""
				matched = true
				break innerloop

			case matchHotkey:
				current.Lines = append(current.Lines, line)
				if e.action == hotkey2 && currentHotkey == "" {
					currentHotkey = e.extract(line)
				} else if e.action == hotkey {
					currentHotkey = e.extract(line)
				}
				matched = true
				break innerloop

			case matchTrue:
				current.Lines = append(current.Lines, line)
				matched = true
				break innerloop
			}
		}

		if !matched {
			current.Lines = append(current.Lines, line)
		}
	}

	current.Apply(rules)
	current.Print(w)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not scan intput: %w", err)
	}

	return nil
}
