//+build example

package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ExampleTestSuite is the structure/interface passed to server specific tests
type ExampleTestSuite struct {
	suite.Suite
	//Shared state goes here
}

// SetupSuite is called once per suite run
func (suite *ExampleTestSuite) SetupSuite() {
	//Setup shared state that won't be altered in tests here
}

// SetupTest is called prior to each test
func (suite *ExampleTestSuite) SetupTest() {
	//Before each test, reset this data to a clean state; useful to ensure DB/mock data is in known state prior to each test, or what have you
}

func (suite *ExampleTestSuite) TearDownTest() {
	//Runs after test, performing whatever teardown is necessary
}

func (suite *ExampleTestSuite) TearDownSuite() {
	//Runs after suite, performing whatever teardown is necessary
}

//Test Section: Below is all the tests run by the suite, which you can tell because it they're all attached to the suite struct we defined above

func (suite *ExampleTestSuite) TestFirstTestSuiteWillRun() {
	suite.Equal(1, 1)
}

func (suite *ExampleTestSuite) TestSecondTestSuiteWillRun() {
	suite.Equal(2, 2)
}

// Run just this one if you want the whole suite to run through
func TestExampleTestSuite(t *testing.T) {
	//This is a wrapper for go test -v or whatever you run with, since this method conforms to the standard library approach
	suite.Run(t, new(ExampleTestSuite))
}
