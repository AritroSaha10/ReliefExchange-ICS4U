package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const SERVICE_ACCOUNT_FILENAME = "ics4u0-project-firebase-key.json"

// firebaseContext and firebaseApp are global variables for firebase context and application.
// firestoreClient and authClient are clients for firestore and authentication respectively.
var (
	firebaseContext context.Context
	firebaseApp     *firebase.App
	firestoreClient *firestore.Client
	authClient      *auth.Client
)

// Donation represents a donation item.
// It includes information about the item like title, description, location, image,
// creation timestamp, owner's id, tags, and reports.
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

// UserData represents a user's data.
// It includes display name, email, registration timestamp, admin status, user's posts,
// UID and count of donations made.
type UserData struct {
	DisplayName           string                   `json:"display_name"`
	Email                 string                   `json:"email"`
	RegistrationTimestamp time.Time                `json:"registered_date"` // In UTC
	Admin                 bool                     `json:"admin"`
	Posts                 []*firestore.DocumentRef `json:"posts"`
	UID                   string                   `json:"uid"`
	DonationsMade         int64                    `json:"donations_made"`
}

// getDonationsListEndpoint handles the endpoint to fetch all donations.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends a list of all donations in the database to the client.
func getDonationsListEndpoint(c *gin.Context) {
	donations, err := getAllDonations(firebaseContext, firestoreClient)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, donations)
	}
}

// getDonationFromIDEndpoint handles the endpoint to fetch a donation by id using the getDonationById function
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends the requested donation to the client.
func getDonationFromIDEndpoint(c *gin.Context) {
	id := c.Param("id")
	donation, err := getDonationByID(firebaseContext, firestoreClient, id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, donation)
	}
}

// getUserDataFromIDEndpoint handles the endpoint to fetch a user's data by id using the getUserDataById Function
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends the requested user's data to the client.
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

// postDonationEndpoint handles the endpoint to post a new donation.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a donation and a user's id token, verifies the token,
// and then uses the addDonation function to add the donation to the database.
func postDonationEndpoint(c *gin.Context) {
	var body struct {
		DonationData Donation `json:"data"`
		IDToken      string   `json:"token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil { // stores request body info into the body varible, so that it matches feild in struct in json format
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // if user not signed in, then will send error
		return
	}

	token, err := authClient.VerifyIDToken(firebaseContext, body.IDToken) // token is for user to verify with the server, after it is decoded, we have access to all feilds
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to make this donation."})
		return
	}

	userUID := token.UID

	docID, err := addDonation(firebaseContext, firestoreClient, body.DonationData, userUID) // create new donation object from struct
	// add to the firestore databse
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusCreated, docID)
	}
}

// deleteDonationEndpoint handles the endpoint to delete a donation by id.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It requires an Authorization header with a bearer token, verifies the token,
// checks if the user is authorized to delete the donation, then deletes it.
func deleteDonationEndpoint(c *gin.Context) {
	id := c.Param("id")

	donationRef := firestoreClient.Collection("donations").Doc(id)
	donationData, err := donationRef.Get(context.Background())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "No authorization header provided"})
		return
	}

	splitToken := strings.Split(authHeader, "Bearer")

	if len(splitToken) != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Incorrect format for authorization header"})
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := authClient.VerifyIDToken(firebaseContext, tokenString)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	userUID := token.UID

	isAdmin, _ := checkIfAdmin(firebaseContext, firestoreClient, userUID)
	if donationData.Data()["owner_id"].(string) != userUID && (!isAdmin) {
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

// addUserEndpoint handles the endpoint to add a new user using the addUser function.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token, verifies the token, and then adds the user to the database.
func addUserEndpoint(c *gin.Context) {
	var body struct {
		IDToken string `json:"token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil { // stores request body info into the body varible, so that it matches feild in struct in json format
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // if user not signed in, then will send error
		return
	}

	token, err := authClient.VerifyIDToken(firebaseContext, body.IDToken) // token is for user to verify with the server, after it is decoded, we have access to all feilds
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create this user"})
		return
	}

	userUID := token.UID

	err = addUser(firebaseContext, firestoreClient, userUID) // create new donation object from struct
	// add to the firestore databse
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User added successfully"})
	}
}

// confirmCAPTCHAToken handles the endpoint to verify a CAPTCHA token.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends a request to Google's reCAPTCHA API and returns whether the token is valid.
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

// banUserEndpoint handles the endpoint to ban a user.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token and the id of the user to be banned, verifies the token,
// checks if the user performing the ban is an admin, then bans the user if authorized using the banUser function and the checkIfAdmin function
func banUserEndpoint(c *gin.Context) {
	var body struct {
		UserData struct { // Define the UserData structure or replace it with your actual structure
			UUID string `json:"uuid"`
		}
		IDToken string `json:"token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get sending user token
	token, err := authClient.VerifyIDToken(firebaseContext, body.IDToken)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create this user"})
		return
	}

	// get uuid of user to ban
	uuidToBan := body.UserData.UUID

	// Check if the user trying to perform the ban is an admin
	isAdmin, err := checkIfAdmin(firebaseContext, firestoreClient, token.UID)
	if err != nil {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "internal server error"})
		return
	}

	if isAdmin {
		// if sending user is an admin delete all the donations of the user to ban including their data and account
		err = banUser(firebaseContext, firestoreClient, uuidToBan)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There was an error processing the ban"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"status": "User banned successfully"})
	} else {
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to ban this user"})
	}
}

// reportDonationEndpoint handles the endpoint to report a donation.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token and a donation id, verifies the token,
// then adds the user's report to the donation.
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

// main function initializes Firebase, Sentry, Firestore client, Auth client, and
// sets up the server routes.
func main() {
	firebaseContext = context.Background()
	firebaseCreds := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS_JSON")))

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("https://4044f25736934d42862ea077a1283931@o924596.ingest.sentry.io/4505213654073344"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("Error initializing Sentry: %s", err)
	}

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
		AllowHeaders:     []string{"*"},
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
	r.POST("/users/ban", banUserEndpoint)
	r.POST("/donations/report", reportDonationEndpoint)
	r.POST("/donations/:id/delete", deleteDonationEndpoint)
	err = r.Run()
	if err != nil {
		return
	}
}

// getAllDonations retrieves all donation records from Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
// Return values:
//   - Slice of all Donation objects retrieved.
//   - error, if any occurred during retrieval.

func getAllDonations(ctx context.Context, client *firestore.Client) ([]Donation, error) {
	var donations []Donation
	iter := client.Collection("donations").Documents(ctx) //.Documents(ctx) returns a iterator
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err // no data was retrieved-nil, but there was an error -err
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
		donation.ID = doc.Ref.ID // sets donation struct id to the one in the firebase
		donations = append(donations, donation)
	}
	return donations, nil // nil-data was retrived without any errors
}

// getDonationByID retrieves a donation record by its ID from Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - id: the ID of the donation to retrieve.
//
// Return values:
//   - Donation object that corresponds to the provided ID.
//   - error, if any occurred during retrieval.
func getDonationByID(ctx context.Context, client *firestore.Client, id string) (Donation, error) {
	var donation Donation
	doc, err := client.Collection("donations").Doc(id).Get(ctx) // get a single donation from its id
	if err != nil {
		return donation, err // returns empty donation struct
	}
	err = doc.DataTo(&donation)
	if err != nil {
		return Donation{}, err
	}

	// Override some attributes that don't work with DataTo
	data := doc.Data()
	donation.Image = data["img"].(string)
	donation.OwnerId = data["owner_id"].(string)
	donation.CreationTimestamp = data["creation_timestamp"].(time.Time)

	// Convert the empty interface types to actual strings
	donation.Reports = make([]string, 0)
	for _, reportRaw := range data["reports"].([]interface{}) {
		donation.Reports = append(donation.Reports, fmt.Sprintf("%+v", reportRaw))
	}

	donation.ID = doc.Ref.ID // ID is stored in the Ref feild, so DataTo, does not store id in the donations object
	return donation, nil
}

// getUserDataByID retrieves user data by the user's ID from Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - id: the ID of the user to retrieve.
//
// Return values:
//   - UserData object that corresponds to the provided ID.
//   - error, if any occurred during retrieval.
func getUserDataByID(ctx context.Context, client *firestore.Client, id string) (UserData, error) {
	var userData UserData
	doc, err := client.Collection("users").Doc(id).Get(ctx) // Get a single user from its id
	if err != nil {
		return userData, err // returns empty user struct
	}
	err = doc.DataTo(&userData)
	if err != nil {
		return UserData{}, err
	}

	// Set values that aren't set in the DataTo function
	var ok bool
	userData.DisplayName, ok = doc.Data()["display_name"].(string)
	userData.RegistrationTimestamp, ok = doc.Data()["registered_date"].(time.Time)
	userData.DonationsMade, ok = doc.Data()["donations_made"].(int64)
	if !ok {
		log.Println("WARN: User data may have not been converted properly.")
	}

	userData.UID = doc.Ref.ID // ID is stored in the Ref feild, so DataTo, does not store id in the user data object
	return userData, nil
}

// addDonation adds a new donation record to Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - donation: the Donation object to add.
//   - userId: the ID of the user making the donation.
//
// Return values:
//   - ID of the new donation record.
//   - error, if any occurred during the operation.
func addDonation(ctx context.Context, client *firestore.Client, donation Donation, userId string) (string, error) {
	docRef, _, err := client.Collection("donations").Add(ctx, map[string]interface{}{
		"title":              donation.Title,
		"description":        donation.Description,
		"location":           donation.Location,
		"img":                donation.Image,
		"owner_id":           userId,
		"creation_timestamp": donation.CreationTimestamp,
		"tags":               donation.Tags,
		"reports":            make([]string, 0),
	})
	// Get current posts and append new post
	// Get the user's document
	userDoc, err := client.Doc("users/" + userId).Get(ctx)
	if err != nil {
		return "", err
	}

	// Extract posts array
	data := userDoc.Data()
	rawPosts := data["posts"].([]interface{})

	// Iterate over the array field to find document references
	var posts []*firestore.DocumentRef
	for _, value := range rawPosts {
		posts = append(posts, value.(*firestore.DocumentRef))
	}
	// Append the new donation reference
	posts = append(posts, docRef)

	// Update the user document
	_, err = client.Doc("users/"+userId).Set(ctx, map[string]interface{}{
		"posts": posts,
	}, firestore.MergeAll) // mergeall ensures that only the posts feild is changed
	if err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}
	return docRef.ID, nil
}

// reportDonation adds a report to a specific donation record.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - donationID: the ID of the donation to report.
//   - userUID: the UID of the user making the report.
//
// Return values:
//   - error, if any occurred during the operation.
func reportDonation(ctx context.Context, client *firestore.Client, donationID string, userUID string) error {
	doc, err := client.Collection("donations").Doc(donationID).Get(ctx) // Get the donation's data
	if err != nil {
		return err
	}

	// Check whether they have already reported this donation
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

// addUser adds a new user to Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - userId: the ID of the user to add.
//
// Return values:
//   - error, if any occurred during the operation.
func addUser(ctx context.Context, client *firestore.Client, userId string) error {
	userData, err := authClient.GetUser(ctx, userId)
	if err != nil {
		return err
	}

	_, err = client.Doc("users/"+userId).Create(ctx, map[string]interface{}{
		"display_name":    userData.DisplayName,
		"email":           userData.Email,
		"admin":           false,
		"posts":           []firestore.DocumentRef{},
		"uid":             userId,
		"donations_made":  0,
		"registered_date": time.Unix(userData.UserMetadata.CreationTimestamp/1000, 0),
	})
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

// banUser bans a user by removing their records from Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - userId: the ID of the user to ban.
//
// Return values:
//   - error, if any occurred during the operation.
func banUser(ctx context.Context, client *firestore.Client, userId string) error {
	// Get user data document reference
	userDataRef := client.Doc("users/" + userId)

	// Get user data
	userDataDoc, err := userDataRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed getting user data: %w", err)
	}

	// get the document references from the user data
	var posts []firestore.DocumentRef
	if err := userDataDoc.DataTo(&posts); err != nil {
		return fmt.Errorf("failed getting posts from user data: %w", err)
	}

	// Delete each post
	for _, postRef := range posts {
		if _, err := postRef.Delete(ctx); err != nil {
			return fmt.Errorf("failed deleting post %v: %w", postRef.ID, err)
		}
	}

	// Delete user data
	if _, err := userDataRef.Delete(ctx); err != nil {
		return fmt.Errorf("failed deleting user data: %w", err)
	}

	return nil
}

// checkIfAdmin checks if a user has admin privileges.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - senderId: the ID of the user to check.
//
// Return values:
//   - true if the user has admin privileges, false otherwise.
//   - error, if any occurred during the check.
func checkIfAdmin(ctx context.Context, client *firestore.Client, senderId string) (bool, error) {
	// Get the user document
	doc, err := client.Doc("users/" + senderId).Get(ctx)
	if err != nil {
		return false, err
	}

	// Get the data and access the 'isAdmin' field
	data := doc.Data()
	isAdmin, ok := data["isAdmin"].(bool)
	if !ok {
		return false, fmt.Errorf("isAdmin field not found or is not a bool")
	}

	return isAdmin, nil
}
