package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"rickandmorty-backend/models"
	"testing"
)

func setupTopPairsTestData() {
	cachedEpisodes = []models.Episode{
		{
			ID: 1,
			Characters: []string{
				"url1",
				"url2",
				"url3",
			},
		},
		{
			ID: 2,
			Characters: []string{
				"url1",
				"url2",
			},
		},
	}

	cachedCharacters = []models.NamedResource{
		{Name: "Rick", Url: "url1"},
		{Name: "Morty", Url: "url2"},
		{Name: "Summer", Url: "url3"},
	}

	cacheLoaded = true
}

func TestTopPairsHandler_Basic(t *testing.T) {
	setupTopPairsTestData()

	req := httptest.NewRequest(http.MethodGet, "/top-pairs", nil)
	rec := httptest.NewRecorder()

	TopPairsHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var results []models.TopPairResult
	err := json.NewDecoder(rec.Body).Decode(&results)
	if err != nil {
		t.Fatal(err)
	}

	if len(results) == 0 {
		t.Fatal("expected some results")
	}

	if results[0].Episodes != 2 {
		t.Fatalf("expected top pair to appear 2 times, got %d", results[0].Episodes)
	}
}

func TestTopPairsHandler_Limit(t *testing.T) {
	setupTopPairsTestData()

	req := httptest.NewRequest(http.MethodGet, "/top-pairs?limit=1", nil)
	rec := httptest.NewRecorder()

	TopPairsHandler(rec, req)

	var results []models.TopPairResult
	json.NewDecoder(rec.Body).Decode(&results)

	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestTopPairsHandler_InvalidMin(t *testing.T) {
	setupTopPairsTestData()

	req := httptest.NewRequest(http.MethodGet, "/top-pairs?min=-1", nil)
	rec := httptest.NewRecorder()

	TopPairsHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
