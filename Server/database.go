package forum

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Posts struct {
	Id          int
	Categorie   string
	Name        string
	Title       string
	Description string
	Date        string
}

type Commentaire struct {
	Id          int
	PostId      int
	Name        string
	Commentaire string
}

func InitDatabase(database string) *sql.DB {
	fmt.Println("-- Creation --")
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
		PRAGMA foreign_keys = ON;
		CREATE TABLE IF NOT EXISTS users (
			id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS post (
			id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			categorie TEXT NOT NULL,
			name TEXT NOT NULL,
			title	TEXT UNIQUE NOT NULL,
			description TEXT UNIQUE NOT NULL,
			date TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS commentaire (
			id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			commentaire TEXT NOT NULL,
			FOREIGN KEY (post_id) REFERENCES post(id)
		);
		`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return nil
	}
	fmt.Println("-- Créer --")
	return db
}

func InsertIntoUsers(db *sql.DB, name string, email string, password string) (int64, error) {
	result, err := db.Exec(`INSERT INTO users (name, email, password) VALUES (?, ?, ?)`, name, email, password)
	if err != nil {
		fmt.Println("Ce nom ou email existe déjà")
		fmt.Println(err)
		return 0, err
	}
	return result.LastInsertId()
}

func InsertIntoPost(db *sql.DB, categorie string, name string, title string, description string, date string) (int64, error) {
	result, err := db.Exec(`INSERT INTO post (categorie, name, title, description, date) VALUES (?, ?, ?, ?, ?)`, categorie, name, title, description, date)
	if err != nil {
		fmt.Println(err)
		// fmt.Println(err)
		return 0, err
	}
	return result.LastInsertId()
}

func InsertIntoComments(db *sql.DB, commentaire string, name string, postId int) (int64, error) {
	result, err := db.Exec(`INSERT INTO commentaire (commentaire, name, post_id) VALUES (?, ?, ?)`, commentaire, name, postId)
	if err != nil {
		fmt.Println(err)
		// fmt.Println(err)
		return 0, err
	}
	return result.LastInsertId()
}

func SelectAllFromTable(db *sql.DB, table string) *sql.Rows {
	query := "SELECT * FROM " + table
	result, _ := db.Query(query)
	return result
}

func SelectAllByCategorie(db *sql.DB, categorie string) *sql.Rows {
	// var u Posts
	res, _ := db.Query(`SELECT * FROM post WHERE lower(categorie) = '` + categorie + "'") //.Scan(&u.Id, &u.Categorie, &u.Title, &u.Description)
	return res
}

// func selectUserById(db *sql.DB, id int) User {
// 	var u User
// 	db.QueryRow(`SELECT * FROM users WHERE id = ?`, id).Scan(&u.Id, &u.Name, &u.Email, &u.Password)
// 	return u
// }

func SelectUserWhenLogin(db *sql.DB, email string) User {
	var u User
	fmt.Println("select user :", email)
	db.QueryRow(`SELECT * FROM users WHERE email = ?`, email).Scan(&u.Id, &u.Name, &u.Email, &u.Password)
	fmt.Println(u)
	return u
}

// func SelectPostWithId(db *sql.DB, id int) Posts {
// 	var u Posts
// 	fmt.Println("select post :", id)
// 	db.QueryRow(`SELECT * FROM post WHERE id  = ?`, id).Scan(&u.Id, &u.Categorie, &u.Title, &u.Description)
// 	fmt.Println(u)
// 	return u
// }

// func SelectAllPost(db *sql.DB, categorie string) []Posts {
// 	var u Posts
// 	rows := SelectAllByCategorie(db, categorie) //SelectAllFromTable(db, "post")
// 	final := make([]Posts, 0)
// 	for rows.Next() {
// 		rows.Scan(&u.Id, &u.Categorie, &u.Name, &u.Title, &u.Description)
// 		final = append(final, u)

// 	}
// 	return final
// }

func SelectAllPost(db *sql.DB, categorie string) []Posts {
	var u Posts
	rows := SelectAllByCategorie(db, categorie) //SelectAllFromTable(db, "post")
	final := make([]Posts, 0)
	for rows.Next() {
		rows.Scan(&u.Id, &u.Categorie, &u.Name, &u.Title, &u.Description, &u.Date)
		final = append(final, u)

	}
	return final
}

// func SelectAllCommentByPost(db *sql.DB) {
// 	var u Zeubi
// 	rows, _ := db.Query("Select * From post Inner JOIN commentaire WHERE commentaire.post_id = post.id")
// 	fmt.Println(rows)
// 	rows.Scan(&u.Id, &u.Categorie, &u.Name, &u.Title, &u.Description, &u.Date, &u.Test, &u.Post_id, &u.Name2, &u.Commentaire)
// 	fmt.Println(u)
// }

func SelectAllComments(db *sql.DB, post_id int) []Commentaire {
	var u Commentaire
	rows := SelectCommentByPost(db, post_id)
	final := make([]Commentaire, 0)
	for rows.Next() {
		rows.Scan(&u.Id, &u.PostId, &u.Name, &u.Commentaire)
		final = append(final, u)

	}
	fmt.Println(final)
	return final
}

func SelectCommentByPost(db *sql.DB, post_id int) *sql.Rows {
	res, _ := db.Query(`SELECT * FROM commentaire WHERE post_id = ?`, post_id)
	return res
}

func SelectUserNameWithPattern(db *sql.DB, pattern string) []User {
	var u User
	query := "SELECT * FROM users WHERE name LIKE '%" + pattern + "%'"
	rows, _ := db.Query(query)
	final := make([]User, 0)
	for rows.Next() {
		rows.Scan(&u.Name, &u.Email, &u.Password)
		final = append(final, u)
	}
	return final
}

// func main() {
// 	db := initDatabase("AAAforum.db")
// 	defer db.Close()

// 	// fmt.Println("-- Creation --")
// 	// fmt.Println("-- Créer --")

// 	// insertIntoUsers(db, "Mathieu", "m.m@gmail.com", "abcde")
// 	// insertIntoUsers(db, "Thomas", "t.t@gmail.com", "fghij")
// 	// insertIntoUsers(db, "Lucas", "l.l@gmail.com", "klmno")

// 	fmt.Println("-- Sélection --")

// 	selectAllFromTable(db, "users")
// 	// user := selectUserById(db, 2)
// 	// fmt.Println(user)

// 	fmt.Println(selectUserNameWithPattern(db, "as"))
// 	// fmt.Println(test)
// }
