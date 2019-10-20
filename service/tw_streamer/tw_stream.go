package tw_streamer

import (
	"github.com/cthulhu/tw-trend/domain"
	"github.com/cthulhu/tw-trend/service/tokenizer"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	log "github.com/sirupsen/logrus"
)

type TwStream struct {
	chTweets chan domain.Tweet
	stream   *twitter.Stream
}

func New(consumerKey, consumerSecret, accessToken, accessSecret string) (*TwStream, error) {

	var err error

	twStream := &TwStream{chTweets: make(chan domain.Tweet)}

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)
	// hack to test authentication
	// twitter stream api doesn't check correctness of tockens
	_, _, err = client.Timelines.HomeTimeline(&twitter.HomeTimelineParams{
		Count: 1,
	})
	if err != nil {
		return nil, err
	}
	// hack to test authentication
	filterParams := &twitter.StreamFilterParams{
		Locations: []string{"4.729242", "52.278174", "5.079162", "52.431064"},
	}

	twStream.stream, err = client.Streams.Filter(filterParams)
	if err != nil {
		return nil, err
	}
	return twStream, nil
}

func (tws *TwStream) Run() error {
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		tw := domain.Tweet{}
		tw.Text = tweet.Text
		tw.Tokens, _ = tokenizer.Tokenize(tweet.Text)
		if tweet.Entities != nil {
			for _, h := range tweet.Entities.Hashtags {
				tw.Hashtags = append(tw.Hashtags, h.Text)
			}
		}

		tws.chTweets <- tw
	}
	log.Info("Starting Stream")
	demux.HandleChan(tws.stream.Messages)
	return nil
}

func (tws *TwStream) Close() {
	log.Info("Stoping stream")
	tws.stream.Stop()
	close(tws.chTweets)
}

func (tws *TwStream) TweetsAsJSONl() chan domain.Tweet {
	return tws.chTweets
}
