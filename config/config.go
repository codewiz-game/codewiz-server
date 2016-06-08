package config

import (
	"github.com/spf13/viper"
	"strings"
)

const (
	envPrefix = "CODEWIZ"
)

var envReplacer *strings.Replacer

func init() {
	envReplacer = strings.NewReplacer(".", "_")

	// Configure Viper to search for all configurations
	// in the system environment variables
	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(envReplacer)
}

func GetString(key string, defaultVal ...string) string {
	if !viper.InConfig(key) && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return viper.GetString(key)
}

func GetBool(key string, defaultVal ...bool) bool {
	if !viper.InConfig(key) && len(defaultVal) > 0 {
		return defaultVal[0]
	}
	return viper.GetBool(key)
}

func GetEnvironmentVariableName(key string) string {
	return envPrefix + "_" + envReplacer.Replace(strings.ToUpper(key))
}
