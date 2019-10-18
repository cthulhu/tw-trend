package tw_trend

import (
	"net/http"

	"github.com/cthulhu/tw-trend/resource"
	"github.com/cthulhu/tw-trend/service/tw_streamer"
	"github.com/cthulhu/tw-trend/store"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type TwTrendApp struct {
	consumerKey, consumerSecret, accessToken, accessSecret string

	httpServer *http.Server
	fileStore  *store.File
	stream     *tw_streamer.TwStream
}

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Method, r.URL.Path)
	l.handler.ServeHTTP(w, r)
}

func New(consumerKey, consumerSecret, accessToken, accessSecret string) *TwTrendApp {
	return &TwTrendApp{consumerKey, consumerSecret, accessToken, accessSecret, nil, nil, nil}
}

func (app *TwTrendApp) Run() error {
	var err error
	log.Info("Starting http")
	router := httprouter.New()

	app.stream, err = tw_streamer.New(app.consumerKey, app.consumerSecret, app.accessToken, app.accessSecret)
	if err != nil {
		return err
	}

	app.fileStore, err = store.New()
	if err != nil {
		return err
	}

	router.GET("/words", resource.Words)
	router.GET("/hashtags", resource.Hashtags)
	router.GET("/data", resource.Data)

	go func() {
		app.fileStore.JSONlStream(app.stream.TweetsAsJSONl())
		if err := app.stream.Run(); err != nil {
			log.Error(err)
		}
	}()

	return http.ListenAndServe(":8000", &Logger{router})
}

func (app *TwTrendApp) Stop() {
	log.Info("Stoping http")
	app.stream.Close()
	app.fileStore.Close()
}
