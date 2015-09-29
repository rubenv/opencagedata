package opencagedata

import (
	"os"
	"strings"
	"testing"
)

type UrlTest struct {
	Query    string
	Params   *GeocodeParams
	Expected string
}

func TestGeocodeUrl(t *testing.T) {
	tests := []UrlTest{
		UrlTest{
			Query:    "Steystraat 30",
			Expected: "https://api.opencagedata.com/geocode/v1/json?key=test&q=Steystraat+30",
		},
		UrlTest{
			Query: "Steystraat 30",
			Params: &GeocodeParams{
				CountryCode: "BE",
			},
			Expected: "https://api.opencagedata.com/geocode/v1/json?countrycode=be&key=test&q=Steystraat+30",
		},
	}

	geocoder := NewGeocoder("test")

	for _, test := range tests {
		url := geocoder.geocodeUrl(test.Query, test.Params)
		if url != test.Expected {
			t.Errorf("Bad result for %#v (%#v), got %#v, expected %#v", test.Query, test.Params, url, test.Expected)
		}
	}
}

func TestErrorRequest(t *testing.T) {
	geocoder := NewGeocoder("test")
	r, err := geocoder.Geocode("Leuven", nil)
	if err == nil {
		t.Error("Expected error")
	}
	if r != nil {
		t.Error("Expected no result")
	}

	geo_err := err.(*GeocodeError)
	if geo_err.Result.Status.Code != 403 {
		t.Error("Expected bad API key")
	}
}

func TestGeocode(t *testing.T) {
	key := os.Getenv("API_KEY")
	if key == "" {
		t.Skip("No API_KEY set, skipping online tests")
	}

	geocoder := NewGeocoder(key)
	r, err := geocoder.Geocode("Fonteinstraat 75, Leuven", nil)
	if err != nil {
		t.Error("Unexpected error")
	}
	if len(r.Results) == 0 {
		t.Error("Expected results")
	}

	if r.Results[0].Geometry.Latitude == 0 {
		t.Error("Expected a coordinate")
	}
	if r.Results[0].Confidence != 10 {
		t.Error("Geocoder suddenly feeling very insecure")
	}
}

func TestParams(t *testing.T) {
	key := os.Getenv("API_KEY")
	if key == "" {
		t.Skip("No API_KEY set, skipping online tests")
	}

	geocoder := NewGeocoder(key)
	r, err := geocoder.Geocode("Grote Markt", &GeocodeParams{
		CountryCode: "be",
	})
	if err != nil {
		t.Error("Unexpected error")
	}
	if len(r.Results) == 0 {
		t.Error("Expected results")
	}
	if !strings.Contains(r.Results[0].Formatted, "Belgium") {
		t.Error("Expected a result in Belgium")
	}

	r, err = geocoder.Geocode("Grote Markt", &GeocodeParams{
		CountryCode: "nl",
	})
	if err != nil {
		t.Error("Unexpected error")
	}
	if len(r.Results) == 0 {
		t.Error("Expected results")
	}
	if !strings.Contains(r.Results[0].Formatted, "The Netherlands") {
		t.Error("Expected a result in The Netherlands")
	}
}
