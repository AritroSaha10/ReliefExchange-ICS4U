package main
//importing
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
//intall go get firebase.google.com/go/v4  go get github.com/gin-gonic/gin

	"firebase.google.com/go/v4"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)
//take proprety from json 
type Donation struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Location    string `json:"location"`
}


func main() {

	ctx :=context.Background()
	//allow go code to interact with firebase safetly
	serviceAccountKeyFile:="path/to/our-service-account-key.json"
	opt := option.WithCredentialsFile(serviceAccountKeyFile)
	//create new firebase app with options
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Error initializing Firebase app: %v\n", err)
		os.Exit(1)
	}
		// Initialize firestore client
	client, err := app.firestore(ctx)
	if err != nil {
		log.Fatalf("Error initializing Firestore client: %v\n", err)
		os.Exit(1)
	}

	// Initialize gin
	r := gin.Default()
	//retrieves donations from getAllDonations() ([]Donation,error) and sends a JSON response to the fronend with a http statusOK and a response body as a donations slice
	r.GET("/donations/donationList",func(c *gin.Context) //passing in the req res objecs
	{
		donations,err :=getAllDonations(ctx,client)
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		else{
			c.JSON(http.StatusOK,donations)
		}
	})
	r.GET("/donations/:id",func(c *gin.Context){
		id=c.param("id")
		donation,err=getDonationByID(ctx,client,id)
		if err!=nil{
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		else{
			return c.JSON(http.StatusOK,donation)
		}

	})
	
	r.run(":4000")
}
//go function can return 2 things, in this case it is a donation struct and the error
func getAllDonations(ctx context.Context,client *firestore.Client) ([]Donation, error) {
	var donations []Donation
	iter :=client.Collection("donations").Documents(ctx) //.Documents(ctx) returns a iterator
	for {
		doc,err:=iter.Next()
		if err==iterator.Done{
			break
		}
		if err!=nil
		{
			return nil,err //no data was retrieved-nil, but there was an error -err
		}
		var donation Donation
		doc.DataTo(&donation)
		donation.ID=doc.Ref.ID //sets new proprety called id to the one in the firebase
		donations=append(donations,donation)
	}
	return donations,nil //nil-data was retrived without any errors
}

func getDonationByID(ctx context.Context, client *firestore.Client,id string) (Donation,error) {
	var donation Donation
	doc,err:=client.Collection("donations").Doc(id).Get(ctx) //get a single donation from its id
	if err!=nil
	{
		return donation,err //returns empty donation struct
	}
	doc.DataTo(&donation)
	donation.ID=doc.Ref.ID
	return donation,nil
}

