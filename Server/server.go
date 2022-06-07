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
	Name            string
	Email           string
	Password        string
	ConfirmPassword string
	// test            string `json:test`
	//Picture  string `json:picture`
}

func HandleFunc(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/HomePage.html", "templates/footer.html", "templates/navbar.html", "templates/login.html", "templates/Signup.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})

	// http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
	// 	template := template.Must(template.ParseFiles("Page/Signup.html"))
	// 	if r.Method != http.MethodPost {
	// 		template.Execute(w, db)
	// 		return
	// 	}
	// })
	http.HandleFunc("/registerApi", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("{\"test\":\"${Users.name}\""))
		var Users MyUsers
		// w.Write([]byte("{\"name\":\"" + Users.Name + "\"}"))
		// w.Write([]byte("{\"email\":\"" + Users.Email + "\"}"))
		// w.Write([]byte("{\"password\":\"" + Users.Password + "\"}"))
		// w.Write([]byte("{\"userConfirmPassword\":\"" + Users.ConfirmPassword + "\"}"))

		body, _ := ioutil.ReadAll(r.Body)
		fmt.Println(body)
		// fmt.Println(r.Body)
		json.Unmarshal(body, &Users)
		fmt.Println(Users)
		// fmt.Println(Users)

		w.Write([]byte("{\"name\":\"" + Users.Name + "\",\""))
		w.Write([]byte("\"email\":\"" + Users.Email))
		w.Write([]byte("\"password\":\"" + Users.Password))
		w.Write([]byte("\"userConfirmPassword\":\"" + Users.ConfirmPassword + "\"}"))

		if Users.Name == "test" {
			// w.Write([]byte("{\"test\":\"\"}"))
			fmt.Println(Users)
			return
		}

		// Requete SQL
		// Scan requête
		// JSON.Marshal
		// fmt.Println(Users.name)
		// w.Write([]byte("{\"test\":\"\"}"))
	})

	http.HandleFunc("/fondateurs", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Fondateur.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			//return
		}
	})

	http.HandleFunc("/drugs", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Drugs.html", "templates/footer.html", "templates/navbar.html", "templates/login.html", "templates/Signup.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			//return
		}
	})

	http.HandleFunc("/homepage", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles(
			"Page/Homepage.html",
		))
		if r.Method != http.MethodPost {
			err := template.Execute(w, "")
			fmt.Println(err)
			return
		}
	})

	fs := http.FileServer(http.Dir("Static/"))
	http.Handle("/Static/", http.StripPrefix("/Static/", fs))
	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
