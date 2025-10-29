package update

import (
	"errors"
	"flag"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/ba58ajbse/envcraft/internal/fs"
)

// UpdateOptions holds the options for updating an environment variable.
type UpdateOptions struct {
	Key      string
	Value    string
	FilePath string
}

// UpdateCmd represents the command for updating an environment variable in a file.
type UpdateCmd struct {
	Options  UpdateOptions
	OrgLines []string
}

// ErrNoUpdated is returned when no matching key is found to update.
var ErrNoUpdated = errors.New("no lines updated")

func Run(args []string) error {
	options, err := ParseUpdateOptions(args)
	if err != nil {
		return err
	}
	cmd, err := NewUpdateCmd(options)
	if err != nil {
		return err
	}
	err = cmd.Exec()
	if err != nil {
		return err
	}
	return nil
}

// NewUpdateCmd creates a new UpdateCmd instance with the specified options.
func NewUpdateCmd(options *UpdateOptions) (*UpdateCmd, error) {
	if options.FilePath == "" {
		return nil, errors.New("file path is required")
	}

	return &UpdateCmd{
		Options:  *options,
		OrgLines: []string{},
	}, nil
}

// Exec executes the update command: reads lines, updates the value, and writes back.
func (c *UpdateCmd) Exec() error {
	err := c.readLines()
	if err != nil {
		return err
	}

	newLines, err := c.makeNewLines()
	if err != nil {
		return err
	}

	err = c.apply(newLines)
	if err != nil {
		return err
	}

	return nil
}

// readLines reads all lines from the file specified in UpdateCmd and stores them in OrgLines.
func (c *UpdateCmd) readLines() error {
	lines, err := fs.ReadLines(c.filePath())
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", c.filePath(), err)
	}
	c.OrgLines = lines

	return nil
}

// makeNewLines returns a new slice of lines with the updated value if the key matches.
// If no matching key is found, returns ErrNoUpdated.
func (c *UpdateCmd) makeNewLines() ([]string, error) {
	if len(c.OrgLines) == 0 {
		return []string{}, errors.New("no lines read from file")
	}

	newLines := slices.Clone(c.OrgLines)

	updatedFlag := false
	for i, line := range c.OrgLines {
		parts := strings.Split(line, "=")
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if c.keyEqual(key) {
			if strings.HasSuffix(line, "\n") {
				newLines[i] = c.keyAndValue() + "\n"
			} else {
				newLines[i] = c.keyAndValue()
			}
			updatedFlag = true
			break
		}
	}

	if !updatedFlag {
		return newLines, ErrNoUpdated
	}

	return newLines, nil
}

// apply writes the new lines to the file, overwriting the original content.
func (c *UpdateCmd) apply(newLines []string) error {
	if err := fs.WriteLines(c.filePath(), newLines); err != nil {
		return fmt.Errorf("error writing to file %s: %w", c.filePath(), err)
	}

	return nil
}

// filePath returns the file path from the options.
func (s *UpdateCmd) filePath() string {
	return s.Options.FilePath
}

// keyEqual checks if the given key matches the update target key.
func (s *UpdateCmd) keyEqual(key string) bool {
	return s.Options.Key == key
}

// keyAndValue returns the key and quoted value in the format KEY="value".
func (s *UpdateCmd) keyAndValue() string {
	return s.Options.Key + "=" + strconv.Quote(s.Options.Value)
}

// ParseUpdateOptions parses command-line arguments and returns an UpdateOptions struct.
func ParseUpdateOptions(opts []string) (*UpdateOptions, error) {
	flagSet := flag.NewFlagSet("update", flag.ContinueOnError)
	file := flagSet.String("f", "", "Path to .env file")

	var key, value string

	if len(opts) >= 2 && !strings.HasPrefix(opts[0], "-") && !strings.HasPrefix(opts[1], "-") {
		key = opts[0]
		value = opts[1]
		if err := flagSet.Parse(opts[2:]); err != nil {
			return nil, err
		}
	} else {
		if err := flagSet.Parse(opts); err != nil {
			return nil, err
		}
		args := flagSet.Args()
		if len(args) < 2 {
			return nil, errors.New("key and value are required")
		}
		key = args[0]
		value = args[1]
		if strings.HasPrefix(key, "-") || strings.HasPrefix(value, "-") {
			return nil, errors.New("key and value are required")
		}
	}
	if *file == "" {
		fmt.Println("Error: -f flag is required")
		flagSet.Usage()
		return nil, errors.New("file path is required")
	}

	return &UpdateOptions{
		Key:      key,
		Value:    value,
		FilePath: *file,
	}, nil
}
