//+build gomock

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

// GoMockTestSuite is the structure/interface passed to server specific tests
type GoMockTestSuite struct {
	suite.Suite
	//Shared state goes here
	ctrl   *gomock.Controller
	store  *MockPlayerStore
	server *PlayerServer
}

// SetupTest is called prior to each test
func (suite *GoMockTestSuite) SetupTest() {
	//Before each test, reset this data to a clean state; useful to ensure DB/mock data is in known state prior to each test, or what have you
	suite.ctrl = gomock.NewController(suite.T())
	suite.store = NewMockPlayerStore(suite.ctrl)
	suite.server = NewPlayerServer(suite.store)
}

func (suite *GoMockTestSuite) TearDownTest() {
	//Runs after test, performing whatever teardown is necessary
	suite.ctrl.Finish()
}

//Test Section: Below is all the tests run by the suite, which you can tell because it they're all attached to the suite struct we defined above
func (suite *GoMockTestSuite) TestGetScoreSuite() {
	suite.store.EXPECT().GetPlayerScore("Pepper").Return(20)

	request := newGetScoreRequest("Pepper")
	response := httptest.NewRecorder()

	suite.server.ServeHTTP(response, request)

	suite.Equal(response.Code, http.StatusOK)
	suite.Equal(response.Body.String(), "20")
}

// Run just this one if you want the whole suite to run through
func TestGoMockTestSuite(t *testing.T) {
	//This is a wrapper for go test -v or whatever you run with, since this method conforms to the standard library approach
	suite.Run(t, new(GoMockTestSuite))
}

// func TestGetLeague(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	store := NewMockPlayerStore(ctrl)

// 	server := NewPlayerServer(store)

// 	wantedLeague := []Player{
// 		{"Pepper", 20},
// 		{"Floyd", 10},
// 	}

// 	request := newLeagueRequest()
// 	response := httptest.NewRecorder()
// 	store.EXPECT().GetLeague().Return(wantedLeague)
// 	server.ServeHTTP(response, request)

// 	assertContentType(t, response, jsonContentType)

// 	got := getLeagueFromResponse(t, response.Body)
// 	assertLeague(t, got, wantedLeague)
// 	assertStatus(t, response.Code, http.StatusOK)
// }

// func TestMockStoreWins(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	store := NewMockPlayerStore(ctrl)

// 	server := NewPlayerServer(store)

// 	request := newPostWinRequest("Pepper")
// 	response := httptest.NewRecorder()
// 	store.EXPECT().RecordWin("Pepper")
// 	store.EXPECT().RecordWin("Floyd")

// 	server.ServeHTTP(response, request)

// 	request = newPostWinRequest("Floyd")

// 	server.ServeHTTP(response, request)

// 	assertStatus(t, response.Code, http.StatusAccepted)

// }

// func TestGetScore(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	store := NewMockPlayerStore(ctrl)

// 	server := NewPlayerServer(store)
// 	//playerName := "Pepper"

// 	store.EXPECT().GetPlayerScore("Pepper").Return(20)

// 	request := newGetScoreRequest("Pepper")
// 	response := httptest.NewRecorder()

// 	server.ServeHTTP(response, request)

// 	assertStatus(t, response.Code, http.StatusOK)
// 	assertResponseBody(t, response.Body.String(), "20")
// }
