package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/oukoutdam/pokedexcli/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache      *pokecache.Cache
}

func NewClient(timeout time.Duration, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
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

func getResource[T any](c *Client, endpoint string) (T, error) {
	var output T
	data, err := c.getBody(endpoint)
	if err != nil {
		return output, err
	}

	if err := json.Unmarshal(data, &output); err != nil {
		return output, err
	}

	return output, nil
}
