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
