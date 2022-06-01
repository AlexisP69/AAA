package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type MyUsers struct {
	name     string `json:name`
	email    string `json:email`
	password string `json:password`
	//Picture  string `json:picture`
}

func HandleFunc(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/HomePage.html", "templates/footer.html", "templates/navbar.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Register.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, db)
			return
		}
	})
	http.HandleFunc("/registerApi", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"test\":\"tata\"}"))
		var Users MyUsers

		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &Users)
		fmt.Println(Users)

		if Users.name == "test" {
			w.Write([]byte("{\"test\":\"tata\"}"))
			return
		}

		// Requete SQL
		// Scan requête
		// JSON.Marshal
		w.Write([]byte("{\"test\":\"toto\"}"))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Login.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})

	fs := http.FileServer(http.Dir("Static/"))
	http.Handle("/Static/", http.StripPrefix("/Static/", fs))
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
