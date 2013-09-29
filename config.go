package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type ReleaseHandlerType string

const (
	EMAIL ReleaseHandlerType = "email"
)

type config struct {
	BindAddress string                `toml:"bind_address"`
	BindPort    int                   `toml:"bind_port"`
	RepoConfigs map[string]repoConfig `toml:"repos"`
}

type repoConfig struct {
	Name                  string                          `toml:"name"`
	FullName              string                          `toml:"full_name"`
	ReleaseHandlerConfigs map[string]releaseHandlerConfig `toml:"release_handlers"`
}

type releaseHandlerConfig struct {
	HandlerType   ReleaseHandlerType `toml:"type"`
	HandlerConfig toml.Primitive     `toml:"config"`
}

func (repoConfig *repoConfig) buildRepo() (*Repo, error) {
	releaseHandlers := make([]NotificationHandler, 0)

	for _, releaseHandlerConfig := range repoConfig.ReleaseHandlerConfigs {
		switch releaseHandlerConfig.HandlerType {
		case EMAIL:
			releaseHandler, err := NewEmailReleaseHandler(releaseHandlerConfig.HandlerConfig)
			if err != nil {
				return nil, err
			}

			releaseHandlers = append(releaseHandlers, releaseHandler)
		}
	}

	return &Repo{repoConfig.Name, repoConfig.FullName, releaseHandlers}, nil
}

func (handler *HookHubHandler) ParseConfig(configFile string) error {
	var config config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return err
	}

	repos := make(map[string]*Repo)
	for _, repoConfig := range config.RepoConfigs {
		repo, err := repoConfig.buildRepo()
		if err != nil {
			return err
		}
		log.Printf("Configuring %s", repoConfig.FullName)
		repos[repoConfig.FullName] = repo
	}

	handler.bindAddress = config.BindAddress
	handler.bindPort = config.BindPort
	handler.repos = repos

	return nil
}
