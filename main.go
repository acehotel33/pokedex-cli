package main

import (
	"bufio"
	"fmt"
	"os"
)

var helpText string = `

Welcome to the Pokedex!
Usage:

help: Displays a help message
exit: Exit the Pokedex

`

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var cliCommandMap = map[string]cliCommand{
	"help": {
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	},
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
}

func commandHelp() error {
	fmt.Println(helpText)
	return nil
}

func commandExit() error {
	fmt.Println("Exiting")
	os.Exit(0)
	return nil
}

func main() {
	// fmt.Println("Hello World")
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")

		scanner.Scan()

		line := scanner.Text()
		switch line {
		case "help", "exit":
			cliCommandMap[line].callback()
		default:
			fmt.Println("Unknown command. Type 'help' for a list of commands")
		}
	}
}
