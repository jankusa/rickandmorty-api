package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rickandmorty-backend/models"
	"strconv"
	"testing"
)

func getAllPagesMock(endpoint string, params map[string]string) []models.NamedResource {
	term := params["name"]
	switch endpoint {
	case "character":
		if term == "rick" {
			return []models.NamedResource{
				{Name: "Rick Sanchez", Url: "url1"},
				{Name: "Cool Rick", Url: "url5"},
			}
		}
		if term == "morty" {
			return []models.NamedResource{{Name: "Morty Smith", Url: "url2"}}
		}
	case "location":
		return []models.NamedResource{{Name: "Earth", Url: "url3"}}
	case "episode":
		return []models.NamedResource{{Name: "Pilot", Url: "url4"}}
	}
	return []models.NamedResource{}
}

func SearchHandlerTest(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	limitParam := r.URL.Query().Get("limit")

	var limit int
	var err error
	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	results := []models.SearchResult{}

	for _, c := range getAllPagesMock("character", map[string]string{"name": term}) {
		results = append(results, models.SearchResult{
			Name: c.Name,
			Type: "character",
			Url:  c.Url,
		})
	}

	for _, l := range getAllPagesMock("location", map[string]string{"name": term}) {
		results = append(results, models.SearchResult{
			Name: l.Name,
			Type: "location",
			Url:  l.Url,
		})
	}

	for _, e := range getAllPagesMock("episode", map[string]string{"name": term}) {
		results = append(results, models.SearchResult{
			Name: e.Name,
			Type: "episode",
			Url:  e.Url,
		})
	}

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func TestSearchHandler_Basic(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?term=rick", nil)
	rec := httptest.NewRecorder()

	SearchHandlerTest(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var results []models.SearchResult
	if err := json.NewDecoder(rec.Body).Decode(&results); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(results) != 4 { // 2 characters + 1 location + 1 episode
		t.Fatalf("expected 4 results, got %d", len(results))
	}

	if results[0].Name != "Rick Sanchez" {
		t.Errorf("expected first result to be 'Rick Sanchez', got %s", results[0].Name)
	}
}

func TestSearchHandler_Limit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?term=rick&limit=2", nil)
	rec := httptest.NewRecorder()

	SearchHandlerTest(rec, req)

	var results []models.SearchResult
	json.NewDecoder(rec.Body).Decode(&results)

	if len(results) != 2 {
		t.Fatalf("expected 2 results due to limit, got %d", len(results))
	}
}

func TestSearchHandler_InvalidLimit(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?term=rick&limit=-1", nil)
	rec := httptest.NewRecorder()

	SearchHandlerTest(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 BadRequest for invalid limit, got %d", rec.Code)
	}
}

func TestSearchHandler_EmptyTerm(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()

	SearchHandlerTest(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
}

func TestSearchHandler_LimitGreaterThanResults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/search?term=rick&limit=10", nil)
	rec := httptest.NewRecorder()

	SearchHandlerTest(rec, req)

	var results []models.SearchResult
	json.NewDecoder(rec.Body).Decode(&results)

	if len(results) != 4 { // 4 wyniki w mocku
		t.Fatalf("expected 4 results, got %d", len(results))
	}
}
