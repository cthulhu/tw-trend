package tw_trend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Shop2market/goping"
	"github.com/cthulhu/tw-trend/resource"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type Server struct {
	port int
}

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info(r.Method, r.URL.Path)
	l.handler.ServeHTTP(w, r)
}

func NewServer(port int) *Server {
	return &Server{port}
}

func (s *Server) Run() error {
	log.Infof("Starting http on port %d", s.port)

	router := httprouter.New()

	router.GET("/words", handleJson(resource.Words))
	router.GET("/hashtags", handleJson(resource.Hashtags))
	router.GET("/data/:time_id", resource.Data)

	router.GET("/ping", goping.Ping())
	router.GET("/", goping.Ping())

	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), &Logger{router})
}
func (s *Server) Close() error {
	return nil
}

type jsonHandler func(*http.Request, httprouter.Params) (interface{}, error)

func handleJson(handler jsonHandler) httprouter.Handle {
	return handleJsonWithStatusCode(handler)
}
func handleJsonWithStatusCode(handler jsonHandler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		response := json.NewEncoder(w)
		resp, err := handler(req, params)
		if err != nil {
			w.WriteHeader(500)
			log.Error(err)
			return
		}
		if resp == nil {
			w.WriteHeader(404)
			return
		}
		if err = response.Encode(resp); err != nil {
			w.WriteHeader(500)
			log.Error(err)
		}
	}
}
