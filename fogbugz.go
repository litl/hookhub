package main

import (
	"fmt"
	"log"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/litl/hookhub/fogbugz"
)

func (hookHubHandler *HookHubHandler) GetFogbugzSession(config fogbugz.Config) (*fogbugz.Session, error) {
	var err error
	if hookHubHandler.fogbugzSession == nil {
		hookHubHandler.fogbugzSession, err = fogbugz.NewSession(config)
	}
	return hookHubHandler.fogbugzSession, err
}

type FogbugzConfig struct {
	Host     string `toml:"host"`
	Email    string `toml:"email"`
	Password string `toml:"password"`
}

func (config *FogbugzConfig) GetHost() string {
	return config.Host
}

func (config *FogbugzConfig) GetEmail() string {
	return config.Email
}

func (config *FogbugzConfig) GetPassword() string {
	return config.Password
}

type FogbugzResolveHandler struct {
	Config FogbugzConfig
}

func (handler FogbugzResolveHandler) Handle(hookHubHandler *HookHubHandler, notification GithubPushEvent, debug bool) error {
	re := regexp.MustCompile(`(?i)Fixes:? http[^?]+[?].*?(\d+)\b`)
	session, err := hookHubHandler.GetFogbugzSession(&handler.Config)
	if err != nil {
		return err
	}

	for _, commit := range notification.Commits {
		matches := re.FindAllStringSubmatch(commit.Message, -1)
		for _, m := range matches {
			log.Printf("Resolving bug #%s\n", m[1])
			session.ResolveBug(m[1], fmt.Sprintf("This was fixed in Github: %s", commit.Url))
		}
	}
	return nil
}

func NewFogbugzResolveHandler(fogbugzDefaultConfig FogbugzConfig, configPrimitive toml.Primitive) (PushHandler, error) {
	// TODO: Parse anything specific out of configPrimitive, save in handler,
	//       so we can make our own session if we need one.
	return &FogbugzResolveHandler{fogbugzDefaultConfig}, nil
}
