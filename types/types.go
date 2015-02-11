package types

import "fmt"

type Loc struct {
	Lat float64
	Lon float64
}
type Station struct {
	Loc  Loc
	name string
	dest []string
	abbr string
	dir  string
}

func (l Station) String() string {
	return fmt.Sprintf("%s", l.name)
}
