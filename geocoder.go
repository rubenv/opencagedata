package opencagedata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const endpoint = "https://api.opencagedata.com/geocode/v1/"

type Geocoder struct {
	Key string
}

type GeocodeParams struct {
	CountryCode string
}

type GeocodeResult struct {
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`

	Rate struct {
		Limit     int   `json:"limit"`
		Remaining int   `json:"remaining"`
		Reset     int64 `json:"reset"`
	} `json:"rate"`

	Results []*GeocodeResultItem `json:"results"`
}

type GeocodeResultItem struct {
}

// Returned when geocoding fails, contains the actual response
type GeocodeError struct {
	Result *GeocodeResult
}

func (err *GeocodeError) Error() string {
	return fmt.Sprintf("%s: %s", err.Result.Status.Code, err.Result.Status.Message)
}

func NewGeocoder(key string) *Geocoder {
	return &Geocoder{
		Key: key,
	}
}

// Geocode a query
//
// The params parameter is optional, feel free to pass nil when no specific options are needed.
func (g *Geocoder) Geocode(query string, params *GeocodeParams) (*GeocodeResult, error) {
	u := g.geocodeUrl(query, params)
	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GeocodeResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	if result.Status.Code != 200 {
		return nil, &GeocodeError{Result: &result}
	}
	return &result, nil
}

// Split out for testing purposes
func (g *Geocoder) geocodeUrl(query string, params *GeocodeParams) string {
	u, _ := url.Parse(endpoint)
	u.Path += "json"

	q := u.Query()
	q.Set("q", query)
	q.Set("key", g.Key)
	if params != nil {
		if params.CountryCode != "" {
			q.Set("countrycode", strings.ToLower(params.CountryCode))
		}
	}

	u.RawQuery = q.Encode()
	return u.String()
}