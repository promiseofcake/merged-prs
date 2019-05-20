package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcl"
)

const (
	configFile = ".merged-prs"
)

const configUsageText = `
Configuration file must exist at $HOME/.merged-prs. It is where your GitHub token and Slack settings are stored.
Example config:

// ~/.merged-prs
GitHub {
	Token = "foo"
	Org = "vsco"
}

Slack {
	WebhookURL = "https://hooks.slack.com/services/foo/bar/baz"
	Channel  = "#platform"
	Emoji  = ":shipit:"
}

Once this is generated the script will work.`

type githubConfig struct {
	Token string
	Org   string
}

type slackConfig struct {
	WebhookURL string
	Channel    string
	Emoji      string
}

// Config is general configuration used throughout the applicaiton
type Config struct {
	Github githubConfig
	Slack  slackConfig
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
