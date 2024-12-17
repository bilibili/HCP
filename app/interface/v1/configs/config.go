package configs

import (
	"errors"

	"github.com/BurntSushi/toml"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
	"github.com/go-kratos/kratos/pkg/net/rpc/warden"
	"github.com/go-kratos/kratos/pkg/time"

	"os"
	"strings"

	bm "github.com/go-kratos/kratos/pkg/net/http/blademaster"
)

var (
	confPath string = "configs/application.toml"
	Conf            = &Config{}
	LogOff          = false
	TraceOff        = false
)

type Config struct {
	HttpServer  *bm.ServerConfig    `toml:"httpServer"`
	GrpcServer  warden.ServerConfig `toml:"grpcServer"`
	Mysql       Mysql               `toml:"mysql"`
	CloudSecret CloudSecret         `toml:"cloud_secret"`
	ServerSpec  ServerSpec          `toml:"server_spec"`
	CloudSync   CloudSync           `toml:"cloud_sync"`
	JobConfig   *JobConfig          `toml:"job"`
	LoginConf   *LoginConfig        `toml:"login"`
	CloudConf   *Cloud              `toml:"cloud"`
}

type Cloud struct {
	TagProjectKey   string `toml:"tag_project_key"`
	TerraformPath   string `toml:"terraform_path"`
	TerraformBin    string `toml:"terraform_bin"`
	TerraformPlugin string `toml:"terraform_plugin"`
}

type BaiduKeyPairConf struct {
	AccountCid  string `toml:"account_cid"`
	KeyPairName string `toml:"key_pair_name"`
	KeyPairId   string `toml:"key_pair_id"`
}

type LoginConfig struct {
	ExpireSecond int    `toml:"expire_second"`
	SecretKey    string `toml:"secret_key"`
	LetterBytes  string `toml:"letter_bytes"`
}

type CloudSync struct {
	SyncInterval       int    `toml:"sync_interval"`
	ConcurrencyAccount int    `toml:"concurrency_account"`
	Sync               bool   `toml:"sync"`
	Clear              bool   `toml:"clear"`
	GCPProxy           string `toml:"gcp_proxy"`
	ScanAt             bool   `toml:"scan_at"`
	SyncCron           bool   `toml:"sync_cron"`
}

type ServerSpec struct {
	Price  Price  `toml:"price"`
	Series Series `toml:"series"`
	Gpu    Gpu    `toml:"gpu"`
}

type Price struct {
	Mask float64 `toml:"mask"`
}

type Gpu struct {
	GpuModels []string `toml:"gpu_models"`
}

type Series struct {
	AliSeries []string `toml:"ali_series"`
	QSeries   []string `toml:"q_series"`
	HwSeries  []string `toml:"hw_series"`
	KsSeries  []string `toml:"ks_series"`
	USeries   []string `toml:"u_series"`
}

type Mysql struct {
	Driver       string        `json:"driver"`
	Addr         string        `json:"addr"`
	Dsn          string        `json:"dsn"`
	ReadDSN      []string      `json:"readDSN"`
	Active       int           `json:"active"`
	Idle         int           `json:"idle"`
	IdleTimeout  time.Duration `json:"idleTimeout"`
	QueryTimeout time.Duration `json:"queryTimeout"`
	ExecTimeout  time.Duration `json:"execTimeout"`
	TranTimeout  time.Duration `json:"tranTimeout"`
}
type CloudSecret struct {
	SecretAesKey string `toml:"secret_aes_key"`
}

type JobConfig struct {
	IssueWait int `toml:"issue_wait"`
}

func init() {
	envs := os.Environ()
	for _, env := range envs {
		if strings.Contains(env, "confPath=") {
			confPath = strings.Replace(env, "confPath=", "", 1)
		}
		if strings.Contains(env, "LogOff") {
			LogOff = true
		}
		if strings.Contains(env, "TraceOff") {
			TraceOff = true
		}
	}
	if err := Init(); err != nil {
		panic(err)
	}
}

func Init() error {
	// 配置初始化
	log.Info("init config>>>")
	if confPath != "" {
		return local()
	}
	return remote()
}

func local() error {
	// 本地配置
	_, err := toml.DecodeFile(confPath, &Conf)
	return err
}

func remote() (err error) {
	if load() != nil {
		return err
	}
	return
}

func load() (err error) {
	err = paladin.Init()
	if err != nil {
		panic(err)
	}
	var tmpConf = &Config{}
	s, err := paladin.Get("application.toml").Raw()
	if err != nil { // 传入文件名，获取该文件内容
		return errors.New("load configs center error")
	}
	if _, err = toml.Decode(s, tmpConf); err != nil {
		log.Info("err: %+v", err)
		return errors.New("could not decode configs")
	}
	Conf = tmpConf
	return
}
