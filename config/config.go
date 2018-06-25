package config

import (
	"github.com/spf13/viper"
	"github.com/v-zhidu/orb/logging"
)

//LoadConfig - Load configuration for application.
func LoadConfig(name string, paths []string, configType string) error {
	viper.SetConfigName(name)
	viper.SetConfigType(configType)
	for _, path := range paths {
		viper.AddConfigPath(path)
	}
	if err := viper.ReadInConfig(); err != nil {
		logging.Error("failed to load config file.", logging.Fields{
			"fileName":  name,
			"filePaths": paths,
			"fileType":  configType,
		}, err)
		return err
	}

	return nil
}

//GetString - load string configuration with key.
func GetString(key string) string {
	return viper.GetString(key)
}

//GetSlice - load string slice configuration with key.
func GetSlice(key string) []string {
	return viper.GetStringSlice(key)
}

//GetInt - load integer configuration with key.
func GetInt(key string) int {
	return viper.GetInt(key)
}

//Unmarshal convert specific value to a struct, map
func Unmarshal(key string, rawVal interface{}) error {
	return viper.UnmarshalKey(key, rawVal)
}

//SetDefaultConfig - Set some default configuration in web application.
func SetDefaultConfig(values map[string]interface{}) {
	for k, v := range values {
		viper.SetDefault(k, v)
	}
}
