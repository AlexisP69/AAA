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
	if auth == nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	json.Unmarshal([]byte(auth.(string)), &data)

	tmpl, _ := template.ParseFiles("Page/HomePage.html", "Page/Signup.html", "templates/footer.html", "templates/navbar.html", "templates/login.html")

	fmt.Println("DATA IN HOME", data)
	tmpl.Execute(w, data)
}

func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB, login *Login) {
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

	// var login Login

	// body, _ := ioutil.ReadAll(r.Body)

	// json.Unmarshal(body, &login)

	// fmt.Println(login.Email)
	// login.Password, _ = CheckPasswordHash(login.Password)
	fmt.Println(login.Password)
	result := SelectUserWhenLogin(db, login.Email, login.Password)
	if result.Id == 0 {
		w.Write([]byte(`{"test": "wrong mail or password"}`))

	} else {
		login.Name = result.Name
		res, _ := json.Marshal(login)
		session, _ := store.Get(r, "cookie-name")
		fmt.Printf("POSTFOR IN LOGIN %v", login)
		session.Values["authenticated"] = string(res)
		session.Save(r, w)
		w.Write([]byte(`{"test": "success"}`))
		// if fil == true {
		// 	document.querySelectorAll(".signup").forEach((e) => {
		// 		e.style.display = "none";
		// 	})
		// 	document.querySelectorAll(".sign-up").forEach((e) => {
		// 		e.style.display = "none";
		// 	})
		// }
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
		// register.Password, _ = HashPassword(register.Password)
		fmt.Println(register.Password)
		_, err := InsertIntoUsers(db, register.Name, register.Email, register.Password)
		if err != nil {
			// if( err == "UNIQUE constraint failed: users.email") {

			// }
			// fmt.Println(err)
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
			//return
		}
	})

	http.HandleFunc("/drugs", func(w http.ResponseWriter, r *http.Request) {
		var postSlice []PostWithComments
		fmt.Println(register.Name)
		posts := SelectAllPost(db, "drugs")
		for _, post := range posts {
			var t PostWithComments
			t.Post = post
			t.EveryComments = SelectAllComments(db, post.Id)
			postSlice = append(postSlice, t)
		}
		fmt.Printf("%v", postSlice)
		template := template.Must(template.ParseFiles("Page/Drugs.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html", "templates/filtre.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/erotica", func(w http.ResponseWriter, r *http.Request) {
		var postSlice []PostWithComments
		fmt.Println(register.Name)
		posts := SelectAllPost(db, "erotica")
		for _, post := range posts {
			var t PostWithComments
			t.Post = post
			t.EveryComments = SelectAllComments(db, post.Id)
			postSlice = append(postSlice, t)
		}
		template := template.Must(template.ParseFiles("Page/Erotica.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html", "templates/filtre.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/counterfeit", func(w http.ResponseWriter, r *http.Request) {
		var postSlice []PostWithComments
		posts := SelectAllPost(db, "counterfeit")
		for _, post := range posts {
			var t PostWithComments
			t.Post = post
			t.EveryComments = SelectAllComments(db, post.Id)
			postSlice = append(postSlice, t)
		}
		template := template.Must(template.ParseFiles("Page/Counterfeit.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html", "templates/filtre.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/tutorials", func(w http.ResponseWriter, r *http.Request) {
		var postSlice []PostWithComments
		posts := SelectAllPost(db, "tutorials")
		for _, post := range posts {
			var t PostWithComments
			t.Post = post
			t.EveryComments = SelectAllComments(db, post.Id)
			postSlice = append(postSlice, t)
		}
		template := template.Must(template.ParseFiles("Page/Tutorials.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/guns", func(w http.ResponseWriter, r *http.Request) {
		var postSlice []PostWithComments
		posts := SelectAllPost(db, "guns")
		for _, post := range posts {
			var t PostWithComments
			t.Post = post
			t.EveryComments = SelectAllComments(db, post.Id)
			postSlice = append(postSlice, t)
		}
		template := template.Must(template.ParseFiles("Page/Guns.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/software", func(w http.ResponseWriter, r *http.Request) {
		var postSlice []PostWithComments
		posts := SelectAllPost(db, "software")
		for _, post := range posts {
			var t PostWithComments
			t.Post = post
			t.EveryComments = SelectAllComments(db, post.Id)
			postSlice = append(postSlice, t)
		}
		template := template.Must(template.ParseFiles("Page/SoftWare.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/games", func(w http.ResponseWriter, r *http.Request) {
		var postSlice []PostWithComments
		posts := SelectAllPost(db, "games")
		for _, post := range posts {
			var t PostWithComments
			t.Post = post
			t.EveryComments = SelectAllComments(db, post.Id)
			postSlice = append(postSlice, t)
		}
		template := template.Must(template.ParseFiles("Page/Games.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/jsp", func(w http.ResponseWriter, r *http.Request) {
		var postSlice []PostWithComments
		posts := SelectAllPost(db, "jsp")
		for _, post := range posts {
			var t PostWithComments
			t.Post = post
			t.EveryComments = SelectAllComments(db, post.Id)
			postSlice = append(postSlice, t)
		}
		template := template.Must(template.ParseFiles("Page/Jsp.html", "templates/footer.html", "templates/navbar.html", "Page/Signup.html", "Page/Login.html", "templates/Post.html", "templates/PostBlock.html", "templates/CompletePost.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, postSlice)
			return
		}
	})

	http.HandleFunc("/newPost", func(w http.ResponseWriter, r *http.Request) {
		var post NewPost
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &post)
		fmt.Println(body)
		fmt.Println(post)
		InsertIntoPost(db, post.Categorie, login.Name, post.Title, post.Description, post.Date)
	})

	http.HandleFunc("/newComments", func(w http.ResponseWriter, r *http.Request) {
		var Commentaire NewComments
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &Commentaire)
		fmt.Println(db)
		fmt.Println(Commentaire)
		x, _ := strconv.Atoi(Commentaire.PostId)
		InsertIntoComments(db, Commentaire.Input, login.Name, x)
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
