package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Post struct { //Standart struct for providing information.
	ID          string `json:"id"`
	Title       string `json:"title"`
	Bodypart    string `json:"bodypart"`
	Musclegroup string `json:"musclegroup"`
	Content     string `json:"content"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "celdrith:Doerak2003@tcp(192.168.178.23:3306)/celdrithdb") //use format: "user:password@tcp(x.x.x.x:3306)/databasename" for correct connection.
	if err != nil {
		log.Fatal("Error connecting to the database, possible helpful info:", err)
	}
	defer db.Close()

	myRouter := mux.NewRouter().StrictSlash(true)
	fmt.Println("API ONLINE") //Tells the user the API succesfully initiated

	myRouter.HandleFunc("/", getPosts).Methods("GET")
	myRouter.HandleFunc("/fitness/", getPosts).Methods("GET")
	myRouter.HandleFunc("/fitness/exercises", handler).Methods("GET")

	log.Fatal(http.ListenAndServe(":420", myRouter))
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var posts []Post
	result, err := db.Query("SELECT * from exercises") //Provides all exercises in the database in table 'exercises'
	if err != nil {
		log.Fatal("Database query failed", err)
	}
	defer result.Close()
	for result.Next() {
		var post Post
		err := result.Scan(&post.ID, &post.Title, &post.Bodypart, &post.Musclegroup, &post.Content)
		if err != nil {
			log.Fatal("Database table scan failed", err)
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

func handler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["key"] //Query key! format like "localhost:8000/yourname/?s=yourkey".

	if !ok || len(keys[0]) < 1 { //if there is no input provided it will tell the user in the terminal.
		log.Println("Url Param 'key' is missing")
		return
	}
	key1 := keys[0]
	log.Println("Url Param 'key' is: " + string(key1)) //Tells the terminal what the provided key is for debugging reasons, comment this line if you don't need this.

	w.Header().Set("Content-Type", "application/json")
	var posts []Post
	result, err := db.Query("SELECT * FROM exercises WHERE Bodypart = '" + key1 + "'") //Selects the chosen variable within the database and get a response.
	if err != nil {
		log.Fatal("Database query failed", err)
	}
	defer result.Close()
	for result.Next() {
		var post Post
		err := result.Scan(&post.ID, &post.Title, &post.Bodypart, &post.Musclegroup, &post.Content)
		if err != nil {
			log.Fatal("Database table scan failed", err)
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}
