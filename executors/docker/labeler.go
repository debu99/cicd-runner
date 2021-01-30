package docker

import "github.com/debu99/cicd-runner/executors/docker/internal/labels"

func (e *executor) createLabeler() error {
	e.labeler = labels.NewLabeler(e.Build)
	return nil
}
