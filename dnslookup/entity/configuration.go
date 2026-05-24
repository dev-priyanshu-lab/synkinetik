package entity

type Configuration struct {
	Address string                 `json:"address" binding:"required"`
	Sources map[string]IntelSource `json:"sources" binding:"required"`
}
