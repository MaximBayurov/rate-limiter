package configuration

import "time"

type ClientConf struct {
	Host    string        `mapstructure:"host"`
	Port    string        `mapstructure:"port"`
	Timeout time.Duration `mapstructure:"timeout"`
}
