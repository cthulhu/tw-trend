package resource

import (
	"net/http"

	"github.com/cthulhu/tw-trend/domain"
	"github.com/cthulhu/tw-trend/service/aggregator"
	"github.com/cthulhu/tw-trend/store"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

var DEFAULT_DAYS_BACK int
var MAX_RESULTS int

func init() {
	DEFAULT_DAYS_BACK = 3
	MAX_RESULTS = 50
}

func Words(r *http.Request, ps httprouter.Params) (interface{}, error) {
	var err error
	readCloser, err := store.TweetsReadCloser()
	if err != nil {
		return nil, err
	}
	aggregated, err := aggregator.Aggregate(readCloser, "words", MAX_RESULTS)
	defer readCloser.Close()
	return domain.WordsReport{aggregated}, err
}

func Hashtags(r *http.Request, ps httprouter.Params) (interface{}, error) {
	var err error
	readCloser, err := store.TweetsReadCloser()
	if err != nil {
		return nil, err
	}
	aggregated, err := aggregator.Aggregate(readCloser, "hashtags", MAX_RESULTS)
	defer readCloser.Close()
	return domain.HashtagsReport{aggregated}, err
}

func Data(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	data, err := store.ReadFileByTimeStamp(ps.ByName("time_id"))
	if err != nil {
		w.WriteHeader(500)
		log.Error(err)
		return
	}
	w.Write(data)
}
