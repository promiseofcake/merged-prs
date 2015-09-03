package main

import (
	"fmt"
	"os"
)

func main() {
	a := os.Args

	if len(a) <= 2 {
		showUsage()
		os.Exit(1)
	}

}

const usage = `
merged-prs <HASH> <HASH>
`

func showUsage() {
	fmt.Println(usage)
}
