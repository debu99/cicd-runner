// +build !windows

package custom

import (
	"errors"

	terminalsession "github.com/debu99/cicd-runner/session/terminal"
)

func (e *executor) Connect() (terminalsession.Conn, error) {
	return nil, errors.New("not yet supported")
}
