package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Artist represents the structure of an artist
type Artist struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Year  int    `json:"year"` // creationDate?
	// concertDates?
	// relations
	FirstAlbum string   `json:"first_album"`
	Members    []string `json:"members"`
}

// Location represents the structure of a location
type Location struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

// ConcertDate represents the structure of a concert date
type ConcertDate struct {
	ID       int    `json:"id"`
	ArtistID int    `json:"artist_id"`
	Date     string `json:"date"`
	Location string `json:"location"`
	IsPast   bool   `json:"is_past"`
}

// Relation represents the structure of the relations between artists, locations, and dates
type Relation struct {
	Artists   []Artist      `json:"artists"`
	Locations []Location    `json:"locations"`
	Concerts  []ConcertDate `json:"concerts"`
}

// entry point to REST server
func main() {
	fs := http.FileServer(http.Dir("templates"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artists", artistsHandler)
	http.HandleFunc("/locations", locationsHandler)
	http.HandleFunc("/dates", datesHandler)
	http.HandleFunc("/relation", relationHandler)

	/*
		We register the atoi function with the template using template.FuncMap and template.New("").Funcs(funcMap).
		This ensures that the custom function is available within the template. Additionally,
		we use template.ParseGlob to parse all HTML template files within the "templates" folder and make them available for execution.
	*/
	funcMap := template.FuncMap{
		"atoi": atoi,
	}
	tmpl := template.New("").Funcs(funcMap)
	tmpl = template.Must(tmpl.ParseGlob("templates/*.html"))

	http.HandleFunc("/artists.html", func(w http.ResponseWriter, r *http.Request) {
		idParam := r.URL.Query().Get("id")
		data := struct {
			IDParam string
		}{
			IDParam: idParam,
		}
		err := tmpl.ExecuteTemplate(w, "artists.html", data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	log.Println("Server is running on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

// Takes a string input, attempts to convert it to an integer using strconv.Atoi, and returns the resulting integer. If the conversion fails, it returns 0.
func atoi(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}

// handles all requests through routh url
func homeHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func artistsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/* A new endpoint handler named artistDetailsHandler added to handle the artists.html endpoint
we retrieve the id query parameter from the request URL using r.URL.Query().Get("id").
We then convert the ID from a string to an integer using strconv.Atoi.*/
func artistDetailsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	// Fetch the artists' data from the API endpoint and decode it into the artists slice.
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	var artists []Artist
	err = json.NewDecoder(response.Body).Decode(&artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find the selected artist by ID
	var selectedArtist Artist
	found := false
	for _, artist := range artists {
		if artist.ID == id {
			selectedArtist = artist
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	// Create a data structure to pass to the template, which includes the selected artist
	data := struct {
		Artist Artist
	}{
		Artist: selectedArtist,
	}

	// Parse the artists.html template file and execute it with the data
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func locationsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var locations []Location
	err = json.NewDecoder(response.Body).Decode(&locations)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles("templates/locations.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, locations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func datesHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var dates []ConcertDate
	err = json.NewDecoder(response.Body).Decode(&dates)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.ParseFiles("templates/dates.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, dates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func relationHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var relation Relation
	err = json.NewDecoder(response.Body).Decode(&relation)
	if err != nil {
		log.Fatal(err)
	}

	data := struct {
		Artists  []Artist
		Concerts []ConcertDate
	}{
		Artists:  relation.Artists,
		Concerts: relation.Concerts,
	}

	tmpl, err := template.ParseFiles("templates/relation.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
