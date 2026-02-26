package models

type TopPairResult struct {
	Character1 NamedResource `json:"character1"`
	Character2 NamedResource `json:"character2"`
	Episodes   int           `json:"episodes"`
}
