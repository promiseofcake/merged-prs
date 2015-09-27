package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"

	"github.com/fatih/color"
	"github.com/google/go-github/github"
	"github.com/rodaine/table"
	"golang.org/x/oauth2"
)

const (
	gc       = "git"
	tokenVar = "MPR_GITHUB_TOKEN"
)

const (
	gho = "vsco"
)

func main() {

	var err error
	var cmdLog string
	var repopath, repo string

	// Get token configuration
	ght := os.Getenv(tokenVar)
	if ght == "" {
		cmdLog = fmt.Sprintf("GitHub token missing, please set %s", tokenVar)
		log.Fatal(cmdLog)
	}

	// Get command line arguments
	wd, _ := os.Getwd()
	flag.StringVar(&repopath, "path", wd, "Path to git repo")
	flag.Parse()

	// Get remaining arguments
	a := flag.Args()
	if len(a) < 2 {
		showUsage()
		os.Exit(1)
	}

	// Store refs
	ref1 := a[0]
	ref2 := a[1]

	// Derive repo from the path
	repo = path.Base(repopath)

	// Determine that Git is installed
	gchk := exec.Command(gc, "--version")
	err = gchk.Run()
	if err != nil {
		cmdLog = fmt.Sprintf("%s is not a valid git application, exiting.", gc)
		log.Fatal(cmdLog)
	}

	// Auth with GitHub
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ght},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	client := github.NewClient(tc)

	// define table output
	columnFmt := color.New(color.FgYellow).SprintfFunc()
	tbl := table.New("PR", "Author", "Description", "URL")
	tbl.WithFirstColumnFormatter(columnFmt)

	// Check to ensure our path is a Git repo, if not exit!
	chkr := exec.Command(gc, "-C", repopath, "status")
	err = chkr.Run()
	if err != nil {
		cmdLog = fmt.Sprintf("%s is not a valid git repository! Exiting", repopath)
		log.Fatal(cmdLog)
	}

	// Output what we are about to do
	cmdLog = fmt.Sprintf("Determining merged branches between the following hashes: %s %s", ref1, ref2)
	fmt.Print(cmdLog)

	// Determine the merged branches between the two hashes
	marg := fmt.Sprintf("%s...%s", ref1, ref2)
	c := exec.Command(gc, "-C", repopath, "log", "--merges", "--pretty=format:\"%s\"", marg)
	out, err := c.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	err = c.Start()
	if err != nil {
		log.Fatal(err)
	}

	// iteratre through matcthes, and pull out the issues id into a slice
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
		pr, _, err := client.PullRequests.Get(gho, repo, iid)
		if err != nil {
			log.Fatal(err)
		}

		i := fmt.Sprintf("#%d", *pr.Number)
		u := fmt.Sprintf("@%s", *pr.User.Login)
		t := fmt.Sprintf("%s", *pr.Title)
		l := fmt.Sprintf("http://github.com/%s/%s/pull/%d", gho, repo, *pr.Number)

		tmpstr := fmt.Sprintf("%s (%s): %s (%s)", i, u, t, l)
		lines = append(lines, tmpstr)

		tbl.AddRow(i, u, t, l)
	}

	// Output results
	tbl.Print()
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
