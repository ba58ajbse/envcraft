package run

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Run parses flags, loads environment variables from a file, and executes the specified command with those variables set.
func Run(args []string) error {
	// 1) Define flags
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	envFile := fs.String("f", ".env", "path to .env file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	// 2) The rest of args after flags is the command to exec
	cmdArgs := fs.Args()
	if len(cmdArgs) == 0 {
		return fmt.Errorf("usage: envcraft run [-f .env] -- <your-command>")
	}
	cmdName := cmdArgs[0]
	cmdParams := cmdArgs[1:]

	file, close, err := openFile(*envFile)
	if err != nil {
		return err
	}
	defer close()

	envMap, err := parseEnv(file)
	if err != nil {
		return fmt.Errorf("parse env error: %w", err)
	}

	// 3) Load env
	loadEnv(envMap)

	// 4) Exec the command, inheriting stdin/stdout/stderr
	cmd := exec.Command(cmdName, cmdParams...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// openFile opens the specified file and returns the file pointer, a close function, and an error if any.
// The returned close function should be called to properly close the file.
func openFile(filePath string) (*os.File, func() error, error) {
	envFile, err := os.Open(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	close := envFile.Close

	return envFile, close, nil
}

// parseEnv parses environment variables from the given io.Reader and returns them as a map.
// Lines starting with '#' or empty lines are ignored. Inline comments are supported.
func parseEnv(r io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(r)
	m := make(map[string]string)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %q", line)
		}
		key := strings.TrimSpace(parts[0])
		val, err := extractVarValue(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("extract var value error: %w", err)
		}

		m[key] = val
	}
	return m, scanner.Err()
}

// loadEnv sets environment variables from the given map if they are not already set in the current environment.
func loadEnv(envMap map[string]string) {
	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	// NOTE: Overwrite feature is not implemented yet
	for key, val := range envMap {
		if !currentEnv[key] {
			_ = os.Setenv(key, val)
		}
	}
}

// extractVarValue extracts the value part of an environment variable line, handling quotes and inline comments.
// Returns the cleaned value or an error if the syntax is invalid.
func extractVarValue(val string) (string, error) {
	val = strings.TrimSpace(val)
	if hasQuote(val) {
		// If quoted, only the quoted string is used as the value
		lidx := lastQuote(val)
		if lidx < 0 {
			return "", fmt.Errorf("invalid syntax env value: %v", val)
		}
		val = strings.Trim(val[:lidx+1], `"'`)
	} else {
		// If not quoted, treat everything after # as a comment
		if idx := strings.Index(val, "#"); idx >= 0 {
			prefix := val[:idx]
			r, _ := utf8.DecodeLastRuneInString(prefix)
			if unicode.IsSpace(r) {
				val = strings.TrimSpace(val[:idx])
			}
		}
	}

	return val, nil
}

// hasQuote returns true if the string starts with a single or double quote.
func hasQuote(str string) bool {
	return len(str) > 0 && (str[0] == '"' || str[0] == '\'')
}

// lastQuote returns the last index of a quote character in the string, or -1 if not found.
func lastQuote(str string) int {
	if lidx := strings.LastIndex(str, `"`); lidx > 0 {
		return lidx
	}
	if lidx := strings.LastIndex(str, `'`); lidx > 0 {
		return lidx
	}

	return -1
}
