package bartparse

import (
	"encoding/xml"
	"fmt"
)

type Station struct {
	Lines []Line `xml:"station>etd"`
}

type Line struct {
	Dest   string  `xml:"destination"`
	Abbr   string  `xml:"abbreviation"`
	Trains []Train `xml:"estimate"`
}

func (l Line) String() string {
	return fmt.Sprintf("[%s] %s", l.Abbr, l.Dest)
}

type Train struct {
	Minutes string `xml:"minutes" json:"minutes"`
	Length  string `xml:"length" json:"length"`
}

func (t Train) String() string {
	return fmt.Sprintf("%s car train in %s minutes", t.Minutes, t.Length)
}

func Get_lines(contents []byte) ([]Line, error) {
	var station Station
	if err := xml.Unmarshal(contents, &station); err != nil {
		return nil, err
	}
	return station.Lines, nil
}

func Filter_lines(lines []Line, dests []string) (notfound Line, err error) {
	// No filter for you.
	for _, dest := range dests {
		for _, line := range lines {
			if line.Abbr == dest {
				return line, nil
			}
		}
	}
	return notfound, fmt.Errorf("No good line found")
}
