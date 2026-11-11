package impl

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/vitor-franz/book-log/entities"
	"github.com/vitor-franz/book-log/infrastructure/settings"
)

func NewTOMLSettingsLoader(filePath string) (settings.Provider, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Join(errors.New("failed to read settings file"), err)
	}

	var applicationSettings entities.ApplicationSettings
	err = toml.Unmarshal(b, &applicationSettings)
	if err != nil {
		return nil, errors.Join(errors.New("failed to unmarshal settings"), err)
	}

	return tomlSettingsLoader{applicationSettings: applicationSettings}, nil
}

type tomlSettingsLoader struct {
	applicationSettings entities.ApplicationSettings
}

func (s tomlSettingsLoader) GetDatabaseSettings() entities.DatabaseSettings {
	return s.applicationSettings.DatabaseSettings
}

func (s tomlSettingsLoader) GetJWTSettings() entities.JWTSettings {
	return s.applicationSettings.JWTSettings
}

func (s tomlSettingsLoader) GetLoggerSettings() entities.LoggerSettings {
	return s.applicationSettings.LoggerSettings
}

func (s tomlSettingsLoader) GetRuntimeSettings() entities.RuntimeSettings {
	return s.applicationSettings.RuntimeSettings
}
