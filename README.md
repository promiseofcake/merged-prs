# merged-prs

merged-prs is a go tool to assit in determining differences between git hashes based upon the GitHub PR as the vechicle for change.

## Requirements

- Go 1.5
- GPM
- Git
- GitHub
- Slack (optional)

## Features

- Built in Go
- Given a set of git hashes, a list of merged GitHub Pull Requests will be retrieved and parsed.
- List will be output to console contining the PR #, Author, Summary, and a link to the PR
- If Slack credentials are configured, a notification will be sent to the channel of your choosing

## Configuration

A `.merged-prs` configuration file must be created in your `$HOME` directory. This configuration uses [HashiCorp's HCL Syntax](https://github.com/hashicorp/hcl#syntax).
Using the example configuration below create a file in `~/.merged-prs`

*Example Config*

```
// GitHub Personal Access Token && Org/Username required for use.
GitHub {
    Token = "foo"
    Org   = "promiseofcake"
}

// Optional config for Slack notifications
Slack {
    WebhookURL = "https://hooks.slack.com/services/foo/bar/baz"
    Channel    = "@lucas"
    Emoji      = ":shipit:"
}

```

## Installation

In order to use the `merged-prs` tool one should use `go get`

```
go get github.com/promiseofcake/merged-prs
cd $GOPATH/src/github.com/promiseofcake/merged-prs
gpm install
go install
```

## Usage

Calling the `merged-prs` tool will act in the current directory's context.

### Flags

```
# Runtime Flags
  -test <do not notify Slack, only output to console>
  -path <Specify path to repo in order to use outside the context of a repo>

```

### Example

```
# merged-prs [-test] [-path /path/to/repo] [Previous Git Hash] [Current Git Hash]
# User should specify the older revision first ie. merging `dev` into `master` would necessitate that `master` is the older commit, and `dev` is the newer

$ merged-prs master dev
Determining merged branches between the following hashes: master dev

REPO: Merged GitHub PRs between the following refs: master dev
PR   Author    Description              URL
#55  @lucas    Typo 100 vs 1000         http://github.com/promiseofcake/foo/pull/55
#54  @lucas    LRU Cache Store Results  http://github.com/promiseofcake/foo/pull/54
```
