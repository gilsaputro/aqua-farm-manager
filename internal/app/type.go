package app

// UrlID denotes a id of every path
type UrlID int

// list defined UrlID
const (
	Farms   UrlID = 1
	FarmsID UrlID = 2
	Ponds   UrlID = 3
	PondsID UrlID = 4
	Limit   UrlID = 5 // this will be flag to stop
	Stat    UrlID = 6 // include stat api for getting metrics
)

// this list define all known of path setting
var (
	UrlIDName = map[UrlID]string{
		Farms:   "/farms",
		FarmsID: "/farms/{id}",
		Ponds:   "/ponds",
		PondsID: "/ponds/{id}",
		Stat:    "/stat",
	}

	UrlIDValue = map[string]UrlID{
		UrlIDName[Farms]:   Farms,
		UrlIDName[Ponds]:   Ponds,
		UrlIDName[Stat]:    Stat,
		UrlIDName[FarmsID]: FarmsID,
		UrlIDName[PondsID]: PondsID,
	}

	UrlIDMethod = map[UrlID][]string{
		Farms:   {"POST", "GET", "PUT", "DELETE"},
		FarmsID: {"GET"},
		Ponds:   {"POST", "GET", "PUT", "DELETE"},
		PondsID: {"GET"},
		Stat:    {"GET"},
	}
)

// Int return int representation of urlID
func (urlID UrlID) Int() int { return int(urlID) }

// string return string representation of urlID
func (urlID UrlID) String() string { return UrlIDName[urlID] }

// GetListMethod return list method representation of urlID
func (urlID UrlID) GetListMethod() []string { return UrlIDMethod[urlID] }
