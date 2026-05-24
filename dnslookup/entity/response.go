package entity

type Response struct {
	Category string   `json:"category"`
	Tag      []string `json:"tag"`
}
