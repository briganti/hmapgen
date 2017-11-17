package hmapgen

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type bingResResources struct {
	Elevations []int `json:"elevations"`
}

type bingResResourceSets struct {
	Resources []bingResResources `json:"resources"`
}

type bingResponse struct {
	StatusCode   int                   `json:"statusCode"`
	ErrorDetails []string              `json:"errorDetails"`
	ResourceSets []bingResResourceSets `json:"resourceSets"`
}

// Call BING Elevation API to retrieve the height of each lat/lng passed
// https://msdn.microsoft.com/fr-fr/library/jj158961.aspx
func getBingEvelations(key string, lat, lng []string) ([]float64, error) {
	heights := []float64{}
	points := []string{}

	for i := 0; i < len(lat); i++ {
		currentPoint := lat[i] + "," + lng[i]
		points = append(points, currentPoint)
	}

	urlStr := "http://dev.virtualearth.net/REST/v1/Elevation/List?points=" + strings.Join(points, ",") + "&heights=sealevel&key=" + key

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
	var data bingResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return heights, err
	}

	if data.StatusCode != 200 {
		err := errors.New(data.ErrorDetails[0])
		return heights, err
	}

	// Convert to an height array
	for _, val := range data.ResourceSets[0].Resources[0].Elevations {
		heights = append(heights, float64(val))
	}

	return heights, nil
}

// Returns the max number of point allowed per API request
func getBingMaxNumberOfPoints() int {
	return 100
}
