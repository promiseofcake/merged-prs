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

	"sync"

	"github.com/google/go-github/github"
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

	// Get info from GitHub (experiment with concurrency)
	pulls := processPullRequests(ids, client, config, repo)

	for _, pull := range pulls {
		i := fmt.Sprintf("#%d", *pull.Number)
		u := fmt.Sprintf("@%s", *pull.User.Login)
		t := fmt.Sprintf("%s", *pull.Title)
		l := fmt.Sprintf("%s/%s/%s/pull/%d", gitHubRoot, config.Github.Org, repo, *pull.Number)
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

// Function to go to GetHub and process the passed pull requests
func processPullRequests(ids []int, client *github.Client, config Config, repo string) []*github.PullRequest {

	var wg sync.WaitGroup
	jobs := make(chan *github.PullRequest, len(ids))
	wg.Add(len(ids))

	// List of Pull Requests
	pulls := []*github.PullRequest{}

	for _, item := range ids {
		go func(client *github.Client, org string, r string, id int) {
			pr, _, err := client.PullRequests.Get(org, r, id)
			if err != nil {
				log.Printf("Error retrieving pull request: %d, %v", id, err)
				wg.Done()
				return
			}

			jobs <- pr
			wg.Done()
		}(client, config.Github.Org, repo, item)
	}

	wg.Wait()
	close(jobs)

	for pull := range jobs {
		pulls = append(pulls, pull)
	}

	return pulls
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
