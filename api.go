package placer

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const OneMileRadius = 1609.34

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
	return response, err
}
