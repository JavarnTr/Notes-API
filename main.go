// forms.go
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "Enterprise"
)

//Struct for the note table
type Note struct {
	NoteID     string
	Name       string
	Text       string
	Status     string
	Delegation string
	Userid     string
	Time       string
}

//Struct for the user table
type User struct {
	UserID string
	Name   string
}

//Empty variable that will later be used to store note data when it is retrieved.
var noteData Note

func main() {

	tmpl := template.Must(template.ParseFiles("forms/forms.html"))

	//Get the current time so it can be added to the note as the time of creation.
	var currentTime = time.Now()

	//Connect to the database
	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		//Finds out what button on the html page was clicked, so that it knows which function to run.
		submit := r.FormValue("submit")

		if submit == "addNote" {
			//-------------------------------Add Note-------------------------------//
			details := Note{
				Name:       r.FormValue("addName"),
				Text:       r.FormValue("addText"),
				Status:     r.FormValue("addStatus"),
				Delegation: r.FormValue("addDelegation"),
				Userid:     r.FormValue("addUserId"),
			}

			//Details
			_ = details

			//Insert the note data into the Postgres database for storage
			insertStatement := `INSERT INTO notes (name, text, status, delegation, userid, date) VALUES ($1, $2, $3, $4, $5, $6)`
			_, err = db.Exec(insertStatement, details.Name, details.Text, details.Status, details.Delegation, details.Userid, currentTime)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("\nRow inserted successfully!")
			}
		} else if submit == "deleteNote" {
			//-------------------------------Remove Note-------------------------------//

			var currentID = r.FormValue("deleteLoggedUser")

			removeDetails := Note{
				NoteID: r.FormValue("deleteNoteID"),
			}

			//Details
			_ = removeDetails

			//Insert the note data into the Postgres database for storage
			deleteStatement := `delete from notes where id = $1 and userid = $2`
			_, err = db.Exec(deleteStatement, removeDetails.NoteID, currentID)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("\nNote removed successfully.")
				fmt.Println(removeDetails.NoteID)
			}
		} else if submit == "createAcc" {
			//-------------------------------Create User-------------------------------//
			accountCreate := User{
				Name: r.FormValue("createName"),
			}

			//Details
			_ = accountCreate

			insertStatement := `INSERT INTO users (name) VALUES ($1)`
			_, err = db.Exec(insertStatement, accountCreate.Name)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("\nRow inserted successfully!")
			}
		} else if submit == "editNote" {
			//-------------------------------Edit Note-------------------------------//
			editDetails := Note{
				NoteID:     r.FormValue("editNoteID"),
				Name:       r.FormValue("editName"),
				Text:       r.FormValue("editText"),
				Status:     r.FormValue("editStatus"),
				Delegation: r.FormValue("editDelegation"),
			}

			var currentlyLogged = r.FormValue("editLoggedUser")

			//Details
			_ = editDetails

			//Insert the note data into the Postgres database for storage
			editStatement := `update notes set (name, text, status, delegation, date) = ($1, $2, $3, $4, $5) where id = $6 and userid = $7`
			_, err = db.Exec(editStatement, editDetails.Name, editDetails.Text, editDetails.Status, editDetails.Delegation, currentTime, editDetails.NoteID, currentlyLogged)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("\nNote removed successfully.")
			}
		} else if submit == "searchNote" {
			//-------------------------------Search Notes-------------------------------//
			//Interact with the database to search for those notes that match user input. This is case insensitive.
			var searchedValue = r.FormValue("searchValue")
			var loggedInUser = r.FormValue("searchLoggedUser")

			searchStatement := `Select * from notes Where LOWER(name) = LOWER($1) or LOWER(text) = LOWER($1) or LOWER(status) = LOWER($1) or LOWER(delegation) = LOWER($1) and userid = $2;`
			rows, err := db.Query(searchStatement, searchedValue, loggedInUser)
			if err != nil {
				log.Fatal(err)
				fmt.Println("An error occurred when querying data!")
			}

			if err == sql.ErrNoRows {
				println("Does not exist")
			}
			defer rows.Close()

			for rows.Next() {
				var id string
				var name string
				var text string
				var status string
				var delegation string
				var userid string
				var date string

				//Print out the recieved data if it is found.
				switch err = rows.Scan(&id, &name, &text, &status, &delegation, &userid, &date); err {
				case sql.ErrNoRows:
					fmt.Println("No rows were returned!")
				case nil:
					noteData = Note{NoteID: id, Name: name, Text: text, Status: status, Delegation: delegation, Userid: userid, Time: date}

					fmt.Println("ID:", id, "| Note Name:", name, "| Note Text:", text, "| Note Status:", status, "| Delegation:", delegation, "| Users:", userid, "| Time of creation:", date)
				default:
					fmt.Println("SQL query error occurred: ")
					panic(err)
				}
			}
			// get any error encountered during iteration
			err = rows.Err()
			if err != nil {
				panic(err)
			}

		} else if submit == "displayNote" {
			//-------------------------------Display Notes-------------------------------//
			//Interact with the database to search for those notes that match user input
			var loggedIn = r.FormValue("displayLoggedUser")

			selectStatement := `Select * from notes where userid = $1;`
			rows, err := db.Query(selectStatement, loggedIn)
			if err != nil {
				log.Fatal(err)
				fmt.Println("An error occurred when querying data!")
			}

			if err == sql.ErrNoRows {
				println("Does not exist")
			}
			defer rows.Close()

			for rows.Next() {
				var id string
				var name string
				var text string
				var status string
				var delegation string
				var userid string
				var date string

				//Print out the recieved data if it is found.
				switch err = rows.Scan(&id, &name, &text, &status, &delegation, &userid, &date); err {
				case sql.ErrNoRows:
					fmt.Println("No rows were returned!")
				case nil:
					noteData = Note{NoteID: id, Name: name, Text: text, Status: status, Delegation: delegation, Userid: userid, Time: date}
					fmt.Println("ID:", id, "| Note Name:", name, "| Note Text:", text, "| Note Status:", status, "| Delegation:", delegation, "| Users:", userid, "| Time of creation:", date)
					tmpl.Execute(w, noteData)
				default:
					fmt.Println("SQL query error occurred: ")
					panic(err)
				}
			}
			// get any error encountered during iteration
			err = rows.Err()
			if err != nil {
				panic(err)
			}
		}
		tmpl.Execute(w, noteData)
	})

	http.ListenAndServe(":8080", nil)
}
