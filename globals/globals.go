package globals

var LocationsAllURL = "https://pokeapi.co/api/v2/location-area/"

type Config struct {
	NextURL     string `json:"next"`
	PreviousURL string `json:"previous"`
}

type CliCommand struct {
	Name        string
	Description string
	Callback    func(*Config) error
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
