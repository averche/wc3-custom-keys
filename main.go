package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/go-multierror"
)

func main() {
	var customKeysPath string

	switch len(os.Args) {
	case 1:
		customKeysPath = "CustomKeys.txt"
	case 2:
		customKeysPath = os.Args[1]
	default:
		log.Fatalf("Usage: %s <path/to/CustomKeys.txt>", os.Args[0])
	}

	if err := run(customKeysPath); err != nil {
		log.Fatal(err)
	}
}

func run(path string) (errs error) {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open %q: %w", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			errs = multierror.Append(errs, fmt.Errorf("could not close %q: %w", path, err))
		}
	}()

	if err := generate(rules(), f, os.Stdout); err != nil {
		return fmt.Errorf("generation error: %w", err)
	}

	return nil
}

func generate(rules []rule, r io.Reader, w io.Writer) error {
	var (
		current       Group
		currentHotkey string
	)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		matched := false

	innerloop:
		for _, r := range rules {
			switch r.matches(line) {

			case matchCommand:
				current.Hotkey = currentHotkey
				if err := current.Apply(rules); err != nil {
					return fmt.Errorf("could not apply rules: %w", err)
				}
				current.Print(w)
				current = Group{Lines: []string{line}}
				currentHotkey = ""
				matched = true
				break innerloop

			case matchHotkey:
				current.Lines = append(current.Lines, line)
				if r.action == actionHotkey2 && currentHotkey == "" {
					currentHotkey = r.extract(line)
				} else if r.action == actionHotkey {
					currentHotkey = r.extract(line)
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

	// the last group
	if err := current.Apply(rules); err != nil {
		return fmt.Errorf("could not apply rules: %w", err)
	}
	current.Print(w)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not scan intput: %w", err)
	}

	return nil
}
