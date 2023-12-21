package artists

type Response struct {
	Response struct {
		GeoObjectCollection struct {
			FeatureMember []struct {
				GeoObject struct {
					Point struct {
						Pos string `json:"pos"`
					} `json:"Point"`
				} `json:"GeoObject"`
			} `json:"featureMember"`
		} `json:"GeoObjectCollection"`
	} `json:"response"`
}

type Relation struct {
	Location string
	Date     string
}

type Artist struct {
	ID             int                 `json:"id"`
	Name           string              `json:"name"`
	Image          string              `json:"image"`
	Members        []string            `json:"members"`
	RelationsApi   string              `json:"relations"`
	DatesLocations map[string][]string `json:"datesLocations"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
	Loc            string              `json:"locations"`
	Geolocation    map[string]Coordinate
	Country        []string
}

type Location struct {
	ID        int      `json:"id"`
	Locations []string `json:"locations"`
}

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type RelationOtYelnar struct {
	LocationDates map[string][]string `json:"datesLocations"`
}

type Data_group struct {
	ID             int                 `json:"id"`
	Image          string              `json:"image"`
	Name           string              `json:"name"`
	Members        []string            `json:"members"`
	DatesLocations map[string][]string `json:"datesLocations"`
	CreationDate   int                 `json:"creationDate"`
	FirstAlbum     string              `json:"firstAlbum"`
}

type Coincidence struct {
	Artist     []Artist
	Data_group []Data_group
}
