package main

import "net/http"
//import "io"
import "io/ioutil"
import "encoding/json"
import "fmt"
import "math"
import "strconv"
import "github.com/danielsoneg/bartbot/bartparse"

type Loc struct {
	lat float64
	lon float64
}
type Station struct {
	loc Loc
	name string
	dest []string
	abbr string
	dir  string
}

const BART_ETD string = "http://api.bart.gov/api/etd.aspx?cmd=etd&orig=%s&dir=%s&key=MW9S-E7SL-26DU-VV8V"
const BART_ADV string = "http://api.bart.gov/api/bsa.aspx?cmd=bsa&orig=%s&key=MW9S-E7SL-26DU-VV8V"

func (l Station) String() string {
	return fmt.Sprintf("%s", l.name)
}

func radians(deg float64) float64 {
	return deg / (math.Pi / 180)
}

func get_estimates(s Station) ([]byte, error) {
  url := fmt.Sprintf(BART_ETD, s.abbr, s.dir)
  resp, err := http.Get(url)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  contents, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  return contents, nil
}

func dist(from, to Loc) float64 {
	from.lat = radians(from.lat)
	to.lat = radians(to.lat)
	d_lat := (to.lat - from.lat) / 2
	d_lon := radians(to.lon-from.lon) / 2
	haversin := math.Pow(math.Sin(d_lat), 2) + math.Cos(to.lat)
	haversin = haversin * math.Cos(from.lat) * math.Pow(math.Sin(d_lon), 2)
	return math.Asin(math.Sqrt(haversin))
}

func nearest(from Loc, a, b Station) Station {
	if dist(from, a.loc) < dist(from, b.loc) {
		return a
	} else {
		return b
	}
}

type Response struct {
  Station string            `json:"station"`
  Abbr    string            `json:"abbr"`
  Dest    string            `json:"dest"`
  Trains  []bartparse.Train `json:"trains"`
}

func fmt_estimates(line bartparse.Line, station Station) []byte {
  resp := Response{station.name, station.abbr, line.Dest, line.Trains}
  js, _ := json.MarshalIndent(resp, "", "  ")
  return js
}

func find_trains(here Loc) (js []byte, err error) {
	cvc := Station{
		Loc{37.779471, -122.413809}, "Civic Center",
    []string{"RICH", "PITT"},
		"civc", "n",
	}
	brk := Station{
		Loc{37.869842, -122.267986}, "Downtown Berkeley",
    []string{"MLBR", "DALY", "FRMT"},
		"dbrk", "s",
	}
  station := nearest(here, cvc, brk)
  xml, err := get_estimates(station)
  if (err != nil) {
    return js, fmt.Errorf("Error talking to Bart: %s", err)
  }
  lines, _ := bartparse.Get_lines(xml)
  if (err != nil) {
    return js, fmt.Errorf("Could not read Bart's response", err)
  }
  line, err := bartparse.Filter_lines(lines, station.dest)
  if (err != nil) {
    return js, fmt.Errorf("Couldn't find an available train")
  }
  js = fmt_estimates(line, station)
  return js, err
}

func read_latlon(r *http.Request) (Loc, error) {
  lat_str, lon_str := r.PostFormValue("lat"), r.PostFormValue("lon")
  if (lat_str == "" || lon_str == "") {
    return Loc{0,0}, fmt.Errorf("Must include lat & lon")
  }
  lat, lat_err := strconv.ParseFloat(lat_str, 32)
  lon, lon_err := strconv.ParseFloat(lon_str, 32)
  if (lat_err != nil || lon_err != nil) {
    errmsg := fmt.Sprintf("Bad values passed for lat & lon: %s, %s",
                  lat_str, lon_str)
    return Loc{0,0}, fmt.Errorf(errmsg)
  }
  return Loc{lat, lon}, nil
}

func get_loc(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json")
  loc, err := read_latlon(r)
  if err != nil {
    http.Error(w, fmt.Sprint(err), http.StatusBadRequest)
    return
  }
  js, err := find_trains(loc)
  if err != nil {
    http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
  }
  w.Write(js)
}

func main() {
  http.HandleFunc("/loc", get_loc)
  http.ListenAndServe(":3000", nil)
}
