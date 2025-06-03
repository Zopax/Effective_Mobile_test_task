package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type EnrichmentClient struct {
	GenderizeURL   string
	AgifyURL       string
	NationalizeURL string
}

func NewEnrichmentClient(genderize, agify, nationalize string) *EnrichmentClient {
	return &EnrichmentClient{
		GenderizeURL:   genderize,
		AgifyURL:       agify,
		NationalizeURL: nationalize,
	}
}

func (c *EnrichmentClient) GetAge(name string) (int, error) {
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", c.AgifyURL, name))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var data struct {
		Age int `json:"age"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}
	return data.Age, nil
}

func (c *EnrichmentClient) GetGender(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", c.GenderizeURL, name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Gender string `json:"gender"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	return data.Gender, nil
}

func (c *EnrichmentClient) GetNationality(name string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s?name=%s", c.NationalizeURL, name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data struct {
		Country []struct {
			CountryID string `json:"country_id"`
		} `json:"country"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if len(data.Country) > 0 {
		return data.Country[0].CountryID, nil
	}
	return "unknown", nil
}
