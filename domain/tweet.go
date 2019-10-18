package domain

type Tweet struct {
	Text     string   `json:"tweet"`
	Hashtags []string `json:"hashtags"`
	Tokens   []string `json:"tokens"`
}
