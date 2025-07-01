package update

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
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
	envFile, err := os.Open(c.filePath())
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", c.filePath(), err)
	}
	defer envFile.Close()

	reader := bufio.NewReader(envFile)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				fmt.Printf("Error reading file %s: %v\n", c.filePath(), err)
			}
			c.OrgLines = append(c.OrgLines, line)
			break
		}
		c.OrgLines = append(c.OrgLines, line)
	}

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
	out, err := os.OpenFile(c.filePath(), os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("error opening file %s for writing: %w", c.filePath(), err)
	}
	defer out.Close()

	writer := bufio.NewWriter(out)
	for _, line := range newLines {
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("error writing to file %s: %w", c.filePath(), err)
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("error flushing writer for file %s: %w", c.filePath(), err)
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
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	file := fs.String("f", "", "Path to .env file")

	if len(opts) < 2 {
		return nil, errors.New("key and value are required")
	}
	key := opts[0]
	value := opts[1]
	flags := opts[2:]

	if err := fs.Parse(flags); err != nil {
		return nil, err
	}
	if *file == "" {
		fmt.Println("Error: -f flag is required")
		fs.Usage()
		return nil, errors.New("file path is required")
	}

	return &UpdateOptions{
		Key:      key,
		Value:    value,
		FilePath: *file,
	}, nil
}
