package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/sessions"
)

type test struct {
	Enregistrer []Register
	Connecter   []Login
}

type Register struct {
	Name                string
	Email               string
	Password            string
	UserConfirmPassword string
}

type Login struct {
	Email    string
	Password string
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func secret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	// Print secret message
	fmt.Fprintln(w, "The cake is a lie!")
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	session, _ := store.Get(r, "cookie-name")

	// Authentication goes here
	// ...

	// Set user as authenticated
	session.Values["authenticated"] = true
	session.Save(r, w)
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}

func HandleFunc(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/HomePage.html", "templates/footer.html", "templates/navbar.html", "templates/login.html", "templates/Signup.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})
	http.HandleFunc("/Therms-of-use", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Therms-of-use.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
		}
	})

	http.HandleFunc("/registerApi", func(w http.ResponseWriter, r *http.Request) {
		// w.Write([]byte("{\"test\":\"${Users.name}\""))
		var register Register
		// w.Write([]byte("{\"name\":\"" + register.Name + "\"}"))
		// w.Write([]byte("{\"email\":\"" + register.Email + "\"}"))
		// w.Write([]byte("{\"password\":\"" + register.Password + "\"}"))
		// w.Write([]byte("{\"userConfirmPassword\":\"" + register.ConfirmPassword + "\"}"))

		body, _ := ioutil.ReadAll(r.Body)
		// fmt.Println(r.Body)
		json.Unmarshal(body, &register)
		// fmt.Println(register)
		// InsertIntoUsers(db, "name", "email", "password")
		// test := SelectUserById(db, 1)
		// fmt.Println(test)
		// fmt.Println(register)
		_, err := InsertIntoUsers(db, register.Name, register.Email, register.Password)

		if err != nil {
			fmt.Println(err)
			w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		} else {
			// w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
			// w.Write([]byte("{\"name\": \"" + register.Name + "\","))
			// w.Write([]byte("\"email\": \"" + register.Email + "\","))
			// w.Write([]byte("\"password\": \"" + register.Password + "\","))
			// w.Write([]byte("\"confirmPassword\": \"" + register.UserConfirmPassword + "\"}"))
		}

		// w.Write([]byte(Users.Email))
		// w.Write([]byte(Users.Password))
		// w.Write([]byte(Users.UserConfirmPassword))
		// w.Write([]byte("{\"name\":\"" + Users.Name + "\",\""))
		// w.Write([]byte("\"email\":\"" + Users.Email))
		// w.Write([]byte("\"password\":\"" + Users.Password))
		// w.Write([]byte("\"userConfirmPassword\":\"" + Users.ConfirmPassword + "\"}"))

		// if Users.Name == "test" {
		// 	// w.Write([]byte("{\"test\":\"\"}"))
		// 	fmt.Println(Users)
		// 	return
		// }

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
