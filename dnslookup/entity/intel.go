package entity

type IntelSource struct {
	Exclusions string   `json:"exclusions" binding:"false"`
	Feeds      []string `json:"feeds" binding:"required"`
}
