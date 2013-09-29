package main

import (
	"time"
)

// The value of the X-Github-Event header
const (
	GITHUB_EVENT_RELEASE string = "release"
)

type GithubNotification struct {
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
	NumberOfCrashes int                  `json:"number_of_crashes"`
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
