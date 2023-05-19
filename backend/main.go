package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const SERVICE_ACCOUNT_FILENAME = "ics4u0-project-firebase-key.json"

var firebaseContext context.Context
var firebaseApp *firebase.App
var firestoreClient *firestore.Client
var authClient *auth.Client

type Donation struct {
	ID                string    `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	Location          string    `json:"location"`
	Image             string    `json:"img"`
	CreationTimestamp time.Time `json:"creation_timestamp"` // In UTC
	OwnerId           string    `json:"owner_id"`
	Tags              []string  `json:"tags"`
	Reports           []string  `json:"reports"` // Includes the UIDs of every person who reported it
}

type UserData struct {
	DisplayName           string                   `json:"display_name"`
	Email                 string                   `json:"email"`
	RegistrationTimestamp time.Time                `json:"registered_date"` // In UTC
	Admin                 bool                     `json:"admin"`
	Posts                 []*firestore.DocumentRef `json:"posts"`
	DonationsMade         int64                    `json:"donations_made"`
	UID                   string
}

type DeleteDonationRequestBody struct {
	IDToken string `json:"token"`
}

func getDonationsListEndpoint(c *gin.Context) {
	donations, err := getAllDonations(firebaseContext, firestoreClient)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, donations)
	}
}

func getDonationFromIDEndpoint(c *gin.Context) {
	id := c.Param("id")
	donation, err := getDonationByID(firebaseContext, firestoreClient, id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, donation)
	}
}

func getUserDataFromIDEndpoint(c *gin.Context) {
	id := c.Param("id")
	userData, err := getUserDataByID(firebaseContext, firestoreClient, id)
	if err != nil {
		log.Println(err.Error())
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, userData)
	}
}

func postDonationEndpoint(c *gin.Context) {
	var body struct {
		DonationData Donation `json:"data"`
		IDToken      string   `json:"token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil { //stores request body info into the body varible, so that it matches feild in struct in json format
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //if user not signed in, then will send error
		return
	}

	token, err := authClient.VerifyIDToken(firebaseContext, body.IDToken) //token is for user to verify with the server, after it is decoded, we have access to all feilds
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to make this donation."})
		return
	}

	userUID := token.UID

	docID, err := addDonation(firebaseContext, firestoreClient, body.DonationData, userUID) //create new donation object from struct
	//add to the firestore databse
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusCreated, docID)
	}
}

func deleteDonationEndpoint(c *gin.Context) {
	id := c.Param("id")
	donationRef := firestoreClient.Collection("donations").Doc(id)
	donationData, err := donationRef.Get(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var body DeleteDonationRequestBody
	if err := c.ShouldBindJSON(&body); err != nil { //transfers request body so that feilds match the donation struct
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := authClient.VerifyIDToken(firebaseContext, body.IDToken)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	userUID := token.UID

	if donationData.Data()["owner_id"] != userUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	_, err = donationRef.Delete(context.Background()) // only need the err return value
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Donation deleted successfully"})
}
func addUserEndpoint(c *gin.Context) {
	var body struct {
		UserData
		IDToken string `json:"token"`
	}

	if err := c.ShouldBindJSON(&body.UserData); err != nil { //stores request body info into the body varible, so that it matches feild in struct in json format
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()}) //if user not signed in, then will send error
		return
	}

	token, err := authClient.VerifyIDToken(firebaseContext, body.IDToken) //token is for user to verify with the server, after it is decoded, we have access to all feilds
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create this user"})
		return
	}

	userUID := token.UID

	err = addUser(firebaseContext, firestoreClient, body.UserData, userUID) //create new donation object from struct
	//add to the firestore databse
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User added successfully"})
	}
}
func confirmCAPTCHAToken(c *gin.Context) {
	var captchaResponseBody struct {
		Success bool `json:"success"`
	}

	token := c.Query("token")

	resp, err := http.Get("https://www.google.com/recaptcha/api/siteverify?secret=" + os.Getenv("RECAPTCHA_SECRET_KEY") + "&response=" + token)

	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&captchaResponseBody)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"human": captchaResponseBody.Success})
}

/*
	func banUser(c *gin.context) {
		var body struct {
			UserData
			IDToken string `json:"token"`
		}
		// get sending user token
		token, err := authClient.VerifyIDToken(firebaseContext, body.IDToken) //token is for user to verify with the server, after it is decoded, we have access to all feilds
		if err != nil {
			c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create this user"})
			return
		}
		token.UID
		// get uuid of user to ban
		// if sending user is a admin delete all the donations of the user to ban including their data and account
	}
*/

func reportDonationEndpoint(c *gin.Context) {
	var body struct {
		DonationID string `json:"donation_id"`
		IDToken    string `json:"token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil { // Transfers request body so that fields match the struct
		log.Println(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := authClient.VerifyIDToken(firebaseContext, body.IDToken)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	userUID := token.UID
	err = reportDonation(firebaseContext, firestoreClient, body.DonationID, userUID)
	if err != nil {
		if err.Error() == "User has already sent a report" {
			c.IndentedJSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			log.Println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusAccepted)
}

func main() {
	firebaseContext = context.Background()
	firebaseCreds := option.WithCredentialsFile(SERVICE_ACCOUNT_FILENAME)

	app, err := firebase.NewApp(firebaseContext, nil, firebaseCreds)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}
	firebaseApp = app

	firestoreClient, err = firebaseApp.Firestore(firebaseContext)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v\n", err)
	}

	authClient, err = firebaseApp.Auth(firebaseContext)
	if err != nil {
		log.Fatalf("Error initializing Firebase Auth client: %v\n", err)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin,Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/donations/list", getDonationsListEndpoint)
	r.GET("/donations/:id", getDonationFromIDEndpoint)
	r.GET("/users/:id", getUserDataFromIDEndpoint)
	r.POST("/confirmCAPTCHA", confirmCAPTCHAToken)
	r.POST("/donations/new", postDonationEndpoint)
	r.POST("/users/new", addUserEndpoint)
	r.POST("/donations/report", reportDonationEndpoint)
	// r.POST("/users/ban", banUserEndpoint)
	r.DELETE("/donations/:id", deleteDonationEndpoint)
	err = r.Run()
	if err != nil {
		return
	}
}

// Gets all donations available
func getAllDonations(ctx context.Context, client *firestore.Client) ([]Donation, error) {
	var donations []Donation
	iter := client.Collection("donations").Documents(ctx) //.Documents(ctx) returns a iterator
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err //no data was retrieved-nil, but there was an error -err
		}
		var donation Donation
		err = doc.DataTo(&donation)

		// Override some attributes that don't work with DataTo
		donation.Image = doc.Data()["img"].(string)
		donation.OwnerId = doc.Data()["owner_id"].(string)
		donation.CreationTimestamp = doc.Data()["creation_timestamp"].(time.Time)

		if err != nil {
			return nil, err
		}
		donation.ID = doc.Ref.ID //sets donation struct id to the one in the firebase
		donations = append(donations, donation)
	}
	return donations, nil //nil-data was retrived without any errors
}

func getDonationByID(ctx context.Context, client *firestore.Client, id string) (Donation, error) {
	var donation Donation
	doc, err := client.Collection("donations").Doc(id).Get(ctx) //get a single donation from its id
	if err != nil {
		return donation, err //returns empty donation struct
	}
	err = doc.DataTo(&donation)
	if err != nil {
		return Donation{}, err
	}

	// Override some attributes that don't work with DataTo
	donation.Image = doc.Data()["img"].(string)
	donation.OwnerId = doc.Data()["owner_id"].(string)
	donation.CreationTimestamp = doc.Data()["creation_timestamp"].(time.Time)

	donation.ID = doc.Ref.ID //ID is stored in the Ref feild, so DataTo, does not store id in the donations object
	return donation, nil
}

func getUserDataByID(ctx context.Context, client *firestore.Client, id string) (UserData, error) {
	var userData UserData
	doc, err := client.Collection("users").Doc(id).Get(ctx) // Get a single user from its id
	if err != nil {
		return userData, err //returns empty user struct
	}
	err = doc.DataTo(&userData)
	if err != nil {
		return UserData{}, err
	}

	// Set values that aren't set in the DataTo function
	userData.DisplayName = doc.Data()["display_name"].(string)
	userData.RegistrationTimestamp = doc.Data()["registered_date"].(time.Time)
	userData.DonationsMade = doc.Data()["donations_made"].(int64)

	userData.UID = doc.Ref.ID // ID is stored in the Ref feild, so DataTo, does not store id in the user data object
	return userData, nil
}

func addDonation(ctx context.Context, client *firestore.Client, donation Donation, userId string) (string, error) {
	docRef, _, err := client.Collection("donations").Add(ctx, map[string]interface{}{
		"title":              donation.Title,
		"description":        donation.Description,
		"location":           donation.Location,
		"img":                donation.Image,
		"owner_id":           userId,
		"creation_timestamp": donation.CreationTimestamp,
		"tags":               donation.Tags,
		"reports":            donation.Reports,
	})
	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

func reportDonation(ctx context.Context, client *firestore.Client, donationID string, userUID string) error {
	doc, err := client.Collection("donations").Doc(donationID).Get(ctx) // Get the donation's data
	if err != nil {
		return err
	}

	// Check whether they have already reported this donation
	log.Println(doc.Data()["reports"].([]interface{}))
	currentReports, ok := doc.Data()["reports"].([]string)
	if !ok {
		// Convert the empty interface types to actual strings
		currentReports = make([]string, 0)
		for _, reportRaw := range doc.Data()["reports"].([]interface{}) {
			currentReports = append(currentReports, fmt.Sprintf("%+v", reportRaw))
		}
	}

	// Check whether they've already made a report
	for _, report := range currentReports {
		log.Println(report)
		if report == userUID {
			return errors.New("User has already sent a report")
		}
	}

	// Add their UID to the donation, and update the doc
	newReports := append(currentReports, userUID)
	_, err = client.Collection("donations").Doc(donationID).Update(ctx, []firestore.Update{
		{
			Path:  "reports",
			Value: newReports,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func addUser(ctx context.Context, client *firestore.Client, userdata UserData, userId string) error {
	_, err := client.Doc("userdata/"+userId).Create(ctx, map[string]interface{}{
		"display_name": userdata.DisplayName,
		"email":        userdata.Email,
		"admin":        false,
		"posts":        []firestore.DocumentRef{},
		"id":           userId,
	})
	if err != nil {
		return err
	}
	return nil

}
