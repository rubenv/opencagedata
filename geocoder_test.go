package opencagedata

import (
	"os"
	"strings"
	"testing"

	"github.com/cheekybits/is"
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
	is := is.New(t)

	geocoder := NewGeocoder("test")
	r, err := geocoder.Geocode("Leuven", nil)
	is.Err(err)
	is.Nil(r)

	geo_err, ok := err.(*GeocodeError)
	is.True(ok)
	is.Equal(geo_err.Result.Status.Code, 403)
}

func TestGeocode(t *testing.T) {
	is := is.New(t)

	key := os.Getenv("API_KEY")
	if key == "" {
		t.Skip("No API_KEY set, skipping online tests")
	}

	geocoder := NewGeocoder(key)
	r, err := geocoder.Geocode("Fonteinstraat 75, Leuven", nil)
	is.NoErr(err)
	is.NotNil(r)
	is.NotEqual(len(r.Results), 0)

	is.NotEqual(r.Results[0].Geometry.Latitude, 0)
	is.Equal(r.Results[0].Confidence, 10)
}

func TestParams(t *testing.T) {
	is := is.New(t)

	key := os.Getenv("API_KEY")
	if key == "" {
		t.Skip("No API_KEY set, skipping online tests")
	}

	geocoder := NewGeocoder(key)
	r, err := geocoder.Geocode("Grote Markt", &GeocodeParams{
		CountryCode: "be",
	})
	is.NoErr(err)
	is.NotNil(r)
	is.NotEqual(len(r.Results), 0)
	is.True(strings.Contains(r.Results[0].Formatted, "Belgium"))

	r, err = geocoder.Geocode("Grote Markt", &GeocodeParams{
		CountryCode: "nl",
	})
	is.NoErr(err)
	is.NotNil(r)
	is.NotEqual(len(r.Results), 0)
	is.True(strings.Contains(r.Results[0].Formatted, "The Netherlands"))
}
