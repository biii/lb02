package yelp

import (
	"errors"
	"fmt"

	"github.com/guregu/null"
)

// CoordinateOptions are used with complex searches for locations.
// The geographic coordinate format is defined as:
// ll=latitude,longitude,accuracy,altitude,altitude_accuracy
type CoordinateOptions struct {
	Latitude         null.Float // Latitude of geo-point to search near (required)
	Longitude        null.Float // Longitude of geo-point to search near (required)
	Accuracy         null.Float // Accuracy of latitude, longitude (optional)
	Altitude         null.Float // Altitude (optional)
	AltitudeAccuracy null.Float // Accuracy of altitude (optional)
}

// getParameters will reflect over the values of the given
// struct, and provide a type appropriate set of querystring parameters
// that match the defined values.
func (o CoordinateOptions) getParameters() (params map[string]string, err error) {
	params = make(map[string]string)
	// coordinate requires at least a latitude and longitude - others are option
	if !o.Latitude.Valid || !o.Longitude.Valid {
		return nil, errors.New("latitude and longitude are required fields for a coordinate based search")
	}
	params["latitude"] = fmt.Sprintf("%v", o.Latitude.Float64)
	params["longitude"] = fmt.Sprintf("%v", o.Longitude.Float64)
	params["locale"] = "zh-TW"

	ll := fmt.Sprintf("%v,%v", o.Latitude.Float64, o.Longitude.Float64)
	if o.Accuracy.Valid {
		ll += fmt.Sprintf(",%v", o.Accuracy.Float64)
	}
	if o.Altitude.Valid {
		ll += fmt.Sprintf(",%v", o.Altitude.Float64)
	}
	if o.AltitudeAccuracy.Valid {
		ll += fmt.Sprintf(",%v", o.AltitudeAccuracy.Float64)
	}

	return params, nil
//	return map[string]string{
//		"ll": ll,
//	}, nil
}
