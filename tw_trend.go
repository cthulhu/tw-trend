package tw_trend

type TwTrendApp struct {
	consumerKey, consumerSecret, accessToken, accessSecret string
}

func New(consumerKey, consumerSecret, accessToken, accessSecret string) *TwTrendApp {
	return &TwTrendApp{consumerKey, consumerSecret, accessToken, accessSecret}
}

func (app *TwTrendApp) Run() error {
	return nil
}
