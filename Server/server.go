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

type Post struct {
	Id          int
	Categorie   string
	Title       string
	Description string
	Date        string
}

type Test struct {
	EveryPost []Posts
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

func HandleHome(w http.ResponseWriter, r *http.Request) {
	var data Login = Login{}

	if r.URL.Path != "/loginApi" {
		http.NotFound(w, r)
		fmt.Println("wowo")
		return
	}
	session, _ := store.Get(r, "cookie-name")
	auth := session.Values["authenticated"]
	fmt.Println(auth)
	// if auth == nil {
	// 	http.Redirect(w, r, "/login", http.StatusFound)
	// 	return
	// }

	json.Unmarshal([]byte(auth.(string)), &data)

	tmpl, _ := template.ParseFiles("Page/HomePage.html", "Page/Signup.html", "templates/footer.html", "templates/navbar.html", "templates/login.html")

	fmt.Println("DATA IN HOME", data)
	tmpl.Execute(w, data)
}

func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.URL.Path != "/loginApi" {
		http.NotFound(w, r)
		return
	}

	// tmpl, _ := template.ParseFiles("Page/HomePage.html", "Page/Signup.html", "templates/footer.html", "templates/navbar.html", "templates/login.html")

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", 500)
		return
	}
	fmt.Println("Welcome")

	// if _, ok := r.PostForm["Submit"]; ok {
	// fmt.Println(string("uv"))

	var login Login

	body, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(body, &login)

	fmt.Println(login.Email)
	result := SelectUserWhenLogin(db, login.Email, login.Password)
	if result.Id == 0 {
		w.Write([]byte(`{"test": "wrong mail or password"}`))

	} else {
		res, _ := json.Marshal(login)
		session, _ := store.Get(r, "cookie-name")
		fmt.Println(session)
		fmt.Printf("POSTFOR IN LOGIN %v", login)
		session.Values["authenticated"] = string(res)
		session.Save(r, w)
		w.Write([]byte(`{"test": "success"}`))
	}
	// http.Redirect(w, r, "/", http.StatusFound)
	// return
	// } else

	// result := SelectUserWhenLogin(db, login.Email, login.Password)
	// if result.Id == 0 {
	// 	w.Write([]byte(`{"test": "wrong mail or password"}`))

	// } else {
	// 	w.Write([]byte(`{"test": "success"}`))
	// }

	// if session.Values["authenticated"] != nil {
	// 	http.Redirect(w, r, "/", http.StatusFound)
	// 	return
	// }

	// tmpl.Execute(w, nil)
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = nil
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func HandleFunc(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/HomePage.html", "Page/Signup.html", "templates/footer.html", "templates/navbar.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})
	http.HandleFunc("/UserPage", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/UserPage.html"))
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
		// fmt.Println(body)
		fmt.Println(db)
		// fmt.Println(register)
		// InsertIntoUsers(db, "name", "email", "password")
		// test := SelectUserById(db, 1)
		// fmt.Println(test)
		fmt.Println(register.Name)
		_, err := InsertIntoUsers(db, register.Name, register.Email, register.Password)
		if err != nil {
			// if( err == "UNIQUE constraint failed: users.email") {

			// }
			// fmt.Println(err)
		}

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
		template := template.Must(template.ParseFiles("Page/Login.html", "templates/navbar.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})

	http.HandleFunc("/loginApi", func(w http.ResponseWriter, r *http.Request) {

		// var login Login
		// fmt.Println(db)

		// body, _ := ioutil.ReadAll(r.Body)

		// json.Unmarshal(body, &login)
		// fmt.Println(body)

		// fmt.Println(login.Email)
		// SelectUserWhenLogin(db, login.Email, login.Password)
		HandleLogin(w, r, db)
		// SelectAllFromTable(db, "users")
	})

	http.HandleFunc("/fondateurs", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Fondateur.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			//return
		}
	})

	http.HandleFunc("/drugs", func(w http.ResponseWriter, r *http.Request) {
		var postSlice Test
		postSlice.EveryPost = SelectAllPost(db, "drugs")
		template := template.Must(template.ParseFiles("Page/Drugs.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/erotica", func(w http.ResponseWriter, r *http.Request) {
		var postSlice Test
		postSlice.EveryPost = SelectAllPost(db, "erotica")
		template := template.Must(template.ParseFiles("Page/Erotica.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/counterfeit", func(w http.ResponseWriter, r *http.Request) {
		var postSlice Test
		postSlice.EveryPost = SelectAllPost(db, "counterfeit")
		template := template.Must(template.ParseFiles("Page/Counterfeit.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/tutorials", func(w http.ResponseWriter, r *http.Request) {
		var postSlice Test
		postSlice.EveryPost = SelectAllPost(db, "tutorials")
		template := template.Must(template.ParseFiles("Page/Tutorials.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/guns", func(w http.ResponseWriter, r *http.Request) {
		var postSlice Test
		postSlice.EveryPost = SelectAllPost(db, "guns")
		template := template.Must(template.ParseFiles("Page/Guns.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/software", func(w http.ResponseWriter, r *http.Request) {
		var postSlice Test
		postSlice.EveryPost = SelectAllPost(db, "software")
		template := template.Must(template.ParseFiles("Page/SoftWare.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		var postSlice Test
		postSlice.EveryPost = SelectAllPost(db, "games")
		template := template.Must(template.ParseFiles("Page/Games.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/jsp", func(w http.ResponseWriter, r *http.Request) {
		var postSlice Test
		postSlice.EveryPost = SelectAllPost(db, "jsp")
		template := template.Must(template.ParseFiles("Page/Jsp.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/newPost", func(w http.ResponseWriter, r *http.Request) {
		var post Post
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &post)
		fmt.Println(body)
		fmt.Println(post)
		InsertIntoPost(db, post.Categorie, post.Title, post.Description, post.Date)
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

	http.HandleFunc("/Userpage", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles(
			"Page/UserPage.html",
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
