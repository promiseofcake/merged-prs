package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
)

const (
	gc = "git"
)

func parseArgs() (string, string) {
	a := flag.Args()
	if len(a) < 2 {
		showUsage()
		os.Exit(1)
	}

	return a[0], a[1]
}

func parseFlags() (string, bool) {
	var path string
	var test bool

	wd, _ := os.Getwd()

	flag.StringVar(&path, "path", wd, "Path to git repo, defaults to working directory.")
	flag.BoolVar(&test, "test", false, "Run command in test mode (do not notify Slack)")
	flag.Parse()

	return path, test
}

func checkForGit() error {
	gchk := exec.Command(gc, "--version")
	return gchk.Run()
}

func checkPathIsGitRepo(repopath string) {
	check := exec.Command(gc, "-C", repopath, "status")
	err := check.Run()
	if err != nil {
		log.Fatalf("%s is not a valid git repository! Exiting", repopath)
	}
}
