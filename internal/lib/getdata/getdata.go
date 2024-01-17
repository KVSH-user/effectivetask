package getdata

import (
	"encoding/json"
	"io/ioutil"
	"log/slog"
	"net/http"
)

type AgeApi struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type SexApi struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type CountryApi struct {
	Count   int    `json:"count"`
	Name    string `json:"name"`
	Country []struct {
		CountryID   string  `json:"country_id"`
		Probability float64 `json:"probability"`
	} `json:"country"`
}

func GetAge(name string) int {
	const op = "internal.lib.getdata.GetAge"

	log := slog.With(slog.String("op", op))

	var age AgeApi

	reqUrl := "https://api.agify.io/?name=" + name

	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Error("no response from request", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &age)

	return age.Age
}

func GetSex(name string) string {
	const op = "internal.lib.getdata.GetSex"

	log := slog.With(slog.String("op", op))

	var gender SexApi

	reqUrl := "https://api.genderize.io/?name=" + name

	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Error("no response from request", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &gender)

	return gender.Gender
}

func GetCountry(name string) string {
	const op = "internal.lib.getdata.GetCountry"

	log := slog.With(slog.String("op", op))

	var country CountryApi

	reqUrl := "https://api.nationalize.io/?name=" + name

	resp, err := http.Get(reqUrl)
	if err != nil {
		log.Error("no response from request", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &country)

	return country.Country[0].CountryID
}
