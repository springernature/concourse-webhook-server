package cws_test

import (
	"cws"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/go-playground/webhooks.v5/github"
)

func testProcessor(repo cws.GitRepo) ([]cws.Resource, []error) {
	return nil, nil
}

func testGithubParser(pushPayload github.PushPayload, err error) cws.GitHubWebHookParser {
	return func(r *http.Request, events ...github.Event) (interface{}, error) {
		return pushPayload, err
	}
}

func Test_ParserError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader("body"))
	res := httptest.NewRecorder()
	parserError := fmt.Errorf("some parser error")
	gitHub := cws.GitHub{
		Parser: testGithubParser(github.PushPayload{}, parserError),
	}

	gitHub.Handler(testProcessor)(res, req)

	assert.Equal(t, 500, res.Code)
	assert.Contains(t, res.Body.String(), parserError.Error())
}

func Test_RefError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader("body"))
	res := httptest.NewRecorder()
	payload := github.PushPayload{}
	payload.Ref = "some/unknown/ref"
	gitHub := cws.GitHub{
		Parser: testGithubParser(payload, nil),
	}

	gitHub.Handler(testProcessor)(res, req)

	assert.Equal(t, 500, res.Code)
	assert.Contains(t, res.Body.String(), payload.Ref)
}

func Test_CallsProcessor(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader("foo"))
	res := httptest.NewRecorder()

	expectedRepo := cws.GitRepo{
		Name:   "some-repo",
		Branch: "some-branch",
	}
	payload := github.PushPayload{}
	payload.Repository.FullName = expectedRepo.Name
	payload.Ref = "refs/heads/" + expectedRepo.Branch

	gitHub := cws.GitHub{
		Parser: testGithubParser(payload, nil),
	}

	actualRepo := cws.GitRepo{}
	processor := func(repo cws.GitRepo) ([]cws.Resource, []error) {
		actualRepo = repo
		return nil, nil
	}

	gitHub.Handler(processor)(res, req)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, expectedRepo, actualRepo)
}

func Test_ProcessorError(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "", strings.NewReader("body"))
	res := httptest.NewRecorder()
	payload := github.PushPayload{}
	payload.Ref = "refs/heads/some-branch"
	gitHub := cws.GitHub{
		Parser: testGithubParser(payload, nil),
	}

	processor := func(repo cws.GitRepo) ([]cws.Resource, []error) {
		return nil, []error{fmt.Errorf("some error")}
	}

	gitHub.Handler(processor)(res, req)

	assert.Equal(t, 500, res.Code)
	assert.Contains(t, res.Body.String(), "error checking resources")
}
