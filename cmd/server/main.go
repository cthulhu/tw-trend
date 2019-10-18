package main

import (
	// "fmt"
	"fmt"
	"os"

	// "os/signal"
	// "syscall"

	log "github.com/sirupsen/logrus"

	// "github.com/dghubble/go-twitter/twitter"
	// "github.com/dghubble/oauth1"
	tw_trend "github.com/cthulhu/tw-trend"
)

type MainedData struct {
	Text     string
	Hashtags []string
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
}

func main() {

	log.Info("Application start")
	defer log.Info("Application stop")

	consumerKey, err := ensureEnvValue("CONSUMERKEY")
	if err != nil {
		log.Panic(err)
	}

	consumerSecret, err := ensureEnvValue("CONSUMERSECRET")
	if err != nil {
		log.Panic(err)
	}

	accessToken, err := ensureEnvValue("ACCESSTOKEN")
	if err != nil {
		log.Panic(err)
	}

	accessSecret, err := ensureEnvValue("ACCESSSECRET")
	if err != nil {
		log.Panic(err)
	}

	app := tw_trend.New(consumerKey, consumerSecret, accessToken, accessSecret)

	// // Wait for SIGINT and SIGTERM (HIT CTRL-C)
	// ch := make(chan os.Signal)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	if err = app.Run(); err != nil {
		log.Panic(err)
	}
}

func ensureEnvValue(envVarName string) (string, error) {
	value := os.Getenv(envVarName)
	if value == "" {
		return "", fmt.Errorf("Env Variable %s is required", envVarName)
	}
	return value, nil
}
