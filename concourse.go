package cws

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //blank import
	"golang.org/x/oauth2"
)

type Resource struct {
	Name     string
	Team     string
	Pipeline string
}

type Concourse struct {
	Endpoint    string
	Username    string
	Password    string
	DBHost      string
	DBUsername  string
	DBPassword  string
	accessToken string
}

func (c Concourse) GetGitResources(repo GitRepo) ([]Resource, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/concourse?sslmode=disable", c.DBUsername, c.DBPassword, c.DBHost)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var resources []Resource
	query := `
		SELECT r.name, t.name team, p.name pipeline
		FROM resources r
		INNER JOIN pipelines p ON p.id = r.pipeline_id
		INNER JOIN teams t ON t.id = p.team_id
		WHERE r.active = TRUE
			AND p.paused = FALSE
			AND config::json->>'type' = 'git'
			AND config::json->'source'->>'uri' IN ($1, $2)
			AND config::json->'source'->>'branch' = $3`

	err = db.Select(&resources, query, repo.URI(), repo.URI()+".git", repo.Branch)
	return resources, err
}

func (c Concourse) CheckResource(resource Resource) error {
	token, err := c.token()
	if err != nil {
		return fmt.Errorf("error getting access token: %s", err)
	}

	checkURL := fmt.Sprintf("%s/api/v1/teams/%s/pipelines/%s/resources/%s/check", c.Endpoint, resource.Team, resource.Pipeline, resource.Name)

	var httpClient = &http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("POST", checkURL, strings.NewReader(`{"from":null}`))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error calling concourse api: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 201 {
		return fmt.Errorf("error calling concourse api: %v %s", res.StatusCode, res.Status)
	}

	return nil
}

func (c Concourse) token() (string, error) {
	if c.accessToken != "" {
		return c.accessToken, nil
	}

	oauth2Config := oauth2.Config{
		ClientID:     "fly",
		ClientSecret: "Zmx5",
		Endpoint:     oauth2.Endpoint{TokenURL: c.Endpoint + "/sky/issuer/token"},
		Scopes:       []string{"openid", "profile", "email", "federated:id", "groups"},
	}

	token, err := oauth2Config.PasswordCredentialsToken(context.Background(), c.Username, c.Password)
	if err != nil {
		c.accessToken = ""
		return "", err
	}

	c.accessToken = token.AccessToken
	return c.accessToken, nil
}
