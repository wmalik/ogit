# ogit

TUI for organizing git repositories.

### Configuration

Add a section in your `~/.gitconfig`:

```
[ogit]
  orgs = padawin, tpope, charmbracelet, wmalik
  gitlabGroups = fdroid
  clonedirpath = /absolute/path/on/disk
  fetchAuthenticatedUserRepos = true
  useSSHAgent = true
  privKeyPath =
```

Please note that `privKeyPath` must be specified in the config, however it
can be an empty string. Also, `privKeyPath` can not point to a private key with
passphrase. In that case, do the following:

* add the private key to ssh-agent
* set `privKeyPath` to an empty string
* set `useSSAgent = true`

### Run

Generate a GitHub personal access token
[here](https://github.com/settings/tokens) with full `repo` access.

```
$ ogit --help

Usage: ogit [OPTION]
Organize git repositories
Sync repositories on startup unless -nosync is specified

  -clear
    	Clear all local repository metadata
  -nosync
    	Disable syncing of repositories metadata at startup
```

#### Examples

```
export GITHUB_TOKEN="yourpersonalaccesstoken"
export GITLAB_TOKEN="yourtokenhere"
go run cmd/ogit/main.go
go run cmd/ogit/main.go --nosync
go run cmd/ogit/main.go --clear
```


Please note that the GitHub API enforces [rate limits](https://docs.github.com/en/developers/apps/building-github-apps/rate-limits-for-github-apps)
(5000 requests per hour) when a personal access token is used.
