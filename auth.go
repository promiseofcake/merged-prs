package main

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func authWithGitHub(tkn string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tkn},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}
