package cws

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //blank import
)

type Resource struct {
	Name     string
	Team     string
	Pipeline string
}

type Concourse struct {
	DbHost     string
	DbUsername string
	DbPassword string
}

func (c Concourse) GetGitResources(repo GitRepo) ([]Resource, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/concourse?sslmode=disable", c.DbUsername, c.DbPassword, c.DbHost)
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
			AND r.paused = FALSE
			AND p.paused = FALSE
			AND config::json->>'type' = 'git'
			AND config::json->'source'->>'uri' IN ($1, $2)
			AND config::json->'source'->>'branch' = $3`

	err = db.Select(&resources, query, repo.URI(), repo.URI()+".git", repo.Branch)
	return resources, err
}

func (c Concourse) CheckResource(resource Resource) error {
	return nil
}
