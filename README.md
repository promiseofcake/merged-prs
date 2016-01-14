# merged-prs

merged-prs is a tool to assit in determining differences between git hashes based upon the GitHub PR as the vechicle for change.

## Features

- Given a set of git hashes, a list of merged GitHub Pull Requests will be retrieved and parsed.
- List will be output to console contining the PR #, Author, Summary, and a link to the PR
- If Slack credentials are configured, a notification will be sent to the channel of your choosing

## Configuration

A `.merged-prs` configuration file must be created in your `$HOME` directory. This configuration uses HCL Syntax.

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

## Usage

Calling the `merged-prs` tool will act in the current directory context.

### Flags

```
  -test <Do not notify Slack>
  -path <Specify path to repo in order to use outside the context of a repo

```

### Example

```
# merged-prs [-test] [-path /path/to/repo] [Previous Git Hash] [Latest Git Hash]
# User should specify the older revision first ie. merging dev into master would necessitate that master is the older commit, and dev is the newer

$ merged-prs master dev
Determining merged branches between the following hashes: master dev

REPO: Merged GitHub PRs between the following refs: master dev
PR   Author    Description              URL
#55  @lucas    Typo 100 vs 1000         http://github.com/promiseofcake/foo/pull/55
#54  @lucas    LRU Cache Store Results  http://github.com/promiseofcake/foo/pull/54
```
