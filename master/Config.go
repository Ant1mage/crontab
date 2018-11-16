package master

type Config struct {
	EtcdEndPoints   []string `json:"etcdEndPoints"`
	EtcdDialTimeOut int      `json:"etcdDialTimeOut"`
}

var (
	G_config *Config
)
