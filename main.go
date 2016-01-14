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

	"github.com/rodaine/table"
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
	err = checkForGit()
	if err != nil {
		cmdLog = fmt.Sprintf("%s is not a valid git application, exiting.", gc)
		log.Fatal(cmdLog)
	}

	// Auth with GitHub
	client := authWithGitHub(config.GithubToken)

	// define table output
	tbl := table.New("PR", "Author", "Description", "URL")

	// Check to ensure our path is a Git repo, if not exit!
	checkPathIsGitRepo(repopath)

	// Output what we are about to do
	cmdLog = fmt.Sprintf("Determining merged branches between the following hashes: %s %s", ref1, ref2)
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
		cmdLog = fmt.Sprintf("No merged PRs / GH issues found between: %s %s", ref1, ref2)
		fmt.Println(cmdLog)
		os.Exit(2)
	}

	// curl github api to get the contents of the PR
	// at present, add it to a table row and output
	var lines []string
	for _, iid := range ids {
		pr, _, err := client.PullRequests.Get(config.GithubOrg, repo, iid)
		if err != nil {
			log.Fatal(err)
		}

		i := fmt.Sprintf("#%d", *pr.Number)
		u := fmt.Sprintf("@%s", *pr.User.Login)
		t := fmt.Sprintf("%s", *pr.Title)
		l := fmt.Sprintf("http://github.com/%s/%s/pull/%d", config.GithubOrg, repo, *pr.Number)

		tmpstr := fmt.Sprintf("%s (%s): %s (%s)", i, u, t, l)
		lines = append(lines, tmpstr)

		tbl.AddRow(i, u, t, l)
	}

	// Output results / Send to Slack
	output := new(bytes.Buffer)
	tbl.WithWriter(output)
	output.WriteString(fmt.Sprintf("Merged PRs between the following refs: %s %s", ref1, ref2))
	tbl.Print()

	// print the merge message
	mergedMessage := output.String()
	fmt.Println(mergedMessage)
	if !testMode {
		notifySlack(mergedMessage, config.SlackWebhook, config.SlackChannel)
	}
}

const usage = `Script can be used within a Git repository between any two hashes, tags, or branches

$ merged-prs <PREV> <NEW>

User should specify the older revision first ie. merging dev into master would necessitate that master is the older commit, and dev is the newer

$ merged-prs master dev
Determining merged branches between the following hashes: master dev
PR   Author    Description              URL
#55  @lucas    Typo 100 vs 1000         http://github.com/promiseofcake/foo/pull/55
#54  @lucas    LRU Cache Store Results  http://github.com/promiseofcake/foo/pull/54
`

func showUsage() {
	fmt.Println(usage)
	os.Exit(1)
}
