package main

import (
	"log"

	"os"

	"io/ioutil"

	"path/filepath"

	"fmt"

	"github.com/hashicorp/hcl"
)

const (
	configFile = ".merged-prs"
)

const configUsageText = `
Configuration file must exist at $HOME/.merged-prs. It is where your GitHub token and Slack settings are stored.
Example config:

// ~/.merged-prs
githubtoken  = "foo"
githuborg    = "vsco"
slackwebhook = "https://hooks.slack.com/services/foo/bar/baz"
slackchannel = "#platform"

Once this is generated the script will work.
`

// Config struct usable through application
type Config struct {
	GithubToken  string
	GithubOrg    string
	SlackWebhook string
	SlackChannel string
}

func initConfig() Config {
	var configPath = os.Getenv("HOME")
	var c Config
	var err error

	config, err := ioutil.ReadFile(filepath.Join(configPath, configFile))
	if err != nil {
		log.Printf("Error parsing config: %s", err)
		configUsage()
	}

	err = hcl.Decode(&c, string(config))
	if err != nil {
		log.Fatalf("Configuration decode error: %s", err)
	}

	return c
}

func configUsage() {
	fmt.Println(configUsageText)
	os.Exit(1)
}
