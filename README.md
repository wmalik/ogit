# ogit

TUI for organizing git repositories.

### Configuration

Add a section in your `~/.gitconfig`:

```
[ogit]
  storagePath = /home/arthur/ogit
  fetchUserRepos = true
  sshAuth = ssh-agent # or sshAuth = /path/to/ssh_private_key
[ogit "github"]
  orgs = tpope
[ogit "gitlab"]
  orgs = fdroid
```

### Usage

Generate a GitHub personal access token
[here](https://github.com/settings/tokens) with full `repo` access.

```
$ ogit --help
NAME:
   ogit - Organize git repositories

USAGE:
   ogit [global options] command [command options] [arguments...]

COMMANDS:
   clone, c  Clone repositories in bulk
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --nosync    Disable syncing of repositories metadata at startup (default: false)
   --clear     Clear all local repository metadata (default: false)
   --help, -h  show help (default: false)
```

#### Examples

```
export GITHUB_TOKEN="yourpersonalaccesstoken_with_full_repo_access"
export GITLAB_TOKEN="yourtoken_with_read_api_scope"
go run cmd/ogit/main.go
go run cmd/ogit/main.go --cached
go run cmd/ogit/main.go --clear
```


Please note that the GitHub API enforces [rate limits](https://docs.github.com/en/developers/apps/building-github-apps/rate-limits-for-github-apps)
(5000 requests per hour) when a personal access token is used.
