package config

import (
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	MainConfig  MainConfig `yaml:"mainConfig"`
	CangMiaoKey string     `yaml:"cangmiaoKey"`
}

var BotConfig *Config

type MainConfig struct {
	NickName    []string `yaml:"nickName"`
	SuperUser   []int64  `yaml:"superUser"`
	NapcatUrl   string   `yaml:"napcatUrl"`
	NapcatToken string   `yaml:"napcatToken"`
}

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
