# ogit

TUI for browsing an organization on GitHub and GitLab.

### Configuration

Add a section in your `~/.gitconfig`:

```
[ogit]
  orgs = padawin, tpope, charmbracelet, wmalik
  clonedirpath = /absolute/path/on/disk
  fetchAuthenticatedUserRepos = true
  useSSHAgent = true
  privKeyPath =
```

Please note that `privKeyPath` must be specified in the config, however it
can be an empty string.

Please note that the `orgs` parameter currently only supports _public_ users and
organisations.

### Run

Generate a GitHub personal access token
[here](https://github.com/settings/tokens) with full `repo` access.

```
export GITHUB_TOKEN="yourpersonalaccesstoken"
go run cmd/browser/main.go
```

Please note that the GitHub API enforces [rate limits](https://docs.github.com/en/developers/apps/building-github-apps/rate-limits-for-github-apps)
(5000 requests per hour) when a personal access token is used.
