// [1]Google Developers. (2021). "reCAPTCHA v2 | Google for Developers," Google Developers [Online].
// Available: https://developers.google.com/recaptcha/docs/display. [Accessed: Day-Month-Year].
package post

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ValidateCAPTCHAToken handles the endpoint to verify a CAPTCHA token.
// Parameters:
//   - c: the gin context, the request and response http.
//
// It sends a request to Google's reCAPTCHA API and returns whether the token is valid.
func ValidateCAPTCHAToken(c *gin.Context) {
	// Define a struct to hold the response body from the reCAPTCHA API.
	var captchaResponseBody struct {
		Success bool `json:"success"`
	}
	// Extract the CAPTCHA token from the request query parameters.
	token := c.Query("token")
	// Send a GET request to the reCAPTCHA API with the secret key and the token.
	// [1]
	resp, err := http.Get("https://www.google.com/recaptcha/api/siteverify?secret=" + os.Getenv("RECAPTCHA_SECRET_KEY") + "&response=" + token)
	// If an error occurred while sending the request, return a 500 Internal Server Error status.
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Ensure the response body is closed after the function returns
	defer resp.Body.Close()
	// If the reCAPTCHA API returned a status code other than 200 OK, return a 500 Internal Server Error status.
	if resp.StatusCode != http.StatusOK {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	// Decode the response body into the captchaResponseBody struct (defined at start of function)
	err = json.NewDecoder(resp.Body).Decode(&captchaResponseBody)
	// If an error occurred while decoding the response body, return a 500 Internal Server Error status.
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Return a 200 OK status and whether the CAPTCHA token was valid.
	c.IndentedJSON(http.StatusOK, gin.H{"human": captchaResponseBody.Success})
}
