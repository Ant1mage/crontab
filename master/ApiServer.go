package master

import (
	"net/http"
)

type ApiServer struct {
	httpServer *http.Server
}

func handleJobSave(w http.ResponseWriter, r *http.Request) {

}

func InitApiServer() (err error) {
	var (
		mux *http.ServeMux
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	return nil
}
