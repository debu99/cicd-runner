package machine

import (
	"errors"

	"github.com/debu99/cicd-runner/session/terminal"
)

func (e *machineExecutor) Connect() (terminal.Conn, error) {
	if term, ok := e.executor.(terminal.InteractiveTerminal); ok {
		return term.Connect()
	}

	return nil, errors.New("executor does not have terminal")
}
