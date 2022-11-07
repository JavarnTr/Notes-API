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

type Note struct {
	NoteID     string
	Name       string
	Text       string
	Status     string
	Delegation string
	Userid     string
	Time       string
}

type User struct {
	UserID string
	Name   string
}

func main() {

	tmpl := template.Must(template.ParseFiles("forms/forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		submit := r.FormValue("submit")

		if submit == "submit1" {
			//Connect to the database
			connStr := fmt.Sprintf("host=%s port=%d user=%s "+
				"password=%s dbname=%s sslmode=disable",
				host, port, user, password, dbname)
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			//-------------------------------Add Note-------------------------------//
			details := Note{
				Name:       r.FormValue("addName"),
				Text:       r.FormValue("addText"),
				Status:     r.FormValue("addStatus"),
				Delegation: r.FormValue("addDelegation"),
				Userid:     r.FormValue("addUserID"),
				Time:       r.FormValue("addTime"),
			}

			// do something with details
			_ = details

			//Insert the note data into the Postgres database for storage
			insertStatement := `INSERT INTO notes (name, text, status, delegation, userid, date) VALUES ($1, $2, $3, $4, $5, $6)`
			_, err = db.Exec(insertStatement, details.Name, details.Text, details.Status, details.Delegation, details.Userid, details.Time)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("\nRow inserted successfully!")
			}
		} else if submit == "submit2" {
			//Connect to the database
			connStr := fmt.Sprintf("host=%s port=%d user=%s "+
				"password=%s dbname=%s sslmode=disable",
				host, port, user, password, dbname)
			db, err := sql.Open("postgres", connStr)
			if err != nil {
				panic(err)
			}
			defer db.Close()

			//-------------------------------Remove Note-------------------------------//
			removeDetails := Note{
				NoteID: r.FormValue("deleteNoteID"),
			}

			// do something with details
			_ = removeDetails

			//Insert the note data into the Postgres database for storage
			deleteStatement := `delete from notes where id = $1`
			_, err = db.Exec(deleteStatement, removeDetails.NoteID)
			if err != nil {
				panic(err)
			} else {
				fmt.Println("\nNote removed successfully.")
			}
		}

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":8080", nil)
}

func add() {
	tmpl := template.Must(template.ParseFiles("forms/forms.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}

		//Connect to the database
		connStr := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//-------------------------------Add Note-------------------------------//
		details := Note{
			Name:       r.FormValue("addName"),
			Text:       r.FormValue("addText"),
			Status:     r.FormValue("addStatus"),
			Delegation: r.FormValue("addDelegation"),
			Userid:     r.FormValue("addUserID"),
			Time:       r.FormValue("addTime"),
		}

		// do something with details
		_ = details

		//Insert the note data into the Postgres database for storage
		insertStatement := `INSERT INTO notes (name, text, status, delegation, userid, date) VALUES ($1, $2, $3, $4, $5, $6)`
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

		//Connect to the database
		connStr := fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		//-------------------------------Remove Note-------------------------------//
		removeDetails := Note{
			NoteID: r.FormValue("deleteNoteID"),
		}

		// do something with details
		_ = removeDetails

		//Insert the note data into the Postgres database for storage
		deleteStatement := `delete from notes where id = $1`
		_, err = db.Exec(deleteStatement, removeDetails.NoteID)
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

		accountCreate := User{
			Name: r.FormValue("createName"),
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

		editDetails := Note{
			NoteID:     r.FormValue("editNoteID"),
			Name:       r.FormValue("editName"),
			Text:       r.FormValue("editText"),
			Status:     r.FormValue("editStatus"),
			Delegation: r.FormValue("editDelegation"),
			Userid:     r.FormValue("editUserID"),
			Time:       r.FormValue("editDate"),
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
		deleteStatement := `update notes set (name, text, status, delegation, userid, date) = ($1, $2, $3, $4, $5, $6) where id = $7`
		_, err = db.Exec(deleteStatement, editDetails.Name, editDetails.Text, editDetails.Status, editDetails.Delegation, editDetails.Userid, editDetails.Time, editDetails.NoteID)
		if err != nil {
			panic(err)
		} else {
			fmt.Println("\nNote removed successfully.")
		}

		tmpl.Execute(w, struct{ Success bool }{true})
	})
}
