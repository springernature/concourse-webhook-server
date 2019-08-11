package main

import (
	"cws"
	"fmt"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	gitHub, err := cws.NewGitHub(os.Getenv("GITHUB_SECRET"))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error creating github webhook parser: %s\n", err)
		return
	}

	concourse := cws.Concourse{
		Endpoint:   os.Getenv("CONCOURSE_ENDPOINT"),
		Username:   os.Getenv("CONCOURSE_USERNAME"),
		Password:   os.Getenv("CONCOURSE_PASSWORD"),
		DBHost:     os.Getenv("CONCOURSE_DB_HOST"),
		DBUsername: os.Getenv("CONCOURSE_DB_USERNAME"),
		DBPassword: os.Getenv("CONCOURSE_DB_PASSWORD"),
	}

	gitHubEventProcessor := cws.Processor{
		GetResources:  concourse.GetGitResources,
		CheckResource: concourse.CheckResource,
	}

	server := cws.Server{
		Port:          port,
		GitHubHandler: gitHub.Handler(gitHubEventProcessor.Process),
	}

	fmt.Println("Running at http://0.0.0.0:" + server.Port)
	err = server.Start()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error serving: %s\n", err)
		return
	}

}
