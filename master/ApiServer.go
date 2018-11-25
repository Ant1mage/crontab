package master

import (
	"encoding/json"
	"github.com/SugarAlex/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	G_apiServer *ApiServer
)

func InitApiServer() (err error) {
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir
		staticHandler http.Handler
	)
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)
	// 静态文件目录
	staticDir = http.Dir(G_config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	mux.Handle("/", http.StripPrefix("/", staticHandler)) //   ./webroot/index.html

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}
	httpServer = &http.Server{
		ReadHeaderTimeout: time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout:      time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
		Handler:           mux,
	}
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}
	// 启动
	go httpServer.Serve(listener)
	return nil
}

// POST Job -> json
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	// 解析表单
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	postJob = r.PostForm.Get("job")
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}
	return
ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

// 删除任务接口
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)

	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	name = r.PostForm.Get("name")
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

// 列举所有crontab任务
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	var (
		jobList []*common.Job
		bytes   []byte
		err     error
	)

	// 获取任务列表
	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}

// 强制杀死某个任务
// POST /job/kill  name=job1
func handleJobKill(resp http.ResponseWriter, req *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)

	// 解析POST表单
	if err = req.ParseForm(); err != nil {
		goto ERR
	}

	// 要杀
	name = req.PostForm.Get("name")
	if err = G_jobMgr.KillJob(name); err != nil {
		goto ERR
	}
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		resp.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		resp.Write(bytes)
	}
}
