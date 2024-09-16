package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/acehotel33/pokedex-cli/globals"
)

type config struct {
	NextURL     string `json:"next"`
	PreviousURL string `json:"previous"`
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

var cliCommandMap map[string]cliCommand

func init() {
	cliCommandMap = map[string]cliCommand{
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
		"map": {
			name:        "map",
			description: "Displays first 20 locations of map, consecutive calls display next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations of map",
			callback:    commandMapB,
		},
	}
}

func commandHelp(conf *config) error {
	fmt.Println()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for key, cmd := range cliCommandMap {
		fmt.Printf("%s: %s\n", key, cmd.description)
	}
	fmt.Println()

	return nil
}

func commandExit(conf *config) error {
	fmt.Println("Exiting")
	os.Exit(0)
	return nil
}

func commandMap(conf *config) error {
	// needs to change to getLocationAreasAll
	nextURL := conf.NextURL

	locations, err := getLocationsAll(nextURL, conf)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB(conf *config) error {
	// needs to change to getLocationAreasAll
	previousURL := conf.PreviousURL

	if previousURL == "" {
		fmt.Println("Already on Page 1")
		return nil
	}

	locations, err := getLocationsAll(previousURL, conf)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	// Id     int            `json:"id"`
	// Region Region         `json:"region"`
	// Areas  []LocationArea `json:"areas"`
}

type LocationsAll struct {
	Count       int        `json:"count"`
	NextURL     string     `json:"next"`
	PreviousURL string     `json:"previous"`
	Results     []Location `json:"results"`
}

type Region struct {
	Id        int        `json:"id"`
	Locations []Location `json:"locations"`
	Name      string     `json:"name"`
}

type LocationArea struct {
	Id       int      `json:"id"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
}

func getLocationsAll(url string, conf *config) ([]Location, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []Location{}, fmt.Errorf("could not create GET request - %w", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []Location{}, fmt.Errorf("could not perform GET request - %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []Location{}, fmt.Errorf("status code of response is not OK - %v", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body - %w", err)
	}

	var locationsAll LocationsAll
	if err := json.Unmarshal(body, &locationsAll); err != nil {
		return nil, fmt.Errorf("could not decode json body into instance of LocationsAll - %w", err)
	}

	conf.NextURL = locationsAll.NextURL
	conf.PreviousURL = locationsAll.PreviousURL

	return locationsAll.Results, nil
}

func main() {
	// Initialize configuration
	conf := &config{
		NextURL:     globals.LocationsAllURL,
		PreviousURL: "",
	}

	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Pokedex > ")

		scanner.Scan()
		line := scanner.Text()

		if command, exists := cliCommandMap[line]; exists {
			if err := command.callback(conf); err != nil {
				fmt.Println("Error")
			}
		} else {
			fmt.Println("Unknown command. Type 'help' for a list of commands.")
		}

	}
}
