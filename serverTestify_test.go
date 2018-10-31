//+build testify

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

// ServerTestSuite is the structure/interface passed to server specific tests
type ServerTestSuite struct {
	suite.Suite
	//Shared state
	store    StubPlayerStore //Could be interface, is not
	server   *PlayerServer
	testName string
}

// SetupSuite is called once per suite run
func (suite *ServerTestSuite) SetupSuite() {
	//Setup shared state that won't be altered in tests
	suite.testName = "Pepper"
}

// SetupTest is called prior to each test
func (suite *ServerTestSuite) SetupTest() {
	//Before each test, reset this data to a clean state; useful to ensure DB/mock data is in known state prior to each test
	suite.store = StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
		[]string{},
		[]Player{
			{"Pepper", 20},
			{"Floyd", 10},
		},
	}
	suite.server = NewPlayerServer(&suite.store)
}

//Test Section: Below is all the tests run by the suite, which you can tell because it they're all attached to the suite struct we defined above

func (suite *ServerTestSuite) TestifyGetPepperScore() {
	request := newGetScoreRequest(suite.testName)
	response := httptest.NewRecorder()

	suite.server.ServeHTTP(response, request)

	suite.Equal(response.Code, http.StatusOK)
	suite.Equal(response.Body.String(), "20")
}

func (suite *ServerTestSuite) TestifyGetFloydScore() {
	request := newGetScoreRequest("Floyd")
	response := httptest.NewRecorder()

	suite.server.ServeHTTP(response, request)

	suite.Equal(response.Code, http.StatusOK)
	suite.Equal(response.Body.String(), "10")
}

func (suite *ServerTestSuite) TestifyGetNobodyScore() {
	request := newGetScoreRequest("Apollo")
	response := httptest.NewRecorder()
	suite.server.ServeHTTP(response, request)

	suite.Equal(response.Code, http.StatusNotFound)
}

func (suite *ServerTestSuite) TestifyStoreWins() {
	request := newPostWinRequest(suite.testName)
	response := httptest.NewRecorder()

	suite.server.ServeHTTP(response, request)

	suite.Equal(response.Code, http.StatusAccepted)
	//TODO: Find more expressive way of testing what we called store with
	suite.Equal(len(suite.store.winCalls), 1, "Unexpected call count to suite.store!")
	suite.Equal(suite.store.winCalls[0], suite.testName, "Unexpected params in call to suite.store!")
}

func (suite *ServerTestSuite) TestifyGetLeague() {
	wantedLeague := []Player{
		{"Pepper", 20},
		{"Floyd", 10},
	}

	var got []Player

	request := newLeagueRequest()
	response := httptest.NewRecorder()

	suite.server.ServeHTTP(response, request)

	suite.Equal(response.Header().Get("content-type"), jsonContentType)

	err := json.NewDecoder(response.Body).Decode(&got)

	suite.Nil(err, "Failed to marshall json for league response!")
	suite.Equal(wantedLeague, got, "League JSON didn't match expectations!")
	suite.Equal(response.Code, http.StatusOK)
}

// Run just this one if you want the whole suite to run through
func TestServerTestSuite(t *testing.T) {
	//This is a wrapper for go test -v or whatever you run with, since this method conforms to the standard library approach
	suite.Run(t, new(ServerTestSuite))
}
