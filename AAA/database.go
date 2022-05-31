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
	Picture  string
}

func InitDatabase(database string) *sql.DB {
	fmt.Println("-- Creation --")
	db, err := sql.Open("sqlite3", database)
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS users (
			id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL
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
	result, _ := db.Exec(`INSERT INTO users (name, email, password) VALUES (?, ?, ?)`, name, email, password)
	return result.LastInsertId()
}

func SelectAllFromTable(db *sql.DB, table string) *sql.Rows {
	query := "SELECT * FROM " + table
	result, _ := db.Query(query)
	return result
}

func SelectUserById(db *sql.DB, id int) User {
	var u User
	db.QueryRow(`SELECT * FROM users WHERE id = ?`, id).Scan(&u.Id, &u.Name, &u.Email, &u.Password)
	return u
}

func SelectUserNameWithPattern(db *sql.DB, pattern string) []User {
	var u User
	query := "SELECT * FROM users WHERE name LIKE '%" + pattern + "%'"
	rows, _ := db.Query(query)
	final := make([]User, 0)
	for rows.Next() {
		rows.Scan(&u.Id, &u.Name, &u.Email, &u.Password)
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
