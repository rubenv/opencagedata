package opencagedata

import "testing"

type Test struct {
	Query       string
	Params      *GeocodeParams
	ExpectedURL string
}

func TestGeocode(t *testing.T) {
	tests := []Test{
		Test{
			Query:       "Steystraat 30",
			ExpectedURL: "https://api.opencagedata.com/geocode/v1/json?key=test&q=Steystraat+30",
		},
		Test{
			Query: "Steystraat 30",
			Params: &GeocodeParams{
				CountryCode: "BE",
			},
			ExpectedURL: "https://api.opencagedata.com/geocode/v1/json?countrycode=be&key=test&q=Steystraat+30",
		},
	}

	geocoder := NewGeocoder("test")

	for _, test := range tests {
		url := geocoder.geocodeUrl(test.Query, test.Params)
		if url != test.ExpectedURL {
			t.Errorf("Bad result for %#v (%#v), got %#v, expected %#v", test.Query, test.Params, url, test.ExpectedURL)
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
