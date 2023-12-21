package artists

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

var Country []string

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		Error(w, http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		Error(w, http.StatusBadRequest)
		return
	}

	body, err := GetUrl("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}

	var artists []Artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	var firstAlbums []string
	var creationDates []string

	for _, artist := range artists {
		firstAlbums = append(firstAlbums, artist.FirstAlbum[6:])
		creationDates = append(creationDates, strconv.Itoa(artist.CreationDate))
	}
	minAlb, maxAlb := minimal(firstAlbums)
	minCreat, maxCreat := minimal(creationDates)
	Country = GetCountry(len(artists))
	type TemplateData struct {
		Artists  []Artist
		Mincreat int
		Maxcreat int
		Minalb   int
		Maxalb   int
		Country  []string
	}
	data := TemplateData{
		Artists:  artists,
		Mincreat: minCreat,
		Maxcreat: maxCreat,
		Minalb:   minAlb,
		Maxalb:   maxAlb,
		Country:  Country,
	}
	startTime := time.Now()

	// Ваш код для запуска сайта

	elapsedTime := time.Since(startTime)
	fmt.Printf("Сайт запущен за %s\n", elapsedTime)
	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, data)
}

func minimal(arr []string) (int, int) {
	res := []int{}
	for _, i := range arr {
		num, _ := strconv.Atoi(i)
		res = append(res, num)
	}
	min := res[0]
	max := res[0]
	for i := 1; i <= len(res)-1; i++ {
		if min > res[i] {
			min = res[i]
		}
		if max < res[i] {
			max = res[i]
		}
	}
	return min, max
}

func GetCountry(numOfArtist int) []string {
	var wg sync.WaitGroup
	var mu sync.Mutex
	CountryMap := make(map[string]bool)
	for i := 1; i <= numOfArtist; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			Country := GetCountryCity("https://groupietrackers.herokuapp.com/api/locations/" + strconv.Itoa(i))
			mu.Lock()
			for _, c := range Country {
				v := strings.ReplaceAll(c, "_", " ")
				if len(v) <= 3 {
					v = strings.ToUpper(v)
				} else {
					v = strings.Title(v)
				}
				CountryMap[v] = true
			}
			mu.Unlock()
		}(i)
	}
	wg.Wait()

	Country := make([]string, 0, len(CountryMap))
	for k := range CountryMap {
		Country = append(Country, k)
	}

	return Country
}

// CountryMap := make(map[string]bool)
// resultChan := make(chan []string)

// for i := 1; i <= c; i++ {
// 	go func(index int) {
// 		k := GetCountryCity("https://groupietrackers.herokuapp.com/api/locations/" + strconv.Itoa(index))
// 		countrySlice := make([]string, 0, len(k))
// 		for _, city := range k {
// 			country := strings.Title(strings.Split(city, ",")[1])
// 			if !CountryMap[country] {
// 				CountryMap[country] = true
// 				countrySlice = append(countrySlice, country)
// 			}
// 		}
// 		resultChan <- countrySlice
// 	}(i)
// }

// for i := 1; i <= c; i++ {
// 	countrySlice := <-resultChan
// 	for _, country := range countrySlice {
// 		Country = append(Country, country)
// 	}
// }

// close(resultChan)
