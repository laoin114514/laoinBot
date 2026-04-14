package config

import (
	"os"

	"github.com/laoin114514/jmapi"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	MainConfig  MainConfig `yaml:"mainConfig"`
	CangMiaoKey string     `yaml:"cangmiaoKey"`
	Db          DbConfig   `yaml:"db"`
}

var BotConfig *Config

type MainConfig struct {
	NickName    []string `yaml:"nickName"`
	SuperUser   []int64  `yaml:"superUser"`
	NapcatUrl   string   `yaml:"napcatUrl"`
	NapcatToken string   `yaml:"napcatToken"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

var JMClient = jmapi.NewClient(jmapi.Config{
	ClientType:        jmapi.ClientTypeAPI,
	AutoUpdateHost:    true,
	AutoEnsureCookies: true,
})

func LoadConfig(path ...string) error {
	if len(path) == 0 {
		path = []string{"config/config.yml"}
	}
	f, err := os.Open(path[0])
	if err != nil {
		return err
	}
	defer f.Close()
	err = yaml.NewDecoder(f).Decode(&BotConfig)
	if err != nil {
		return err
	}
	return nil
}
