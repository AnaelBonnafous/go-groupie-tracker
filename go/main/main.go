package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string
	Email    string
	Password string
	Height   string
	Gender   string
}

var user User
var db *sql.DB

func main() {
	fmt.Println("Server started on port 8080")

	var err error
	db, err = sql.Open("sqlite3", "groupie_tracker.db")
	logError(err)
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (username text, email text, password text, height text, gender text)")
	logError(err)

	http.HandleFunc("/", MainPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", RegisterPage)

	fs := http.FileServer(http.Dir("css/"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	http.ListenAndServe(":8080", nil)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/index.html")
	logError(err)
	template.Execute(w, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/login-page.html")
	logError(err)
	template.Execute(w, nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/register-page.html", "./templates/register-form.html")
	logError(err)

	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Password = r.FormValue("password")
	user.Height = r.FormValue("height")
	user.Gender = r.FormValue("gender")

	switch r.Method {
	case http.MethodGet:
		template.Execute(w, user)
	case http.MethodPost:
		stmt, err := db.Prepare("INSERT INTO users (username, email, password, height, gender) VALUES (?, ?, ?, ?, ?)")
		logError(err)
		_, err = stmt.Exec(user.Username, user.Email, user.Password, user.Height, user.Gender)
		logError(err)
		fmt.Fprintln(w, "L'utilisateur", user.Username, "a été enregistré avec succès !")
	}
}

func logError(err error) {
	if err != nil {
		panic(err)
	}
}
