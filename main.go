package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gogits/git"
)

// var cwd string

func main() {

	wd, _ := os.Getwd()
	fmt.Println(wd)

	r, err := git.OpenRepository(wd)
	if err != nil {
		log.Fatal("Could not open repository")
	}

	a := os.Args

	if len(a) <= 2 {
		showUsage()
		os.Exit(1)
	}

	fmt.Println("arguments:", a)

	refA, err := r.GetCommit(a[1])
	if err != nil {
		log.Fatal(a[1] + "not a valid hash")
	}
	refB, err := r.GetCommit(a[2])
	if err != nil {
		log.Fatal(a[2] + "not a valid hash")
	}

	fmt.Println(refA.Message())
	fmt.Println(refB.Message())
}

const usage = `
merged-prs <HASH> <HASH>
`

func showUsage() {
	fmt.Println(usage)
}
