/*
Go bindings for OpenCage Geocoder

http://geocoder.opencagedata.com/

Usage:

	geocoder := opencagedata.NewGeocoder("my-api-key")

Simple queries:

	result, err := geocoder.Geocode("Fonteinstraat, Leuven", nil)
	if err != nil {
		// Handle error
	}
	// Do something with result

Extra options can be passed as well:

	result, err := geocoder.Geocode("Fonteinstraat, Leuven", &opencagedata.GeocodeParams{
		CountryCode: "be",
	})
	if err != nil {
		// Handle error
	}
	// Do something with result

*/
package opencagedata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const endpoint = "https://api.opencagedata.com/geocode/v1/"

type Geocoder struct {
	// API Key, found on the developer dashboard: https://developer.opencagedata.com/
	Key string

	// Disable rate limiting (not recommended).
	//
	// This library will sleep automatically to avoid hitting the rate limit
	// when using the free plan.
	DisableRateLimitSleep bool

	lock  sync.Mutex
	sleep time.Time
}

type GeocodeParams struct {
	// Country hint
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

	Results []GeocodeResultItem `json:"results"`
}

type GeocodeResultItem struct {
	Confidence int      `json:"confidence"`
	Formatted  string   `json:"formatted"`
	Geometry   Geometry `json:"geometry"`

	Bounds struct {
		NorthEast Geometry `json:"northeast"`
		SouthWest Geometry `json:"southwest"`
	} `json:"bounds"`
}

type Geometry struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
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
	g.lock.Lock()
	defer g.lock.Unlock()

	sleep := g.sleep.Sub(time.Now())
	if sleep > 0 {
		time.Sleep(sleep)
	}

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

	if !g.DisableRateLimitSleep && result.Rate.Limit > 0 {
		reset := time.Unix(result.Rate.Reset, 0)
		untilReset := reset.Sub(time.Now())
		delay := time.Duration(float64(untilReset+1) / (float64(result.Rate.Remaining) + 1))
		g.sleep = time.Now().Add(delay)
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
