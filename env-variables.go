package main

import (
	"github.com/golang/glog"
	"github.com/kelseyhightower/envconfig"
)

type envSettings struct {
	DictionaryPath string `envconfig:"dictionary_path"`
}

var ProjectEnvSettings *envSettings

//environment variables must be in following format
//SETTINGS_URL_CHANNEL_API
//SETTINGS_QUEUE_HOST
func (settings *envSettings) readEnvironmentVariables() error {
	err := envconfig.Process("settings", settings)
	if err != nil {
		return err
	}
	return nil
}

//Always use this function to initialize a Settings struct
func EnvSettingsInit() {
	ProjectEnvSettings = &envSettings{
		DictionaryPath: "test.json",
	}
	ProjectEnvSettings.readEnvironmentVariables()
	glog.Infof("Environment settings are %v\n", ProjectEnvSettings)
}
