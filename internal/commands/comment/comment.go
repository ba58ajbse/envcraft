package comment

import (
	"errors"
	"flag"
	"fmt"
	"slices"

	"github.com/ba58ajbse/envcraft/internal/fs"
	"github.com/ba58ajbse/envcraft/internal/utils"
)

// CommentOptions holds the options for adding a new environment variable.
type CommentOptions struct {
	Value    string
	FilePath string
	Line     int
}

// CommentCmd represents the command for adding a new environment variable to a file.
type CommentCmd struct {
	Options  CommentOptions
	OrgLines []string
}

func Run(args []string) error {
	options, err := ParseCommentOptions(args)
	if err != nil {
		return err
	}
	cmd, err := NewCommentCmd(options)
	if err != nil {
		return err
	}
	err = cmd.Exec()
	if err != nil {
		return err
	}
	return nil
}

// NewCommentCmd creates a new CommentCmd instance with the specified file path.
func NewCommentCmd(options *CommentOptions) (*CommentCmd, error) {
	if options.FilePath == "" {
		return nil, errors.New("file path is required")
	}

	return &CommentCmd{
		Options:  *options,
		OrgLines: []string{},
	}, nil
}

// Exec is the main function that processes the add command using the provided options.
func (a *CommentCmd) Exec() error {
	err := a.readLines()
	if err != nil {
		return err
	}

	newLines, err := a.makeNewLines()
	if err != nil {
		return err
	}

	err = a.apply(newLines)
	if err != nil {
		return err
	}

	return nil
}

// readLines reads all lines from the file specified in CommentCmd and stores them in OrgLines.
func (c *CommentCmd) readLines() error {
	lines, err := fs.ReadLines(c.filePath())
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", c.filePath(), err)
	}
	c.OrgLines = lines

	return nil
}

// makeNewLines generates the new lines to be written to the file after adding the new variable.
func (c *CommentCmd) makeNewLines() ([]string, error) {
	newLines := slices.Clone(c.OrgLines)
	if c.insertLineNum() == 0 {
		if utils.IsEmptyOrBlank(newLines) {
			newLines = []string{c.value()}
			return newLines, nil
		}
		if utils.EndsWithoutNewline(newLines) {
			newLines[len(newLines)-1] += "\n" // Add a newline if the last line does not end with a newline
		}
		newLines = slices.Insert(newLines, len(newLines), c.value())
		return newLines, nil
	}

	if c.insertLineNum() > len(c.OrgLines) {
		if utils.EndsWithoutNewline(newLines) {
			newLines[len(newLines)-1] += "\n" // Add a newline if the last line does not end with a newline
		}
		emplyLines := slices.Repeat([]string{"\n"}, c.insertLineNum()-len(c.OrgLines)-1)
		newLines = slices.Concat(newLines, emplyLines, []string{c.value()})
		return newLines, nil
	}

	if len(newLines) == 1 && newLines[0] == "" {
		// If the file is empty, add the new line
		newLines = []string{c.value()}
		return newLines, nil
	}
	newLines = slices.Insert(newLines, c.insertLineNum()-1, c.value()+"\n")
	return newLines, nil
}

// apply writes the new lines to the file, overwriting the original content.
func (c *CommentCmd) apply(newLines []string) error {
	if err := fs.WriteLines(c.filePath(), newLines); err != nil {
		return fmt.Errorf("error writing to file %s: %w", c.filePath(), err)
	}

	return nil
}

func (a *CommentCmd) filePath() string {
	return a.Options.FilePath
}

func (a *CommentCmd) value() string {
	return "# " + a.Options.Value
}

func (a *CommentCmd) insertLineNum() int {
	return a.Options.Line
}

// ParseCommentOptions parses command-line arguments and returns an CommentOptions struct.
func ParseCommentOptions(opts []string) (*CommentOptions, error) {
	fs := flag.NewFlagSet("comment", flag.ContinueOnError)
	file := fs.String("f", "", "Path to .env file")
	line := fs.Int("l", 0, "Line number to insert comment (optional)")

	if len(opts) < 1 {
		return nil, errors.New("comment required")
	}
	value := opts[0]
	flags := opts[1:]

	if err := fs.Parse(flags); err != nil {
		return nil, err
	}
	if *file == "" {
		fmt.Println("Error: -f flag is required")
		fs.Usage()
		return nil, errors.New("file path is required")
	}

	if *line < 0 {
		fmt.Println("Error: -l must be a non-negative integer")
		fs.Usage()
		return nil, errors.New("line number must be a non-negative integer")
	}

	return &CommentOptions{
		Value:    value,
		FilePath: *file,
		Line:     *line,
	}, nil
}
