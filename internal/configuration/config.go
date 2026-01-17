package configuration

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf `mapstructure:"logger"`
	Database DbConf     `mapstructure:"database"`
	App      AppConf    `mapstructure:"app"`
	Server   ServerConf `mapstructure:"server"`
}

func New(filePath string) (Config, error) {
	pathInfo := getConfigPathInfoFor(filePath)

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigName(pathInfo.Name)
	viper.AddConfigPath(pathInfo.Path)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			return Config{}, fmt.Errorf("config file not found: %w", err)
		}
		return Config{}, fmt.Errorf("error reading config: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return Config{}, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return config, nil
}
