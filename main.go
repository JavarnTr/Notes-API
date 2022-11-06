// forms.go
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "Enterprise"
)

type CreateNote struct {
	Name       string
	Text       string
	Status     string
	Delegation string
	Userid     string
	Time       string
}

type RemoveNote struct {
	NoteID string
	Name   string
}

type CreateUser struct {
	UserID string
	Name   string
}

type EditNote struct {
	NoteID string
}

func main() {

	remove()

	http.ListenAndServe(":8082", nil)
}

func add() {
	tmpl := template.Must(template.ParseFiles("forms/forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := CreateNote{
			Name:       r.FormValue("name"),
			Text:       r.FormValue("text"),
			Status:     r.FormValue("status"),
			Delegation: r.FormValue("delegation"),
			Userid:     r.FormValue("userid"),
			Time:       r.FormValue("time"),
		}

		// do something with details
		_ = details

		//Connect to the database
		connStr := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//Insert the note data into the Postgres database for storage
		insertStatement := `INSERT INTO notes (name, text, status, delegation, userid, date) 
		VALUES ($1, $2, $3, $4, $5, $6)`
		_, err = db.Exec(insertStatement, details.Name, details.Text, details.Status, details.Delegation, details.Userid, details.Time)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("\nRow inserted successfully!")
		}

		tmpl.Execute(w, struct{ Success bool }{true})
	})
}

func remove() {
	tmpl := template.Must(template.ParseFiles("forms/forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		removeDetails := RemoveNote{
			NoteID: r.FormValue("noteid"),
			Name:   r.FormValue("named"),
		}

		// do something with details
		_ = removeDetails

		//Connect to the database
		connStr := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//Insert the note data into the Postgres database for storage
		deleteStatement := `update notes set name = $2 where id = $1`
		_, err = db.Exec(deleteStatement, removeDetails.NoteID, removeDetails.Name)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("\nNote removed successfully.")
		}

		tmpl.Execute(w, struct{ Success bool }{true})
	})
}

func createUser() {
	tmpl := template.Must(template.ParseFiles("forms/forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		accountCreate := CreateUser{
			Name: r.FormValue("name"),
		}

		// do something with details
		_ = accountCreate

		//Connect to the database
		connStr := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		insertStatement := `INSERT INTO users (name) VALUES ($1)`
		_, err = db.Exec(insertStatement, accountCreate.Name)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("\nRow inserted successfully!")
		}

		tmpl.Execute(w, struct{ Success bool }{true})
	})

}

func edit() {
	tmpl := template.Must(template.ParseFiles("forms/forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		editDetails := RemoveNote{
			NoteID: r.FormValue("noteid"),
		}

		// do something with details
		_ = editDetails

		//Connect to the database
		connStr := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//Insert the note data into the Postgres database for storage
		deleteStatement := `Delete from notes WHERE id = $1`
		_, err = db.Exec(deleteStatement, editDetails.NoteID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("\nNote removed successfully.")
		}

		tmpl.Execute(w, struct{ Success bool }{true})
	})
}
