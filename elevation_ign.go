package hmapgen

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type ignResElevation struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
	Z   float64 `json:"z"`
	Acc float64 `json:"acc"`
}

type ignResponse struct {
	Elevations []ignResElevation `json:"elevations"`
	Http       struct {
		Status int    `json:"status"`
		Error  string `json:"error"`
	} `json:"http"`
}

// Call IGN Elevation API to retrieve the height of each lat/lng passed
func getIGNEvelations(key string, lat, lng []string) ([]float64, error) {
	heights := []float64{}
	urlStr := "https://wxs.ign.fr/" + key + "/alti/rest/elevation.json?lon=" + strings.Join(lng, ",") + "&lat=" + strings.Join(lat, ",") + "&delimiter=,&output=json&zonly=false"

	// Do the request
	res, err := http.Get(urlStr)
	if err != nil {
		return heights, err
	}

	// Read json response
	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return heights, err
	}

	// Parse the json response
	var data ignResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return heights, err
	}

	if data.Http.Error != "" {
		err := errors.New(data.Http.Error)
		return heights, err
	}

	// Convert to an height array
	for _, val := range data.Elevations {
		z := val.Z
		// IGN use -9999 Z value for unknown height. Default to 0
		if z <= 0 {
			z = 0
		}
		heights = append(heights, z)
	}

	return heights, nil
}

// Returns the max number of point allowed per API request
func getIGNMaxNumberOfPoints() int {
	return 50
}
