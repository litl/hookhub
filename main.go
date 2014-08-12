package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/litl/hookhub/fogbugz"
	"log"
	"net/http"
	"os"
	"strconv"
)

type HookHubHandler struct {
	bindAddress    string
	bindPort       int
	repos          map[string]*Repo
	debug          bool
	fogbugzSession *fogbugz.Session
}

type ReleaseHandler interface {
	Handle(repo *Repo, notification GithubReleaseEvent, debug bool) error
}

type PushHandler interface {
	Handle(hookHubHandler *HookHubHandler, notification GithubPushEvent, debug bool) error
}

type Repo struct {
	Name            string
	FullName        string
	releaseHandlers []ReleaseHandler
	pushHandlers    []PushHandler
}

func (handler *HookHubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-Github-Event")
	if event == "" {
		log.Println("No X-Github-Event header")
		return
	}

	defer r.Body.Close()
	var err error
	if err = r.ParseForm(); err != nil {
		fmt.Println("Failed to parse form values from request body", err)
	}
	var jsonStr = r.PostFormValue("payload")

	if handler.debug {
		log.Println("Received notification", jsonStr)
	}

	switch event {
	case GITHUB_EVENT_RELEASE:
		var notification GithubReleaseEvent
		if err = json.Unmarshal([]byte(jsonStr), &notification); err != nil {
			log.Println("Got errors parsing notification", err)
			log.Println(jsonStr)
			return
		}

		repo := handler.repos[notification.Repository.FullName]
		if repo == nil {
			log.Println("No repo configured for", notification.Repository.FullName)
			return
		}
		for _, eventHandler := range repo.releaseHandlers {
			if err = eventHandler.Handle(repo, notification, handler.debug); err != nil {
				log.Println("Error when handling release", err)
			} else {
				log.Println("Successfully handled release")
			}
		}
		break
	case GITHUB_EVENT_PUSH:
		var notification GithubPushEvent
		if err = json.Unmarshal([]byte(jsonStr), &notification); err != nil {
			log.Println("Got errors parsing notification", err)
			log.Println(jsonStr)
			return
		}

		repo := handler.repos[notification.Repository.FullName]
		if repo == nil {
			log.Println("No repo configured for", notification.Repository.FullName)
			return
		}
		for _, eventHandler := range repo.pushHandlers {
			if err = eventHandler.Handle(handler, notification, handler.debug); err != nil {
				log.Println("Error when handling push", err)
			} else {
				log.Println("Successfully handled push")
			}
		}
		break
	}
}

func main() {
	var configFile string
	handler := new(HookHubHandler)
	flag.BoolVar(&handler.debug, "debug", false, "Run in debug mode")
	flag.StringVar(&configFile, "config", "hookhub.toml", "Path to configuration file")
	flag.Parse()

	if handler.debug {
		log.Println("Debug mode enabled")
	}

	if err := handler.ParseConfig(configFile); err != nil {
		log.Fatalln("Failed to initialize from config", err)
		return
	}

	portEnv := os.Getenv("PORT")
	if portEnv != "" {
		if port, err := strconv.Atoi(portEnv); err == nil {
			handler.bindPort = port
		}
	}

	http.Handle("/github_webhook", handler)

	log.Printf("Listening on %s:%d\n", handler.bindAddress, handler.bindPort)
	http.ListenAndServe(fmt.Sprintf("%s:%d", handler.bindAddress, handler.bindPort), nil)
}
