package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/rs/cors"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model

	Name  string
	Email string
	Notes []Note
}

type Note struct {
	gorm.Model

	Content string
	OwnerID int
}

var db *gorm.DB
var err error

func main() {
	// Open the database
	db, err = gorm.Open(
		"postgres", "host=localhost port=5432 user=postgres dbname=notesprototype sslmode=disable password=passsword")
	if err != nil {
		panic(err)
	}

	// Close the database connection when the program stops running
	defer db.Close()

	// Make automatic migrations for the Users and Notes if they haven't been migrated yet
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Note{})

	/*---------------- Route Handling ----------------*/

	// Initialize router
	router := mux.NewRouter()

	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/user/{id}", getUser).Methods("GET") // Includes notes related to user
	router.HandleFunc("/note/{id}", getNote).Methods("GET")

	// Cors middleware
	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

/*------------- API Getters and Setters --------------*/

func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User

	db.Find(&users)

	err = json.NewEncoder(w).Encode(&users)
	if err != nil {
		fmt.Println(err)
	}
}

func getNote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var note Note

	db.First(&note, params["id"])

	err = json.NewEncoder(w).Encode(&note)
	if err != nil {
		fmt.Println(err)
	}
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user User
	var notes []Note

	db.First(&user, params["id"])
	db.Model(&user).Related(&notes)

	user.Notes = notes

	err = json.NewEncoder(w).Encode(&user)
	if err != nil {
		fmt.Println(err)
	}
}
