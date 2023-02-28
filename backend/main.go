package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/db"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

type Donation struct {
	ID          string `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}

var (
	dbClient *db.Client
)

func main() {
	// Initialize Firebase Realtime Database client.
	opt := option.WithCredentialsFile("firebase-credentials.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}
	dbClient, err = app.DatabaseWithURL(context.Background(), "https://YOUR_PROJECT_ID.firebaseio.com/")
	if err != nil {
		log.Fatalf("Error initializing Firebase Realtime Database client: %v\n", err)
	}

	// Initialize router and routes.
	r := mux.NewRouter()
	r.HandleFunc("/donations", getDonations).Methods("GET")
	r.HandleFunc("/donations/{id}", getDonationByID).Methods("GET")
	r.HandleFunc("/donations", createDonation).Methods("POST") //write function for this 
	r.HandleFunc("/donations/{id}", deleteDonation).Methods("DELETE") //write function for this too

	// Serve static files from the "build" directory.
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./build/")))

	// Start server.
	log.Println("Starting server on port 4000...")
	log.Fatal(http.ListenAndServe(":4000", r))
}

func getDonations(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Retrieve all donations from the database.
	donationsRef := dbClient.NewRef("donations")
	snapshot, err := donationsRef.Get(context.Background())
	if err != nil {
		log.Printf("Error retrieving donations from Firebase Realtime Database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Convert snapshot to []Donation and return as JSON response.
	var donations []Donation
	if err := snapshot.ForEach(func(donationSnapshot *db.Ref) error {
		var donation Donation
		if err := donationSnapshot.Value(&donation); err != nil {
			return err
		}
		donations = append(donations, donation)
		return nil
	}); err != nil {
		log.Printf("Error converting snapshot to []Donation: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(donations)
}

func getDonationByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract id from URL path parameter.
	vars := mux.Vars(r)
	id := vars["id"]

	// Retrieve donation from the database.
	donationRef := dbClient.NewRef(fmt.Sprintf("donations/%s", id))
	donationSnapshot, err := donationRef.Get(context.Background())
	if err != nil {
		if err == db.ErrNotFound {
			http.NotFound(w, r)
		} else {
			log.Printf("Error retrieving donation from Firebase Realtime Database: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Parse the donation snapshot into a Donation struct.
	var donation Donation
	if err := donationSnapshot.Unmarshal(&donation); err != nil {
		log.Printf("Error unmarshaling donation snapshot: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Return the donation as JSON.
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(donation); err != nil {
		log.Printf("Error encoding donation as JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// package main

// import (
//     "encoding/json"
//     "fmt"
//     "log"
//     "net/http"
// 	// firebase "firebase.google.com/go"
// 	// "firebase.google.com/go/db"
// 	// "google.golang.org/api/option"
// )

// type Donation struct {
//     ID          int    `json:"id"`
//     Title       string `json:"title"`
//     Description string `json:"description"`
//     Location    string `json:"location"`
//     Date        string `json:"date"`
// }

// var donations = []Donation{
//     {ID: 1, Title: "Bookshelf", Description: "Brown wooden bookshelf", Location: "Los Angeles", Date: "2023-02-16"},
//     {ID: 2, Title: "Table", Description: "Glass top dining table", Location: "San Francisco", Date: "2023-02-18"},
//     {ID: 3, Title: "Sofa", Description: "Red velvet sofa with cushions", Location: "New York City", Date: "2023-02-20"},
// }

// func main() {
//     http.HandleFunc("/donations/", getDonationByID)
// 	http.HandleFunc("/donations", getAllDonations)
//     log.Fatal(http.ListenAndServe(":8080", nil)) //start listening on port and if any erros log them and exit
// }
// func getAllDonations(w http.ResponseWriter, r *http.Request) {
//     json.NewEncoder(w).Encode(donations)
// }

// func getDonationByID(w http.ResponseWriter, r *http.Request) { //Request contains information about the link 
//     id := r.URL.Path[len("/donations/"):] //extract id from the link 
// 	//donations is the object literal defined above 
//     for _, donation := range donations { //loop through donations to check if it matches the link
//         if fmt.Sprintf("%d", donation.ID) == id { //sprinf is used to convert the id (which is an int) to a string so it can be compared with the id varible 
//             json.NewEncoder(w).Encode(donation)
//             return
//         }
//     }
//     http.NotFound(w, r)
// }

