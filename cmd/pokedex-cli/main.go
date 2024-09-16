package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/acehotel33/pokedex-cli/globals"
	"github.com/acehotel33/pokedex-cli/internal/api"
)

var cliCommandMap map[string]globals.CliCommand

func init() {
	cliCommandMap = map[string]globals.CliCommand{
		"help": {
			Name:        "help",
			Description: "Displays a help message",
			Callback:    commandHelp,
		},
		"exit": {
			Name:        "exit",
			Description: "Exit the Pokedex",
			Callback:    commandExit,
		},
		"map": {
			Name:        "map",
			Description: "Displays first 20 locations of map, consecutive calls display next 20 locations",
			Callback:    commandMap,
		},
		"mapb": {
			Name:        "mapb",
			Description: "Displays the previous 20 locations of map",
			Callback:    commandMapB,
		},
	}
}

func commandHelp(conf *globals.Config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for key, cmd := range cliCommandMap {
		fmt.Printf("%s: %s\n", key, cmd.Description)
	}
	fmt.Println()

	return nil
}

func commandExit(conf *globals.Config) error {
	fmt.Println("Exiting")
	os.Exit(0)
	return nil
}

func commandMap(conf *globals.Config) error {
	nextURL := conf.NextURL

	locations, err := api.GetLocationAreasAll(nextURL, conf)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB(conf *globals.Config) error {
	previousURL := conf.PreviousURL

	if previousURL == "" {
		fmt.Println("Already on Page 1")
		return nil
	}

	locations, err := api.GetLocationAreasAll(previousURL, conf)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

func main() {
	// Initialize configuration
	conf := &globals.Config{
		NextURL:     globals.LocationsAllURL,
		PreviousURL: "",
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")

		scanner.Scan()
		line := scanner.Text()

		if command, exists := cliCommandMap[line]; exists {
			if err := command.Callback(conf); err != nil {
				fmt.Println("Error")
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}

	}
}
