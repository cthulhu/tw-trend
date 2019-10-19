package tw_trend

import (
	"os"

	"github.com/cthulhu/tw-trend/service/tw_streamer"
	"github.com/cthulhu/tw-trend/store"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

type TwTrendApp struct {
	httpServer *Server
	fileStore  *store.File
	stream     *tw_streamer.TwStream
}

func New(consumerKey, consumerSecret, accessToken, accessSecret string, httpPort int) (*TwTrendApp, error) {
	var err error
	app := &TwTrendApp{nil, nil, nil}
	app.stream, err = tw_streamer.New(consumerKey, consumerSecret, accessToken, accessSecret)
	if err != nil {
		return app, err
	}
	app.fileStore, err = store.New()
	if err != nil {
		return app, err
	}
	app.httpServer = NewServer(httpPort)
	return app, nil
}

func (app *TwTrendApp) Run() error {
	go func() {
		app.fileStore.JSONlStream(app.stream.TweetsAsJSONl())
		if err := app.stream.Run(); err != nil {
			log.Error(err)
		}
	}()
	return app.httpServer.Run()
}

func (app *TwTrendApp) Stop() {
	app.httpServer.Close()
	app.stream.Close()
	app.fileStore.Close()
}
