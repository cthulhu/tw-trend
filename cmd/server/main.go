package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type MainedData struct {
	Text     string
	Hashtags []string
}

func main() {

  consumerKey := ensureEnvValue("CONSUMERKEY")
  consumerSecret := ensureEnvValue("CONSUMERSECRET")
  accessToken := ensureEnvValue("ACCESSTOKEN")
  accessSecret := ensureEnvValue("ACCESSSECRET")

  config := oauth1.NewConfig(consumerKey, consumerSecret)
  token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	filterParams := &twitter.StreamFilterParams{
		Locations: []string{"4.729242", "52.278174", "5.079162", "52.431064"},
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		mainedData := MainedData{}
		mainedData.Text = tweet.Text
		if tweet.Entities != nil {
			for _, h := range tweet.Entities.Hashtags {
				mainedData.Hashtags = append(mainedData.Hashtags, h.Text)
			}
		}
		fmt.Println(mainedData)
	}
	demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	demux.Event = func(event *twitter.Event) {
		fmt.Printf("%#v\n", event)
	}
	fmt.Println("Starting Stream...")

	go demux.HandleChan(stream.Messages)
	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()
}

func ensureEnvValue(envVarName string) string {
	value := os.Getenv(envVarName)
	if value == "" {
		log.Panicf("Env Variable %s is required\n", envVarName)
	}
	return value
}
