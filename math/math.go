package math

import "github.com/danielsoneg/bartbot/types"
import "math"

func dist(from, to types.Loc) float64 {
	from.Lat = radians(from.Lat)
	to.Lat = radians(to.Lat)
	d_lat := (to.Lat - from.Lat) / 2
	d_lon := radians(to.Lon-from.Lon) / 2
	haversin := math.Pow(math.Sin(d_lat), 2) + math.Cos(to.Lat)
	haversin = haversin * math.Cos(from.Lat) * math.Pow(math.Sin(d_lon), 2)
	return math.Asin(math.Sqrt(haversin))
}

func Nearest(from types.Loc, a, b types.Station) types.Station {
	if dist(from, a.Loc) < dist(from, b.Loc) {
		return a
	} else {
		return b
	}
}

func radians(deg float64) float64 {
	return deg / (math.Pi / 180)
}
