// @file main.go initializes the log objects, sentry, and routes used in the backend
// @authors Aritro Saha, Joshua Chou
// @cite [1] Stack Overflow.(2015). "Go Gin framework CORS," Stack Overflow [Online].
// Available: https://stackoverflow.com/questions/29418478/go-gin-framework-cors.  [Accessed: 16-May-2023].
package main

import (
	endpointsGet "relief_exchange_backend/endpoints/get"
	endpointsPost "relief_exchange_backend/endpoints/post"
	globals "relief_exchange_backend/globals"

	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

// init function sets up logging and Firebase connections
func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

	// Initialize Firebase globals
	err := globals.InitializeFirebaseGlobals()
	if err != nil {
		log.Error(err)
	}
}

// main function initializes Firebase, Sentry, Firestore client, Auth client, and
// sets up the server routes.
func main() {
	// Set up Sentry
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://4044f25736934d42862ea077a1283931@o924596.ingest.sentry.io/4505213654073344",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("Error initializing Sentry: %s", err)
	}
	// Initialize web server
	r := gin.Default()

	// Set up CORS middleware for all requests
	//citations.txt: [2]
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Set up all GET endpoints
	r.GET("/donations/list", endpointsGet.GetDonationsList)
	r.GET("/donations/:id", endpointsGet.GetDonationByID)
	r.GET("/users/:id", endpointsGet.GetUserDataByID)
	r.GET("/users/banned", endpointsGet.GetIfBanned)
	r.GET("/users/admin", endpointsGet.GetIfAdmin)

	// Set up all POST endpoints
	r.POST("/confirmCAPTCHA", endpointsPost.ValidateCAPTCHAToken)
	r.POST("/donations/new", endpointsPost.AddDonation)
	r.POST("/users/new", endpointsPost.AddUser)
	r.POST("/users/delete", endpointsPost.DeleteUser)
	r.POST("/users/ban", endpointsPost.BanUser)
	r.POST("/donations/report", endpointsPost.ReportDonation)
	r.POST("/donations/edit", endpointsPost.EditDonation)
	r.POST("/donations/:id/delete", endpointsPost.DeleteDonation)

	// Start the server
	err = r.Run()
	if err != nil {
		log.Error(err)
		return
	}
}
