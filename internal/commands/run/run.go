package run

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
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

	// 3) Load env
	err := godotenv.Load(*envFile)
	if err != nil {
		return err
	}

	// 4) Exec the command, inheriting stdin/stdout/stderr
	cmd := exec.Command(cmdName, cmdParams...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
