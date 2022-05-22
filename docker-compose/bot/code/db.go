package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"strings"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

//Collecting data from bot
func collectData(username string, chatid int64, message string, answer []string) error {

	//Connecting to database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Converting slice with answer to string
	answ := strings.Join(answer, ", ")

	//Creating SQL command
	data := `INSERT INTO users(username, chat_id, message, answer) VALUES($1, $2, $3, $4);`

	//Execute SQL command in database
	if _, err = db.Exec(data, `@`+username, chatid, message, answ); err != nil {
		return err
	}

	return nil
}

//Creating users table in database
func createTable() error {

	//Connecting to database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Creating users Table
	if _, err = db.Exec(`CREATE TABLE users(ID SERIAL PRIMARY KEY, TIMESTAMP TIMESTAMP DEFAULT CURRENT_TIMESTAMP, USERNAME TEXT, CHAT_ID INT, MESSAGE TEXT, ANSWER TEXT);`); err != nil {
		return err
	}
	fmt.Println("table users created")

	//Creating answer type Table
	// ANSWER_TYPE: 1 - wiki, 2 - jokes
	if _, err = db.Exec(`CREATE TABLE answer_type(ID SERIAL PRIMARY KEY, CHAT_ID INT UNIQUE, ANSWER_TYPE INT);`); err != nil {
		return err
	}
	fmt.Println("table answer_type created")

	return nil
}

//Getting number of users who using bot
func getNumberOfUsers() (int64, error) {

	var count int64

	//Connecting to database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	//Counting number of users
	row := db.QueryRow("SELECT COUNT(DISTINCT username) FROM users;")
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// по умолчанию тип ответа при ошибке - wiki
func getAnswerType(chatId int64) (int, error) {
	var answerType int = wikiAnswer

	//Connecting to database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return answerType, err
	}
	defer db.Close()

	////Counting number of users
	//ID SERIAL PRIMARY KEY, CHAT_ID INT,  INT)
	//data := `INSERT INTO users(username, chat_id, message, answer) VALUES($1, $2, $3, $4);`

	////Execute SQL command in database
	//if _, err = db.Exec(data, `@`+username, chatid, message, answ); err != nil {
	//	return answerType, err
	//}
	//

	row := db.QueryRow("SELECT ANSWER_TYPE FROM answer_type WHERE CHAT_ID = $1;", chatId)
	err = row.Scan(&answerType)
	if err != nil {
		return answerType, nil
	}

	fmt.Println("get answer type: ", answerType)
	return answerType, nil
}

func setAnswerType(chatId int64, answerType int) error {

	//Connecting to database
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()

	//Creating SQL command
	data := `INSERT INTO answer_type(CHAT_ID, ANSWER_TYPE) VALUES($1, $2)
			 ON CONFLICT (CHAT_ID)
             DO UPDATE SET
			 ANSWER_TYPE = $2;`

	//Execute SQL command in database
	if result, err := db.Exec(data, chatId, answerType); err != nil {
		fmt.Println(result)
		return err
	}

	fmt.Println("set answer type: ", answerType)

	return nil
}
