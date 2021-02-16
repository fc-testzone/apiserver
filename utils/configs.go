package utils

import (
	"encoding/json"
	"io/ioutil"
)

type DBCfg struct {
	IP     string
	Port   int
	User   string
	Passwd string
	DB     string
}

type ServerCfg struct {
	IP   string
	Port int
}

type AppCfg struct {
	Server ServerCfg
	DB     DBCfg
	Domain string
}

type Configs struct {
	configs AppCfg
}

func NewConfigs() *Configs {
	return &Configs{}
}

func (c *Configs) LoadFromFile(fileName string) error {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c.configs)
	if err != nil {
		return err
	}

	return nil
}

func (c *Configs) Settings() *AppCfg {
	return &c.configs
}
