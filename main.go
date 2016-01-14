package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"

	"bytes"

	"strings"

	"github.com/rodaine/table"
)

const (
	gitHubRoot = "http://github.com"
	gc         = "git"
)

func main() {

	var err error
	var cmdLog string

	config := initConfig()

	// Get command line arguments and store args as refs
	repopath, testMode := parseFlags()
	ref1, ref2 := parseArgs()

	// Derive repo from the path
	repo := path.Base(repopath)

	// Determine that Git is installed
	checkForGit()

	// Auth with GitHub
	client := authWithGitHub(config.Github.Token)

	// define table output
	tbl := table.New("PR", "Author", "Description", "URL")

	// Check to ensure our path is a Git repo, if not exit!
	checkPathIsGitRepo(repopath)

	// Output what we are about to do
	cmdLog = fmt.Sprintf("Determining merged branches between the following hashes: %s %s \n", ref1, ref2)
	fmt.Println(cmdLog)

	// Determine the merged branches between the two hashes
	marg := fmt.Sprintf("%s..%s", ref1, ref2)
	c := exec.Command(gc, "-C", repopath, "log", "--merges", "--pretty=format:\"%s\"", marg)
	out, err := c.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = c.Start()
	if err != nil {
		log.Fatal(err)
	}

	// iterate through matches, and pull out the issues id into a slice
	var ids []int
	s := bufio.NewScanner(out)
	for s.Scan() {
		t := s.Text()
		r, _ := regexp.Compile("#([0-9]+)")
		sm := r.FindStringSubmatch(t)
		if len(sm) > 0 {
			i, err := strconv.Atoi(sm[1])
			if err == nil {
				ids = append(ids, i)
			}
		}
	}
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	// if the list is empty, error out
	if len(ids) == 0 {
		cmdLog = fmt.Sprintf("No merged PRs / GitHub issues found between: %s %s", ref1, ref2)
		fmt.Println(cmdLog)
		os.Exit(2)
	}

	// Call GitHub API to get the contents of the individual PRs and add as Rows to the Table
	var lines []string
	for _, iid := range ids {
		pr, _, err := client.PullRequests.Get(config.Github.Org, repo, iid)
		if err != nil {
			log.Fatal(err)
		}

		i := fmt.Sprintf("#%d", *pr.Number)
		u := fmt.Sprintf("@%s", *pr.User.Login)
		t := fmt.Sprintf("%s", *pr.Title)
		l := fmt.Sprintf("%s/%s/%s/pull/%d", gitHubRoot, config.Github.Org, repo, *pr.Number)

		tmpstr := fmt.Sprintf("%s (%s): %s (%s)", i, u, t, l)
		lines = append(lines, tmpstr)

		tbl.AddRow(i, u, t, l)
	}

	// Generate results output
	output := new(bytes.Buffer)
	tbl.WithWriter(output)
	output.WriteString(fmt.Sprintf("%s: Merged GitHub PRs between the following refs: %s %s", strings.ToUpper(repo), ref1, ref2))
	tbl.Print()

	// Print the merge message to console and notifies slack
	mergedMessage := output.String()
	fmt.Println(mergedMessage)
	if !testMode {
		notifySlack(mergedMessage, config.Slack)
	}
}

const usage = `Script can be used within a Git repository between any two hashes, tags, or branches

$ merged-prs <PREV> <NEW>

User should specify the older revision first ie. merging dev into master would necessitate that master is the older commit, and dev is the newer

$ merged-prs master dev
Determining merged branches between the following hashes: master dev

REPO: Merged GitHub PRs between the following refs: master dev
PR   Author    Description              URL
#55  @lucas    Typo 100 vs 1000         http://github.com/promiseofcake/foo/pull/55
#54  @lucas    LRU Cache Store Results  http://github.com/promiseofcake/foo/pull/54
`

func showUsage() {
	fmt.Println(usage)
	os.Exit(1)
}
