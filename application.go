package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

// db is a global to this file.
var db *sql.DB

var dbSetup = [...]string{"DROP DATABASE IF EXISTS test;",
	"CREATE DATABASE test;",
	"USE test;",
	"CREATE TABLE restaurants (ID int NOT NULL AUTO_INCREMENT, Name varchar(255) NOT NULL, Image varchar(255) NOT NULL, PRIMARY KEY (ID));",
	"INSERT INTO restaurants (Name,Image) VALUES ('Leftys', 'http://leftystaverncoralsprings.com/images/sharklogo.png');",
	"INSERT INTO restaurants (Name,Image) VALUES ('The Whale Tale', 'http://www.thewhalerawbar.com/images/logo.png');",
	"CREATE TABLE items (ID int NOT NULL AUTO_INCREMENT, RestaurantID int NOT NULL, Name varchar(255) NOT NULL, Image varchar(255), Price float, PRIMARY KEY (ID));",
	"ALTER TABLE items ADD FOREIGN KEY (RestaurantID) REFERENCES restaurants(ID);",
	"INSERT INTO items(RestaurantID,Name,Image,Price) VALUES (1, 'Chicken Wings', '', 9.99);",
	"INSERT INTO items(RestaurantID,Name,Image,Price) VALUES (1, 'Big Pretzel', '', 8.99);",
	"INSERT INTO items(RestaurantID,Name,Image,Price) VALUES (2, 'Sushi', '', 19.99);",
	"INSERT INTO items(RestaurantID,Name,Image,Price) VALUES (2, 'Clams', '', 12.99);"}

func getRestaurants(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	list, err := GetAllRestaurants(db)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	b, err := json.Marshal(list)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	w.Write(b)
}

func getItems(w http.ResponseWriter, req *http.Request) {

	var restaurantID int
	p := req.URL.Path[1:]
	_, e := fmt.Sscanf(p, "restaurants/%d/items", &restaurantID)
	if e != nil {
		fmt.Println(e.Error())
	}

	w.Header().Set("Content-Type", "application/json")

	items, err := GetItemsForRestaurantID(db, restaurantID)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}

	b, err := json.Marshal(items)
	if err != nil {
		fmt.Printf("Error during marshal: %s\n", err.Error())
		return
	}
	w.Write(b)
}

func root(w http.ResponseWriter, req *http.Request) {

	var filename string
	if filename = req.URL.Path[1:]; filename == "" {
		filename = "index.html"
	}
	filename = "public/" + filename
	fmt.Printf("Filename is: [%s]\n", filename)
	dat, err := ioutil.ReadFile(filename)

	acceptHeader := req.Header.Get("Accept")

	if strings.Contains(acceptHeader, "text/css") {
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	}

	if err != nil {
		io.WriteString(w, err.Error())
	} else {
		w.Write(dat)
	}
}

func pong(w http.ResponseWriter, req *http.Request) {
	fmt.Print("pong handler")
	io.WriteString(w, "pong v05")
}

func initializeDB() (*sql.DB, error) {

	dbConn := os.Getenv("DBCONN")
	if dbConn == "" {
		log.Printf("Missing dbConn environment variable")
		return nil, errors.New("Missing dbConn environment variable")
	}
	fmt.Printf("dbConn: %s\n", dbConn)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Printf("Error during sql.Open\n")
		log.Printf("error connecting to db: ", err.Error())
		return nil, err
	}
	fmt.Printf("sql.Open was successful\n")

	for i := 0; i < len(dbSetup); i++ {
		fmt.Printf("sql: %s\n", dbSetup[i])
		_, err := db.Exec(dbSetup[i])
		if err != nil {
			log.Printf("Error: %s", err.Error())
		}
	}
	return db, nil
}

func main() {

	fmt.Printf("Starting Application\n")

	var err error
	db, err = initializeDB()
	if err != nil {
		fmt.Print("Error initializing DB: %s\n", err.Error())
		return
	}
	defer db.Close()

	http.HandleFunc("/", root)
	http.HandleFunc("/restaurants", getRestaurants)
	http.HandleFunc("/restaurants/", getItems)
	http.HandleFunc("/ping", pong)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("ListenAndServe: %s", err.Error())
	}
}
