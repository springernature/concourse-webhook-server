package cws_test

import (
	"cws"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EventProcessorChecksAllResources(t *testing.T) {
	repo := cws.GitRepo{
		Name:   "some-repo",
		Branch: "some-branch",
	}

	resources := []cws.Resource{
		{Name: "r1", Team: "t1", Pipeline: "p1"},
		{Name: "r2", Team: "t2", Pipeline: "p2"},
	}

	getResourcesRepo := cws.GitRepo{}
	getResources := func(repo cws.GitRepo) ([]cws.Resource, error) {
		getResourcesRepo = repo
		return resources, nil
	}

	checkResource := func(resource cws.Resource) error {
		return nil
	}

	ep := cws.Processor{GetResources: getResources, CheckResource: checkResource}

	checkedResources, errs := ep.Process(repo)

	assert.Empty(t, errs)
	assert.Equal(t, repo, getResourcesRepo, "repo was not passed to resource getter")
	assert.Len(t, checkedResources, len(resources), "resources were not all passed to resource checker")
	assert.Contains(t, checkedResources, resources[0])
	assert.Contains(t, checkedResources, resources[1])
}

func Test_EventProcessorContinuesWhenACheckErrors(t *testing.T) {
	repo := cws.GitRepo{
		Name:   "some-repo",
		Branch: "some-branch",
	}

	resources := []cws.Resource{
		{Name: "r1"},
		{Name: "r2Error"},
		{Name: "r3"},
		{Name: "r4Error"},
		{Name: "r5"},
	}

	getResources := func(repo cws.GitRepo) ([]cws.Resource, error) {
		return resources, nil
	}

	checkResource := func(resource cws.Resource) error {
		if strings.Contains(resource.Name, "Error") {
			return fmt.Errorf("test error %s", resource.Name)
		}
		return nil
	}

	ep := cws.Processor{GetResources: getResources, CheckResource: checkResource}

	checkedResources, checkResourceErrors := ep.Process(repo)

	assert.Len(t, checkedResources, 3, "expected 3 resources to have checked successfully")
	assert.Len(t, checkResourceErrors, 2, "expected 2 resources to have errored")

}
