package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type openWeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  int     `json:"temp_min"`
		TempMax  int     `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type RESTCountriesResponse struct {
	Name           string    `json:"name"`
	TopLevelDomain []string  `json:"topLevelDomain"`
	Alpha2Code     string    `json:"alpha2Code"`
	Alpha3Code     string    `json:"alpha3Code"`
	CallingCodes   []string  `json:"callingCodes"`
	Capital        string    `json:"capital"`
	AltSpellings   []string  `json:"altSpellings"`
	Region         string    `json:"region"`
	Subregion      string    `json:"subregion"`
	Population     int       `json:"population"`
	Latlng         []float64 `json:"latlng"`
	Demonym        string    `json:"demonym"`
	Area           float64   `json:"area"`
	Gini           float64   `json:"gini"`
	Timezones      []string  `json:"timezones"`
	Borders        []string  `json:"borders"`
	NativeName     string    `json:"nativeName"`
	NumericCode    string    `json:"numericCode"`
	Currencies     []struct {
		Code   string `json:"code"`
		Name   string `json:"name"`
		Symbol string `json:"symbol"`
	} `json:"currencies"`
	Languages []struct {
		Iso6391    string `json:"iso639_1"`
		Iso6392    string `json:"iso639_2"`
		Name       string `json:"name"`
		NativeName string `json:"nativeName"`
	} `json:"languages"`
	Translations struct {
		De string `json:"de"`
		Es string `json:"es"`
		Fr string `json:"fr"`
		Ja string `json:"ja"`
		It string `json:"it"`
		Br string `json:"br"`
		Pt string `json:"pt"`
		Nl string `json:"nl"`
		Hr string `json:"hr"`
		Fa string `json:"fa"`
	} `json:"translations"`
	Flag          string `json:"flag"`
	RegionalBlocs []struct {
		Acronym       string        `json:"acronym"`
		Name          string        `json:"name"`
		OtherAcronyms []interface{} `json:"otherAcronyms"`
		OtherNames    []string      `json:"otherNames"`
	} `json:"regionalBlocs"`
	Cioc string `json:"cioc"`
}

type Result struct {
	country1       string
	country2       string
	temperature1   float64
	temperature2   float64
	currency1      string
	currency2      string
	conversionRate float64
}

var apiOpenWeatherKey string = "54fc61ef1a4d3dca104fb3aafbdd2a83"
var apiOpenWeatherURL string = "https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&appid=%s"

var apiRESTCountries string = "https://restcountries.eu/rest/v2/alpha/%s"

var apiCurrencyConverterKey string = "8a5769626b3b9bff0d45"
var apiCurrencyConverterURL string = "https://free.currconv.com/api/v7/convert?q=%s_%s&compact=ultra&apiKey=%s"

var myClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	var city1 string = os.Args[1]
	var city2 string = os.Args[2]
	var err error
	res := new(Result)
	weatherResponse := new(openWeatherResponse) // or &Foo{}
	countriesResponse := new(RESTCountriesResponse)

	//Get Country and Temperature
	getJson(fmt.Sprintf(apiOpenWeatherURL, city1, apiOpenWeatherKey), weatherResponse)
	res.country1 = weatherResponse.Sys.Country
	res.temperature1 = weatherResponse.Main.Temp
	getJson(fmt.Sprintf(apiOpenWeatherURL, city2, apiOpenWeatherKey), weatherResponse)
	res.country2 = weatherResponse.Sys.Country
	res.temperature2 = weatherResponse.Main.Temp
	//Get Currency
	getJson(fmt.Sprintf(apiRESTCountries, res.country1), countriesResponse)
	res.currency1 = countriesResponse.Currencies[0].Code
	getJson(fmt.Sprintf(apiRESTCountries, res.country2), countriesResponse)
	res.currency2 = countriesResponse.Currencies[0].Code
	getJson(fmt.Sprintf(apiRESTCountries, res.country2), countriesResponse)

	res.conversionRate, err = getCurrency(res.currency1, res.currency2)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\t\t\t%s\n----------------------------\n", city1, city2)
	fmt.Printf("%2.2fº\t\t\t%2.2fº\n", res.temperature1, res.temperature2)
	fmt.Printf("%s\t\t\t%s\n", res.currency1, res.currency2)
	fmt.Printf("\t%f\n", res.conversionRate)
}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func getCurrency(currency1, currency2 string) (float64, error) {
	r, err := myClient.Get(fmt.Sprintf(apiCurrencyConverterURL, currency1, currency2, apiCurrencyConverterKey))
	if err != nil {
		return 0.00, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	m := map[string]float64{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		panic(err)
	}

	return m[fmt.Sprintf("%s_%s", currency1, currency2)], nil
}
