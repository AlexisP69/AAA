package forum

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
)

type NewPost struct {
	Id          int
	Categorie   string
	Title       string
	Description string
	Date        string
}

type PostWithComments struct {
	Post          Posts
	EveryComments []Commentaire
}

type Register struct {
	Name                string
	Email               string
	Password            string
	UserConfirmPassword string
}

type NewComments struct {
	Input  string
	Name   string
	PostId string
}

type Login struct {
	Name     string
	Email    string
	Password string
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

//look if the users is already connected and redirect him if not
func HandleHome(w http.ResponseWriter, r *http.Request) {
	var data Login = Login{}

	if r.URL.Path != "/loginApi" {
		http.NotFound(w, r)
		return
	}
	session, _ := store.Get(r, "cookie-name")
	auth := session.Values["authenticated"]
	if auth == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	json.Unmarshal([]byte(auth.(string)), &data)

	tmpl, _ := template.ParseFiles("Page/HomePage.html", "Page/Signup.html", "templates/footer.html", "templates/navbar.html", "templates/login.html")

	tmpl.Execute(w, data)
}

//handle sessions creation
func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB, login *Login) {
	if r.URL.Path != "/loginApi" {
		http.NotFound(w, r)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form", 500)
		return
	}
	result := SelectUserWhenLogin(db, login.Email)
	pwd2 := []byte(login.Password)
	//compare password with hash in database
	pwdMatch := comparePasswords(result.Password, pwd2)
	if result.Id == 0 || !pwdMatch {
		w.Write([]byte(`{"test": "wrong mail or password"}`))

	} else {
		login.Name = result.Name
		res, _ := json.Marshal(login)
		session, _ := store.Get(r, "cookie-name")
		fmt.Printf("POSTFOR IN LOGIN %v", login)
		session.Values["authenticated"] = string(res)
		session.Save(r, w)
		w.Write([]byte(`{"test": "success"}`))
	}
}

// handle session deletion
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/log-out" {
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
	var register Register
	var login Login
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/HomePage.html", "Page/Signup.html", "templates/footer.html", "templates/navbar.html", "templates/Post.html", "templates/PostBlock.html"))
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
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &register)
		pwd := []byte(register.Password)
		//hash and salt password of the users
		hash := hashAndSalt(pwd)
		validEmail := valid(register.Email)

		// check if the email is write correctly with @ and looks like a valid email
		if validEmail {
			InsertIntoUsers(db, register.Name, register.Email, hash)
		}
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Login.html", "templates/navbar.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})

	http.HandleFunc("/loginApi", func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &login)
		HandleLogin(w, r, db, &login)
	})

	http.HandleFunc("/fondateurs", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Fondateur.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
		}
	})

	http.HandleFunc("/drugs", func(w http.ResponseWriter, r *http.Request) {
		posts := SelectAllPost(db, "drugs")
		postSlice := FilterByCategory(db, posts)
		fmt.Printf("%v", postSlice)
		template := template.Must(template.ParseFiles("Page/Drugs.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html", "templates/filtre.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/erotica", func(w http.ResponseWriter, r *http.Request) {
		posts := SelectAllPost(db, "erotica")
		postSlice := FilterByCategory(db, posts)
		template := template.Must(template.ParseFiles("Page/Erotica.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html", "templates/filtre.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/counterfeit", func(w http.ResponseWriter, r *http.Request) {
		posts := SelectAllPost(db, "counterfeit")
		postSlice := FilterByCategory(db, posts)
		template := template.Must(template.ParseFiles("Page/Counterfeit.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html", "templates/filtre.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/tutorials", func(w http.ResponseWriter, r *http.Request) {
		posts := SelectAllPost(db, "tutorials")
		postSlice := FilterByCategory(db, posts)
		template := template.Must(template.ParseFiles("Page/Tutorials.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/guns", func(w http.ResponseWriter, r *http.Request) {
		posts := SelectAllPost(db, "guns")
		postSlice := FilterByCategory(db, posts)
		template := template.Must(template.ParseFiles("Page/Guns.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/software", func(w http.ResponseWriter, r *http.Request) {
		posts := SelectAllPost(db, "software")
		postSlice := FilterByCategory(db, posts)
		template := template.Must(template.ParseFiles("Page/SoftWare.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		posts := SelectAllPost(db, "games")
		postSlice := FilterByCategory(db, posts)
		template := template.Must(template.ParseFiles("Page/Games.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		posts := SelectAllPost(db, "services")
		postSlice := FilterByCategory(db, posts)
		template := template.Must(template.ParseFiles("Page/Services.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/newPost", func(w http.ResponseWriter, r *http.Request) {
		var post NewPost
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &post)
		InsertIntoPost(db, post.Categorie, login.Name, post.Title, post.Description, post.Date)
	})

	http.HandleFunc("/newComments", func(w http.ResponseWriter, r *http.Request) {
		var Commentaire NewComments
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &Commentaire)
		post_id, _ := strconv.Atoi(Commentaire.PostId)
		InsertIntoComments(db, Commentaire.Input, login.Name, post_id)
	})

	http.HandleFunc("/homepage", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles(
			"Page/Homepage.html",
		))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})

	http.HandleFunc("/Userpage", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles(
			"Page/UserPage.html",
		))
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
