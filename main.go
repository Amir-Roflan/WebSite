package main

import (
	"html/template"
	"net/http"
)

func homepage(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/home_page.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	err = tmpl.ExecuteTemplate(w, "home_page", nil)
	if err != nil{ 
		http.Error(w, err.Error(), http.StatusInternalServerError) 
	}
}

func handleRequests(){
	http.HandleFunc("/", homepage)
	http.ListenAndServe(":8080", nil)
}

func main(){
	handleRequests()
}
