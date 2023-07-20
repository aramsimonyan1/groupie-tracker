package main

/* This code sets up a basic HTTP server that serves different HTML templates and handles different endpoints
to display information about artists, their locations, concert dates, and relations (dates and locations). The server
communicates with the GroupieTrackers API to retrieve artist data and display it on the corresponding HTML pages. */

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// Define the Artist struct to represent the structure of an artist with various fields like ID, Image, Name, Members, etc.
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

// The Location struct represents the structure of a location with fields ID, Locations, and Dates.
type Location struct {
	ID        int    `json:"id"`
	Locations string `json:"locations"`
	Dates     string `json:"dates"`
}

// Struct represents the structure of the response from the /api/locations endpoint.
// It contains an array of locations along with their respective IDs.
type LocationResponse struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

// Struct consists artists last and/or upcoming concert dates.
type Date struct {
	ID    int    `json:"id"`
	Dates string `json:"dates"`
}

// Struct represents the structure of concert dates. It contains an array of dates for each artist.
type ConcertDate struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

// Relation represents the structure of the relations between artists, locations, and dates
// It contains an array of artists along with a map that stores the dates and locations for each artist.
type Relations struct {
	Index []struct {
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

// We set up the HTTP server and handle different endpoints:
func main() {
	// Serves static files (CSS, JS, etc.) from the static directory.
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", homeHandler)                      // Handles the root endpoint and renders the index.html template with the list of artists.
	http.HandleFunc("/artists", artistsHandler)            // Redirects to the root endpoint
	http.HandleFunc("/artists.html", artistDetailsHandler) // Handles the endpoint for displaying the details of a specific artist using the artists.html template.
	http.HandleFunc("/locations.html", locationsHandler)   // Handles the endpoint for displaying the locations of a specific artist using the locations.html template.
	http.HandleFunc("/dates.html", datesHandler)           // Handles the endpoint for displaying the concert dates of a specific artist using the dates.html template.
	http.HandleFunc("/relations.html", relationsHandler)   // endpoint for displaying the relations (dates and locations) of a specific artist using the relations.html
	http.HandleFunc("/404", notFoundHandler)               // Handles the endpoint for displaying a "Page not found" message for invalid routes.

	log.Println("Server is running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Function handles the "Page not found" scenario and responds with a 404 status code.
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Page not found", http.StatusNotFound)
}

/* The homeHandler function retrieves the list of artists from the JSON API, decodes
the JSON response, and renders the index.html template to display the list of artists. */
func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFoundHandler(w, r)
		return
	}

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

// The artistsHandler function is a placeholder handler that redirects to the home page ("/").
func artistsHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

/* Retrieves the details of a specific artist based on the provided ID, decodes the
JSON response, and renders the artists.html template with the artist's information. */
func artistDetailsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}
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
	var artist Artist
	for _, a := range artists {
		if a.ID == id {
			artist = a
			break
		}
	}
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

/* The locationsHandler function retrieves the locations for a specific artist based on the
provided ID, decodes the JSON response, and renders the locations.html template to display them. */
func locationsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var locationResponse LocationResponse
	err = json.NewDecoder(response.Body).Decode(&locationResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var locations []string
	for _, loc := range locationResponse.Index {
		if loc.ID == id {
			locations = loc.Locations
			break
		}
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

/* The datesHandler function retrieves the concert dates for a specific artist based on the
provided ID, decodes the JSON response, and renders the dates.html template to display them. */
func datesHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var concertDates ConcertDate
	err = json.NewDecoder(response.Body).Decode(&concertDates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var dates []string
	for _, dat := range concertDates.Index {
		if dat.ID == id {
			dates = dat.Dates
			break
		}
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

/* Function retrieves the relations (dates and locations) for a specific artist based on the
provided ID, decodes the JSON response, and renders the relations.html template to display them. */
func relationsHandler(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}
	response, err := http.Get("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var relationsData Relations
	err = json.NewDecoder(response.Body).Decode(&relationsData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var datesLocations map[string][]string
	for _, rel := range relationsData.Index {
		if rel.ID == id {
			datesLocations = rel.DatesLocations
			break
		}
	}
	tmpl, err := template.ParseFiles("templates/relations.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, datesLocations)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
