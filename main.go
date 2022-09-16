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
	var err error

	switch len(os.Args) {
	case 1:
		err = run("", "")
	case 2:
		err = run(os.Args[1], "")
	case 3:
		err = run(os.Args[1], os.Args[2])
	default:
		log.Fatalf("Usage: %s <path/to/CustomKeys.txt>", os.Args[0])
	}

	if err != nil {
		log.Fatal(err)
	}
}

func run(pathIn, pathOut string) (errs error) {
	var (
		in  = os.Stdin
		out = os.Stdout
		err error
	)

	if pathIn != "" {
		in, err = os.Open(pathIn)
		if err != nil {
			return fmt.Errorf("could not open input %q: %w", pathIn, err)
		}
		defer func() {
			if err := in.Close(); err != nil {
				errs = multierror.Append(errs, fmt.Errorf("could not close input %q: %w", pathIn, err))
			}
		}()
	}

	if pathOut != "" {
		out, err = os.Create(pathOut)
		if err != nil {
			return fmt.Errorf("could not open output %q: %w", pathOut, err)
		}
		defer func() {
			if err := out.Close(); err != nil {
				errs = multierror.Append(errs, fmt.Errorf("could not close output %q: %w", pathOut, err))
			}
		}()
	}

	if err := generate(rules(), in, out); err != nil {
		return fmt.Errorf("generation error: %w", err)
	}

	return nil
}

func generate(rules []rule, in io.Reader, out io.Writer) error {
	var (
		current       Group
		currentHotkey string
	)

	scanner := bufio.NewScanner(in)

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
				current.Print(out)
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
	current.Print(out)

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("could not scan intput: %w", err)
	}

	return nil
}
