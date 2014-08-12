package main

import (
	"time"
)

// The value of the X-Github-Event header
const (
	GITHUB_EVENT_RELEASE string = "release"
	GITHUB_EVENT_PUSH    string = "push"
)

type GithubPushEvent struct {
	Ref        string           `json:"ref"`
	Commits    []GithubCommit   `json:"commits"`
	HeadCommit GithubCommit     `json:"head_commit"`
	Repository GithubRepository `json:"repository"`
	Pusher     GithubUser       `json:"pusher"`
}

type GithubUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GithubCommit struct {
	Id      string     `json:"id"`
	Message string     `json:"message"`
	Url     string     `json:"url"`
	Author  GithubUser `json:"author"`
}

type GithubReleaseEvent struct {
	Action     string           `json:"action,omitempty"`
	Release    GithubRelease    `json:"release,omitempty"`
	Repository GithubRepository `json:"repository"`
	Event      string
}

type GithubRelease struct {
	Url             string               `json:"url"`
	AssetsUrl       string               `json:"assets_url"`
	UploadUrl       string               `json:"upload_url"`
	HtmlUrl         string               `json:"html_url"`
	Id              int                  `json:"id"`
	TagName         string               `json:"tag_name"`
	TargetCommitish string               `json:"target_commitish"`
	Name            string               `json:"name"`
	Body            string               `json:"body"`
	Draft           bool                 `json:"draft"`
	Prerelease      bool                 `json:"prerelease"`
	CreatedAt       time.Time            `json:"created_at"`
	PublishedAt     time.Time            `json:"published_at"`
	Assets          []GithubReleaseAsset `json:"assets,omitempty"`
}

type GithubReleaseAsset struct {
	Url           string    `json:"url"`
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Label         string    `json:"label"`
	ContentType   string    `json:"content_type"`
	State         string    `json:"state"`
	Size          int       `json:"size"`
	DownloadCount int       `json:"updated_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GithubRepository struct {
	Id       int    `json:"id"`
	FullName string `json:"full_name"`
}
