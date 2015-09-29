# opencagedata - Go bindings for OpenCage Geocoder

[![Build Status](https://travis-ci.org/rubenv/opencagedata.svg?branch=master)](https://travis-ci.org/rubenv/opencagedata) [![GoDoc](https://godoc.org/github.com/rubenv/opencagedata?status.png)](https://godoc.org/github.com/rubenv/opencagedata)

http://geocoder.opencagedata.com/

## Installation
```
go get github.com/rubenv/opencagedata
```

Import into your application with:

```go
import "github.com/rubenv/opencagedata"
```

## Usage

```go
geocoder := opencagedata.NewGeocoder("my-api-key")
```

Simple queries:
```go
result, err := geocoder.Geocode("Fonteinstraat, Leuven", nil)
if err != nil {
    // Handle error
}
// Do something with result
```

Extra options can be passed as well:
```go
result, err := geocoder.Geocode("Fonteinstraat, Leuven", &opencagedata.GeocodeParams{
    CountryCode: "be",
})
if err != nil {
    // Handle error
}
// Do something with result
```

## License

    (The MIT License)

    Copyright (C) 2015 by Ruben Vermeersch <ruben@rocketeer.be>

    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:

    The above copyright notice and this permission notice shall be included in
    all copies or substantial portions of the Software.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
    THE SOFTWARE.
