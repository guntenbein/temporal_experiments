package workflow_error

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
)

type MoveByProductsWorkflowTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
}

func TestMoveByProductsWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(MoveByProductsWorkflowTestSuite))
}

func (s *MoveByProductsWorkflowTestSuite) Test_WorkflowError() {
	env := s.NewTestWorkflowEnvironment()
	activityHandler := ActivityHandler{}

	// 1. why cannot we do it like this?
	// we have "test panicked: activity "Activity1" is not registered with the TestWorkflowEnvironment"
	// env.OnActivity("Activity1", mock.Anything).Return(BusinessError{}).Once()
	env.OnActivity(activityHandler.Activity1, mock.Anything).Return(BusinessError{}).Once()
	// 2. how to ensure that the Activity2 was not called?
	// env.OnActivity(activityHandler.Activity2, mock.Anything).Return(nil).Once()

	env.ExecuteWorkflow(ExperimentalWorkflow)

	s.True(env.IsWorkflowCompleted())
	s.Error(env.GetWorkflowError())

	env.AssertExpectations(s.T())
}

// dummy declaration
type ActivityHandler struct{}

func (ah ActivityHandler) Activity1(_ context.Context) error {
	return nil
}

func (ah ActivityHandler) Activity2(_ context.Context) error {
	return nil
}
