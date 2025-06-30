package main

import (
	"fmt"
	"os"

	"github.com/ba58ajbse/envcraft/cmd"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: envcraft [command] [flags]")
		os.Exit(1)
	}

	command := os.Args[1]
	opts := os.Args[2:]
	switch command {
	case "add":
		cmdAdd(opts)
	case "set":
		cmdUpdate()
	case "comment":
		cmdComment()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func cmdAdd(args []string) {
	options, err := cmd.ParseAddOptions(args)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	addCmd, err := cmd.NewAddCmd(options)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	err = addCmd.Exec()
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("✅ Successfully added.")
}

func cmdUpdate() {
	// TODO: フラグ定義
	fmt.Println("insert command executed")
}

func cmdComment() {
	// TODO: フラグ定義
	fmt.Println("comment command executed")
}
