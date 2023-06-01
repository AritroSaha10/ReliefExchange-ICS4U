/*
 * File: ban_user.go
 * -------------
 * This module handles the ban user endpoint in the server. It bans a user
 * from the platform if the request maker has admin privileges.
 * It takes a gin context as a parameter, extracts the token and userToBan from the request,
 * checks if the token is valid and if the user who made the request is an admin.
 * If everything checks out, it calls the BanUser helper function to ban the user with the given userToBan id.
 */
// @authors Joshua Chou,Aritro Saha
package post

import (
    "net/http"
    "relief_exchange_backend/globals"
    "relief_exchange_backend/helpers"

    "github.com/gin-gonic/gin"
    log "github.com/sirupsen/logrus"
)

// BanUser handles the endpoint to ban a user.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It accepts a user's id token and the id of the user to be banned, verifies the token,
// checks if the user performing the ban is an admin, then bans the user if authorized using the banUser function and the checkIfAdmin function
func BanUser(c *gin.Context) {
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
    isAdmin, err := helpers.CheckIfAdmin(token.UID)
    if err != nil {
        log.Error(err.Error())
        c.IndentedJSON(http.StatusForbidden, gin.H{"error": "internal server error"})
        return
    }

    if isAdmin {
        // if sending user is an admin delete all the donations of the user to ban including their data and account
        err = helpers.BanUser(uuidToBan)
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
