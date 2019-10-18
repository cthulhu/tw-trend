package domain

type Aggregated struct {
	TokensWithCounts []TokenWithCount `json:"aggregated"`
}

type TokenWithCount struct {
	Token string `json:"token"`
	Count int    `json:"rank"`
}
