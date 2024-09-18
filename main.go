package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/acehotel33/pokedex-cli/globals"
	"github.com/acehotel33/pokedex-cli/internal/api"
)

var cliCommandMap map[string]globals.CliCommand

func main() {
	// Initialize configuration
	conf := &globals.Config{
		NextURL:     globals.LocationsAllURL,
		PreviousURL: "",
		Pokedex:     make(map[string]globals.Pokemon),
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")

		scanner.Scan()
		line := scanner.Text()
		words := strings.Split(line, " ")

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
		"catch": {
			Name:        "catch",
			Description: "Try to catch the specified pokemon",
			Callback:    commandCatch,
		},
		"pokedex": {
			Name:        "pokedex",
			Description: "Display Pokedex of current Pokemon",
			Callback:    commandPokedex,
		},
		"inspect": {
			Name:        "inspect",
			Description: "Inspect Pokemon's attributes if already caught",
			Callback:    commandInspect,
		},
	}
}

func commandHelp(conf *globals.Config, params []string) error {
	defer fmt.Println(".\n.")
	fmt.Println(".\n.")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for key, cmd := range cliCommandMap {
		fmt.Printf("%s: %s\n", key, cmd.Description)
	}

	return nil
}

func commandExit(conf *globals.Config, params []string) error {
	defer fmt.Println(".\n.")
	fmt.Println(".\n.")
	fmt.Println("Exiting")
	os.Exit(0)
	return nil
}

func commandMap(conf *globals.Config, params []string) error {
	defer fmt.Println(".\n.")
	fmt.Println(".\n.")
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
	defer fmt.Println(".\n.")
	fmt.Println(".\n.")

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
	defer fmt.Println(".\n.")
	fmt.Println(".\n.")

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
	fmt.Printf("Exploring %s...\n", location)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range pokemonSplice {
		fmt.Printf("- %s\n", pokemon)
	}
	return nil
}

func commandCatch(conf *globals.Config, params []string) error {
	fmt.Println(".\n.")

	if len(params) < 1 {
		return fmt.Errorf("catch command missing arguments")
	}
	toCatch := params[0]
	fullURL := globals.PokemonURL + toCatch
	pokemon, err := api.GetPokemon(fullURL, conf)
	if err != nil {
		return fmt.Errorf("could not find pokemon - %w", err)
	}

	if _, exists := conf.Pokedex[pokemon.Name]; exists {
		return fmt.Errorf("pokemon %s already in pokedex", pokemon.Name)
	}

	if result := helperCatch(pokemon); result {
		if err := addToPokedex(conf, pokemon); err != nil {
			return fmt.Errorf("pokemon %s already caught", pokemon.Name)
		}
		fmt.Printf("Result: Success! You caught %v!\n", pokemon.Name)
	} else {
		fmt.Printf("Result: Oh no! %v slipped away!\n", pokemon.Name)
	}

	time.Sleep(time.Second)
	fmt.Println(".\n.")
	fmt.Println("Current Pokedex:")
	if err := commandPokedex(conf, params); err != nil {
		return fmt.Errorf("could not display pokedex - %w", err)
	}
	return nil
}

func commandInspect(conf *globals.Config, params []string) error {
	fmt.Println(".\n.")
	if len(params) < 1 {
		return fmt.Errorf("missing parameter")
	}
	pokemonToInspect := params[0]
	if poke, exists := conf.Pokedex[pokemonToInspect]; !exists {
		return fmt.Errorf("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %s\n", poke.Name)
		fmt.Printf("Height: %v\n", poke.Height)
		fmt.Printf("Weight: %v\n", poke.Weight)

		fmt.Println("Stats:")
		stats := poke.Stats
		for _, stat := range stats {
			statName := stat.Stat.Name
			statValue := stat.BaseStat
			fmt.Printf("  -%s: %v\n", statName, statValue)
		}

		fmt.Println("Types:")
		pTypes := poke.Types
		for _, pType := range pTypes {
			pTypeName := pType.Type.Name
			fmt.Printf("  - %s\n", pTypeName)
		}
	}
	fmt.Println(".\n.")
	return nil
}

func commandPokedex(conf *globals.Config, params []string) error {
	fmt.Println(".\n.")
	if len(conf.Pokedex) == 0 {
		fmt.Println("Pokedex is empty!")
		return nil
	}
	for key := range conf.Pokedex {
		fmt.Printf("- %s -\n", key)
	}
	return nil
}

func helperCatch(pokemon globals.Pokemon) bool {

	baseExperience := pokemon.BaseExperience
	fmt.Printf("Base Experience of %s: %v\n", pokemon.Name, baseExperience)

	// Normalize baseExperience to a minimum of 100 and maximum of 500
	if baseExperience < 100 {
		baseExperience = 100
	} else if baseExperience > 500 {
		baseExperience = 500
	}

	// Calculate a realistic chance using baseExperience
	// Higher baseExperience should result in a lower chance
	chance := (500.0 - float64(baseExperience)) / 5.0

	fmt.Printf("Chance of success: %.1f percent\n", chance)

	// Use a random seed
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	// Determine success based on random number
	result := r.Intn(100)

	time.Sleep(2 * time.Second)
	fmt.Println(".\n.")
	fmt.Printf("Throwing a Pokeball at %s...", pokemon.Name)
	fmt.Println("")
	for i := 0; i < 4; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println(".")
	}

	return result < int(chance)

}

func addToPokedex(conf *globals.Config, pokemon globals.Pokemon) error {
	if _, exists := conf.Pokedex[pokemon.Name]; exists {
		return fmt.Errorf("pokemon %s already in pokedex", pokemon.Name)
	}

	conf.Pokedex[pokemon.Name] = pokemon
	return nil
}
