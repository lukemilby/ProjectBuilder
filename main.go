package main

import (
	"bytes"
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/google/go-github/v45/github"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed templates/README.md
var t []byte

type Content struct {
	Name string
}

var (
	private     = flag.Bool("private", false, "set repo to private")
	description = flag.String("description", "", "set description of repo")
	msg         = "add README.md"
	pba         = "Project Builder"
	aEmail      = "luke.milby@gmail.com"
)

// feature: profiles custom pre and post
func main() {

	flag.Parse()

	args := flag.Args()

	if len(args) <= 0 {
		fmt.Println("no arguments provided, use \"new\" to create a project")
	}

	switch args[0] {
	case "new":
		//		buildProject(flag.Args()[1:])
		fmt.Println("You caught me at a good time. Ill make that now.")
	default:
		fmt.Println("yeah I need you to tell me \"new\" if you want me to make that")
		os.Exit(1)
	}

	viper.SetConfigName("config")            // name of config file (without extension)
	viper.SetConfigType("yaml")              // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.probuilder") // call multiple times to add many search paths
	viper.AddConfigPath(".")                 // optionally look for config in the working directory
	err := viper.ReadInConfig()              // Find and read the config file
	if err != nil {                          // Handle errors reading the config file
		fmt.Println("Shit boss we hit an error...")
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	// Args are done
	// Setup Viper to pull config.yml
	// Contains Github app creds

	// Create project in current location

	// Github client
	// Github check name, create repo private
	// Push README.md thats a template
	// Parse Template
	tmp, err := template.New("readme").Parse(string(t))

	var pp = viper.GetString("PROPATH")
	var content Content
	content.Name = strings.Join(args[1:], " ")
	var buf bytes.Buffer
	err = tmp.Execute(&buf, content)
	if err != nil {
		panic(err)
	}
	fmt.Println("template loaded")

	// We dont want to do this.

	// ============================= Setup Git ===========

	if err != nil {
		panic(err)
	}

	// =============================== Create File structure ===============

	// ===================== GITHUB MAKES THE REPO ======================
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: viper.GetString("PAT")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// Get all repos
	repos, _, err := client.Repositories.ListAll(context.Background(), &github.RepositoryListAllOptions{})
	if err != nil {
		panic(err)
	}
	for _, r := range repos {
		if *r.Name == strings.Join(args[1:], " ") {
			fmt.Println("Yo, we have some bad news. You already created a repo with that name")
			os.Exit(1)
		}
	}

	name := strings.Join(args[1:], " ")
	repo, _, err := client.Repositories.Create(
		context.Background(),
		"",
		&github.Repository{
			Name:        &name,
			Description: description,
			Private:     private,
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(*repo.URL, " created")

	fmt.Println("Calling Build Utility")
	switch viper.GetString("LANG") {
	case "rust":
		cmd := exec.Command("cargo", "new", pp+name)
		err := cmd.Run()
		if err != nil {
			fmt.Println("how the fuck did this get made already")
			panic(err)
		}
	default:
		fmt.Println("I dont support that")
	}

	fmt.Println("configuring git")

	cmd := exec.Command("git", "init", pp+name)
	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	gitRepo, err := git.PlainOpen(pp + name)
	if err != nil {
		panic(err)
	}

	fmt.Println("adding README.md")
	f, err := os.Create(pp + name + "/README.md")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}

	// getting  Config
	cfg, err := gitRepo.Config()
	if err != nil {
		panic(err)
	}

	// setting remote to the repo we created
	cfg.Remotes["origin"] = &config.RemoteConfig{
		Name: "origin",
		URLs: []string{"git@github.com:" + viper.GetString("GITHUB_NAME") + "/" + name + ".git"},
	}

	//from repo.GitURL = git://github.com/lukemilby/more40.git
	// need to remove and edit to git@github.com:lukemilby/more40.git

	gitRepo.Storer.SetConfig(cfg)
	fmt.Println("I think we did it. reap created here ", pp+name)

}
