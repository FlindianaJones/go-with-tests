package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type PostgresPlayerStore struct {
}

func (store *PostgresPlayerStore) GetPlayerScore(name string) int {

	dbname := "gwt_dev"

	psqlInfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)

	db, err := sql.Open("postgres", psqlInfo)

	defer db.Close()

	if err != nil {
		fmt.Printf("Got error when connecting to db! %s\n", err)
		return 0
	}

	rows, err := db.Query(`SELECT score FROM player_scores WHERE name = $1`, name)
	if err != nil {
		fmt.Printf("Got error when querying to db! %s\n", err)
		return 0
	}

	var score int

	if rows.Next() {
		if err := rows.Scan(&score); err != nil {
			fmt.Printf("Got error scanning score! %s", err.Error())
		}
	}

	return score
}

// RecordWin creates a new user if name passed in doesn't match; otherwise increments count
func (store *PostgresPlayerStore) RecordWin(name string) {

	dbname := "gwt_dev"

	psqlInfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)

	db, err := sql.Open("postgres", psqlInfo)

	defer db.Close()

	if err != nil {
		fmt.Printf("Got error when connecting to db! %s\n", err)
		return
	}
	rows, err := db.Query(`SELECT score FROM player_scores WHERE name = $1`, name)
	if err != nil {
		fmt.Printf("Got error when querying to db! %s\n", err)
		return
	}
	var currentScore int

	if rows.Next() {
		if err := rows.Scan(&currentScore); err != nil {
			fmt.Printf("Got error scanning score! %s", err.Error())
		}
		sqlStatement := `UPDATE player_scores SET score = $1 WHERE name = $2`
		_, err = db.Exec(sqlStatement, currentScore+1, name)
		if err != nil {
			fmt.Printf("Got error when updating db! %s\n", err)
			return
		}
	} else {
		sqlStatement := `INSERT INTO player_scores VALUES ($1, 1)`
		_, err = db.Exec(sqlStatement, name)
		if err != nil {
			fmt.Printf("Got error when inserting into db! %s\n", err)
			return
		}
	}

}

// GetLeague returns all players/scores into a slice of Players
func (store *PostgresPlayerStore) GetLeague() []Player {
	dbname := "gwt_dev"

	psqlInfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)

	db, err := sql.Open("postgres", psqlInfo)

	defer db.Close()

	if err != nil {
		fmt.Printf("Got error when connecting to db! %s\n", err)
		return nil
	}
	rows, err := db.Query(`SELECT name, score FROM player_scores`)
	if err != nil {
		fmt.Printf("Got error when querying the db! %s\n", err)
		return nil
	}

	players := []Player{}

	for rows.Next() {
		var pName = ""
		var pScore = 0
		if err := rows.Scan(&pName, &pScore); err != nil {
			fmt.Printf("Got error scanning score/name! %s", err.Error())
			break
		}
		players = append(players, Player{Name: pName, Wins: pScore})
	}

	return players
}

// TestConnect ensures we have the connection information correct
func (store *PostgresPlayerStore) TestConnect(environment string) error {
	dbname := "gwt_dev"

	psqlInfo := fmt.Sprintf("dbname=%s sslmode=disable", dbname)

	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()

	err = db.Ping()
	return err
}
