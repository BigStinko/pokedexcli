package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/BigStinko/pokedexcli/pokecache"
)

type Client struct {
	cache pokecache.Cache
	httpClient http.Client
}

const (
	baseURL = "https://pokeapi.co/api/v2"
)

func NewClient(timeout time.Duration, cacheInterval time.Duration) Client {
	return Client{
		cache: pokecache.NewCache(cacheInterval),
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	url := baseURL + "/location-area?offset=0&limit=20"
	if pageURL != nil {
		url = *pageURL
	}

	if value, ok := c.cache.Get(url); ok {
		locations := RespShallowLocations{}
		err := json.Unmarshal(value, &locations)
		if err != nil {
			return RespShallowLocations{}, err
		}

		fmt.Println("---cached value---")
		return locations, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespShallowLocations{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return RespShallowLocations{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return RespShallowLocations{}, err
	}

	locations := RespShallowLocations{}
	err = json.Unmarshal(data, &locations)
	if err != nil {
		return RespShallowLocations{}, err
	}

	c.cache.Add(url, data)

	return locations, nil
}

func (c *Client) GetLocation(locationName string) (Location, error) {
	url := baseURL + "/location-area/" + locationName

	if value, ok := c.cache.Get(url); ok {
		locationResp := Location{}
		err := json.Unmarshal(value, &locationResp)
		if err != nil {
			return Location{}, err
		}

		fmt.Println("---cached value---")
		return locationResp, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Location{}, err
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Location{}, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return Location{}, err
	}

	locationResponse := Location{}
	err = json.Unmarshal(data, &locationResponse)
	if err != nil {
		return Location{}, err
	}

	c.cache.Add(url, data)

	return locationResponse, nil
}
