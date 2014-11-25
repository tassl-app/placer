package placer

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

type Endpoint interface {
	Url() (string, error)
	Fetch() error
	ResultLen() int
	json.Unmarshaler
}

const (
	nearbySearchBase = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"
	autocompleteBase = "https://maps.googleapis.com/maps/api/place/autocomplete/json"
)

type NearbySearchEndpoint struct {
	Request  *NearbySearchRequest
	Response *SearchResponse
}

func NewNearbySearchEndpoint(key string, lat, lng, radius float64) *NearbySearchEndpoint {
	request := &NearbySearchRequest{
		Key: key,
		Location: &Location{
			Latitude:  lat,
			Longitude: lng,
		},
		Radius: radius,
	}
	return &NearbySearchEndpoint{
		Request: request,
	}
}

func (e *NearbySearchEndpoint) Url() (string, error) {
	u, err := url.Parse(nearbySearchBase)
	if err != nil {
		return "", err
	}
	q := u.Query()
	req := e.Request
	q.Set("key", req.Key)
	if req.Location == nil {
		return "", errors.New("Location required")
	} else {
		q.Set("location", req.Location.Formatted())
	}
	q.Set("radius", fmt.Sprintf("%f", req.Radius))
	if len(req.Keyword) > 0 {
		q.Set("keyword", req.Keyword)
	}
	if len(req.Language) > 0 {
		q.Set("language", req.Language)
	}
	if req.MinpriceValid && req.Minprice >= 0 && req.Minprice <= 4 {
		q.Set("minprice", strconv.Itoa(int(req.Minprice)))
	}
	if req.MaxpriceValid && req.Maxprice >= 0 && req.Maxprice <= 4 {
		q.Set("maxprice", strconv.Itoa(int(req.Maxprice)))
	}
	if len(req.Name) > 0 {
		q.Set("name", req.Name)
	}
	if req.OpenNowValid {
		q.Set("opennow", strconv.FormatBool(req.OpenNow))
	}
	if len(req.Rankby) > 0 {
		q.Set("rankby", req.Rankby)
	}
	if len(req.Types) > 0 {
		typeStr := ""
		for i, searchType := range req.Types {
			if i > 0 {
				typeStr += "|"
			}
			typeStr += searchType
		}
		q.Set("types", typeStr)
	}
	if len(req.Pagetoken) > 0 {
		q.Set("pagetoken", req.Pagetoken)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (e *NearbySearchEndpoint) Fetch() error {
	return Fetch(e)
}

func (e *NearbySearchEndpoint) UnmarshalJSON(data []byte) error {
	resp := new(SearchResponse)
	err := json.Unmarshal(data, resp)
	if err != nil {
		return err
	}
	e.Response = resp
	return nil
}

func (e *NearbySearchEndpoint) ResultLen() int {
	return len(e.Response.Results)
}

type NearbySearchRequest struct {
	// Required
	Key      string
	Location *Location
	Radius   float64
	// Optional
	Keyword       string
	Language      string
	Minprice      int64
	MinpriceValid bool
	Maxprice      int64
	MaxpriceValid bool
	Name          string
	OpenNow       bool
	OpenNowValid  bool
	Rankby        string
	Types         []string
	Pagetoken     string
}

type SearchResponse struct {
	Attributions []string       `json:"html_attributions"`
	Results      []SearchResult `json:"results"`
	Status       string         `json:"status"`
}

type SearchResult struct {
	Geometry     Geometry     `json:"geometry"`
	Icon         string       `json:"icon"`
	Id           string       `json:"id"`
	Name         string       `json:"name"`
	OpeningHours OpeningHours `json:"opening_hours"`
	Photos       []Photo      `json:"photos"`
	PlaceId      string       `json:"place_id"`
	Scope        string       `json:"scope"`
	Reference    string       `json:"reference"`
	Types        []string     `json:"types"`
	Vicinity     string       `json:"vicinity"`
}

type Geometry struct {
	Location Location `json:"location"`
}

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

func (l Location) Formatted() string {
	return fmt.Sprintf("%f,%f", l.Latitude, l.Longitude)
}

type OpeningHours struct {
	OpenNow bool `json:"open_now"`
}

type Photo struct {
	Reference    string   `json:"photo_reference"`
	Height       int64    `json:"height"`
	Width        int64    `json:"width"`
	Attributions []string `json:"html_attributions"`
}

type AutocompleteEndpoint struct {
	Request  *AutocompleteRequest
	Response *AutocompleteResponse
}

func NewAutocompleteEndpoint(key, input string, lat, lng, radius float64) *AutocompleteEndpoint {
	request := &AutocompleteRequest{
		Input: input,
		Key:   key,
		Location: &Location{
			Latitude:  lat,
			Longitude: lng,
		},
		Radius: radius,
	}
	return &AutocompleteEndpoint{
		Request: request,
	}
}

func (e *AutocompleteEndpoint) Url() (string, error) {
	u, err := url.Parse(autocompleteBase)
	if err != nil {
		return "", err
	}
	q := u.Query()
	req := e.Request
	q.Set("input", req.Input)
	q.Set("key", req.Key)
	if req.Offset > 0 {
		q.Set("offset", strconv.Itoa(int(req.Offset)))
	}
	if req.Location != nil {
		q.Set("location", req.Location.Formatted())
	}
	if req.Radius > 0 {
		q.Set("radius", fmt.Sprintf("%f", req.Radius))
	}
	if len(req.Language) > 0 {
		q.Set("language", req.Language)
	}
	if len(req.Types) > 0 {
		q.Set("types", req.Types)
	}
	if len(req.Components) > 0 {
		q.Set("components", req.Components)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (e *AutocompleteEndpoint) Fetch() error {
	return Fetch(e)
}

func (e *AutocompleteEndpoint) UnmarshalJSON(data []byte) error {
	resp := new(AutocompleteResponse)
	err := json.Unmarshal(data, resp)
	if err != nil {
		return err
	}
	e.Response = resp
	return nil
}

func (e *AutocompleteEndpoint) ResultLen() int {
	return len(e.Response.Predictions)
}

type AutocompleteRequest struct {
	// Required
	Input string
	Key   string
	// Optional
	Offset     int64
	Location   *Location
	Radius     float64
	Language   string
	Types      string
	Components string
}

type AutocompleteResponse struct {
	Status      string       `json:"status"`
	Predictions []Prediction `json:"predictions"`
}

type Prediction struct {
	Description       string             `json:"description"`
	Id                string             `json:"id"`
	MatchedSubstrings []MatchedSubstring `json:"matched_substrings"`
	PlaceId           string             `json:"place_id"`
	Reference         string             `json:"reference"`
	Terms             []Term             `json:"terms"`
	Types             []string           `json:"types"`
}

type MatchedSubstring struct {
	Length int64 `json:"length"`
	Offset int64 `json:"offset"`
}

type Term struct {
	Offset int64  `json:"offset"`
	Value  string `json:"value"`
}
