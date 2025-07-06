package delete

import (
	"errors"
	"flag"
	"fmt"
	"slices"
	"strings"

	"github.com/ba58ajbse/envcraft/internal/fs"
)

// DeleteOptions holds the options for updating an environment variable.
type DeleteOptions struct {
	Key      string
	FilePath string
}

// DeleteCmd represents the command for updating an environment variable in a file.
type DeleteCmd struct {
	Options  DeleteOptions
	OrgLines []string
}

// ErrNoUpdated is returned when no matching key is found to update.
var ErrNoUpdated = errors.New("no lines updated")

func Run(args []string) error {
	options, err := ParseDeleteOptions(args)
	if err != nil {
		return err
	}
	cmd, err := NewDeleteCmd(options)
	if err != nil {
		return err
	}
	err = cmd.Exec()
	if err != nil {
		return err
	}
	return nil
}

// NewDeleteCmd creates a new DeleteCmd instance with the specified options.
func NewDeleteCmd(options *DeleteOptions) (*DeleteCmd, error) {
	if options.FilePath == "" {
		return nil, errors.New("file path is required")
	}

	return &DeleteCmd{
		Options:  *options,
		OrgLines: []string{},
	}, nil
}

// Exec executes the update command: reads lines, updates the value, and writes back.
func (c *DeleteCmd) Exec() error {
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

// readLines reads all lines from the file specified in DeleteCmd and stores them in OrgLines.
func (c *DeleteCmd) readLines() error {
	lines, err := fs.ReadLines(c.filePath())
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", c.filePath(), err)
	}
	c.OrgLines = lines

	return nil
}

// makeNewLines returns a new slice of lines with the updated value if the key matches.
// If no matching key is found, returns ErrNoUpdated.
func (c *DeleteCmd) makeNewLines() ([]string, error) {
	if len(c.OrgLines) == 0 {
		return []string{}, errors.New("no lines read from file")
	}

	newLines := slices.Clone(c.OrgLines)

	updatedFlag := false
	for i, line := range c.OrgLines {
		parts := strings.Split(line, "=")
		key := strings.TrimSpace(parts[0])
		if c.keyEqual(key) {
			newLines = slices.Delete(newLines, i, i+1) // Remove the line with the matching key
			updatedFlag = true
			break
		}
	}

	lastLineIndex := len(newLines) - 1
	if lastLineIndex >= 0 {
		// Ensure the last line does not end with a newline
		newLines[lastLineIndex] = strings.TrimSuffix(newLines[lastLineIndex], "\n")
	}

	if !updatedFlag {
		return newLines, ErrNoUpdated
	}

	return newLines, nil
}

// apply writes the new lines to the file, overwriting the original content.
func (c *DeleteCmd) apply(newLines []string) error {
	if err := fs.WriteLines(c.filePath(), newLines); err != nil {
		return fmt.Errorf("error writing to file %s: %w", c.filePath(), err)
	}

	return nil
}

// filePath returns the file path from the options.
func (s *DeleteCmd) filePath() string {
	return s.Options.FilePath
}

// keyEqual checks if the given key matches the update target key.
func (s *DeleteCmd) keyEqual(key string) bool {
	return s.Options.Key == key
}

// ParseDeleteOptions parses command-line arguments and returns an DeleteOptions struct.
func ParseDeleteOptions(opts []string) (*DeleteOptions, error) {
	fs := flag.NewFlagSet("delete", flag.ContinueOnError)
	file := fs.String("f", "", "Path to .env file")

	if len(opts) < 2 {
		return nil, errors.New("key and value are required")
	}
	key := opts[0]
	flags := opts[1:]

	if err := fs.Parse(flags); err != nil {
		return nil, err
	}
	if *file == "" {
		fmt.Println("Error: -f flag is required")
		fs.Usage()
		return nil, errors.New("file path is required")
	}

	return &DeleteOptions{
		Key:      key,
		FilePath: *file,
	}, nil
}
