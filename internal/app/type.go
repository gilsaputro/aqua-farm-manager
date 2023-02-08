package app

// UrlID denotes a id of every path
type UrlID int

// list defined UrlID
const (
	Farms UrlID = 1
	Ponds UrlID = 2
	Limit UrlID = 3 // this will be flag to stop
	Stat  UrlID = 4 // include stat api for getting metrics
)

// this list define all known of path setting
var (
	UrlIDName = map[UrlID]string{
		Farms: "/farms",
		Ponds: "/ponds",
		Stat:  "/stat",
	}

	UrlIDValue = map[string]UrlID{
		UrlIDName[Farms]: Farms,
		UrlIDName[Ponds]: Ponds,
		UrlIDName[Stat]:  Stat,
	}

	UrlIDMethod = map[UrlID][]string{
		Farms: {"POST", "GET", "PUT", "DELETE"},
		Ponds: {"POST", "GET", "PUT", "DELETE"},
		Stat:  {"GET"},
	}
)

// Int return int representation of urlID
func (urlID UrlID) Int() int { return int(urlID) }

// string return string representation of urlID
func (urlID UrlID) String() string { return UrlIDName[urlID] }

// GetListMethod return list method representation of urlID
func (urlID UrlID) GetListMethod() []string { return UrlIDMethod[urlID] }
