package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/ba58ajbse/envcraft/cmd/add"
	"github.com/ba58ajbse/envcraft/cmd/delete"
	"github.com/ba58ajbse/envcraft/cmd/update"
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
	case "update":
		cmdUpdate(opts)
	case "delete":
		cmdDelete(opts)
	case "comment":
		cmdComment()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func cmdAdd(args []string) {
	options, err := add.ParseAddOptions(args)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	addCmd, err := add.NewAddCmd(options)
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

func cmdUpdate(args []string) {
	options, err := update.ParseUpdateOptions(args)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	updateCmd, err := update.NewUpdateCmd(options)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	err = updateCmd.Exec()
	if err != nil {
		if errors.Is(err, update.ErrNoUpdated) {
			fmt.Println("No matching key found. No update performed.")
			os.Exit(0)
		}
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("✅ Successfully updated.")
}

func cmdDelete(args []string) {
	options, err := delete.ParseDeleteOptions(args)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	deleteCmd, err := delete.NewDeleteCmd(options)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	err = deleteCmd.Exec()
	if err != nil {
		if errors.Is(err, update.ErrNoUpdated) {
			fmt.Println("No matching key found. No deletion performed.")
			os.Exit(0)
		}
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	fmt.Println("✅ Successfully deleted.")
}

func cmdComment() {
	// TODO: フラグ定義
	fmt.Println("comment command executed")
}
