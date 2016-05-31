package config

import (
	"errors"
	"strings"
	"github.com/spf13/viper"
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

func GetString(key string) (string, error) {
	value := viper.GetString(key)
	if value == "" {
		return value, getMissingVariableError(key)
	}

	return value, nil
}

func getMissingVariableError(key string) error {
	return errors.New("Missing environment variable " + getEnvironmentVariableName(key))
}

func getEnvironmentVariableName(key string) string {
	return envPrefix + "_" + envReplacer.Replace(strings.ToUpper(key))
}