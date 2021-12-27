# ogit

TUI for browsing an organization on GitHub and GitLab.

### Configuration

Add a section in your `~/.gitconfig`:

```
[ogit]
  orgs = padawin, tpope, charmbracelet, wmalik
  clonedirpath = /absolute/path/on/disk
```

Please note that the `orgs` parameter currently only supports _public_ users and
organisations.

### Run

```
go run cmd/browser/main.go
```

Please note that the current implementation uses the public GitHub API and is
therefore subject to rate limits. In the case of an error returned by the GitHub
API, the error text is shown in the status bar (on the top right).
