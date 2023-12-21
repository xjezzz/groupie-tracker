package artists

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"text/template"
)

var CountryCity []string

func Artists(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		Error(w, http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	var artist Artist

	url, err := GetUrl("https://groupietrackers.herokuapp.com/api/artists/" + id)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	json.Unmarshal(url, &artist)
	apiKey := "c9ad80bd-38f5-4f83-b299-5cb205880e94"
	CountryCity = GetCountryCity("https://groupietrackers.herokuapp.com/api/locations/" + id)
	url, err = GetUrl(artist.RelationsApi)
	if err != nil {
		Error(w, http.StatusBadRequest)
		return
	}
	json.Unmarshal(url, &artist)
	var coordinates []Coordinate

	for _, i := range CountryCity {
		geocodeURL := fmt.Sprintf("https://geocode-maps.yandex.ru/1.x/?apikey=%s&format=json&geocode=%s", apiKey, i)
		latitude, longitude, err := GetCoordinates(geocodeURL)
		if err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}
		coords := Coordinate{
			Latitude:  latitude,
			Longitude: longitude,
		}
		coordinates = append(coordinates, coords)
	}
	Geolocation := map[string]Coordinate{}
	for i := 0; i <= len(coordinates)-1; i++ {
		CountryCity[i] = strings.Title(CountryCity[i])
		Geolocation[CountryCity[i]] = coordinates[i]
	}
	artist.Geolocation = Geolocation
	funcMap := template.FuncMap{
		"title":    strings.Title,
		"replace1": ReplaceAllString,
		"replace2": ReplaceAllString2,
	}
	tmpl := template.New("templates/artist.html").Funcs(funcMap)
	tmpl, err = tmpl.ParseFiles("templates/artist.html")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "artist.html", artist)
}

func ReplaceAllString(s, old, new string) string {
	return regexp.MustCompile("-").ReplaceAllString(s, " ")
}

func ReplaceAllString2(s, old, new string) string {
	return regexp.MustCompile("_").ReplaceAllString(s, " ")
}

func GetCountryCity(url string) []string {
	cc, _ := GetUrl(url)
	hash := ""
	flag := false
	flag2 := false
	CountryCity := []string{}
	for _, i := range cc {
		if i == '[' {
			flag = true
		}
		if flag2 && hash != "" && i == '"' {
			CountryCity = append(CountryCity, hash)
			hash = ""
		}
		if flag {
			if i == '"' {
				flag2 = true
			}
		}
		if flag2 {
			hash += string(i)
		}
		if i == ']' && flag {
			break
		}
	}
	for i, s := range CountryCity {
		for _, char := range "\"" {
			s = strings.ReplaceAll(s, string(char), "")
		}
		CountryCity[i] = s
	}
	for i, s := range CountryCity {
		for _, char := range "," {
			s = strings.ReplaceAll(s, string(char), "")
		}
		CountryCity[i] = s
	}
	for i, s := range CountryCity {
		s = strings.ReplaceAll(s, " ", "") // Удаляем пробелы
		CountryCity[i] = s
	}
	CountryCity2 := []string{}
	for _, i := range CountryCity {
		if i != "" {
			CountryCity2 = append(CountryCity2, i)
		}
	}
	for i, s := range CountryCity2 {
		s = strings.ReplaceAll(s, "-", ",")
		CountryCity2[i] = s
	}
	return CountryCity2
}

func GetCoordinates(url string) (float64, float64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return 0, 0, err
	}

	pos := data.Response.GeoObjectCollection.FeatureMember[0].GeoObject.Point.Pos
	coords := parseCoordinates(pos)
	return coords[0], coords[1], nil
}

func parseCoordinates(pos string) []float64 {
	coords := make([]float64, 2)
	fmt.Sscanf(pos, "%f %f", &coords[1], &coords[0])
	return coords
}
