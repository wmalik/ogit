# ogit

TUI and CLI for organizing git repositories across multiple providers (e.g.
GitHub, GitLab).

### Install

#### Install via go-install

```
go install -ldflags="-X main.version=installed-with-go-install" github.com/wmalik/ogit/cmd/ogit@latest
```

#### Install via source code

```
go build -ldflags="-X main.version=local-build" -o ogit cmd/ogit/main.go
```

#### Install via GitHub Releases

Download a pre-built binary for your platform [here](https://github.com/wmalik/ogit/releases/latest)

### Configuration

Add `[ogit]` sections to `~/.gitconfig`.

<details>
  <summary>Config for GitHub repositories (using ssh-agent)</summary>

```
[ogit]
  storagePath = /absolute/path/on/disk
  fetchUserRepos = false
  sshAuth = ssh-agent
[ogit "github"]
  orgs = tpope, charmbracelet
```
</details>

<details>
  <summary>Config for GitHub and GitLab repositories (using ssh-agent)</summary>

```
[ogit]
  storagePath = /absolute/path/on/disk
  fetchUserRepos = false
  sshAuth = ssh-agent
[ogit "github"]
  orgs = tpope, charmbracelet
[ogit "gitlab"]
  orgs = fdroid
```
</details>

<details>
  <summary>Config for user's repositories only (using ssh-agent)</summary>

```
[ogit]
  storagePath = /absolute/path/on/disk
  fetchUserRepos = true
  sshAuth = ssh-agent
```
</details>

<details>
  <summary>Config for GitHub and GitLab repositories (using private SSH key)</summary>

```
[ogit]
  storagePath = /absolute/path/on/disk
  fetchUserRepos = false
  sshAuth = /absolute/path/to/privatekey
[ogit "github"]
  orgs = tpope
[ogit "gitlab"]
  orgs = fdroid
```
</details>

#### Authentication

##### SSH Auth

The `sshAuth` attribute in `~/.gitconfig` is used to perform the "git clone"
operation for both Github and GitLab repositories.
An SSH key pair must be available on the host machine and associated with the
GitHub and GitLab accounts. The SSH key pair can be fetched from either an
ssh-agent or from a file on disk. If the private key is protected with
a passphrase, the only way to use it is through ssh-agent.

##### GitHub/GitLab API Auth

Personal access tokens for GitHub/GitLab must be configured via the following
environment variables:

* `GITHUB_TOKEN` (with `repo` scope)
* `GITLAB_TOKEN` (with `read_api` scope)

The tokens can be generated [here](https://github.com/settings/tokens/new) and
[here](https://gitlab.com/-/profile/personal_access_tokens).

### Usage

#### Setup credentials

```
ssh-add ~/.ssh/your_private_key
export GITHUB_TOKEN="yourpersonalaccesstoken_with_full_repo_access"
export GITLAB_TOKEN="yourtoken_with_read_api_scope"
```

#### Fetch repository metadata and launch TUI

```
ogit fetch && ogit
```

#### Clone all repositories belonging to an org

```
ogit clone --org tpope
```

#### Clone some repositories belonging to an org

```
ogit clone --org tpope --filter vim
```
