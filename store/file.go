package store

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/cthulhu/tw-trend/domain"
	log "github.com/sirupsen/logrus"
)

type File struct {
}

func New() (*File, error) {
	return &File{}, nil
}

func (s *File) JSONlStream(tweets chan domain.Tweet) {
	go func() {
		for tweet := range tweets {
			b, err := json.Marshal(tweet)
			if err != nil {
				log.Error(err)
				return
			}
			if err = s.Println(string(b)); err != nil {
				log.Error(err)
				return
			}
		}
		log.Infof("Stopped...")
	}()
}

func (s *File) Println(jsonl string) error {
	f, err := os.OpenFile(fileName(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(jsonl + "\n")
	return err
}

func (s *File) Close() {

}

func fileName() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d.jsonl", now.Year(), now.Month(), now.Day())
}
