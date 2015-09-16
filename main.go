package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

const (
	gc = "/usr/local/bin/git"
)

func main() {

	var cmdLog string
	var err error

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
	var ids []string

	s := bufio.NewScanner(out)
	for s.Scan() {
		t := s.Text()
		fmt.Println(t)
		r, _ := regexp.Compile("#([0-9]+)")
		sm := r.FindStringSubmatch(t)
		ids = append(ids, sm[1])
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

	// pull out a lot of values, content, name

	// push output into an array

	// notify slack
}

const usage = `
merged-prs <HASH> <HASH>
`

func showUsage() {
	fmt.Println(usage)
}
