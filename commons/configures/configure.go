package configures

import (
	"flag"
	"im-server/commons/tools"
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
	DefaultPort    int    `yaml:"defaultPort"`
	NodeName       string `yaml:"nodeName"`
	NodeHost       string `yaml:"nodeHost"`
	MsgStoreEngine string `yaml:"msgStoreEngine"`

	Log struct {
		LogPath string `yaml:"logPath"`
		LogName string `yaml:"logName"`
	} `ymal:"log"`

	Kvdb struct {
		IsOpen   bool   `yaml:"isOpen"`
		DataPath string `yaml:"dataPath"`
	} `ymal:"kvdb"`

	MsgLogs struct {
		LogPath    string `yaml:"logPath"`
		MaxBackups int    `yaml:"maxBackups"`
		IsCompress bool   `yaml:"isCompress"`
	}

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
	configFile := flag.String("config", "conf/config.yml", "Path to the configuration file")
	flag.Parse()
	cfBytes, err := os.ReadFile(*configFile)
	if err == nil {
		var conf ImConfig
		yaml.Unmarshal(cfBytes, &conf)
		Config = conf
		//check
		if Config.NodeName == "" {
			Config.NodeName = tools.GenerateUUIDShort11()
		}
		if Config.NodeHost == "" {
			Config.NodeHost = "127.0.0.1"
		}
		if Config.MsgStoreEngine == "" {
			Config.MsgStoreEngine = MsgStoreEngine_MySQL
		}
		if Config.ConnectManager.WsPort <= 0 {
			Config.ConnectManager.WsPort = Config.DefaultPort
		}
		if Config.ConnectManager.WsPort <= 0 {
			Config.ConnectManager.WsPort = 9003
		}
		if Config.AdminGateway.HttpPort <= 0 {
			Config.AdminGateway.HttpPort = 8090
		}
		return nil
	} else {
		return err
	}
}
