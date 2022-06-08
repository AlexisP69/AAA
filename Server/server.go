package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
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

func HandleFunc(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/HomePage.html", "Page/Signup.html", "templates/footer.html", "templates/navbar.html"))
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

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Signup.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
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
		fmt.Println(body)
		// fmt.Println(register)
		// InsertIntoUsers(db, "name", "email", "password")
		// test := SelectUserById(db, 1)
		// fmt.Println(test)
		fmt.Println(register.Name)
		InsertIntoUsers(db, register.Name, register.Email, register.Password)

		// if err != nil {
		// 	// fmt.Println(err)
		// 	// w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		// } else {
		// w.Write([]byte("{\"error\": \"" + err.Error() + "\"}"))
		// w.Write([]byte("{\"name\": \"" + register.Name + "\","))
		// w.Write([]byte("\"email\": \"" + register.Email + "\","))
		// w.Write([]byte("\"password\": \"" + register.Password + "\","))
		// w.Write([]byte("\"confirmPassword\": \"" + register.UserConfirmPassword + "\"}"))
		// }

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

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Login.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})

	http.HandleFunc("/loginApi", func(w http.ResponseWriter, r *http.Request) {

	})

	http.HandleFunc("/fondateurs", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Fondateur.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			//return
		}
	})

	http.HandleFunc("/drugs", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Drugs.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
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
