package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
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

	data, err := c.getBody(url)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	var areas LocationAreasResponse
	if err := json.Unmarshal(data, &areas); err != nil {
		return LocationAreasResponse{}, err
	}

	return areas, nil
}

func (c *Client) getBody(url string) ([]byte, error) {
	if cacheData, ok := c.cache.Get(url); ok {
		return cacheData, nil
	}

	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return nil, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !json.Valid(data) {
		return nil, fmt.Errorf("invalid JSON response from %s", url)
	}

	c.cache.Add(url, data)
	return data, nil
}
