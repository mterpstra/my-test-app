package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// db is a global to this file.
var db *sql.DB

var dbSetup = [...]string{
	//"DROP DATABASE IF EXISTS ebdb;",
	//"CREATE DATABASE ebdb IF NOT EXISTS;",
	//"USE ebdb;",
	"CREATE TABLE restaurants (ID int NOT NULL AUTO_INCREMENT, Name varchar(255) NOT NULL, Image varchar(255) NOT NULL, PRIMARY KEY (ID));",
	"INSERT INTO restaurants (Name,Image) VALUES ('Leftys', 'http://leftystaverncoralsprings.com/images/sharklogo.png');",
	"INSERT INTO restaurants (Name,Image) VALUES ('The Whale Tale', 'http://www.thewhalerawbar.com/images/logo.png');",
	"CREATE TABLE items (ID int NOT NULL AUTO_INCREMENT, RestaurantID int NOT NULL, Name varchar(255) NOT NULL, Image varchar(255), Price float, PRIMARY KEY (ID));",
	"ALTER TABLE items ADD FOREIGN KEY (RestaurantID) REFERENCES restaurants(ID);",
	"INSERT INTO items(RestaurantID,Name,Image,Price) VALUES (1, 'Chicken Wings', 'http://cookdiary.net/wp-content/uploads/images/Spicy-Chicken_12957.jpg', 9.99);",
	"INSERT INTO items(RestaurantID,Name,Image,Price) VALUES (1, 'Big Pretzel', 'http://www.thedeliciouslife.com/wp-content/plugins/hot-linked-image-cacher/upload/photos1.blogger.com/img/98/3385/640/foodies_pretzel.jpg', 8.99);",
	"INSERT INTO items(RestaurantID,Name,Image,Price) VALUES (2, 'Sushi', 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRTAAxaDS1TRufhvUF9zWq0IS5skcYEngPj7M-NQcj8lCREdV67', 19.99);",
	"INSERT INTO items(RestaurantID,Name,Image,Price) VALUES (2, 'Clams', 'http://img1.sunset.timeinc.net/sites/default/files/image/2011/09/razor-clams-colander-l.jpg', 12.99);"}

func getRestaurants(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	list, err := GetAllRestaurants(db)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		b, _ := json.Marshal(err)
		w.Write(b)
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
		fmt.Printf("Error: %s", err.Error())
		b, _ := json.Marshal(err)
		w.Write(b)
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

	fmt.Printf("Entered initializeDB, reading ENV Var\n")

	dbConn := os.Getenv("DBCONN")
	if dbConn == "" {
		fmt.Printf("Missing DBCONN environment variable\n")
		return nil, errors.New("Missing dbConn environment variable")
	}
	fmt.Printf("Found an ENV Var DBCONN %s\n", dbConn)

	fmt.Printf("Opening DB next\n")
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Printf("Error during sql.Open %s\n", err.Error())
		return nil, err
	}

	fmt.Printf("DB is open, lets run these SQL statements next\n")
	for i := 0; i < len(dbSetup); i++ {
		fmt.Printf("Running SQL: %s\n", dbSetup[i])
		_, err := db.Exec(dbSetup[i])
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
	}

	fmt.Printf("Wow, all worked. Ending initializeDB()\n")
	return db, nil
}

// SomeFunc is just a func
func SomeFunc() {
}

func main() {

	fmt.Printf("Starting Application\n")

	var err error
	db, err = initializeDB()
	if err != nil {
		fmt.Print("Error initializing DB: %s\n", err.Error())
		return
	}
	//defer db.Close()

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
