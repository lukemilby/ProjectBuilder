# Project Builder

A tool for building projects
	
- Creates Repo on Github and your local host
- Create project with language specific tooling
  - Support:
    - Rust
- Initiates a git repo in a target build path
- Set a README.md file from template

### Config

probuild will look at the current working directory and $HOME/.probuilder/config.yaml

check out example.conf.yaml

PAT = Personal Access Token from Github

Configured to talk to git with an sshkey

### Run
```go
go run main.go new <repo name>
```
