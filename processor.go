package cws

import (
	"fmt"
	"time"
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
	type Result struct {
		Res Resource
		Err error
	}

	results := make(chan Result, len(resources))

	for _, resource := range resources {
		go func(r Resource) {
			results <- Result{Res: r, Err: p.CheckResource(r)}
		}(resource)
	}

	for range resources {
		select {
		case result := <-results:
			if result.Err != nil {
				checkErrors = append(checkErrors, fmt.Errorf("error checking resourse: %+v %s", result.Res, result.Err))
			} else {
				checkedResources = append(checkedResources, result.Res)
			}
		case <-time.After(5 * time.Second):
			//timeout
		}
	}

	return checkedResources, checkErrors
}
