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
	NewLines []string
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
		NewLines: []string{},
	}, nil
}

// Exec is the main function that processes the add command using the provided options.
func (a *AddCmd) Exec() error {
	err := a.readLines()
	if err != nil {
		return err
	}

	err = a.makeNewLines()
	if err != nil {
		return err
	}

	// TODO: print diff & masked lines

	// TODO: Uncomment the confirmation prompt when ready
	// if !a.confirmSave() {
	// 	return
	// }

	err = a.apply()
	if err != nil {
		return err
	}

	return nil
}

// readLines reads all lines from the file specified in AddCmd and stores them in OrgLines.
func (a *AddCmd) readLines() error {
	lines, err := fs.ReadLines(a.filePath())
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", a.filePath(), err)
	}
	a.OrgLines = lines

	return nil
}

// makeNewLines generates the new lines to be written to the file after adding the new variable.
func (a *AddCmd) makeNewLines() error {
	a.NewLines = slices.Clone(a.OrgLines)
	if a.insertLineNum() == 0 {
		if utils.IsEmptyOrBlank(a.NewLines) {
			a.NewLines = []string{a.keyAndValue()}
			return nil
		}
		if utils.EndsWithoutNewline(a.NewLines) {
			a.NewLines[len(a.NewLines)-1] += "\n" // Add a newline if the last line does not end with a newline
		}
		a.NewLines = slices.Insert(a.NewLines, len(a.NewLines), a.keyAndValue())
		return nil
	}

	if a.insertLineNum() > len(a.OrgLines) {
		if utils.EndsWithoutNewline(a.NewLines) {
			a.NewLines[len(a.NewLines)-1] += "\n" // Add a newline if the last line does not end with a newline
		}
		emplyLines := slices.Repeat([]string{"\n"}, a.insertLineNum()-len(a.OrgLines)-1)
		a.NewLines = slices.Concat(a.NewLines, emplyLines, []string{a.keyAndValue()})
		return nil
	}

	if len(a.NewLines) == 1 && a.NewLines[0] == "" {
		// If the file is empty, add the new line
		a.NewLines = []string{a.keyAndValue()}
		return nil
	}
	a.NewLines = slices.Insert(a.NewLines, a.insertLineNum()-1, a.keyAndValue()+"\n")
	return nil
}

// apply writes the new lines to the file, overwriting the original content.
func (a *AddCmd) apply() error {
	if err := fs.WriteLines(a.filePath(), a.NewLines); err != nil {
		return fmt.Errorf("error writing to file %s: %w", a.filePath(), err)
	}

	return nil
}

func (a *AddCmd) filePath() string {
	return a.Options.FilePath
}

func (a *AddCmd) keyAndValue() string {
	return a.Options.Key + "=" + strconv.Quote(a.Options.Value)
}

func (a *AddCmd) insertLineNum() int {
	return a.Options.Line
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
