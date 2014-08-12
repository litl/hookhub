package main

import (
	"github.com/BurntSushi/toml"
	"log"
)

type HandlerType string

const (
	EMAIL           HandlerType = "email"
	FOGBUGZ_RESOLVE HandlerType = "fogbugz_resolve"
)

type config struct {
	BindAddress          string                `toml:"bind_address"`
	BindPort             int                   `toml:"bind_port"`
	RepoConfigs          map[string]repoConfig `toml:"repos"`
	FogbugzDefaultConfig FogbugzConfig         `toml:"fogbugz_default_config"`
}

type repoConfig struct {
	Name                  string                   `toml:"name"`
	FullName              string                   `toml:"full_name"`
	ReleaseHandlerConfigs map[string]handlerConfig `toml:"release_handlers"`
	PushHandlerConfigs    map[string]handlerConfig `toml:"push_handlers"`
}

type handlerConfig struct {
	HandlerType   HandlerType    `toml:"type"`
	HandlerConfig toml.Primitive `toml:"config"`
}

func (repoConfig *repoConfig) buildRepo(fogbugzDefaultConfig FogbugzConfig) (*Repo, error) {
	releaseHandlers := make([]ReleaseHandler, 0)

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

	pushHandlers := make([]PushHandler, 0)

	for _, pushHandlerConfig := range repoConfig.PushHandlerConfigs {
		switch pushHandlerConfig.HandlerType {
		case FOGBUGZ_RESOLVE:
			pushHandler, err := NewFogbugzResolveHandler(fogbugzDefaultConfig, pushHandlerConfig.HandlerConfig)
			if err != nil {
				return nil, err
			}

			pushHandlers = append(pushHandlers, pushHandler)
		}
	}

	return &Repo{repoConfig.Name, repoConfig.FullName, releaseHandlers, pushHandlers}, nil
}

func (handler *HookHubHandler) ParseConfig(configFile string) error {
	var config config
	if _, err := toml.DecodeFile(configFile, &config); err != nil {
		return err
	}

	repos := make(map[string]*Repo)
	for _, repoConfig := range config.RepoConfigs {
		repo, err := repoConfig.buildRepo(config.FogbugzDefaultConfig)
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
