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

type ContactDetails struct {
	Name       string
	Text       string
	Status     string
	Delegation string
	Userid     string
	Time       string
}

func main() {

	add()

	http.ListenAndServe(":8080", nil)
}

func add() {
	tmpl := template.Must(template.ParseFiles("forms/forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		details := ContactDetails{
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
