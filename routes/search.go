package routes

import (
	"encoding/json"
	"net/http"
	"rickandmorty-backend/models"
	"rickandmorty-backend/services"
	"strconv"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
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

	characters, err := services.GetAllPages[models.NamedResource]("character", map[string]string{"name": term})
	if err == nil {
		for _, c := range characters {
			results = append(results, models.SearchResult{
				Name: c.Name,
				Type: "character",
				Url:  c.Url,
			})
		}
	}

	locations, err := services.GetAllPages[models.NamedResource]("location", map[string]string{"name": term})
	if err == nil {
		for _, l := range locations {
			results = append(results, models.SearchResult{
				Name: l.Name,
				Type: "location",
				Url:  l.Url,
			})
		}
	}

	episodes, err := services.GetAllPages[models.NamedResource]("episode", map[string]string{"name": term})
	if err == nil {
		for _, e := range episodes {
			results = append(results, models.SearchResult{
				Name: e.Name,
				Type: "episode",
				Url:  e.Url,
			})
		}
	}

	if limit > 0 && len(results) > limit {
		results = results[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
