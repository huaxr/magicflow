// Author: XinRui Hua
// Time:   2022/3/17 下午4:49
// Git:    huaxr

package confutil

import (
	"sync"

	"github.com/huaxr/magicflow/pkg/confutil/manager"
	"github.com/ghodss/yaml"
	"github.com/spf13/cast"
)

type env string

const (
	Local  env = "local"
	Test   env = "tmp"
	Online env = "online"
)

type Conf struct {
	Port          Port          `yaml:"port"`
	Dcc           Dcc           `yaml:"dcc"`
	Queue         Queue         `yaml:"queue"`
	Db            Db            `yaml:"db"`
	Configuration Configuration `yaml:"configuration"`
	Switch        Switch        `yaml:"switch"`
}

type Port struct {
	Api     string `yaml:"api"`
	Service string `yaml:"service"`
}
type Dcc struct {
	Basepath string `yaml:"basepath"`
	Hosts    string `yaml:"hosts"`
}

type Queue struct {
	Nsq struct {
		Brokers string `yaml:"brokers"`
		Lookups string `yaml:"lookups"`
		Admin   string `yaml:"admin"`
		Secret  string `yaml:"secret"`
	} `yaml:"nsq"`

	Kafka struct {
		Brokers string `yaml:"brokers"`
	} `yaml:"kafka"`
}

type Db struct {
	Mysql struct {
		Slaves       []string `yaml:"slaves"`
		Master       string   `yaml:"master"`
		MaxConn      int      `yaml:"maxConn"`
		MaxIdle      int      `yaml:"maxIdle"`
		LogLevel     int      `yaml:"logLevel"` // 0 1 2 ..
		ShowSql      bool     `yaml:"showSql"`
		SlowDuration int      `yaml:"slowDuration"`
	} `yaml:"mysql"`

	Redis struct {
		Host        string `yaml:"host"`
		Password    string `yaml:"password"`
		Idletimeout int    `yaml:"idletimeout"`
		Readtimeout int    `yaml:"readtimeout"`
		MaxRetry    int    `yaml:"maxretry"`
		Poolsize    int    `yaml:"poolsize"`
		Db          int    `yaml:"db"`
	} `yaml:"redis"`
}

type Configuration struct {
	DispatchThreadCount     string `yaml:"dispatchThreadCount"`
	ChannelReportInterval   string `yaml:"channelReportInterval"`
	BrokerHeartbeatInterval string `yaml:"brokerHeartbeatInterval"`

	WatchKeyPrefix  string `yaml:"watchKeyPrefix"`
	ServicesPrefix  string `yaml:"servicesPrefix"`
	ElectionPrefix  string `yaml:"electionPrefix"`
	MaxNodeIDPrefix string `yaml:"maxNodeIDPrefix"`
	NodeIDPrefix    string `yaml:"nodeIDPrefix"`

	Appid      string `yaml:"appid"`
	Appkey     string `yaml:"appkey"`
	Env        env    `yaml:"env"`
	Superadmin string `yaml:"superadmin"`

	None string `yaml:"none"`
}

type Switch struct {
	EnableMonitor      bool `yaml:"enableMonitor"`
	EnableHealthyCheck bool `yaml:"enableHealthyCheck"`
	EnableMasterElect  bool `yaml:"enableMasterElect"`
}

var (
	conf *Conf
	once sync.Once
)

func GetConf() *Conf {
	if conf == nil {
		initConf()
	}
	return conf
}

func (c *Conf) IsLocalEnv() bool {
	return c.Configuration.Env == Local
}

func initConf() {
	once.Do(func() {
		manager.Init(*confDir, []string{"yml"}, nil)
		conf = new(Conf)
		yamlText := manager.GetConfigByKey("conf.yml")
		if err := yaml.Unmarshal(yamlText, conf); err != nil {
			panic(err)
		}

		log = new(Log)
		yamlText = manager.GetConfigByKey("log.yml")
		if err := yaml.Unmarshal(yamlText, log); err != nil {
			panic(err)
		}

		prom = new(Prom)
		yamlText = manager.GetConfigByKey("prom.yml")
		if err := yaml.Unmarshal(yamlText, prom); err != nil {
			panic(err)
		}
	})
}

func (c *Conf) GetNodeIDPrefix() string {
	return c.Configuration.NodeIDPrefix
}

func (c *Conf) GetMaxNodeIDPrefix() string {
	return c.Configuration.MaxNodeIDPrefix
}

func (c *Conf) GetDispatchThreadCount() int {
	return cast.ToInt(c.Configuration.DispatchThreadCount)
}
