[![Go Report Card](https://goreportcard.com/badge/github.com/michaelvial/hmapgen)](https://goreportcard.com/report/github.com/michaelvial/hmapgen)

# HMAPGEN

A simple library in Go to generate heightmap images of a given location.

## Prerequisite

In order to generate the heightmap, you must have at least an Api key for the following services:
- [IGN Elevation API](https://geoservices.ign.fr/documentation/geoservices/alti.html)
- [BING Elevation API](https://msdn.microsoft.com/fr-fr/library/jj158961.aspx)

## USAGE

```
import (
  "github.com/michaelvial/hmapgen"
)

// South, West, North and East Lat/Lng coordinates 
area := ["45.142724","6.0888359","45.7316515","7.1368806"]

// Available options
options := hmapgen.Options{
    Service: "bing",  // Either bing or ign
    Key: "xxx",  // The API key for the elevation service
    Precision: 1000,  // Interval between each points (in meter)
    File: "./images/output.png", // Heightmap filename
}

// Heightmap generation
hmapinfo, err := hmapgen.GenerateHeightMap(area, options)

// Filename
fmt.Printf("Output file:", hmapinfo.Filename)
// Diff between the lowest point and the highest point (in meter)
fmt.Printf("Height offset:", hmapinfo.HeightOffset)

```

## Licenses

All source code is licensed under the [MIT License](https://raw.github.com/michaelvial/hmapgen/master/LICENSE).