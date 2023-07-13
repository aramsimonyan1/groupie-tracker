package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
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

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artists", artistsHandler)
	http.HandleFunc("/locations", locationsHandler)
	http.HandleFunc("/dates", datesHandler)
	http.HandleFunc("/relation", relationHandler)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
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
