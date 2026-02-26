package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type PaginatedResponse[T any] struct {
	Info struct {
		Next string `json:"next"`
	} `json:"info"`
	Results []T `json:"results"`
}

func GetAllPages[T any](endpoint string, params map[string]string) ([]T, error) {
	var all []T
	baseURL := fmt.Sprintf("https://rickandmortyapi.com/api/%s", endpoint)
	url := baseURL

	if len(params) > 0 {
		q := "?"
		for k, v := range params {
			q += fmt.Sprintf("%s=%s&", k, v)
		}
		url += q[:len(q)-1]
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for url != "" {
		resp, err := client.Get(url)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode == 429 {
			resp.Body.Close()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if resp.StatusCode >= 400 {
			resp.Body.Close()
			return nil, fmt.Errorf("API returned status: %s (url: %s)", resp.Status, url)
		}

		var page PaginatedResponse[T]
		err = json.NewDecoder(resp.Body).Decode(&page)
		resp.Body.Close()
		if err != nil {
			return nil, err
		}

		all = append(all, page.Results...)
		url = page.Info.Next
	}

	return all, nil
}
