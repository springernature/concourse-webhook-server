package cws

type Resource struct {
	Name     string
	Team     string
	Pipeline string
}

type Concourse struct {
}

func (c Concourse) GetGitResources(gitBranch GitRepo) ([]Resource, error) {
	return []Resource{
		{Name: "some-resource", Pipeline: "some-pipeline", Team: "some-team"},
	}, nil
}

func (c Concourse) CheckResource(resource Resource) error {
	return nil
}
