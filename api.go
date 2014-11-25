package placer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const OneMileRadius = 1609.34

var ErrNoResults = errors.New("No results found")

func FetchEndpoint(endpoint RequestEndpoint) (*Response, error) {
	url, err := endpoint.Url()
	if err != nil {
		return nil, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	response := new(Response)
	err = json.Unmarshal(body, response)
	if err != nil {
		return nil, err
	}
	if len(response.Results) == 0 {
		return nil, ErrNoResults
	}
	return response, nil
}
