package main

import (
	"fmt"
	"os"

	"github.com/ba58ajbse/envcraft/internal/commands/add"
	"github.com/ba58ajbse/envcraft/internal/commands/comment"
	"github.com/ba58ajbse/envcraft/internal/commands/delete"
	"github.com/ba58ajbse/envcraft/internal/commands/update"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: envcraft [command] [flags]")
		os.Exit(1)
	}

	command := os.Args[1]
	opts := os.Args[2:]
	commands := map[string]func([]string) error{
		"add":     add.Run,
		"update":  update.Run,
		"delete":  delete.Run,
		"comment": comment.Run,
	}
	cmd, ok := commands[command]
	if !ok {
		fmt.Println("Usage: envcraft [add|update|delete|comment] [flags]")
		os.Exit(1)
	}

	if err := cmd(opts); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅", command, "completed.")
}
