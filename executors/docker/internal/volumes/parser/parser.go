package parser

import "github.com/debu99/cicd-runner/helpers/path"

type Parser interface {
	ParseVolume(spec string) (*Volume, error)
	Path() path.Path
}
