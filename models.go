package placer

import (
	"fmt"
	"net/url"
	"strconv"
)

type RequestEndpoint interface {
	Url() (string, error)
	Fetch() (*Response, error)
}

const nearbySearchEndpoint = "https://maps.googleapis.com/maps/api/place/nearbysearch/json"

type NearbySearchRequest struct {
	// Required
	Key      string
	Location Location
	Radius   int64
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

func NewNearbySearchRequest(key string, lat, lng float64, radius int64) *NearbySearchRequest {
	return &NearbySearchRequest{
		Key: key,
		Location: Location{
			Latitude:  lat,
			Longitude: lng,
		},
		Radius: radius,
	}
}

func (ns *NearbySearchRequest) Url() (string, error) {
	u, err := url.Parse(nearbySearchEndpoint)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("key", ns.Key)
	q.Set("location", ns.Location.Formatted())
	q.Set("radius", strconv.Itoa(int(ns.Radius)))
	if len(ns.Keyword) > 0 {
		q.Set("keyword", ns.Keyword)
	}
	if len(ns.Language) > 0 {
		q.Set("language", ns.Language)
	}
	if ns.MinpriceValid && ns.Minprice >= 0 && ns.Minprice <= 4 {
		q.Set("minprice", strconv.Itoa(int(ns.Minprice)))
	}
	if ns.MaxpriceValid && ns.Maxprice >= 0 && ns.Maxprice <= 4 {
		q.Set("maxprice", strconv.Itoa(int(ns.Maxprice)))
	}
	if len(ns.Name) > 0 {
		q.Set("name", ns.Name)
	}
	if ns.OpenNowValid {
		q.Set("opennow", strconv.FormatBool(ns.OpenNow))
	}
	if len(ns.Rankby) > 0 {
		q.Set("rankby", ns.Rankby)
	}
	if len(ns.Types) > 0 {
		typeStr := ""
		for i, searchType := range ns.Types {
			if i > 0 {
				typeStr += "|"
			}
			typeStr += searchType
		}
		q.Set("types", typeStr)
	}
	if len(ns.Pagetoken) > 0 {
		q.Set("pagetoken", ns.Pagetoken)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (ns *NearbySearchRequest) Fetch() (*Response, error) {
	return FetchEndpoint(ns)
}

type Response struct {
	Attributions []string `json:"html_attributions"`
	Results      []Result `json:"results"`
	Status       string   `json:"status"`
}

type Result struct {
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
