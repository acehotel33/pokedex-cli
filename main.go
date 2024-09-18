package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

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
	defer fmt.Println(".\n.")
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

func helperCatch(pokemon globals.Pokemon) bool {

	baseExperience := pokemon.BaseExperience
	fmt.Printf(".\nBase experience: %v\n", baseExperience)
	if baseExperience < 100 {
		baseExperience = 100
	} else if baseExperience > 500 {
		baseExperience = 500
	}
	fmt.Printf("Adjusted Base experience: %v\n", baseExperience)

	time.Sleep(2 * time.Second)

	var chance float64
	chance = math.Sqrt(80.0) / math.Sqrt(float64(baseExperience))
	fmt.Printf(".\nChance before multiplier: %.0f percent\n", chance*100)
	mult := (600.0 - float64(baseExperience)) / 600.0
	fmt.Printf("Multiplier: %.2f\n", mult)
	chance = chance * float64(mult) * 100
	chanceInt := int(chance)
	fmt.Printf("Chance after multiplier: %v percent\n", chanceInt)
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	time.Sleep(2 * time.Second)

	fmt.Printf(".\nChance of success: %v percent\n", chanceInt)
	fmt.Printf(".\nAttempting to catch %s...\n", pokemon.Name)

	for i := 0; i < 4; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println(".")
	}
	return r.Intn(100) < chanceInt
}

func addToPokedex(conf *globals.Config, pokemon globals.Pokemon) error {
	if _, exists := conf.Pokedex[pokemon.Name]; exists {
		return fmt.Errorf("pokemon %s already in pokedex", pokemon.Name)
	}

	conf.Pokedex[pokemon.Name] = pokemon
	return nil
}

func commandPokedex(conf *globals.Config, params []string) error {
	defer fmt.Println(".\n.")
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
