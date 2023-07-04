package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)


type CragsList struct {
	CragsList []Crag `json:"crags"`
}

type Crag struct {
	Name        string `json:"name"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Disciplines []string `json:"disciplines"`
	Weather     CragWeather `json:"weather"`
}

type CragWeather struct {
	Temperature 			 map[string][]int `json:"temperature"`
	PrecipitationProbability map[string][]int `json:"precipitationProbability"`
	WindSpeed   			 map[string][]int `json:"windSpeed"`
	CloudsCover              map[string][]int `json:"cloudCover"`
}

func InitServer() {
    http.HandleFunc("/crags", cragsHandler)

    log.Fatal(http.ListenAndServe(":9990", nil))
}

func cragsHandler(w http.ResponseWriter, r *http.Request) {
	go logRequest(r)

	crags, isCached, loadErr := loadCrags()

	if loadErr != nil {
		log.Println(loadErr.Error())
	}

	if !isCached {
		go CacheCragsResponse(crags)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(crags)
}

func loadCrags() (list CragsList, isCached bool, err error) {
	cachedCrags, isValidCache := GetCachedCrags()

	if isValidCache {
		return cachedCrags, true, nil
	}

	dataPath := os.Getenv("DATA_PATH")
	file, _ := ioutil.ReadFile(fmt.Sprintf("%s/crags.json", dataPath))
	crags := CragsList{}

	unmarshalError := json.Unmarshal([]byte(file), &crags)

	var wg sync.WaitGroup
    for i := 0; i < len(crags.CragsList); i++ {
		wg.Add(1)
		crag := &crags.CragsList[i]

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			hydrateCragWithWeather(crag)
		}(&wg)
    }

	wg.Wait()

	return crags, false, unmarshalError
}

func hydrateCragWithWeather(crag *Crag) {
	weatherResponse := FetchForecast(crag.Lat, crag.Lon)

	crag.Weather = CragWeather{
		Temperature: weatherResponse.Temperatures(),
		PrecipitationProbability: weatherResponse.PrecipitationProbabilties(),
		WindSpeed: weatherResponse.WindSpeeds(),
		CloudsCover: weatherResponse.CloudCovers(),
	}
}

func logRequest(r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
}
