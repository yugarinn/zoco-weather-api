package internal

import (
    "fmt"
    "errors"
    "encoding/json"
    "io/ioutil"
    "net/http"
	"strconv"
)


type RequestQueryParameter struct {
    name  string
    value string
}

type FetchForecastResult struct {
    Forecast struct {
        Time                    []string `json:"time"`
        Temperature             []float64 `json:"temperature_2m"`
        CloudCover              []int `json:"cloudcover"`
        PrecipitationProbabilty []int `json:"precipitation_probability"`
        WindSpeed               []float64 `json:"windspeed_10m"`
    } `json:"hourly"`
}

func (result FetchForecastResult) PrecipitationProbabilties() map[string][]int {
    return CastHourlyValuesToWeek(result.Forecast.PrecipitationProbabilty)
}

func (result FetchForecastResult) Temperatures() map[string][]int {
    temperatures := make([]int, len(result.Forecast.Temperature))

    for i := range result.Forecast.Temperature {
        temperatures[i] = int(result.Forecast.Temperature[i])
    }

    return CastHourlyValuesToWeek(temperatures)
}

func (result FetchForecastResult) WindSpeeds() map[string][]int {
    speeds := make([]int, len(result.Forecast.Temperature))

    for i := range result.Forecast.Temperature {
        speeds[i] = int(result.Forecast.WindSpeed[i])
    }

    return CastHourlyValuesToWeek(speeds)
}

func (result FetchForecastResult) CloudCovers() map[string][]int {
    return CastHourlyValuesToWeek(result.Forecast.CloudCover)
}

func FetchForecast(lat float64, lon float64) FetchForecastResult {
    fmt.Println("fetching!")
    params := []RequestQueryParameter{
        {name: "latitude", value: strconv.FormatFloat(lat, 'f', -1, 64)},
        {name: "longitude", value: strconv.FormatFloat(lon, 'f', -1, 64)},
        {name: "hourly", value: "temperature_2m"},
        {name: "hourly", value: "precipitation_probability"},
        {name: "hourly", value: "cloudcover"},
        {name: "hourly", value: "windspeed_10m"},
    }

    forecastResponse, _ := fetch("GET", "https://api.open-meteo.com/v1/forecast", params)

    var forecastResult FetchForecastResult
    json.Unmarshal(forecastResponse, &forecastResult)

    return forecastResult
}

func fetch(method string, url string, params []RequestQueryParameter) ([]byte, error) {
    client := &http.Client{}
    req, err := http.NewRequest(method, url, nil)

    if len(params) > 0 && method == "GET" {
        query := req.URL.Query()

        for _, param := range params {
            query.Add(param.name, param.value)
        }

        req.URL.RawQuery = query.Encode()
    }

    resp, _ := client.Do(req)

    if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
        err = errors.New(http.StatusText(resp.StatusCode))
    }

    if err != nil {
        return nil, err
    } else {
        payload, _ := ioutil.ReadAll(resp.Body)

        return payload, nil
    }
}
