package user

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/debu99/cicd-runner/executors/docker/internal/exec"
	"github.com/debu99/cicd-runner/helpers/docker"
)

func TestDefaultInspect_IsRoot(t *testing.T) {
	containerID := "container-id"

	tests := map[string]struct {
		setupDockerClientMock func(t *testing.T, clientMock *docker.MockClient, expectedCtx context.Context)
		expectedIsRoot        bool
		expectedError         error
	}{
		"ImageInspectWithRaw error": {
			setupDockerClientMock: func(t *testing.T, clientMock *docker.MockClient, expectedCtx context.Context) {
				clientMock.On("ImageInspectWithRaw", expectedCtx, containerID).
					Return(types.ImageInspect{}, nil, assert.AnError).
					Once()
			},
			expectedIsRoot: true,
			expectedError:  assert.AnError,
		},
		"empty Config": {
			setupDockerClientMock: func(t *testing.T, clientMock *docker.MockClient, expectedCtx context.Context) {
				clientMock.On("ImageInspectWithRaw", expectedCtx, containerID).
					Return(types.ImageInspect{}, nil, nil).
					Once()
			},
			expectedIsRoot: true,
			expectedError:  nil,
		},
		"empty user entry in Config": {
			setupDockerClientMock: func(t *testing.T, clientMock *docker.MockClient, expectedCtx context.Context) {
				clientMock.On("ImageInspectWithRaw", expectedCtx, containerID).
					Return(types.ImageInspect{Config: &container.Config{User: ""}}, nil, nil).
					Once()
			},
			expectedIsRoot: true,
			expectedError:  nil,
		},
		"user entry in Config set to root": {
			setupDockerClientMock: func(t *testing.T, clientMock *docker.MockClient, expectedCtx context.Context) {
				clientMock.On("ImageInspectWithRaw", expectedCtx, containerID).
					Return(types.ImageInspect{Config: &container.Config{User: "root"}}, nil, nil).
					Once()
			},
			expectedIsRoot: true,
			expectedError:  nil,
		},
		"user entry in Config set to non-root": {
			setupDockerClientMock: func(t *testing.T, clientMock *docker.MockClient, expectedCtx context.Context) {
				clientMock.On("ImageInspectWithRaw", expectedCtx, containerID).
					Return(types.ImageInspect{Config: &container.Config{User: "non-root"}}, nil, nil).
					Once()
			},
			expectedIsRoot: false,
			expectedError:  nil,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			clientMock := new(docker.MockClient)
			defer clientMock.AssertExpectations(t)

			execMock := new(exec.MockDocker)
			defer execMock.AssertExpectations(t)

			ctx := context.Background()

			tt.setupDockerClientMock(t, clientMock, ctx)

			inspect := NewInspect(clientMock, execMock)
			isRoot, err := inspect.IsRoot(ctx, containerID)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedIsRoot, isRoot, "user root-status mismatch")
		})
	}
}

type uidAndGidTestCase struct {
	assertExecMock func(t *testing.T, clientMock *exec.MockDocker, expectedCtx context.Context)
	expectedID     int
	assertError    func(t *testing.T, err error)
}

func TestDefaultInspect_UID(t *testing.T) {
	testDefaultInspectUIDandGID(
		t,
		commandIDU,
		func(inspect Inspect, ctx context.Context, containerID string) (int, error) {
			return inspect.UID(ctx, containerID)
		},
	)
}

func TestDefaultInspect_GID(t *testing.T) {
	testDefaultInspectUIDandGID(
		t,
		commandIDG,
		func(inspect Inspect, ctx context.Context, containerID string) (int, error) {
			return inspect.GID(ctx, containerID)
		},
	)
}

func testDefaultInspectUIDandGID(
	t *testing.T,
	expectedCommand string,
	testCall func(inspect Inspect, ctx context.Context, containerID string) (int, error),
) {
	containerID := "container-id"

	assertCommand := func(t *testing.T, args mock.Arguments) {
		input, ok := args.Get(2).(io.Reader)
		require.True(t, ok)

		data, err := ioutil.ReadAll(input)
		require.NoError(t, err)

		assert.Equal(t, expectedCommand, string(data))
	}
	mockOutput := func(t *testing.T, args mock.Arguments, data string) {
		output, ok := args.Get(3).(io.Writer)
		require.True(t, ok)

		_, err := fmt.Fprintln(output, data)
		require.NoError(t, err)
	}

	tests := map[string]uidAndGidTestCase{
		"Exec error": {
			assertExecMock: func(t *testing.T, clientMock *exec.MockDocker, expectedCtx context.Context) {
				clientMock.On("Exec", expectedCtx, containerID, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						assertCommand(t, args)
					}).
					Return(assert.AnError).
					Once()
			},
			expectedID: 0,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"ID parsing error": {
			assertExecMock: func(t *testing.T, clientMock *exec.MockDocker, expectedCtx context.Context) {
				clientMock.On("Exec", expectedCtx, containerID, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						assertCommand(t, args)
						mockOutput(t, args, "\n\ntest\n\n")
					}).
					Return(nil).
					Once()
			},
			expectedID: 0,
			assertError: func(t *testing.T, err error) {
				var e *strconv.NumError
				assert.ErrorAs(t, err, &e)
			},
		},
		"proper ID received from output": {
			assertExecMock: func(t *testing.T, clientMock *exec.MockDocker, expectedCtx context.Context) {
				clientMock.On("Exec", expectedCtx, containerID, mock.Anything, mock.Anything).
					Run(func(args mock.Arguments) {
						assertCommand(t, args)
						mockOutput(t, args, "\n\n123\n\n")
					}).
					Return(nil).
					Once()
			},
			expectedID:  123,
			assertError: nil,
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			clientMock := new(docker.MockClient)
			defer clientMock.AssertExpectations(t)

			execMock := new(exec.MockDocker)
			defer execMock.AssertExpectations(t)

			ctx := context.Background()

			tt.assertExecMock(t, execMock, ctx)

			inspect := NewInspect(clientMock, execMock)
			id, err := testCall(inspect, ctx, containerID)

			assert.Equal(t, tt.expectedID, id)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
