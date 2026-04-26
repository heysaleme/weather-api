package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"weather-api/internal/errs"
)

type CountryClient struct {
	httpClient       *http.Client
	countryBaseURL   string
	countryCitiesURL string
}

type CountryInfo struct {
	Name string
	Code string
}

func NewCountryClient(httpClient *http.Client) *CountryClient {
	return &CountryClient{
		httpClient:       httpClient,
		countryBaseURL:   "https://restcountries.com/v3.1/name",
		countryCitiesURL: "https://countriesnow.space/api/v0.1/countries/cities",
	}
}

type restCountryResponse []struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	CCA2 string `json:"cca2"`
}

func (c *CountryClient) GetCountry(ctx context.Context, country string) (*CountryInfo, error) {
	country = strings.TrimSpace(country)
	if country == "" {
		return nil, errs.InvalidInput("country is required")
	}

	u, err := url.Parse(fmt.Sprintf("%s/%s", c.countryBaseURL, url.PathEscape(country)))
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	q := u.Query()
	q.Set("fullText", "true")
	q.Set("fields", "name,cca2")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.Upstream("resolve country: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errs.NotFound("country not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errs.Upstream("resolve country status %d", resp.StatusCode)
	}

	var result restCountryResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errs.Upstream("decode country response: %v", err)
	}

	if len(result) == 0 || result[0].Name.Common == "" || result[0].CCA2 == "" {
		return nil, errs.NotFound("country not found")
	}

	return &CountryInfo{
		Name: result[0].Name.Common,
		Code: result[0].CCA2,
	}, nil
}

type countryCitiesRequest struct {
	Country string `json:"country"`
}

type countryCitiesResponse struct {
	Error bool     `json:"error"`
	Msg   string   `json:"msg"`
	Data  []string `json:"data"`
}

func (c *CountryClient) GetCities(ctx context.Context, country string) ([]string, error) {
	country = strings.TrimSpace(country)
	if country == "" {
		return nil, errs.InvalidInput("country is required")
	}

	body, err := json.Marshal(countryCitiesRequest{Country: country})
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.countryCitiesURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.Upstream("load country cities: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errs.Upstream("load country cities status %d", resp.StatusCode)
	}

	var result countryCitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, errs.Upstream("decode country cities response: %v", err)
	}

	if result.Error {
		return nil, errs.NotFound(strings.TrimSpace(result.Msg))
	}
	if len(result.Data) == 0 {
		return nil, errs.NotFound("no cities found for country")
	}

	return result.Data, nil
}
