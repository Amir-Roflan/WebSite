package main

import (

	"html/template"
	"net/http"
)
type User struct
{
	Name string
	Age int
}
func homepage(w http.ResponseWriter, r *http.Request){
	person := User{"Alan", 18}
	tmpl, err := template.ParseFiles("templates/home_page.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	
	err = tmpl.Execute(w, person)
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
