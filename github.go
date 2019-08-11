package cws

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"gopkg.in/go-playground/webhooks.v5/github"
)

type GitRepo struct {
	Name   string
	Branch string
}

func (r GitRepo) URI() string {
	return fmt.Sprintf("git@github.com:%s", r.Name)
}

type GitHubWebHookParser = func(r *http.Request, events ...github.Event) (interface{}, error)

type GitHub struct {
	Parser GitHubWebHookParser
}

func NewGitHub(secret string) (GitHub, error) {
	hook, err := github.New(github.Options.Secret(secret))
	if err != nil {
		return GitHub{}, fmt.Errorf("error creating github hook parser, check secret is set")
	}
	return GitHub{hook.Parse}, nil
}

func (gh GitHub) Handler(process GitProcessor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := gh.Parser(r, github.PushEvent, github.PingEvent)
		if err != nil && err != github.ErrEventNotFound {
			gh.httpError(w, fmt.Errorf("error parsing event: %s", err))
			return
		}

		switch p := payload.(type) {
		case github.PingPayload:
			fmt.Fprintf(w, string(p.HookID))
		case github.PushPayload:
			repo, err := gh.branchFromPayload(p)
			if err != nil {
				gh.httpError(w, err)
				return
			}
			checkedResources, errors := process(repo)

			fmt.Fprintf(os.Stdout, "processed event for repo %+v\n", repo)
			for _, resource := range checkedResources {
				fmt.Fprintf(w, "%+v\n", resource)
				fmt.Fprintf(os.Stdout, "checked resource %+v\n", resource)
			}
			for _, err := range errors {
				fmt.Fprintln(os.Stderr, err)
			}

			if len(errors) > 0 {
				gh.httpError(w, fmt.Errorf("error checking resources"))
				return
			}
		}
	}
}

func (gh GitHub) branchFromPayload(p github.PushPayload) (GitRepo, error) {
	branchRefPrefix := "refs/heads/"
	if strings.HasPrefix(p.Ref, branchRefPrefix) {
		return GitRepo{
			Name:   p.Repository.FullName,
			Branch: strings.TrimPrefix(p.Ref, branchRefPrefix),
		}, nil
	}
	return GitRepo{}, fmt.Errorf("could not get branch from event ref: %s", p.Ref)
}

func (gh GitHub) httpError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	fmt.Fprintln(w, err)
}
