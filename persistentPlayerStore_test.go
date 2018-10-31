//+build standard

package main

import (
	"os/exec"
	"reflect"
	"testing"
)

func TestPersistence(t *testing.T) {
	store := PostgresPlayerStore{}

	t.Run("it connects to the database", func(t *testing.T) {
		err := store.TestConnect("dev")
		if err != nil {
			t.Errorf("Expected to connect without an error! Got: %s", err.Error())
		}
	})

	t.Run("it returns the score for specified user", func(t *testing.T) {
		databaseReset(t)

		testScore := 20
		name := "Pepper"
		got := store.GetPlayerScore(name)
		expected := testScore

		if got != expected {
			t.Errorf("Unexpected score return for %s! Got %d, expected %d", name, got, expected)
		}
	})

	t.Run("it updates the score for specified user (when exists)", func(t *testing.T) {
		databaseReset(t)

		name := "Pepper"
		expected := store.GetPlayerScore(name) + 1
		store.RecordWin(name)

		returned := store.GetPlayerScore(name)

		if returned != expected {
			t.Errorf("Unexpected score return for %s! Got %d, expected %d", name, returned, expected)
		}
	})

	t.Run("it updates the score for specified user (when not exists)", func(t *testing.T) {
		databaseReset(t)

		name := "Jay"
		expected := 1
		store.RecordWin(name)

		returned := store.GetPlayerScore(name)

		if returned != expected {
			t.Errorf("Unexpected score return for %s! Got %d, expected %d", name, returned, expected)
		}
	})

	t.Run("it returns all scores in the league", func(t *testing.T) {
		databaseReset(t)

		expected := []Player{{Name: "Pepper", Wins: 20}}
		got := store.GetLeague()

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Unexpected league return! Got %v, expected %v", got, expected)
		}
	})
}

func databaseReset(t *testing.T) {
	err := migrateDown()
	if err != nil {
		t.Errorf("Got error migrating down! %s\n", err.Error())
	}
	err = migrateUp()
	if err != nil {
		t.Errorf("Got error migrating up! %s\n", err.Error())
	}
}

func migrateUp() error {
	cmd := exec.Command("migrate", "-source", "file://./migrations", "-database", "postgres://localhost:5432/gwt_dev?sslmode=disable", "up")
	return cmd.Run()
}

func migrateDown() error {
	cmd := exec.Command("migrate", "-source", "file://./migrations", "-database", "postgres://localhost:5432/gwt_dev?sslmode=disable", "down")
	return cmd.Run()

}
