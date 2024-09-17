package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/acehotel33/pokedex-cli/globals"
)

func GetLocationAreasAll(url string, conf *globals.Config) ([]globals.LocationArea, error) {
	if body, exists := globals.Cache.Get(url); exists {
		var locationAreasAll globals.LocationAreasAll
		if err := json.Unmarshal(body, &locationAreasAll); err != nil {
			return nil, fmt.Errorf("could not decode json body into instance of LocationsAll - %w", err)
		}
		conf.NextURL = locationAreasAll.NextURL
		conf.PreviousURL = locationAreasAll.PreviousURL

		return locationAreasAll.Results, nil

	} else {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create GET request - %w", err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("could not perform GET request - %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("status code of response is not OK - %v", res.Status)
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read response body - %w", err)
		}

		globals.Cache.Add(url, body)

		var locationAreasAll globals.LocationAreasAll
		if err := json.Unmarshal(body, &locationAreasAll); err != nil {
			return nil, fmt.Errorf("could not decode json body into instance of LocationsAll - %w", err)
		}

		conf.NextURL = locationAreasAll.NextURL
		conf.PreviousURL = locationAreasAll.PreviousURL

		return locationAreasAll.Results, nil
	}

}

func ExploreArea(url string, conf *globals.Config) ([]string, error) {
	if body, exists := globals.Cache.Get(url); exists {
		var area globals.Area
		if err := json.Unmarshal(body, &area); err != nil {
			return nil, fmt.Errorf("could not unmarshal cached body into Area struct - %w", err)
		}
		pokemonEncounters := area.PokemonEncounters
		pokemonSlice := []string{}
		for _, item := range pokemonEncounters {
			pokemonSlice = append(pokemonSlice, item.Pokemon.Name)
		}

		return pokemonSlice, nil

	} else {

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create GET request - %w", err)
		}

		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("could not perform request - %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			if res.StatusCode == 404 {
				return nil, fmt.Errorf("location not found")
			} else {
				return nil, fmt.Errorf("non-okay response status: %v", res.Status)
			}
		}

		var area globals.Area
		if err := json.NewDecoder(res.Body).Decode(&area); err != nil {
			return nil, fmt.Errorf("could not decode JSON into Area struct - %w", err)
		}

		pokemonEncounters := area.PokemonEncounters
		pokemonSlice := []string{}
		for _, item := range pokemonEncounters {
			pokemonSlice = append(pokemonSlice, item.Pokemon.Name)
		}

		return pokemonSlice, nil
	}
}

func GetPokemon(url string, conf *globals.Config) (globals.Pokemon, error) {
	res, err := http.Get(url)
	if err != nil {
		return globals.Pokemon{}, err
	}
	defer res.Body.Close()

	var pokemon globals.Pokemon
	if err := json.NewDecoder(res.Body).Decode(&pokemon); err != nil {
		return globals.Pokemon{}, err
	}

	return pokemon, nil
}
