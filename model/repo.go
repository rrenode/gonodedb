package model

type Repo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
