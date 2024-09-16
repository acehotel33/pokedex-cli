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
	nextURL := conf.NextURL

	locations, err := getLocationAreasAll(nextURL, conf)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

func commandMapB(conf *config) error {
	previousURL := conf.PreviousURL

	if previousURL == "" {
		fmt.Println("Already on Page 1")
		return nil
	}

	locations, err := getLocationAreasAll(previousURL, conf)
	if err != nil {
		return err
	}

	for _, location := range locations {
		fmt.Println(location.Name)
	}
	return nil
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreasAll struct {
	Count       int            `json:"count"`
	NextURL     string         `json:"next"`
	PreviousURL string         `json:"previous"`
	Results     []LocationArea `json:"results"`
}

func getLocationAreasAll(url string, conf *config) ([]LocationArea, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return []LocationArea{}, fmt.Errorf("could not create GET request - %w", err)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []LocationArea{}, fmt.Errorf("could not perform GET request - %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return []LocationArea{}, fmt.Errorf("status code of response is not OK - %v", res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read response body - %w", err)
	}

	var locationAreasAll LocationAreasAll
	if err := json.Unmarshal(body, &locationAreasAll); err != nil {
		return nil, fmt.Errorf("could not decode json body into instance of LocationsAll - %w", err)
	}

	conf.NextURL = locationAreasAll.NextURL
	conf.PreviousURL = locationAreasAll.PreviousURL

	return locationAreasAll.Results, nil
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
