package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/acehotel33/pokedex-cli/globals"
)

func GetLocationAreasAll(url string, conf *globals.Config) ([]globals.LocationArea, error) {
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

	var locationAreasAll globals.LocationAreasAll
	if err := json.Unmarshal(body, &locationAreasAll); err != nil {
		return nil, fmt.Errorf("could not decode json body into instance of LocationsAll - %w", err)
	}

	conf.NextURL = locationAreasAll.NextURL
	conf.PreviousURL = locationAreasAll.PreviousURL

	return locationAreasAll.Results, nil
}
