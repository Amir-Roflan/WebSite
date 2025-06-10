package main

import (
	"html/template"
	"net/http"
	"encoding/json"
	"io"
	"fmt"
	"log"
	"github.com/joho/godotenv"
	"os"
)	

const (
	apiBaseURL  = "https://catalog.api.2gis.com/3.0/items"
	staticMapURL = "https://static.maps.2gis.com/1.0"
	imageSize = "600x400"
	zoom = "17"
)

type Point struct{
	Lat 	 float64  `json:"lat"`	// создали структуру для координат потом засунули в Location 
	Lon 	 float64  `json:"lon"`
}

type Location struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address_name"`
	Point 	 Point	  `json:"point"` // все координаты
}

type APIResponse struct {
	Result struct {
		Items []Location `json:"items"`
	} `json:"result"`
}

func searchPage(w http.ResponseWriter, r *http.Request){
	tmpl, err := template.ParseFiles("templates/search_page.html", "templates/header.html", "templates/footer.html")
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} // инициализация страницы

	query := r.URL.Query().Get("q") // считывает данные со строки поиска
	if query == "" {
		http.Error(w, "Нужно что-то ввести", http.StatusBadRequest)
		return
	}

	locations, err := searchLocations(query) 
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type ViewLocation struct{
		ID 		string 		
		Name 	string 	
		Address string	
		Lat		float64 // обьявление координат
		Lon 	float64
		MapURL  string // статитческий URL карты
	}

	var viewLocations []ViewLocation
	
	for _, loc := range locations{ 
		mapURL := generateStaticMapURL(loc.Point.Lat, loc.Point.Lon) // функция которая генерирует URL для статической карты

		vl := ViewLocation{
			ID: 		loc.ID,
			Name: 		loc.Name,
			Address: 	loc.Address,
			Lat: 		loc.Point.Lat,
			Lon: 		loc.Point.Lon,
			MapURL:     mapURL,
		}
		viewLocations = append(viewLocations, vl)
	}

	data := struct {
		Query     string
		Locations []ViewLocation
	}{
		Query:     query,
		Locations: viewLocations,
	}

	err = tmpl.ExecuteTemplate(w, "search_page", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func generateStaticMapURL(lat, lon float64) string {
	apiKey := os.Getenv("API_KEY")
	return fmt.Sprintf("%s?zoom=%s&size=%s&center=%f,%f&key=%s",
		staticMapURL, zoom, imageSize, lon, lat, apiKey)
}

func searchLocations(query string) ([]Location, error) { // добавил lon lat в функцию
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	} // получение API 
	
	apiKey := os.Getenv("API")

	url := fmt.Sprintf("%s?q=%s&key=%s&fields=items.point,items.address_name,items.photo_ids", apiBaseURL, query, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call 2GIS API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("2GIS API returned status %d: %s", resp.StatusCode, body)
	}
	var apiResponse APIResponse

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %v", err)
	}

	log.Printf("Parsed %d locations", len(apiResponse.Result.Items))
	return apiResponse.Result.Items, nil
}


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

	// обьявление страниц на хэдере
	http.HandleFunc("/", homePage)
	http.HandleFunc("/settings_page", settingsPage)
	http.HandleFunc("/about_page", aboutPage)
	http.HandleFunc("/help_page", helpPage) 
	http.HandleFunc("/search_page", searchPage)

	http.ListenAndServe(":8080", nil)
}

func main(){
	handleRequests()
}
