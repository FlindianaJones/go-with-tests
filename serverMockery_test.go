//+build mockery

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

// MockeryTestSuite is the structure/interface passed to server specific tests
type MockeryTestSuite struct {
	suite.Suite
	//Shared state goes here
	store  *MockeryPlayerStore
	server *PlayerServer
}

// SetupTest is called prior to each test
func (suite *MockeryTestSuite) SetupTest() {
	//Before each test, reset this data to a clean state; useful to ensure DB/mock data is in known state prior to each test, or what have you
	suite.store = &MockeryPlayerStore{}
	suite.server = NewPlayerServer(suite.store)
}

//Test Section: Below is all the tests run by the suite, which you can tell because it they're all attached to the suite struct we defined above

func (suite *MockeryTestSuite) TestGetScoreMockery() {
	request := newGetScoreRequest("Pepper")
	response := httptest.NewRecorder()

	suite.store.On("GetPlayerScore", "Pepper").Return(20)

	suite.server.ServeHTTP(response, request)

	suite.Equal(response.Code, http.StatusOK)
	suite.Equal(response.Body.String(), "20")
	suite.store.AssertCalled(suite.T(), "GetPlayerScore", "Pepper")
	suite.store.AssertNumberOfCalls(suite.T(), "GetPlayerScore", 1)
}

// Run just this one if you want the whole suite to run through
func TestMockeryTestSuite(t *testing.T) {
	//This is a wrapper for go test -v or whatever you run with, since this method conforms to the standard library approach
	suite.Run(t, new(MockeryTestSuite))
}
