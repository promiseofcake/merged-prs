package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func authWithGitHub(tkn string) *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: tkn},
	)
	tc := oauth2.NewClient(context.TODO(), ts)
	return github.NewClient(tc)
}
