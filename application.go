package main

import (
	"encoding/json"
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
	fmt.Println("getRestaurants")

	dbConn := os.Getenv("DBCONN")
	if dbConn == "" {
		log.Printf("Missing dbConn environment variable")

	}
	fmt.Printf("dbConn: %s\n", dbConn)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Printf("Error during sql.Open\n")
		log.Printf("error connecting to db: ", err.Error())
		return
	}
	fmt.Printf("sql.Open was successful\n")
	defer db.Close()

	_, err = db.Exec("use test")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	println("Use worked...")

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("["))

	r, err := db.Query("select * from restaurants")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	defer r.Close()

	count := 0
	for r.Next() {
		var restaurant = Restaurant{}
		err = r.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Image)
		if err != nil {
			fmt.Printf("Error scaning row: %s\n", err.Error())
		} else {
			fmt.Printf("restaurant: %v\n", restaurant)
			b, err := json.Marshal(restaurant)
			if err != nil {
				fmt.Printf("Error during marshal: %s\n", err.Error())
			} else {
				fmt.Printf("Cleaner: %s", string(b))
				if count > 0 {
					w.Write([]byte(","))
				}
				w.Write(b)
				count++
			}

		}
	}

	w.Write([]byte("]"))
	println(r)
}

func getItems(w http.ResponseWriter, req *http.Request) {
	var restaurantID int
	fmt.Println("getItems")
	p := req.URL.Path[1:]
	fmt.Printf("p: %s\n", p)
	c, e := fmt.Sscanf(p, "restaurants/%d/items", &restaurantID)

	println(c, e, restaurantID)
	if e != nil {
		fmt.Println(e.Error())
	}

	dbConn := os.Getenv("DBCONN")
	if dbConn == "" {
		log.Printf("Missing dbConn environment variable")

	}
	fmt.Printf("dbConn: %s\n", dbConn)

	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Printf("Error during sql.Open\n")
		log.Printf("error connecting to db: ", err.Error())
		return
	}
	fmt.Printf("sql.Open was successful\n")
	defer db.Close()

	_, err = db.Exec("use test")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	println("Use worked...")

	w.Header().Set("Content-Type", "application/json")

	w.Write([]byte("["))

	statement := fmt.Sprintf("select * from items where RestaurantID = %d", restaurantID)
	r, err := db.Query(statement)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
	defer r.Close()

	count := 0
	for r.Next() {
		var item = Item{}
		err = r.Scan(&item.ID, &item.RestaurantID, &item.Name, &item.Image, &item.Price)
		if err != nil {
			fmt.Printf("Error scaning row: %s\n", err.Error())
		} else {
			fmt.Printf("item: %v\n", item)
			b, err := json.Marshal(item)
			if err != nil {
				fmt.Printf("Error during marshal: %s\n", err.Error())
			} else {
				fmt.Printf("Cleaner: %s", string(b))
				if count > 0 {
					w.Write([]byte(","))
				}
				w.Write(b)
				count++
			}

		}
	}

	w.Write([]byte("]"))
	println(r)

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

func initializeDB() {
	dbConn := os.Getenv("DBCONN")
	if dbConn == "" {
		log.Printf("Missing dbConn environment variable")

	}
	fmt.Printf("dbConn: %s\n", dbConn)
	db, err := sql.Open("mysql", dbConn)
	if err != nil {
		fmt.Printf("Error during sql.Open\n")
		log.Printf("error connecting to db: ", err.Error())
	} else {
		fmt.Printf("sql.Open was successful\n")
		defer db.Close()

		for i := 0; i < len(dbSetup); i++ {
			fmt.Printf("sql: %s\n", dbSetup[i])
			_, err := db.Exec(dbSetup[i])
			if err != nil {
				log.Printf("Error: %s", err.Error())
			}
		}
	}
}

func main() {

	fmt.Printf("Starting Application\n")

	initializeDB()

	http.HandleFunc("/", root)
	http.HandleFunc("/restaurants", getRestaurants)
	http.HandleFunc("/restaurants/", getItems)
	http.HandleFunc("/ping", pong)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("ListenAndServe: %s", err.Error())
	}
}
