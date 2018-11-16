package master

import (
	"encoding/json"
	"fmt"
	"github.com/SugarAlex/crontab/common"
	"net/http"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer
)

func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
	)
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	postJob = r.PostForm.Get("job")
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

ERR:
	fmt.Println(err)
}

func InitApiServer() (err error) {
	var (
		mux *http.ServeMux
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	return nil
}
