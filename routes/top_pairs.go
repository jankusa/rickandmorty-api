package routes

import (
	"encoding/json"
	"math"
	"net/http"
	"rickandmorty-backend/models"
	"rickandmorty-backend/services"
	"sort"
	"strconv"
	"strings"
)

var cachedEpisodes []models.Episode
var cachedCharacters []models.NamedResource
var cacheLoaded bool

func LoadCache() error {
	if cacheLoaded {
		return nil
	}

	episodes, err := services.GetAllPages[models.Episode]("episode", nil)
	if err != nil {
		return err
	}

	characters, err := services.GetAllPages[models.NamedResource]("character", nil)
	if err != nil {
		return err
	}

	cachedEpisodes = episodes
	cachedCharacters = characters
	cacheLoaded = true
	return nil
}

func TopPairsHandler(w http.ResponseWriter, r *http.Request) {
	minParam := r.URL.Query().Get("min")
	maxParam := r.URL.Query().Get("max")
	limitParam := r.URL.Query().Get("limit")

	min := 0
	max := math.MaxInt
	limit := 20
	var err error

	if minParam != "" {
		min, err = strconv.Atoi(minParam)
		if err != nil || min < 0 {
			http.Error(w, "Invalid min parameter", http.StatusBadRequest)
			return
		}
	}

	if maxParam != "" {
		max, err = strconv.Atoi(maxParam)
		if err != nil || max < 0 {
			http.Error(w, "Invalid max parameter", http.StatusBadRequest)
			return
		}
	}

	if limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	err = LoadCache()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	episodes := cachedEpisodes

	pairCounts := make(map[string]int)

	for _, ep := range episodes {
		for i := 0; i < len(ep.Characters); i++ {
			for j := i + 1; j < len(ep.Characters); j++ {
				a := ep.Characters[i]
				b := ep.Characters[j]

				if a > b {
					a, b = b, a
				}

				key := a + "|" + b
				pairCounts[key]++
			}
		}
	}

	charactersMap := make(map[string]models.NamedResource)
	for _, c := range cachedCharacters {
		charactersMap[c.Url] = c
	}

	results := []models.TopPairResult{}

	for key, count := range pairCounts {
		if count < min || count > max {
			continue
		}

		parts := strings.Split(key, "|")

		char1, ok1 := charactersMap[parts[0]]
		char2, ok2 := charactersMap[parts[1]]

		if !ok1 || !ok2 {
			continue
		}

		results = append(results, models.TopPairResult{
			Character1: char1,
			Character2: char2,
			Episodes:   count,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Episodes > results[j].Episodes
	})

	if len(results) > limit {
		results = results[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
