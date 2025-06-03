package main

import (
	"html/template"
	"net/http"
)

func helpPage(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/help_page.html", "templates/header.html", "templates/footer.html")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "help_page", nil)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func settingsPage(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/settings_page.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	err = tmpl.ExecuteTemplate(w, "settings_page", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) 
	}
}

func aboutPage(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/about_page.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}

	err = tmpl.ExecuteTemplate(w, "about_page", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) 
	}
}

func homePage(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		http.NotFound(w, r) // добавлено только что 
		return
	}
	

	tmpl, err := template.ParseFiles("templates/home_page.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	
	err = tmpl.ExecuteTemplate(w, "home_page", nil)
	if err != nil{ 
		http.Error(w, err.Error(), http.StatusInternalServerError) 
	}
}

func handleRequests(){
	fs := http.FileServer(http.Dir("static")) // хранение статитечских данных
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homePage)
	http.HandleFunc("/settings_page", settingsPage)
	http.HandleFunc("/about_page", aboutPage)
	http.HandleFunc("/help_page", helpPage)
	http.ListenAndServe(":8080", nil)
}

func main(){
	handleRequests()
}
