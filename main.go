package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type HookHubHandler struct {
	bindAddress string
	bindPort    int
	repos       map[string]*Repo
	debug       bool
}

type NotificationHandler interface {
	Handle(repo *Repo, notification GithubNotification, debug bool) error
}

type Repo struct {
	Name            string
	FullName        string
	releaseHandlers []NotificationHandler
}

func (handler *HookHubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	event := r.Header.Get("X-Github-Event")
	if event == "" {
		log.Println("No X-Github-Event header")
		return
	}

	defer r.Body.Close()
	var err error
	var bytes []byte
	if bytes, err = ioutil.ReadAll(r.Body); err != nil {
		fmt.Println("Failed to read request body", err)
		return
	}

	if handler.debug {
		log.Println("Received notification", string(bytes))
	}

	var notification GithubNotification
	if err = json.Unmarshal(bytes, &notification); err != nil {
		log.Println("Got errors parsing notification", err)
		return
	}
	notification.Event = event

	repo := handler.repos[notification.Repository.FullName]
	if repo == nil {
		log.Println("No repo configured for", notification.Repository.FullName)
		return
	}

	switch notification.Event {
	case GITHUB_EVENT_RELEASE:
		for _, eventHandler := range repo.releaseHandlers {
			if err = eventHandler.Handle(repo, notification, handler.debug); err != nil {
				log.Println("Error when handling release", err)
			} else {
				log.Println("Successfully handled release")
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
