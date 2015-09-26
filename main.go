package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const (
	gc = "/usr/local/bin/git"
)

const (
	gho = "vsco"
	ghr = "image"
	ght = "foo"
)

func main() {

	var err error

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ght},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	client := github.NewClient(tc)

	var cmdLog string

	// Set or get the current working directory
	// wd, _ := os.Getwd()
	wd := "/Users/lucas/Workspace/Work/vsco/image"
	fmt.Println(wd)

	a := os.Args

	if len(a) <= 2 {
		showUsage()
		os.Exit(1)
	}

	// Check to ensure we are in a git repo, if we are not, exit!
	chkr := exec.Command(gc, "-C", wd, "status")
	err = chkr.Run()
	if err != nil {
		cmdLog = fmt.Sprintf("%s is not a valid git repository! Exiting", wd)
		log.Fatal(cmdLog)
	}

	// Output what we are about to do
	cmdLog = fmt.Sprintf("Determining merged branches between the following hashes: %s %s", a[1], a[2])
	fmt.Println(cmdLog)

	// Determine the merged branches between the two hashes
	c := exec.Command(gc, "-C", wd, "log", "--merges", "--pretty=format:\"%s\"", a[1]+"..."+a[2])

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
		// fmt.Println(t)
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
		cmdLog = fmt.Sprintf("No merged PRs / GH issues found between: %s %s", a[1], a[2])
		log.Fatal(cmdLog)
	}

	// curl github api to get the contents of the PR

	var lines []string

	for _, iid := range ids {
		fmt.Println(iid)
		pr, _, err := client.PullRequests.Get(gho, ghr, iid)
		if err != nil {
			log.Fatal(err)
		}

		// pull out a lot of values, content, name
		i := *pr.Number
		u := *pr.User.Login
		t := *pr.Title

		// pr link
		l := fmt.Sprintf("http://github.com/%s/%s/pull/%d", gho, ghr, i)

		tmpstr := fmt.Sprintf("#%d (@%s): %s (%s)", i, u, t, l)

		// push output into an array
		// tmpstr := "#" + strconv.Itoa(i) + " (@" + u + "): " + t + " (" + l + ")"
		lines = append(lines, tmpstr)
	}

	fmt.Println(lines)

	// notify slack
}

const usage = `
merged-prs <HASH> <HASH>
`

func showUsage() {
	fmt.Println(usage)
}
