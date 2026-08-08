package pokeapi

import (
	"encoding/json"
	"fmt"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationAreasResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

func (c *Client) ListLocationAreas(pageURL string) (LocationAreasResponse, error) {
	url := baseURL + "/location-area"
	if pageURL != "" {
		url = pageURL
	}
	res, err := c.httpClient.Get(url)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return LocationAreasResponse{}, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	var areas LocationAreasResponse
	if err := json.NewDecoder(res.Body).Decode(&areas); err != nil {
		return LocationAreasResponse{}, err
	}
	return areas, nil
}
