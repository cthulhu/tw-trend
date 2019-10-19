package domain

type WordsReport struct {
	*Aggregated `json:"words"`
}

type HashtagsReport struct {
	*Aggregated `json:"hashtags"`
}

type Aggregated struct {
	TokensWithCounts []TokenWithCount `json:"aggregated"`
}

type TokenWithCount struct {
	Token string `json:"token"`
	Count int    `json:"rank"`
}
