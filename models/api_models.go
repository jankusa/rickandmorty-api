package models

// for /top-pairs
type Episode struct {
	ID         int      `json:"id"`
	Characters []string `json:"characters"`
}

type NamedResource struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
