package cws

import (
	"fmt"
)

type GitProcessor = func(repo GitRepo) ([]Resource, []error)
type ResourceGetter = func(repo GitRepo) ([]Resource, error)
type ResourceChecker = func(resource Resource) error

type Processor struct {
	GetResources  ResourceGetter
	CheckResource ResourceChecker
}

func (p Processor) Process(repo GitRepo) ([]Resource, []error) {
	resources, err := p.GetResources(repo)
	if err != nil {
		return nil, []error{err}
	}

	var checkedResources []Resource
	var checkErrors []error

	for _, resource := range resources {
		err := p.CheckResource(resource)
		if err != nil {
			checkErrors = append(checkErrors, fmt.Errorf("error checking resourse: %+v  error: %s", resource, err))
		} else {
			checkedResources = append(checkedResources, resource)
		}
	}

	return checkedResources, checkErrors
}
