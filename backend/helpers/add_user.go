// @authors Joshua Chou,Aritro Saha
// @cite "Add Data | Firebase." Google, 2023. [Online].
// Available: https://firebase.google.com/docs/firestore/manage-data/add-data. [Accessed: 20- May- 2023].
// This is a file in the package-"helpers" that contains the AddUser function.
package helpers

import (
    "fmt"
    "relief_exchange_backend/globals"
    "time"

    "cloud.google.com/go/firestore"
    log "github.com/sirupsen/logrus"
)

// addUser adds a new user to Firestore.
// Parameters:
//   - ctx: the context in which the function is invoked.
//   - client: the Firestore client.
//   - userId: the ID of the user to add.
//
// Return values:
//   - error, if any occurred during the operation.
func AddUser(userId string) error {
    // Check if they're already banned
    banned, err := CheckIfBanned(userId)
    if err != nil {
        err = fmt.Errorf("err while checking if banned: %w", err)
        log.Error(err.Error())
        return err
    }
    if banned {
        // If the user is already banned, log and return an error
        err := fmt.Errorf("user is already banned")
        log.Error(err.Error())
        return err
    }
    //Get user data from auth server
    userData, err := globals.AuthClient.GetUser(globals.FirebaseContext, userId)
    if err != nil {
        err = fmt.Errorf("failed getting user data from auth server: %w", err)
        log.Error(err.Error())
        return err
    }
    // Create a new document in Firestore for the user with the provided data
    _, err = globals.FirestoreClient.Doc("users/"+userId).Create(globals.FirebaseContext, map[string]interface{}{
        "display_name":    userData.DisplayName,
        "email":           userData.Email,
        "admin":           false,
        "posts":           []firestore.DocumentRef{}, //the posts made by the user
        "uid":             userId,
        "donations_made":  0,
        "registered_date": time.Unix(userData.UserMetadata.CreationTimestamp/1000, 0),
    })
    if err != nil {
        // Log and return the error if there was a problem creating the user's document
        err = fmt.Errorf("failed creating user data doc: %w", err)
        log.Error(err.Error())
        return err
    }
    // If everything went well, return nil indicating no errors
    return nil
}
