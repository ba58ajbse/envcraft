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
	var key string
	var val string
	var insideQuote bool
	var err error
	var quote rune

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if insideQuote {
			runes := []rune(line)
			last := runes[len(runes)-1]
			if last == quote {
				v := strings.TrimRight(line, string(quote))
				m[key] = val + "\n" + v
				quote = 0
				insideQuote = false
				continue
			}
			m[key] = val + "\n" + line
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			if _, ok := m[key]; !ok {
				return nil, fmt.Errorf("invalid line: %q", line)
			}
			m[key] = val + "\n" + line
			continue
		}
		key = extractVarKey(parts[0])
		quote, insideQuote = betweenQuote(parts[1])
		val, insideQuote, err = extractVarValue(parts[1], quote, insideQuote)
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

func extractVarKey(key string) string {
	if strings.Contains(key, "export") {
		key = strings.TrimPrefix(key, "export")
	}

	return strings.TrimSpace(key)
}

// extractVarValue extracts the value part of an environment variable line, handling quotes and inline comments.
// Returns the cleaned value or an error if the syntax is invalid.
func extractVarValue(val string, quote rune, insideQuote bool) (string, bool, error) {
	val = strings.TrimSpace(val)
	if insideQuote {
		runes := []rune(val)
		len := len(runes)
		for i := 1; i < len; i++ {
			if runes[i] == quote {
				v := runes[:i]
				return strings.TrimLeft(string(v), string(quote)), false, nil
			}
		}
		val = strings.TrimLeft(val, string(quote))
		return val, true, nil
	}
	if hasQuote(val) {
		fmt.Println("\ncall has quote!!!", val)
		// If quoted, only the quoted string is used as the value
		lidx := lastQuote(val)
		if lidx < 0 {
			return "", false, fmt.Errorf("invalid syntax env value: %v", val)
		}
		val = strings.Trim(val[:lidx+1], "\"'`")
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

	return val, false, nil
}

// hasQuote returns true if the string starts with a single or double quote.
func hasQuote(str string) bool {
	return len(str) > 0 && (str[0] == '"' || str[0] == '\'' || str[0] == '`')
}

func betweenQuote(str string) (rune, bool) {
	if len(str) == 0 {
		return 0, false
	}

	if !hasQuote(str) {
		return 0, false
	}

	runes := []rune(str)
	last := runes[len(runes)-1]
	if str[0] == '"' {
		ret := last != '"'
		return '"', ret
	}

	if str[0] == '\'' {
		return '\'', last != '\''
	}

	if str[0] == '`' {
		return '`', last != '`'
	}

	return 0, false
}

// lastQuote returns the last index of a quote character in the string, or -1 if not found.
func lastQuote(str string) int {
	quote := str[0]
	if quote == '\'' {
		return strings.LastIndex(str, `'`)
	}
	if quote == '"' {
		return strings.LastIndex(str, `"`)
	}
	if quote == '`' {
		return strings.LastIndex(str, "`")
	}

	return -1
}
