package service_helpers

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/debu99/cicd-runner/helpers/service/mocks"
)

var errExample = errors.New("example error")

func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mi := &mocks.Interface{}
	s := &SimpleService{i: mi}

	mi.On("Start", s).Return(errExample)

	err := s.Run()
	assert.Equal(t, errExample, err)
	mi.AssertExpectations(t)
}
