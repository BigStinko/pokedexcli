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

type resource interface {
	Location | Pokemon | RespShallowLocations
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

func getResource[R resource](url string, c *Client) (R, error) {
	var zero R
	
	if value, ok := c.cache.Get(url); ok {
		var res R
		err := json.Unmarshal(value, &res)
		if err != nil {	return zero, err }

		fmt.Println("---cached value---")
		return res, nil
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {	return zero, err }

	response, err := c.httpClient.Do(request)
	if err != nil {	return zero, err }
	defer response.Body.Close()

	dat, err := io.ReadAll(response.Body)
	if err != nil { return zero, err }

	var res R
	err = json.Unmarshal(dat, &res)
	if err != nil { return zero, err }

	c.cache.Add(url, dat)
	
	return res, nil
}

func (c *Client) ListLocations(pageURL *string) (RespShallowLocations, error) {
	url := baseURL + "/location-area?offset=0&limit=20"
	if pageURL != nil {
		url = *pageURL
	}
	
	return getResource[RespShallowLocations](url, c)
}

func (c *Client) GetLocation(locationName string) (Location, error) {
	url := baseURL + "/location-area/" + locationName

	return getResource[Location](url, c)
}

func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := baseURL + "/pokemon/" + pokemonName

	return getResource[Pokemon](url, c)
}
