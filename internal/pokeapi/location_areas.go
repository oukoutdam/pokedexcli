package pokeapi

import (
	"net/url"
)

type NamedResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreasResponse struct {
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []NamedResource `json:"results"`
}

type LocationArea struct {
	ID                   int           `json:"id"`
	Name                 string        `json:"name"`
	GameIndex            int           `json:"game_index"`
	Location             NamedResource `json:"location"`
	EncounterMethodRates []struct {
		EncounterMethod NamedResource `json:"encounter_method"`
		VersionDetails  []struct {
			Rate    int           `json:"rate"`
			Version NamedResource `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Names []struct {
		Language NamedResource `json:"language"`
		Name     string        `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon        NamedResource `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int             `json:"chance"`
				ConditionValues []NamedResource `json:"condition_values"`
				MaxLevel        int             `json:"max_level"`
				Method          NamedResource   `json:"method"`
				MinLevel        int             `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int           `json:"max_chance"`
			Version   NamedResource `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func (c *Client) ListLocationAreas(pageURL string) (LocationAreasResponse, error) {
	endpoint := baseURL + "/location-area"
	if pageURL != "" {
		endpoint = pageURL
	}

	return getResource[LocationAreasResponse](c, endpoint)
}

func (c *Client) GetLocationArea(areaName string) (LocationArea, error) {
	return getResource[LocationArea](c, baseURL+"/location-area/"+url.PathEscape(areaName))
}
