package config

import (
	"errors"

	"github.com/benk-techworld/www-backend/internal/validator"
	"github.com/spf13/viper"
)

var config *viper.Viper

func Load(configType, configFilePath string) error {

	if !validator.PermittedValues(configType, "json", "yaml", "toml", "ini") {
		return errors.New("config file type is not supported")
	}

	config = viper.New()
	config.AutomaticEnv()
	config.SetConfigType(configType)
	config.SetConfigFile(configFilePath)

	return config.ReadInConfig()

}

func Get() *viper.Viper {
	return config
}
