package store

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/cthulhu/tw-trend/domain"
	log "github.com/sirupsen/logrus"
)

var VolumeDir string
var FetchDaysForReports func() []string

func init() {
	VolumeDir = "fixtures"
	FetchDaysForReports = func() []string {
		now := time.Now()
		return []string{
			fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day()),
			fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day()-1),
			fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day()-2),
		}
	}
}

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
func ReadFileByTimeStamp(time_id string) ([]byte, error) {
	id, err := strconv.Atoi(time_id)
	if err != nil {
		return []byte{}, err
	}
	file, err := os.Open(fmt.Sprintf("%s/%d.jsonl", VolumeDir, id))
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func FilesForReport() []string {
	fileNames := []string{}
	for _, day := range FetchDaysForReports() {
		fileNames = append(fileNames, fmt.Sprintf("%s/%s.jsonl", VolumeDir, day))
	}
	return fileNames
}

func TweetsReadCloser() (io.ReadCloser, error) {
	mf := &MultiFile{}
	for _, fileName := range FilesForReport() {
		f, err := os.Open(fileName)
		if err != nil {
			return nil, err
		}
		mf.files = append(mf.files, f)
	}
	return mf, nil
}
