package main

import (
	"fmt"
	"os"
)

var cwd string

func main() {

	cwd = getCwd()
	fmt.Println(cwd)

	a := os.Args

	if len(a) <= 2 {
		showUsage()
		os.Exit(1)
	}

	fmt.Println("arguments:", a)

}

const usage = `
merged-prs <HASH> <HASH>
`

func showUsage() {
	fmt.Println(usage)
}

func getCwd() string {
	return os.Getenv("PWD")
}
