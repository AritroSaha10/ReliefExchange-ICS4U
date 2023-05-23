package main

import (
	"context"
	"fmt"
	endpointsGet "relief_exchange_backend/endpoints/get"
	endpointsPost "relief_exchange_backend/endpoints/post"
	globals "relief_exchange_backend/globals"
	types "relief_exchange_backend/types"

	"net/http"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

const SERVICE_ACCOUNT_FILENAME = "ics4u0-project-firebase-key.json"

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

// postDonationEndpoint handles the endpoint to post a new donation.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a donation and a user's id token, verifies the token,
// and then uses the addDonation function to add the donation to the database.

func postDonationEndpoint(c *gin.Context) {
	var body struct {
		DonationData types.Donation `json:"data"`
		IDToken      string         `json:"token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.IDToken)
	if err != nil {
		log.Warn("Failed to verify ID token")
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to make this donation."})
		return
	}

	userUID := token.UID
	docID, err := addDonation(globals.FirebaseContext, globals.FirestoreClient, body.DonationData, userUID)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Info("Post donation successful.")
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

	donationRef := globals.FirestoreClient.Collection("donations").Doc(id)
	donationData, err := donationRef.Get(context.Background())
	if err != nil {
		log.Error(err.Error())
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

	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, tokenString)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	userUID := token.UID

	isAdmin, _ := checkIfAdmin(globals.FirebaseContext, globals.FirestoreClient, userUID)
	if donationData.Data()["owner_id"].(string) != userUID && (!isAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	_, err = donationRef.Delete(context.Background()) // only need the err return value
	if err != nil {
		log.Error(err.Error())
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

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Error(err.Error())
		// stores request body info into the body varible, so that it matches feild in struct in json format
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()}) // if user not signed in, then will send error
		return
	}

	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.IDToken) // token is for user to verify with the server, after it is decoded, we have access to all feilds
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create this user"})
		return
	}

	userUID := token.UID

	err = addUser(globals.FirebaseContext, globals.FirestoreClient, userUID) // create new donation object from struct
	// add to the firestore databse
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		log.Info("user added successfully")
		c.IndentedJSON(http.StatusCreated, gin.H{"message": "User added successfully"})
	}
}

// banUserEndpoint handles the endpoint to ban a user.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token and the id of the user to be banned, verifies the token,
// checks if the user performing the ban is an admin, then bans the user if authorized using the banUser function and the checkIfAdmin function
func banUserEndpoint(c *gin.Context) {
	var body struct {
		UserToBan string `json:"userToBan"`
		Token     string `json:"token"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// get sending user token
	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.Token)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to create this user"})
		return
	}

	// get uuid of user to ban
	uuidToBan := body.UserToBan
	log.Info(uuidToBan)
	// Check if the user trying to perform the ban is an admin
	isAdmin, err := checkIfAdmin(globals.FirebaseContext, globals.FirestoreClient, token.UID)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "internal server error"})
		return
	}

	if isAdmin {
		// if sending user is an admin delete all the donations of the user to ban including their data and account
		err = banUser(globals.FirebaseContext, globals.FirestoreClient, uuidToBan)
		if err != nil {
			log.Error(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There was an error processing the ban"})
			return
		}

		c.IndentedJSON(http.StatusOK, gin.H{"status": "User banned successfully"})
	} else {
		log.Error(err.Error())
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

	token, err := globals.AuthClient.VerifyIDToken(globals.FirebaseContext, body.IDToken)
	if err != nil {
		log.Error(err.Error())
		c.IndentedJSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this donation."})
		return
	}

	userUID := token.UID
	err = reportDonation(globals.FirebaseContext, globals.FirestoreClient, body.DonationID, userUID)
	if err != nil {
		log.Error(err.Error())
		if err.Error() == "User has already sent a report" {

			c.IndentedJSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			log.Println(err.Error())
			c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusAccepted)
	log.Info("report successful")
}

// main function initializes Firebase, Sentry, Firestore client, Auth client, and
// sets up the server routes.
func main() {

	globals.FirebaseContext = context.Background()
	firebaseCreds := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS_JSON")))

	err := sentry.Init(sentry.ClientOptions{
		Dsn:              os.Getenv("https://4044f25736934d42862ea077a1283931@o924596.ingest.sentry.io/4505213654073344"),
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("Error initializing Sentry: %s", err)
	}

	globals.FirebaseApp, err = firebase.NewApp(globals.FirebaseContext, nil, firebaseCreds)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
	}

	globals.FirestoreClient, err = globals.FirebaseApp.Firestore(globals.FirebaseContext)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v\n", err)
	}

	globals.AuthClient, err = globals.FirebaseApp.Auth(globals.FirebaseContext)
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

	r.GET("/donations/list", endpointsGet.GetDonationsList)
	r.GET("/donations/:id", endpointsGet.GetDonationByID)
	r.GET("/users/:id", endpointsGet.GetUserDataByID)
	r.GET("/users/banned", endpointsGet.GetIfBanned)

	r.POST("/confirmCAPTCHA", endpointsPost.ValidateCAPTCHAToken)
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
func addDonation(ctx context.Context, client *firestore.Client, donation types.Donation, userId string) (string, error) {
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
	if err != nil {
		err = fmt.Errorf("error while adding donation: %w", err)
		log.Error(err.Error())
		return "", err
	}

	// Get current posts and append new post
	// Get the user's document
	userDoc, err := client.Doc("users/" + userId).Get(ctx)
	if err != nil {
		err = fmt.Errorf("error while getting user document (addDonation): %w", err)
		log.Error(err.Error())
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
		"posts":          posts,
		"donations_made": userDoc.Data()["donations_made"].(int64),
	}, firestore.MergeAll) // mergeall ensures that only the posts feild is changed
	if err != nil {
		err = fmt.Errorf("error while updating user document (addDonation): %w", err)
		log.Error(err.Error())
		return "", err
	}

	log.Info("ID of new donation: %v", docRef.ID)
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
		log.Error(err.Error())
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
			err := fmt.Errorf("user has already sent a report")
			log.Error(err)
			return err
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
		err = fmt.Errorf("failed adding report to donation doc: %w", err)
		log.Error(err.Error())
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
	userData, err := globals.AuthClient.GetUser(ctx, userId)
	if err != nil {
		err = fmt.Errorf("failed getting user data from auth server: %w", err)
		log.Error(err.Error())
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
		err = fmt.Errorf("failed creating user data doc: %w", err)
		log.Error(err.Error())
		return err
	}
	return nil
}

// banUser bans a user by removing their records from Firestore and flagging their UID.
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
		err = fmt.Errorf("failed getting user data: %w", err)
		log.Error(err.Error())
		return err
	}

	// Extract posts field from user data
	rawPosts, ok := userDataDoc.Data()["posts"].([]interface{})
	if !ok {
		err = fmt.Errorf("failed extracting posts field from user data")
		log.Error(err.Error())
		return err
	}

	// Convert each raw post to a *firestore.DocumentRef
	for _, rawPost := range rawPosts {
		postRef, ok := rawPost.(*firestore.DocumentRef)
		if !ok {
			log.Warn("failed converting raw post to *firestore.DocumentRef")
			continue
		}
		if _, err := postRef.Delete(ctx); err != nil {
			log.Warn("failed deleting post: %w", err)
			continue
		}
	}

	// Delete user data
	if _, err := userDataRef.Delete(ctx); err != nil {
		err = fmt.Errorf("failed deleting user data: %w", err)
		log.Error(err.Error())
		return err
	}

	// Add them to the banned list
	banDocRef := client.Doc("config/bans")
	var banDocSnapshot *firestore.DocumentSnapshot
	if banDocSnapshot, err = banDocRef.Get(ctx); err != nil {
		err = fmt.Errorf("failed getting ban list: %w", err)
		log.Error(err.Error())
		return err
	}

	// Get ban list
	var banList []string
	if banList, ok = banDocSnapshot.Data()["users"].([]string); !ok {
		err = fmt.Errorf("could not convert banned users list to []string")
		log.Error(err.Error())
		return err
	}
	banList = append(banList, userId)

	// Update the document with new banned user
	banDocRef.Update(ctx, []firestore.Update{
		{
			Path:  "users",
			Value: banList,
		},
	})

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
		log.Error(err.Error())
		return false, err
	}

	// Get the data and access the 'isAdmin' field
	data := doc.Data()
	isAdmin, ok := data["admin"].(bool)
	if !ok {
		err = fmt.Errorf("failed getting admin field from user doc: %w", err)
		log.Error(err.Error())
		return false, err
	}

	log.Info("isAdmin: %v", isAdmin)
	return isAdmin, nil
}
