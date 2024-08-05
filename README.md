## Overview
The Groupie Trackers web application is designed to provide users with detailed information about various music artists, including their band members, creation dates, first albums, concert locations, and dates. This application utilizes the Groupie Trackers API to fetch and display data in a user-friendly interface.

## Key Features:
###
    Home Page:
    Displays a list of artists fetched from the Groupie Trackers API.
    Each artist is presented with their image, name, and a link to view more details.
    
    Artist Details Page:
    Provides comprehensive information about a specific artist.
    Includes artist's image, name, band members, creation date, first album, and links to their locations and concert dates.

    Locations Page:
    Displays the locations where a specific artist has performed concerts.
    The data is fetched from the Groupie Trackers API and presented in a list format.

    Concert Dates Page:
    Shows the dates of concerts for a specific artist.
    The concert dates are retrieved from the Groupie Trackers API.

    Relations Page:
    Displays the relationship between an artist's concert locations and dates.
    Provides a comprehensive view of where and when an artist has performed.
    
## Technical Details:
###
    Backend:
    Written in Go, utilizing the net/http package for handling HTTP requests and responses.
    Parses JSON responses from the Groupie Trackers API to extract relevant data.

    Frontend:
    Utilizes HTML templates to render dynamic content on the web pages.
    CSS for styling and creating a visually appealing user interface.
    
    API Integration:
    Fetches data from the Groupie Trackers API endpoints (/api/artists, /api/locations, /api/dates, /api/relation).
    Decodes JSON responses into Go structs for easy manipulation and display.
    
    Endpoints and Handlers:
    1. HomeHandler:
    Endpoint: /
    Fetches and displays a list of artists.
    2. ArtistDetailsHandler:
    Endpoint: /artists.html
    Displays details of a specific artist based on the provided ID.
    3. LocationsHandler:
    Endpoint: /locations.html
    Shows the concert locations for a specific artist.
    4. DatesHandler:
    Endpoint: /dates.html
    Displays concert dates for a specific artist.
    5. RelationsHandler:
    Endpoint: /relations.html
    Presents the relations between an artist's concert dates and locations.
    6. NotFoundHandler:
    Handles invalid routes and displays a "Page not found" message.

## Steps to Run the Project
###
    1. Ensure Go is installed
    2. Clone this repository to your local machine using the following command:
        git clone https://github.com/aramsimonyan1/groupie-tracker.git
        
    3. Navigate to the project directory
    cd <project-directory>    
    
    4. Execute the following command to run the Go application:
    go run main.go

    5. Open your preferred web browser (Chrome is recommended over Edge) and go to the following URL: 
    http://localhost:8080/ 