package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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
		"explore": {
			Name:        "explore",
			Description: "Explore the specifed location for pokemon",
			Callback:    commandExploreArea,
		},
	}
}

func commandHelp(conf *globals.Config, params []string) error {
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

func commandExit(conf *globals.Config, params []string) error {
	fmt.Println("Exiting")
	os.Exit(0)
	return nil
}

func commandMap(conf *globals.Config, params []string) error {
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

func commandMapB(conf *globals.Config, params []string) error {
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

func commandExploreArea(conf *globals.Config, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("missing argument")
	}
	location := params[0]
	if location == " " {
		return fmt.Errorf("empty location given")
	}
	fullURL := globals.LocationsAllURL + location
	pokemonSplice, err := api.ExploreArea(fullURL, conf)
	if err != nil {
		return fmt.Errorf("could not explore area - %w", err)
	}
	for _, pokemon := range pokemonSplice {
		fmt.Println(pokemon)
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
		words := strings.Split(line, " ")

		// fmt.Println("Here are the words:")
		// for _, word := range words {
		// 	fmt.Printf("arg: %s\n", word)
		// }
		// fmt.Println("End of words")

		if command, exists := cliCommandMap[words[0]]; exists {
			params := []string{}
			if len(words) > 1 {
				params = words[1:]
			}
			if err := command.Callback(conf, params); err != nil {
				fmt.Printf("Could not perform command: %v\n", err)
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}

	}
}
