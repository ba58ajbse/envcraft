package add

import (
	"errors"
	"flag"
	"fmt"
	"slices"
	"strconv"

	"github.com/ba58ajbse/envcraft/internal/fs"
	"github.com/ba58ajbse/envcraft/internal/utils"
)

// AddOptions holds the options for adding a new environment variable.
type AddOptions struct {
	Key      string
	Value    string
	FilePath string
	Line     int
}

// AddCmd represents the command for adding a new environment variable to a file.
type AddCmd struct {
	Options  AddOptions
	OrgLines []string
}

func Run(args []string) error {
	options, err := ParseAddOptions(args)
	if err != nil {
		return err
	}
	cmd, err := NewAddCmd(options)
	if err != nil {
		return err
	}
	err = cmd.Exec()
	if err != nil {
		return err
	}
	return nil
}

// NewAddCmd creates a new AddCmd instance with the specified file path.
func NewAddCmd(options *AddOptions) (*AddCmd, error) {
	if options.FilePath == "" {
		return nil, errors.New("file path is required")
	}

	return &AddCmd{
		Options:  *options,
		OrgLines: []string{},
	}, nil
}

// Exec is the main function that processes the add command using the provided options.
func (c *AddCmd) Exec() error {
	err := c.readLines()
	if err != nil {
		return err
	}

	newlines, err := c.makeNewLines()
	if err != nil {
		return err
	}

	// TODO: print diff & masked lines

	// TODO: Uncomment the confirmation prompt when ready
	// if !a.confirmSave() {
	// 	return
	// }

	err = c.apply(newlines)
	if err != nil {
		return err
	}

	return nil
}

// readLines reads all lines from the file specified in AddCmd and stores them in OrgLines.
func (c *AddCmd) readLines() error {
	lines, err := fs.ReadLines(c.filePath())
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", c.filePath(), err)
	}
	c.OrgLines = lines

	return nil
}

// makeNewLines generates the new lines to be written to the file after adding the new variable.
func (c *AddCmd) makeNewLines() ([]string, error) {
	newLines := slices.Clone(c.OrgLines)
	if c.insertLineNum() == 0 {
		if utils.IsEmptyOrBlank(newLines) {
			return []string{c.keyAndValue()}, nil
		}
		if utils.EndsWithoutNewline(newLines) {
			newLines[len(newLines)-1] += "\n" // Add a newline if the last line does not end with a newline
		}
		newLines = slices.Insert(newLines, len(newLines), c.keyAndValue())
		return newLines, nil
	}

	if c.insertLineNum() > len(c.OrgLines) {
		if utils.EndsWithoutNewline(newLines) {
			newLines[len(newLines)-1] += "\n" // Add a newline if the last line does not end with a newline
		}
		emplyLines := slices.Repeat([]string{"\n"}, c.insertLineNum()-len(c.OrgLines)-1)
		newLines = slices.Concat(newLines, emplyLines, []string{c.keyAndValue()})
		return newLines, nil
	}

	if len(newLines) == 1 && newLines[0] == "" {
		// If the file is empty, add the new line
		return []string{c.keyAndValue()}, nil
	}
	newLines = slices.Insert(newLines, c.insertLineNum()-1, c.keyAndValue()+"\n")
	return newLines, nil
}

// apply writes the new lines to the file, overwriting the original content.
func (c *AddCmd) apply(newLines []string) error {
	if err := fs.WriteLines(c.filePath(), newLines); err != nil {
		return fmt.Errorf("error writing to file %s: %w", c.filePath(), err)
	}

	return nil
}

func (c *AddCmd) filePath() string {
	return c.Options.FilePath
}

func (c *AddCmd) keyAndValue() string {
	return c.Options.Key + "=" + strconv.Quote(c.Options.Value)
}

func (c *AddCmd) insertLineNum() int {
	return c.Options.Line
}

// ParseAddOptions parses command-line arguments and returns an AddOptions struct.
func ParseAddOptions(opts []string) (*AddOptions, error) {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	file := fs.String("f", "", "Path to .env file")
	line := fs.Int("l", 0, "Line number to insert the variable (optional)")

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

	if *line < 0 {
		fmt.Println("Error: -l must be a non-negative integer")
		fs.Usage()
		return nil, errors.New("line number must be a non-negative integer")
	}

	return &AddOptions{
		Key:      key,
		Value:    value,
		FilePath: *file,
		Line:     *line,
	}, nil
}

// func (a *AddCmd) confirmSave() bool {
// 	var input string
// 	fmt.Print("Do you want to save the changes? (y/N): ")
// 	_, err := fmt.Scanln(&input)
// 	if err != nil {
// 		fmt.Printf("Error reading input: %v\n", err)
// 		os.Exit(1)
// 	}
// 	input = strings.ToLower(input)
// 	return input == "y"
// }
