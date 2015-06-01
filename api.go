package placer

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

const OneMileRadius = 1609.34

var ErrNoResults = errors.New("No results found")

func Fetch(endpoint Endpoint) error {
	url, err := endpoint.Url()
	log.Printf("url: %+v\n", url)
	if err != nil {
		return err
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = endpoint.UnmarshalJSON(body)
	if err != nil {
		return err
	}
	if endpoint.ResultLen() == 0 {
		return ErrNoResults
	}
	return nil
}

type PlacesAPIManager struct {
	Key string
}

func NewPlacesAPIManager(key string) *PlacesAPIManager {
	return &PlacesAPIManager{key}
}

func (m *PlacesAPIManager) NearbySearch(lat, lng, radius float64) (*SearchResponse, error) {
	endpoint := NewNearbySearchEndpoint(m.Key, lat, lng, radius)
	err := endpoint.Fetch()
	return endpoint.Response, err
}

func (m *PlacesAPIManager) Autocomplete(input string, offset int64, lat, lng, radius float64, language, types, components string) (*AutocompleteResponse, error) {
	endpoint := NewAutocompleteEndpoint(m.Key, input, lat, lng, radius)
	err := endpoint.Fetch()
	return endpoint.Response, err
}

func (m *PlacesAPIManager) PlaceDetails(placeId string) (*PlaceDetailsResponse, error) {
	endpoint := NewPlaceDetailsEndpoint(m.Key, placeId)
	err := endpoint.Fetch()
	return endpoint.Response, err
}
