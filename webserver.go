package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/go-playground/webhooks.v5/github"
)

type Server struct {
	Port         string
	GithubSecret string
}

func NewServer() Server {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Server{
		Port:         port,
		GithubSecret: os.Getenv("GITHUB_SECRET"),
	}
}

func (s Server) Start() error {
	hook, err := github.New(github.Options.Secret(s.GithubSecret))
	if err != nil {
		return err
	}

	githubHandler := func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.PushEvent, github.PingEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				// ok event wasn't one of the ones asked to be parsed
			} else {
				fmt.Printf("error parsing: %s\n", err)
			}
		}

		switch p := payload.(type) {
		case github.PushPayload:
			branch := branchFromRef(p.Ref)
			if branch != "" {
				repo := p.Repository.FullName
				fmt.Printf("push: repo=%s branch=%s\n", repo, branch)
			}
		case github.PingPayload:
			fmt.Printf("ping: id=%d\n", p.HookID)
		}
	}

	http.HandleFunc("/github", githubHandler)

	fmt.Println("Running on port " + s.Port)
	return http.ListenAndServe(":"+s.Port, nil)

}

func branchFromRef(ref string) string {
	branchRefPrefix := "refs/heads/"
	if strings.HasPrefix(ref, branchRefPrefix) {
		return strings.TrimPrefix(ref, branchRefPrefix)
	} else {
		return ""
	}
}

func main() {
	server := NewServer()
	err := server.Start()
	if err != nil {
		fmt.Printf("error serving: %s\n", err)
	}
}
