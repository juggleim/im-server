package configures

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	EnvDev  = "dev"
	EnvProd = "prod"

	MsgStoreEngine_MySQL = "mysql"
	MsgStoreEngine_Mongo = "mongo"

	CmdMsgExpired int64 = 7 * 24 * 60 * 60 * 1000
	MsgExpired    int64 = 24 * 60 * 60 * 1000
)

type ImConfig struct {
	NodeName       string `yaml:"nodeName"`
	NodeHost       string `yaml:"nodeHost"`
	MsgStoreEngine string `yaml:"msgStoreEngine"`

	Log struct {
		LogPath      string `yaml:"logPath"`
		LogName      string `yaml:"logName"`
		Visual       bool   `yaml:"visual"`
		VLogHttpPort int    `yaml:"vloghttpPort"`
	} `ymal:"log"`

	Mysql struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
		DbName   string `yaml:"name"`
		Debug    bool   `yaml:"debug"`
	} `yaml:"mysql"`

	MongoDb struct {
		Address string `yaml:"address"`
		DbName  string `yaml:"name"`
	} `yaml:"mongodb"`

	ConnectManager struct {
		WsPort int `yaml:"wsPort"`
	} `yaml:"connectManager"`

	ApiGateway struct {
		HttpPort int `yaml:"httpPort"`
	} `yaml:"apiGateway"`

	NavGateway struct {
		HttpPort int `yaml:"httpPort"`
	} `yaml:"navGateway"`

	AdminGateway struct {
		HttpPort int `yaml:"httpPort"`
	} `yaml:"adminGateway"`
}

var Config ImConfig

func InitConfigures() error {
	cfBytes, err := os.ReadFile("conf/config.yml")
	if err == nil {
		var conf ImConfig
		yaml.Unmarshal(cfBytes, &conf)
		Config = conf
		return nil
	} else {
		return err
	}
}
