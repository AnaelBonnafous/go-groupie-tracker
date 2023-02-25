package main

import (
	"html/template"
	"log"
	"net/http"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string
	Email string
	Sexe string
	Taille string
	Password string
}

var user User
var db *sql.DB


func main() {
	fmt.Println("Server started on port 8080")
	http.HandleFunc("/", MainPage)
	http.HandleFunc("/login", LoginPage)
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		RegisterPage(w,r, user)})


	fs := http.FileServer(http.Dir("css/"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	http.ListenAndServe(":8080", nil)
}

func MainPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/index.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	template, err := template.ParseFiles("./pages/login-page.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, nil)
}

func RegisterPage(w http.ResponseWriter, r *http.Request, user User) {
	template, err := template.ParseFiles("./pages/register-page.html", "./templates/register-form.html")
	if err != nil {
		log.Fatal(err)
	}

	user.Username = r.FormValue("username")
	user.Email = r.FormValue("email")
	user.Sexe = r.FormValue("sexe")
	user.Taille = r.FormValue("taille")
	user.Password = r.FormValue("password")

	if r.Method != http.MethodPost {
		template.Execute(w, user)
		return
	}

	db, err = sql.Open("sqlite3", "./groupie_tracker.db")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO users (username, email) VALUES (?, ?)")
	if err != nil {
		panic(err)
	}

	_, err = stmt.Exec(user.Username, user.Email)
	if err != nil {
		panic(err)
	}
}

func init() {
	var err error

	db, err = sql.Open("sqlite3", "../../groupie_tracker.db")
	
	if err != nil {
		panic(err)
	}
	sqlStmt := `CREATE TABLE IF NOT EXISTS utilisateurs (username text, email text, password text);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	email := r.FormValue("email")
	// password := r.FormValue("password")

	stmt, err := db.Prepare("INSERT INTO utilisateurs (nom, email) VALUES (?, ?)")
	if err != nil {
		panic(err)
	}
	db.Close()

	_, err = stmt.Exec(username, email)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, "L'utilisateur", username, "a été enregistré avec succès !")
}
