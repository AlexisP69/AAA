package forum

import (
	"fmt"
	"html/template"
	"net/http"
)

func HandleFunc() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/HomePage.html", "templates/footer.html", "templates/navbar.html"))
		if r.Method != http.MethodPost {
			template.Execute(w, "")
			return
		}
	})

	http.HandleFunc("/Therms-of-use", func(w http.ResponseWriter, r *http.Request) {
		template := template.Must(template.ParseFiles("Page/Therms-of-use.html"))
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
