package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const (
	// 配置路径是相对路径，相对于进程而言
	// 编译之后的可执行文件需要跟 conf 文件夹在同一级目录
	configRoot   = "./conf"
	configDb     = configRoot + "/db.json"
	configDbTest = configRoot + "/dbtest.json"

	logPath = "../log"
)

type Config struct {
	DB     *DBConfig
	Common *CommonConfig
}

type DBConfig struct {
	Host     string `json:"Host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbName"`
	Timeout  string `json:"timeout"`
}

type CommonConfig struct {
	LogPath string
}

func NewConfig() (*Config, error) {
	c := Config{}
	_, err := os.Stat(configDbTest)
	if err == nil {
		if err := c.getConfig(configDbTest); err != nil {
			return nil, err
		}
	} else {
		if err := c.getConfig(configDb); err != nil {
			return nil, err
		}
	}
	c.Common = &CommonConfig{logPath}
	return &c, nil
}

func (c *Config) getConfig(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	switch name {
	case configDb:
		err = decoder.Decode(&c.DB) // it needs to be a pointer
	case configDbTest:
		err = decoder.Decode(&c.DB) // it needs to be a pointer
	default:
		err = errors.New(fmt.Sprintf("no such option %s", name))
	}
	if err != nil {
		return err
	}
	return nil
}
