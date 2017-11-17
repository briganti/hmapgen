package hmapgen

import (
	"fmt"
	"github.com/paulmach/go.geo"
	"gopkg.in/cheggaaa/pb.v1"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strconv"
)

// An Options structure is used for generating the heighmap
type Options struct {
	// The elevation service used (ign|bing)
	Service string
	// The API key for the elevation service
	Key string
	// Distance in meter between grid points
	Precision float64
	// Image file name
	File string
}

// Response struct  is returned for a successful GenerateHeightMap call
type Response struct {
	// Filename of the heightmap
	Filename string
	// Height offset bewenn the lowest point and the highest point (in meter)
	HeightOffset int
}

// GenerateHeightMap generates a heightmap image of a given area location.
func GenerateHeightMap(area []string, options Options) (Response, error) {
	var response Response
	var err error

	heights := []float64{}
	bufferLng, bufferLat := []string{}, []string{}

	// Get Surface object
	surface := getSurface(area, options.Precision)
	nbPoints := surface.Width * surface.Height
	var PointsPerCall int
	if options.Service == "bing" {
		PointsPerCall = getBingMaxNumberOfPoints()
	} else {
		PointsPerCall = getIGNMaxNumberOfPoints()
	}
	nbRequiredCalls := int(math.Ceil(float64(nbPoints) / float64(PointsPerCall)))

	fmt.Printf("Surface: %dx%d - %d calls required\n", surface.Width, surface.Height, nbRequiredCalls)

	// Init progress bar
	bar := pb.StartNew(nbRequiredCalls)

	// Reverse loop => this makes it a lot easier while batching points
	for i := (nbPoints - 1); i >= 0; i-- {
		x, y := getCoordinateByIndex((nbPoints-1)-i, surface)
		point := surface.PointAt(x, y)
		bufferLng = append(bufferLng, strconv.FormatFloat(point.Lng(), 'f', 6, 64))
		bufferLat = append(bufferLat, strconv.FormatFloat(point.Lat(), 'f', 6, 64))

		// Enough points for 1 request to IGN / NE point
		if i%PointsPerCall == 0 {
			var newHeights []float64
			if options.Service == "bing" {
				newHeights, err = getBingEvelations(options.Key, bufferLat, bufferLng)
			} else {
				newHeights, err = getIGNEvelations(options.Key, bufferLat, bufferLng)
			}

			if err != nil {
				return response, err
			}

			for _, val := range newHeights {
				heights = append(heights, val)
			}

			bar.Increment()
			bufferLng, bufferLat = nil, nil
		}
	}
	bar.Finish()

	img := image.NewGray16(image.Rect(0, 0, surface.Width, surface.Height))

	// Convert heights to gray and drawing them
	minZ, maxZ := getMinMax(heights)
	// img.set(0,0) => top left corner
	for i, z := range heights {
		x, y := getCoordinateByIndex(i, surface)
		grayColor := getGrayColor(z, minZ, maxZ)

		img.SetGray16(x, (surface.Height-y)-1, color.Gray16{Y: grayColor})
	}

	file := "./output.png"
	if options.File != "" {
		file = options.File
	}
	out, err := os.Create(file)
	if err != nil {
		return response, err
	}

	// Setup Encoder
	var Enc png.Encoder
	Enc.CompressionLevel = -1 // no compression

	// Write hmap
	err = Enc.Encode(out, img)
	if err != nil {
		return response, err
	}

	response = Response{
		Filename:     file,
		HeightOffset: int(maxZ - minZ),
	}

	return response, nil
}

// Create a Geo Surface object from the input lng/lat
func getSurface(areaPtr []string, gridInterval float64) *geo.Surface {
	south, _ := strconv.ParseFloat(areaPtr[0], 64)
	west, _ := strconv.ParseFloat(areaPtr[1], 64)
	north, _ := strconv.ParseFloat(areaPtr[2], 64)
	east, _ := strconv.ParseFloat(areaPtr[3], 64)

	bound := geo.NewBound(west, east, south, north)

	width := int(bound.GeoWidth() / gridInterval)
	height := int(bound.GeoHeight() / gridInterval)

	return geo.NewSurface(bound, width, height)
}

// Return the x/y coordinates of the surface grid, index 0 being the sw point
// and last index being the ne point.
func getCoordinateByIndex(index int, surface *geo.Surface) (int, int) {
	x := int(math.Mod(float64(index), float64(surface.Width)))
	y := (surface.Height - int(index/surface.Width)) - 1

	return x, y
}

// Return the main/max values of given array
func getMinMax(array []float64) (float64, float64) {
	var min float64 = array[0]
	var max float64 = array[0]

	for _, val := range array {
		if min > val {
			min = val
		}
		if max < val {
			max = val
		}
	}

	fmt.Printf("min max: %.2fx%.2f\n", min, max)
	return min, max
}

// Get the gray color of val on a floor-ceil scale
func getGrayColor(val, floor, ceil float64) uint16 {
	var scale float64 = math.MaxUint16 / (ceil - floor)
	var color float64 = (val - floor) * scale

	return uint16(color)
}
