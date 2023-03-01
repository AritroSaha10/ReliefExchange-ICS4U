package main
//importing
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
	ID          string `json:"id,omitempty"` //take proprety from json 
	Description string `json:"description,omitempty"`
	Location    string `json:"location,omitempty"`
}
//global varible
var (
	dbClient *db.Client //create real time client pointer
)

func main() {
// (credentials for firebase)
	opt := option.WithCredentialsFile("firebase-credentials.json")
	//create new firebase app with options
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}
		// Initialize Firebase Realtime Database client. 
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
	snapshot, err := donationsRef.Get(context.Background()) //represents data at specific time 
	if err != nil {
		log.Printf("Error retrieving donations from Firebase Realtime Database: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return 
	}

	// Convert snapshot to []Donation and return as JSON response.
	//snapshot is defined as the reference to the donations
	var donations []Donation
	if err := snapshot.ForEach(func(donationSnapshot *db.Ref) error { //donationSnapshot type *db.Ref points to the donation struct in the child nodes of the snapshot
		var donation Donation
		if err := donationSnapshot.Value(&donation); err != nil { //Store value of snapshot (each donation) to the donation varible
			return err
		}
		donations = append(donations, donation)
		return nil
	}); err != nil {
		log.Printf("Error converting snapshot to []Donation: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(donations) //return all donations 
}

func getDonationByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}


	vars := mux.Vars(r) //return map of the URL info
		// Extract id from URL path parameter.
	id := vars["id"]

	// Retrieve donation from the database.
	donationRef := dbClient.NewRef(fmt.Sprintf("donations/%s", id)) //get donation by path.  Sprintf is used to construct path using id. 
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
	if err := donationSnapshot.Unmarshal(&donation); err != nil { //Stores value of snapshot to donation varible by unmarshalling (even if different names/types)
		log.Printf("Error unmarshaling donation snapshot: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	// Return the donation as JSON.
	w.Header().Set("Content-Type", "application/json") //tells client server is sending a JSon response 
	if err := json.NewEncoder(w).Encode(donation); err != nil {
		log.Printf("Error encoding donation as JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

