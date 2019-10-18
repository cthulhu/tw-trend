package resource

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var DEFAULT_DAYS_BACK int

func init() {
	DEFAULT_DAYS_BACK = 3
}

func Words(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func Hashtags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func Data(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}
