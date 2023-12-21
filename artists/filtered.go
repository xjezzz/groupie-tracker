package artists

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

func Filtered(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		Error(w, http.StatusBadRequest)
		return
	}
	var artists []Artist
	body, err := GetUrl("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &artists)

	err = r.ParseForm()
	if err != nil {
		Error(w, http.StatusBadRequest)
		return
	}

	creatStart, err := strconv.Atoi(r.Form.Get("fromSlider"))
	creatEnd, err := strconv.Atoi(r.Form.Get("toSlider"))
	albumStart, err := strconv.Atoi(r.Form.Get("fromSlider2"))
	albumEnd, err := strconv.Atoi(r.Form.Get("toSlider2"))
	checkboxes := []string{"1", "2", "3", "4", "5", "6", "7", "8"}
	values := make(map[string]string)
	count := 0
	for _, checkbox := range checkboxes {
		values[checkbox] = r.Form.Get(checkbox)
		if r.Form.Get(checkbox) != "" {
			count++
		}
	}
	checkBoxWorks := false
	if count > 0 {
		checkBoxWorks = true
	}
	reqCountry := r.Form.Get("country")
	flagCountry := false
	if reqCountry == "" {
		flagCountry = true
	}
	var artists2 []Artist
	var wg sync.WaitGroup
	if !flagCountry {
		for i := 1; i <= len(artists); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				Countryy := GetCountryCity("https://groupietrackers.herokuapp.com/api/locations/" + strconv.Itoa(i))
				for _, j := range Countryy {
					v := strings.ReplaceAll(j, "_", " ")
					if len(v) <= 3 {
						v = strings.ToUpper(v)
					} else {
						v = strings.Title(v)
					}
					if reqCountry == v {
						artists2 = append(artists2, artists[i-1])
						break
					}
				}
			}(i)
		}
		wg.Wait()
	} else {
		artists2 = artists
	}
	artists3 := []Artist{}
	for j := range artists2 {
		for i := creatStart; i <= creatEnd; i++ {
			if i == artists2[j].CreationDate {
				artists3 = append(artists3, artists2[j])
				break
			}
		}
	}
	artists4 := []Artist{}
	for j := range artists3 {
		for i := albumStart; i <= albumEnd; i++ {
			alb, _ := strconv.Atoi(artists3[j].FirstAlbum[6:])
			if alb == i {
				if checkBoxWorks {
					for v, k := range values {
						num, _ := strconv.Atoi(v)
						if num == len(artists3[j].Members) && k != "" {
							artists4 = append(artists4, artists3[j])
						}
					}
				} else {
					artists4 = append(artists4, artists3[j])
				}
				break
			}
		}
	}
	cearch_artists, err := Add_stuckt(w, body)
	find := r.FormValue("search")
	if find != "" {
		selections, check_struct := Check_coincidence(w, find, cearch_artists)
		if check_struct == 0 || len(find) == 0 {
			Error(w, http.StatusBadRequest)
			return
		}
		endStruct := Coincidence{
			cearch_artists,
			selections,
		}
		fmt.Println(endStruct)
		tmpl, err := template.ParseFiles("templates/search.html")
		if err != nil {
			Error(w, http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, endStruct)
		return
	}

	tmpl, err := template.ParseFiles("templates/filter.html")
	if err != nil {
		Error(w, http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "filter.html", artists4)
}

func Check_coincidence(w http.ResponseWriter, find string, all_data_group []Artist) ([]Data_group, int) {
	res := []Data_group{}
	flag := false
	check_struct := 0
	for _, v := range all_data_group {
		if strings.Contains(strings.ToLower(v.Name), strings.ToLower(find)) ||
			strings.Contains(strconv.Itoa(v.CreationDate), find) ||
			strings.Contains(v.FirstAlbum, find) {
			check_struct++
			// fmt.Println(v.LOCATION_AND_DATES.LocationDates)
			res = append(res, Data_group{
				ID:             v.ID,
				Image:          v.Image,
				Name:           v.Name,
				Members:        v.Members,
				DatesLocations: v.DatesLocations,
				CreationDate:   v.CreationDate,
				FirstAlbum:     v.FirstAlbum,
			})
			continue
		}

		for _, j := range v.Members {
			if strings.Contains(strings.ToLower(j), strings.ToLower(find)) {
				check_struct++
				flag = true
				res = append(res, Data_group{
					ID:             v.ID,
					Image:          v.Image,
					Name:           v.Name,
					Members:        v.Members,
					DatesLocations: v.DatesLocations,
					CreationDate:   v.CreationDate,
					FirstAlbum:     v.FirstAlbum,
				})
				break
			}
		}
		if flag {
			flag = false
			continue
		}
		for key := range v.DatesLocations {
			if strings.Contains(strings.ToLower(key), strings.ToLower(find)) {
				check_struct++
				flag = true
				// fmt.Println(v)
				res = append(res, Data_group{
					ID:             v.ID,
					Image:          v.Image,
					Name:           v.Name,
					Members:        v.Members,
					DatesLocations: v.DatesLocations,
					CreationDate:   v.CreationDate,
					FirstAlbum:     v.FirstAlbum,
				})
				break
			}
		}
		if flag {
			flag = false
			continue
		}

	}
	return res, check_struct
}

func Add_stuckt(w http.ResponseWriter, jsonData1 []byte) ([]Artist, error) {
	res := []Artist{}
	// groups := []Artists2{}

	var wg sync.WaitGroup
	var mu sync.Mutex

	var res_stuckt []Artist
	err := json.Unmarshal(jsonData1, &res_stuckt)
	// fmt.Println(jsonData1)
	if err != nil {
		return nil, err
	}

	for _, v := range res_stuckt {
		wg.Add(1)
		go func(v Artist) {
			defer wg.Done()

			jsonData1, err := GetUrl(v.RelationsApi)
			if err != nil {
				// Обработка ошибок
				return
			}

			var delete RelationOtYelnar
			err = json.Unmarshal(jsonData1, &delete)
			if err != nil {
				// Обработка ошибок
				return
			}
			lock := delete.LocationDates

			mu.Lock()
			res = append(res, Artist{
				ID:             v.ID,
				Image:          v.Image,
				Name:           v.Name,
				Members:        v.Members,
				DatesLocations: lock,
				CreationDate:   v.CreationDate,
				FirstAlbum:     v.FirstAlbum,
				RelationsApi:   v.RelationsApi,
				Loc:            v.Loc,
				Geolocation:    v.Geolocation,
				Country:        v.Country,
			})

			mu.Unlock()
		}(v)
	}
	wg.Wait()

	return res, nil
}
